# Multi-Tenancy & AuthZ Architecture Plan

## Status: v0.2 — decisions confirmed (2026-09-09). Phases 0–1 implemented; Phase 2 identity layer implemented, query-scoping sweep pending.

Implements GitHub **issue #74** (Multi-tenancy support) and **issue #148** (owner scoping) against the product strategy in `plans/CLOUD_HOSTING_PLAN.md` (open SaaS: "Google Photos convenience, without Google", EU-hosted, open signup with approval tooling, 10GB free quota, B2 object storage, Cloudflare CDN).

### Implementation status (branch `issue-74-multitenancy`)

| Phase | Item | Status |
|---|---|---|
| 0 | golang-migrate adoption (embedded, advisory-locked; `000001_baseline` adopts AutoMigrate DBs without data change) | ✅ shipped |
| 1 | `workspaces` + `workspace_members`; `000002_workspaces` backfills all pre-existing content + users into the default workspace | ✅ shipped |
| 1 | Signup creates an isolated personal workspace; admin-created users join the shared default workspace | ✅ shipped |
| 1 | `thumb_cap` capability column + unauth `/t/{cap}/{size}.webp` route (immutable caching) | ✅ shipped |
| 1 | `ObjectStore` + workspace-bound `WorkspaceStore` (structural prefix isolation) | ✅ shipped |
| 1 | S3 originals adapter (minio-go; Wasabi/AWS/MinIO) behind `ORIGINALS_STORAGE` | ✅ shipped |
| 2 | Workspace CRUD + membership API; `workspaceMiddleware` (server-side resolution, fail-closed) | ✅ shipped |
| 2 | **Query-scoping sweep**: every repository query filtered by workspace | ⬜ pending |
| 2 | HTTP upload route; originals/regen cutover to S3; cache-key namespacing | ⬜ pending |
| 2 | Webapp `X-Workspace-ID` header + `<img src>` capability thumbnails | ⬜ pending |
| 2 | Cross-tenant content matrix (403/404 on every endpoint) | ⬜ pending |
| 3 | Per-workspace cron/WS; orphan sweep; backup key split | ⬜ pending |
| 4 | Postgres RLS hardening | ⬜ pending |
| 5 | Multi-member product UI, invitations, paid tiers | ⬜ pending |

**Enforcement boundary:** the workspace *identity* layer is complete and tested, but content queries are not yet filtered by workspace. Until the Phase 2 sweep lands, a caller cannot act outside a workspace they belong to, but a query inside their own workspace can still return rows belonging to another workspace. Do not enable open signup until the sweep and its cross-tenant matrix are green.

Supersedes the tenancy sketches in `plans/ideas-implementation-plan.md` §1 and `plans/CLOUD_HOSTING_PLAN.md` §7. Design direction from Rich: **authN/authZ built through the product so users only ever access their own media, with credentials applied at the resource layer** (the AWS Cognito pattern: identity pool + scoped credentials against `s3:prefix/{customer-id}/*`).

---

## 1. Codebase reality check (verified against source)

Issue #74's premises are stale in ways that change the design — mostly for the better:

