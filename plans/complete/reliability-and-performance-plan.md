# Remediation Plan: Reliability & Performance

## Root Cause Analysis

Three interrelated problems surfaced:

| # | Problem | Root Cause |
|---|---------|------------|
| 1 | **Backend unavailable** | No process supervision — API runs as a bare `&` background job in `entrypoint.sh`. If the Go binary panics/OOMs/exits, nginx keeps running, Docker sees a "healthy" container, but the API is dead. No restart mechanism. |
| 2 | **Slow thumbnail pop-in** | Thumbnails generated on-demand (first request pays full cost) + N+1 HTTP request pattern (100 photos = 100 individual blob fetches) + no server-side cache headers on the `/api/photo/thumbnail/` path means zero browser/nginx caching + no placeholder/LQIP strategy. |
| 3 | **Optional cache not present** | Redis/Valkey is `CACHE_ENABLED`-gated. Without it, every cold thumbnail request re-reads the filesystem or re-generates entirely. |

---

## Webhook vs WebSocket

- **Webhook** is server→external-service. Useful for notifying Discord/Slack/other APIs when photos are indexed. Does **not** help the frontend or address either problem above.
- **WebSocket** is server↔client persistent connection. Useful for real-time UI updates. Partially helps problem #2 (push "thumbnail ready" events) and builds infrastructure for future live features. Does **not** fix the process supervision gap (problem #1) or the thumbnail generation bottleneck.

**Recommendation**: A webhook is a nice-to-have for external integrations but solves none of the current problems. A WebSocket is worth building but only *after* the reliability and generation bottlenecks are fixed — otherwise it's a real-time pipe to a flaky backend serving slow thumbnails.

---

## Phase 0: Immediate Stability (hours, not days)

### 0.1 — Process Supervisor
**Problem**: #1. **Effort**: Low. **Risk**: None.

Replace the bare `&` in `entrypoint.sh` with a proper supervisor. The container already uses Alpine.

**Options**:
- A: **s6-overlay** — battle-tested, used in linuxserver.io images, adds ~2MB
- B: **supervisord** — simple, Python-based, already in many Alpine images
- C: **dumb-init + shell watchdog loop** — simplest, zero new dependencies

**Recommended**: Option C for immediacy, migrate to s6-overlay as a follow-up. The shell loop watches the PID and restarts on exit:

```sh
#!/bin/sh
set -e

# Start API with a restart watchdog
while true; do
    /app/api
    echo "[watchdog] API exited, restarting in 2s..." >&2
    sleep 2
done &
API_WATCHDOG_PID=$!

# Cleanup handler
trap "kill $API_WATCHDOG_PID 2>/dev/null" EXIT

# Start nginx in foreground
exec nginx -g "daemon off;"
```

This ensures: API crash → 2s delay → auto-restart. Docker health check (already configured in `docker-compose.local.yml` at `/api/health`) handles the 2s gap.

### 0.2 — API Startup Readiness Gate
**Problem**: #1 (race condition). **Effort**: Low.

In `entrypoint.sh`, add a readiness loop before starting nginx:

```sh
# Wait for API to be ready before starting nginx
for i in $(seq 1 30); do
    if wget -qO- http://localhost:8080/api/health 2>/dev/null; then
        break
    fi
    echo "[entrypoint] waiting for API... ($i/30)" >&2
    sleep 1
done
```

This prevents the 502-flood window where nginx proxies to a not-yet-started API.

### 0.3 — React Query Retry Configuration
**Problem**: #1 (resilience). **Effort**: Low.

The current `QueryClient` uses TanStack defaults (3 retries with exponential backoff). Explicitly configure:

```typescript
// webapp/src/... (QueryClient setup)
retry: 3,
retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 30000),
```

This pairs with the auto-restart watchdog: API crashes → clients retry with backoff → API restarts within 2s → requests succeed on retry 2 or 3. Transient blip, not user-visible error.

---

## Phase 1: Structural Reliability (days)

### 1.1 — Decouple API from nginx (Separate Containers)
**Problem**: #1 (architectural). **Effort**: Medium.

The single-container design couples API and nginx lifecycles. Split into two containers in `docker-compose.local.yml`:

```yaml
services:
  api:
    build:
      context: .
      target: api-runtime    # new Dockerfile stage: just the Go binary
    healthcheck:
      test: ["CMD", "wget", "-qO-", "http://localhost:8080/api/health"]

  nginx:
    image: nginx:alpine
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
      - ./webapp/build:/usr/share/nginx/html
    depends_on:
      api:
        condition: service_healthy
```

Benefits:
- Independent health checks & restarts
- `docker compose restart api` doesn't bounce nginx
- Scale API independently if needed
- Nginx `proxy_pass http://api:8080/api/` (Docker DNS) instead of localhost

### 1.2 — Implement `panic` Recovery Middleware
**Problem**: #1 (silent crashes). **Effort**: Low.

Add a Gin recovery middleware that logs the panic, returns 500, and keeps the process alive:

```go
router.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
    // log the panic, return 500, do NOT crash the process
    c.AbortWithStatusJSON(500, gin.H{"error": "internal server error"})
}))
```

Go panics in HTTP handlers are the #1 cause of API process death. This keeps the process alive even when a handler panics.

---

## Phase 2: Thumbnail Performance (days)

### 2.1 — Generate Thumbnails at Index Time, Not Request Time
**Problem**: #2 (on-demand generation). **Effort**: Medium.

Currently: first viewer of a photo triggers thumbnail generation. Change to: thumbnail generated at photo upload/index time, stored alongside the photo. The `GenerateThumbnail` call moves from `GetPhotoThumbnail` handler into the photo ingestion pipeline. After indexing, all thumbnails exist. First viewer gets instant response.

### 2.2 — Batch Thumbnail API
**Problem**: #2 (N+1 requests). **Effort**: Medium.

Add `POST /api/photos/thumbnails` that accepts `{ ids: string[], size: "s"|"m"|"l" }` and returns a multipart response or a JSON payload with base64-encoded thumbnails. This collapses 100 HTTP requests into 1. The frontend `PhotosAdapter` gets a `getThumbnailsBatch(ids)` method. For the photo grid, fetch all visible thumbnails in a single round-trip.

### 2.3 — Nginx Caching for Thumbnails
**Problem**: #2 (no caching). **Effort**: Low.

The nginx config caches static assets but not API-served thumbnails. Add a cache zone:

```nginx
proxy_cache_path /var/cache/nginx/thumbnails levels=1:2 keys_zone=thumbnail_cache:50m max_size=1g inactive=24h;

location /api/photo/thumbnail/ {
    proxy_cache thumbnail_cache;
    proxy_cache_key "$uri$is_args$args";
    proxy_cache_valid 200 24h;
    proxy_cache_use_stale error timeout updating;
    proxy_pass http://api:8080;
}
```

Combined with 2.1 (thumbnails pre-generated), this means the Go API is only hit once per thumbnail per 24h. After that, nginx serves from disk at wire speed.

### 2.4 — LQIP / Blurhash Placeholder Strategy
**Problem**: #2 (nothing shown during load). **Effort**: Medium.

Store a tiny (32px) blurhash or base64-encoded JPEG thumbnail as a database column on `Photo`. When the photo grid renders, show the blurhash immediately, then swap in the full thumbnail when loaded. This eliminates the "empty boxes that slowly populate" UX. Libraries: `blurhash` (encode in Go at index time, decode in React).

---

## Phase 3: Real-Time Infrastructure (week+)

### 3.1 — WebSocket for Live Events
**Problem**: #2 (partial — progressive thumbnail delivery), future-proofing. **Effort**: High.

Add a WebSocket endpoint at `GET /api/ws`. The Go side uses `gorilla/websocket` (or the lighter `nhooyr.io/websocket`). The client side uses a React context provider that manages connection lifecycle (auto-reconnect with exponential backoff, auth via token in initial message).

Event types the server pushes:

| Event | Payload | Trigger |
|-------|---------|---------|
| `thumbnail:ready` | `{ photoId, size, url }` | Thumbnail generation completes (async worker) |
| `index:progress` | `{ current, total, phase }` | Photo indexing job |
| `index:complete` | `{ total, errors[] }` | Photo indexing finishes |
| `photo:new` | `{ photo }` | New photo detected in watch dir |
| `system:health` | `{ version, uptime }` | Periodic heartbeat (replaces HEAD polling in NetworkStatusBanner) |

### 3.2 — Progressive Thumbnail Loading via WebSocket
**Problem**: #2 (subjective performance). **Effort**: Medium (depends on 3.1).

The photo grid initially renders with blurhash placeholders (from 2.4). A single batch request fetches thumbnails for visible photos (from 2.2). If the batch response is large or some thumbnails weren't pre-generated, the server pushes `thumbnail:ready` events as they complete. The `useBlobUrl` hook listens to these events for the keys it cares about.

Result: **First paint is instant** (blurhash), **batch fills most within one RTT**, **stragglers stream in as they're ready**.

---

## Summary Matrix

| Phase | Item | Solves Problem | Effort | Impact | Depends On |
|-------|------|---------------|--------|--------|------------|
| **0.1** | Process watchdog | #1 Backend crashes | 30 min | Critical | — |
| **0.2** | Readiness gate | #1 Startup race | 15 min | Critical | — |
| **0.3** | QueryClient retry config | #1 Transient errors | 15 min | High | — |
| **1.1** | Separate containers | #1 Architectural | 2-4 hrs | High | 0.1, 0.2 |
| **1.2** | Panic recovery middleware | #1 Handler panics | 30 min | Medium | — |
| **2.1** | Index-time thumbnails | #2 First-load latency | 4-8 hrs | Critical | — |
| **2.2** | Batch thumbnail API | #2 N+1 requests | 2-4 hrs | High | 2.1 |
| **2.3** | Nginx thumbnail cache | #2 Repeat loads | 1 hr | High | 2.1 |
| **2.4** | Blurhash placeholders | #2 Perceived performance | 4-8 hrs | Medium | — |
| **3.1** | WebSocket infrastructure | #2 progressive, future | 1-2 weeks | Medium | Phase 0-1 done |
| **3.2** | Progressive loading via WS | #2 straggler UX | 2-4 hrs | Medium | 2.4, 3.1 |

**Recommended execution order**: 0.1 → 0.2 → 0.3 in one sitting (under 2 hours, permanently fixes the "backend unavailable" error). Then 2.1 → 2.2 → 2.3 for the thumbnail performance win. Then 1.1 for architectural hygiene. 2.4 and Phase 3 can follow based on whether the performance is "good enough" after Phase 2.