| Issue #74 assumes | Reality (verified) |
|---|---|
| "Maps well to the existing album/photo ownership model" | **There is no ownership model.** `photos`, `albums`, `shared_links`, `jobs`, `photo_analysis`, `photo_tags`, `settings` have no ownership column of any kind. All content is global. |
| "Add workspace claims to JWT" | Auth is **PASETO v2.local** (symmetric, encrypted), not JWT. `TokenPayload` carries `UserID`/`Username`/`Role` only (`api/internal/core/domain/tokenPayload.go`). |
| Implies uploads exist | **There is no HTTP upload route.** Ingestion is filesystem-scan only; the webapp's legacy base64 `POST /photo` 404s. The SaaS upload path is net-new. |
| (unmentioned) | **Photo IDs are `base64(filesystem_path)`** (`api/internal/core/service/photo.go:289`) — path-derived and guessable, not UUIDs. Breaks the naive "UUID-as-capability" idea in `CLOUD_HOSTING_PLAN.md` §13.4; see §6. |
| (unmentioned) | Webapp fetches thumbnails via **authenticated XHR → blob URLs** (`webapp/src/Components/PhotoGrid/PhotoGrid.tsx`) — no `<img src>`, so zero CDN caching is possible today. |
| (unmentioned) | **WebSocket hub broadcasts every event to every client** (`api/internal/components/websocket/hub.go:44-53`) — cross-tenant leak of photo IDs (= filenames) the moment tenancy exists. |
| (unmentioned) | **Live security bugs:** share listing/revocation IDOR (`ListShares`/`RevokeShare` are global), share passwords stored **plaintext** (`api/internal/core/service/share.go`), `GET /api/album/:id` is a no-auth public endpoint. All become cross-tenant breaches under multitenancy. |
| (unmentioned) | **No migration framework** — GORM `AutoMigrate` only (`api/internal/adapter/storage/database/database.go:62`). The required backfill (nullable column → UPDATE → NOT NULL → constraint swaps) is beyond AutoMigrate. |
| (unmentioned) | `albums.name` has a **global unique index** — two tenants can't both have "Holidays". Cache keys (`thumbnail:{photoId}:{size}`) have no tenant namespace. Search/tags are global. Cron jobs are global. Storage layout is flat (`/photos/<album>/`, `{s,m,l}/{photoId}.webp`). |

The good news: no legacy ownership model means no reconciliation — the scoping pass is a clean addition, and the repository layer is centralized (`api/internal/adapter/storage/database/repository/`).

---

## 2. Design principles

1. **Two credential layers (the Cognito translation).** Identity layer: PASETO authN + server-side workspace membership authZ. Resource layer: credentials fused to the resource — capability URLs for reads at scale, prefix-bound storage access for writes/originals, per-process scoped storage keys for operator surfaces.
2. **Authorization must be structural, not just procedural.** Wherever possible, make cross-tenant access impossible by construction (prefix-bound key building, per-request bound accessors), not merely checked by convention.
3. **Fail closed.** Every check defaults to deny: unknown capability → 404 (no existence oracle), no membership → 403, no workspace context → request rejected.
4. **The self-hosted single-tenant deployment never breaks.** Home instance keeps the filesystem adapter, one `default` workspace, unchanged routes. Every phase is independently mergeable and deployable to `develop`.
5. **Membership is data, not schema.** `workspace_id` is the single isolation dimension; whether workspaces are personal or shared is a product decision (still open, `CLOUD_HOSTING_PLAN.md` §10.2) that must not be baked into the schema.
6. **Explicit over implicit where error = breach.** Workspace is a parameter on repository ports (compiler-enforced), transported via context — never read from context deep inside a repository where a forgotten read is invisible.

---

## 3. Layer A — Identity (authN + membership authZ)

Keep PASETO v2.local. **Do not put `workspace_id` in the token** — stale-claim bug: a user removed from a workspace still holds a valid token carrying the claim, and stateless symmetric tokens can't be reliably revoked.

Instead, workspace is resolved **server-side per request**:

1. `authMiddleware` verifies PASETO → user identity (unchanged).
2. New workspace middleware: read `X-Workspace-ID` header (query param for WebSocket upgrade, which can't set headers in browsers) → one indexed `workspace_members` lookup → build `WorkspaceContext{WorkspaceID, Role}` where **Role is the workspace role, orthogonal to the token's global `UserRole`** → 403 if no membership.
3. Handlers pass `WorkspaceContext` into services → repositories.

Properties: membership removal takes effect immediately; workspace switching is a header change (no token reissue); the lookup is a microsecond PK read — don't cache it.

**Roles (confirmed 2026-09-09): four** — `owner`, `admin`, `member`, `viewer` (per issue #74). Membership row stores one role; workspace owner may manage members. Enforcement in Phase 2: `viewer` = read-only, `member` = read + upload/tag, `admin` = full manage (add/remove members, edit workspace), `owner` = superset (transfer/delete workspace). Global `ADMINISTRATOR` remains a deployment-level (ops) concept, not a workspace role.

---

## 4. Layer B — Resource credentials (authZ at the resource layer)

### 4.1 Mechanism per resource class

| Resource | Mechanism | Rationale |
|---|---|---|
| Metadata, mutations, uploads, jobs | **API-mediated**, workspace-scoped at repository ports | Strong authZ, low volume |
| **Thumbnails (~80% of reads)** | **Capability URL**: `GET /t/{thumb_cap}/{size}.webp`, unauth, `Cache-Control: immutable` | The only model that unlocks Cloudflare edge caching (§6); simplifies frontend to `<img src>` |
| Originals / full-res / ZIP download | **API-proxied streams** at launch; presigned 302 redirect as a config-flip escape hatch | Low volume; revocation immediate; nothing cacheable leaks to a shared CDN; netcup traffic is included (⚠️ verify fair-use) |
| Chunked upload (#98, Phase 5) | **Presigned multipart direct-to-B2** | Canonical Cognito pattern; deferred — at launch, API-mediated streaming multipart behind the 50MB body cap is simpler |
| Share links | Share token **is** a capability already — fix plaintext passwords + IDOR scoping | Phase 0 blocker |
| Operator processes | **3-key scheme**: API key (full bucket), backup key (read-only, pull-only per `CLOUD_HOSTING_PLAN.md` §6.2), optional sweep key (`originals/*` only) | Cognito principle applied where processes are genuinely separable |

### 4.2 Why per-tenant storage keys are rejected

B2 application keys support `namePrefix` restriction; MinIO policies and S3 condition policies do too — the literal Cognito analog is technically available. But the API process would hold all N per-tenant keys in one `SealedSecret`, so an app-layer authZ bug in that process still reaches every tenant's objects: key lifecycle, rotation, and re-provisioning infrastructure for zero real isolation. The structural guarantee only materializes when *separate processes* hold separate keys — hence the 3-key process scheme. The prefix layout (§4.4) keeps per-tenant keys a drop-in later if a multi-service architecture emerges. (⚠️ B2's ~1,000 application-key limit per account would also cap tenants; moot if rejected.)

### 4.3 `WorkspaceStore` — structural prefix binding

The failure mode to design against: a handler resolving `photo_id → object key` without going through workspace binding. Kill it by construction:

1. One process-level `ObjectStore` holds the full-access client and exposes **no public key-building API**.
2. `WorkspaceStore` is created per request in middleware, *after* membership resolution: `NewWorkspaceStore(store, workspaceID)` with a validated UUID that exists in `workspaces` and whose membership was checked.
3. Every method (`Get`, `Put`, `Delete`, `Presign`, `ListForSweep`) builds keys exclusively inside the struct: `originals/{wsID}/{photoID}`. `photoID` is boundary-validated (UUID regex — no `/`, no `..`, no escape).
4. Services receive the already-bound `WorkspaceStore` from request scope — never a free `workspace_id` string they could misuse. Cross-tenant access requires deliberately resolving a second workspace through the membership check: one auditable path.

This fuses the credential to the accessor before any byte is touched — the in-process analog of Cognito's scoped STS credentials. Cheap, unit-testable, zero key-provisioning infrastructure.

### 4.4 Storage layout (S3 adapter, SaaS)

```
originals/{workspace_id}/{photo_uuid}                 # UUIDv4 on SaaS uploads
thumbnails/{workspace_id}/{photo_uuid}/{s|m|l}.webp
trash/{workspace_id}/{photo_uuid}                     # copy+delete semantics (§6.5 of cloud plan)
edited/{workspace_id}/{photo_uuid}/{edit_params_hash}
```

- Home instance: `filesystem` originals adapter untouched; path-hash photo IDs retained (they exist for filesystem-scan dedup, meaningless for uploads). SaaS: `ORIGINALS_STORAGE=s3`, UUIDv4 IDs — also removes the `[1:]`-strip / unescape hacks in `http/photo.go`.
- Thumbnail cutover: **one-time regen** (the regen job exists and is already the documented restore path, cloud plan §6.6) rather than permanent dual-lookup logic.
- Decide the S3 key pattern now (this section) so the layout is stable; the S3 originals adapter ships in Phase 1.

---

## 5. Layer A data model changes (Phase 1 schema)

New tables:

```sql
workspaces (id uuid PK, name, slug, owner_user_id FK, created_at,
            storage_used_bytes bigint, storage_limit_bytes bigint, settings jsonb)
workspace_members (workspace_id FK, user_id FK, role text,  -- owner|member
                   PRIMARY KEY (workspace_id, user_id))
```

Columns added (nullable → backfill → `SET NOT NULL`):

| Table | Add | Notes |
|---|---|---|
| `photos` | `workspace_id` FK, `thumb_cap` (random 128-bit, immutable, unique) | composite indexes `(workspace_id, created_epoch)` etc. in the same migration |
| `albums` | `workspace_id` FK | **drop global unique on `name`** → composite unique `(workspace_id, name)` |
| `shared_links` | `workspace_id` FK | fixes ListShares/RevokeShare scoping |
| `jobs` | `workspace_id` | PK becomes composite `(workspace_id, name)` |
| `photo_analysis` | `workspace_id` | |
| `photo_tags` | `workspace_id` | or join-through-photos; column is simpler |
| `users` | — | **no** `workspace_id` column on users — membership table only |

**Backfill (one ordered, idempotent golang-migrate migration):** create tables → insert `default` workspace → insert every existing user as a member → add nullable columns → `UPDATE ... SET workspace_id = <default>` → `SET NOT NULL` → composite indexes → album constraint swap. All existing users land in **one default workspace** = today's shared-library semantics, bit-for-bit. Gate: migration up/down tested against a `pg_dump` of real production data; backup before first deploy.

**Migration framework is a Phase 0 prerequisite:** adopt **golang-migrate** with embedded migrations (plain SQL, single-binary deploy preserved), run on startup with an advisory lock (single-replica K3s) or via Helm init container. AutoMigrate cannot express this backfill.

---

## 6. The thumbnail capability model (resolves cloud plan §13.4)

§13.4's caveat: Cloudflare doesn't cache requests carrying `Authorization` — and thumbnails are authenticated today, so CF absorbs 0% of reads. Of its three options:

| Option | Verdict |
|---|---|
| **(a) capability route** | ✅ **adopt — but not as written.** Photo IDs are `base64(path)` = enumerable. The capability must be a **separate random 128-bit `thumb_cap` column**, immutable per photo, and the URL contains only the cap — never the photo ID: `GET /t/{thumb_cap}/{s|m|l}.webp` → `Cache-Control: public, max-age=31536000, immutable`. Unknown cap → 404 (same as missing; no existence oracle). |
| (b) signed URLs | ❌ signature in query string → every new signature is a CF cache miss; TTLs contradict `immutable`; kills the ~80% edge offload it exists to protect. |
| (c) CF Worker header-strip | ❌ requires the PASETO symmetric key inside Cloudflare (US company holding the identity key — against the EU positioning); CF-specific code against the §9 guardrails; doesn't beat (a) on security. |

**Privacy tradeoff (decision D2):** thumbnails (≤600px, 20–50KB) will transiently live in CF's global edge cache. Originals, metadata, and the DB remain strictly EU. Recommendation: accept + disclose in ToS/privacy policy ("thumbnail representations may be served from a global CDN edge; originals and metadata never leave the EU"); keep an EU-only thumbnail-caching config flag (VPS nginx + Valkey chain) as the alternative. Share links already expose thumbnails to non-owners today, so capability access is not a new privacy surface — it's the same surface with better hygiene.

Residual risks to name: caps leak via browser history/logs; `Referrer-Policy: strict-origin-when-cross-origin` (already set) limits referer leakage; share revocation does not rotate caps of shared photos (acceptable — same bytes remain owner-visible; rotation sweep is a Phase 5 option).

CF page rules (cache everything `/t/*`, bypass `/api/*`) land in OpenTofu in the Phase 1 deploy.

---

## 7. AuthZ matrix — end-to-end flows

Vending: **T** = PASETO at login · **C** = capability URL embedded in API responses · **P** = presigned URL per request · **S** = WorkspaceStore prefix binding · **DB** = workspace-scoped query.

| Flow | Credential | Validated at | Fails closed |
|---|---|---|---|
| Login | T | argon2id verify → token issue; `authMiddleware` on every authed route | 401 |
| Upload | T + quota | membership → **quota pre-check inside the photo-insert transaction** (`storage_used_bytes + size ≤ storage_limit_bytes` atomic with row insert) → S via WorkspaceStore | 413/507 before any object write |
| Grid thumbnails | C | unauth route; cap → key under `thumbnails/{ws}/...` | 404, unenumerable (128-bit) |
| Full-res view | T + DB + S | owner-scoped `GetPhoto(photoID, ws)`; key built only from `(ws, uuid)` | 403 cross-tenant |
| ZIP download | T + DB + S | each ID resolved within caller's workspace | 403/404 per photo |
| Shared link (unauth) | share token | token lookup + **argon2 compare** (fix plaintext) + expiry | 404/401/410; revoked = dead URL |
| WebSocket events | T at upgrade + per-workspace hub bucket | membership at upgrade; hub delivers within workspace only | no cross-tenant events |
| Cron (AI, quality, trash, regen, index) | process-level | iterate per workspace; trash retention stays a global deployment setting | per-ws isolation |
| Orphan sweep | sweep key (optional) + per-ws LIST | objects >24h old with no DB row, within one ws prefix | can't touch other prefixes |
| Backup | read-only pull key | B2 key capabilities | backup bug can't write/delete primary |

**Definition of done for the scoping work: a cross-tenant integration matrix** (extend `api/test/integration`, `-tags=integration` harness + `docker-compose.test.yml`) — two tenants A and B; user from A accessing B's resource must get 403/404 (never 200, never the resource) on **every** endpoint: photos, albums, shares, tags, search, thumbnails, bin, trash, edit/rotate.

---

## 8. Signup → first-upload lifecycle

Open signup (approval flag remains an admin tool):

1. `POST /register` — single transaction: `users` row (UUID, argon2id) + `workspaces` row (personal, 10GB default limit, used=0) + `workspace_members` row (role `owner`).
2. Approved-but-never-uploaded user = those three rows. **No B2 objects, no key provisioning** — a workspace is pure DB + prefix namespace. Cost ≈ zero. Deletion = delete user → cascade memberships → prefix sweep reclaims objects.
3. First upload: transactional quota check; object to `originals/{ws}/{uuid}`; over-quota rejected pre-PUT.
4. **Orphan sweep** (upload succeeds, DB insert fails): daily job LISTs per-workspace prefixes (bounded, tenant-fair by construction), deletes objects older than 24h with no matching `photos` row, reconciles `storage_used_bytes` drift. Orphaned objects are unreachable (no DB row → no URL, no cap) — quota leaks until the sweep, data doesn't.

---

## 9. Phased implementation plan

Each phase is independently mergeable to `develop` (PR with ADR / Test Plan / Test Outcome per `plans/IMPLEMENTATION_DECISIONS.md` §5). Validation layers per AGENTS.md: **L1** = `golangci-lint`, `go test -race`, webapp lint/test/build (always) · **L2** = integration suite · **L3** = local compose + chrome-devtools browser verification with screenshots · **L4** = helm lint/template + K3s.

### Phase 0 — Security blockers & migration infrastructure
*Live bugs today; nothing depends on S3.*
- Adopt golang-migrate (embedded, startup runner + advisory lock); migration test harness.
- Fix share listing/revocation IDOR (scope by creator); argon2id-hash share passwords.
- Resolve no-auth `GET /api/album/:id` (delete if unused outside shared flow, else capability-ize — D3).
- WebSocket hub per-workspace routing groundwork (or interim: authenticated-only delivery) — currently broadcasts photo IDs (= filenames) to all clients.
- Signup hardening: email verification, auth rate limits.
- **Gate:** L1 + L2.

### Phase 1 — Schema, storage & credential model
- `workspaces` + `workspace_members`; register creates personal workspace (§8); default-workspace backfill (§5) — app ignores new columns, zero behavior change for the existing deployment.
- S3 originals adapter behind `ORIGINALS_STORAGE=filesystem|s3`; prefix layout (§4.4); `WorkspaceStore` bound accessor + key-construction unit tests.
- `thumb_cap` column + `/t/{cap}/{size}.webp` route with immutable headers.
- Quota columns + transactional enforcement.
- CF cache rules in OpenTofu (`photobox-infra`).
- **Gate:** migration up/down vs prod-copy data; full existing suite green; L1 + L2 + L4; self-hosted deploy behaves identically.

### Phase 2 — Request scoping (the big pass)
- Workspace middleware (header + membership → `WorkspaceContext`); explicit workspace params on repository ports; GORM tenant scope; every query filtered.
- HTTP upload route (streaming multipart, API-mediated) with quota pre-check; UUIDv4 photo IDs on SaaS.
- `GetPhotoBin`/`DownloadPhotos` → S3 streams; presigned-302 seam behind config.
- Cache-key namespacing (`thumbnail:{ws}:{id}:{size}`); search/tags scoped; thumbnail regen cutover.
- Webapp: `<img src>` thumbnails via cap URLs; `X-Workspace-ID` header in `RestApiAdapter`.
- **Gate:** **cross-tenant 403/404 matrix green on every endpoint (definition of done)**; L1 + L2 + L3 (upload→grid→full-res→share on two tenants, screenshots on PR); self-hosted deploy unchanged.

### Phase 3 — Background & realtime surfaces
- Cron jobs iterate per-workspace (AI, quality, trash purge, regen, index); jobs composite PK; per-workspace job status.
- Orphan sweep (§8); backup key split (read-only rclone key).
- WS events carry workspace end-to-end.
- **Gate:** integration tests asserting event isolation + per-workspace cron behavior; L1 + L2 + L3.

### Phase 4 — RLS hardening (before open SaaS scale)
- Postgres RLS: `SET LOCAL app.workspace_id` per request transaction; fail-closed policies (NULL GUC → zero rows); policy tests proving filterless queries return nothing; admin bypass role for cross-tenant ops.
- Not earlier: half-baked RLS either fails open or breaks every read. No-op for the single-workspace self-hosted deploy.
- **Gate:** L1 + L2 + policy test suite.

### Phase 5 — Product (post §10.2 product-model decision)
- Multi-member workspaces (UI + invitations — now a UI feature, not a migration: queries already scoped); roles 3/4; workspace switcher UI; share-revoke → cap rotation sweep; presigned multipart upload (#98); paid tiers when quotas pinch.

---

## 10. Deferred (YAGNI) — explicit non-goals for now

Per-tenant storage keys · invitations/email flows (direct-add only) · `viewer`/`admin` workspace roles · workspace switcher UI · per-workspace settings (retention etc. stay deployment-global) · presigned originals at launch · RLS (Phase 4) · billing/Stripe · quotas beyond the single storage-bytes limit · cross-workspace shares.

---

## 11. Decision register — for Rich

Confirmed 2026-09-09 (Rich). D1/D2/D5/D6/D7/D8 adopted as recommended; D3 delete (verified unused outside the shared flow); D4 rejected per §4.2.

| # | Decision | Recommendation | Status |
|---|---|---|---|
| D1 | §13.4 → thumbnail capability route with separate random `thumb_cap` (never the path-hash photo ID) | Adopt | ✅ |
| D2 | Thumbnails transiently in CF global edge; originals/metadata/DB strictly EU — accept + disclose in ToS, EU-only caching as config flag | Accept + disclose | ✅ |
| D3 | No-auth `GET /api/album/:id`: delete (if unused) or capability-ize | Verify webapp usage, likely delete | ✅ delete — verified unused outside shared flow |
| D4 | Per-tenant B2 keys rejected; 3-key process scheme (API / backup-readonly / sweep) | Confirm to close | ✅ rejected per §4.2 |
| D5 | `workspaces` + membership table from Phase 1 (not a `users.workspace_id` column) | Adopt | ✅ |
| D6 | Backfill: all existing users into one default workspace (preserves shared-library semantics) | Adopt | ✅ |
| D7 | Originals API-proxied at launch; presigned 302 as config seam | Adopt | ✅ |
| D8 | Workspace via server-side resolution (`X-Workspace-ID` + membership), not token claims | Adopt | ✅ |

---

## 12. Risk register

| # | Risk | L/I | Mitigation |
|---|---|---|---|
| 1 | Missed WHERE in the 40-endpoint sweep → cross-tenant leak | High/Critical | Explicit port params + GORM scope; authZ matrix as merge gate; RLS in Phase 4; PR checklist |
| 2 | Backfill migration fails on real data (NOT NULL violation, constraint swap) | Med/High | golang-migrate first; prod-copy dry run; idempotent steps; pre-deploy backup |
| 3 | Public surfaces leak under tenancy (share IDOR, plaintext passwords, public album endpoint, WS broadcast, global cron) | High/Critical | Phase 0 blockers; capability validation on share resolution; Phase 3 surfaces |
| 4 | Scope creep (invitations/roles/quotas/switcher in the same effort) | High/High | Strict phasing; §10 enforced at review |
| 5 | Scoped queries degrade to seq scans on large libraries | Med/Med | Composite indexes shipped with the schema change; EXPLAIN checks in matrix tests |
| 6 | `thumb_cap` treated as derivable from photo ID (regression to enumerable URLs) | Low/Critical | Cap is a random column; code review + unit test asserting independence |

---

## 13. File anchor index

| Area | File |
|---|---|
| Photo ID scheme (base64 path) | `api/internal/core/service/photo.go:289` |
| Token payload | `api/internal/core/domain/tokenPayload.go` |
| PASETO maker | `api/internal/adapter/handler/auth/paseto.go` |
| Auth middleware | `api/internal/adapter/handler/http/middleware.go:36-69` |
| Router / public routes | `api/internal/adapter/handler/http/router.go` (`:140` public album, `:206-207` shared) |
| Share plaintext password | `api/internal/core/service/share.go` |
| WS hub broadcast | `api/internal/components/websocket/hub.go:44-53` |
| GORM AutoMigrate | `api/internal/adapter/storage/database/database.go:62` |
| Thumbnail S3 key builder | `api/internal/adapter/storage/thumbnail/s3/storage.go` |
| Storage knobs | `api/internal/appconfig/appconfig.go` |
| Frontend auth header | `webapp/src/Adapters/RestApiAdapter.ts:133-137` |
| XHR thumbnail pattern | `webapp/src/Components/PhotoGrid/PhotoGrid.tsx` |
