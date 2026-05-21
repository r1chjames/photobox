# Backend Improvements Plan: Performance & Security

## Priority Summary

| Priority | Area | Impact | Effort | Status |
|----------|------|--------|--------|--------|
| **P0** | Thumbnail storage & serving | High perf gain, reduces DB load | Medium | ✅ Done |
| **P0** | Photo indexing batching & I/O | Dramatically faster re-index | Medium | ✅ Done |
| **P1** | Database query & index optimization | Faster list views, less memory | Low | ✅ Done |
| **P1** | Security hardening | Auth, path traversal, headers | Low-Medium | ⚠️ Partial |
| **P2** | Serialization format | Marginal gain unless huge payloads | Medium-High | ⚠️ Partial |
| **P2** | Thumbnail generation pipeline | Better video support, progressive JPEG | Medium | ⚠️ Partial |

---

## 1. Photo Indexing (Highest Impact)

### Current Issues
- **Individual saves**: `PerformPhotoIndex` calls `save(photo)` per file. The callback `SavePhoto` does a single `INSERT` per photo with album lookup. No batching during the actual index run.
- **Redundant file I/O**: `getMetaData` opens the file once for dimensions, then again for MD5/MIME/EXIF/thumbnail (`utils.OpenFile` at line 132).
- **Album lookup N+1**: Every single photo triggers an album lookup or creation, even though `SavePhotos` exists with caching.
- **Thumbnail always regenerated**: Even if photo hasn't changed, the full thumbnail pipeline runs.

### Recommendations

**a. ✅ Batch save during indexing**
`PerformPhotoIndex` now batches saves (100 per flush via `CreatePhotosInfo`).

*Relevant files:* `api/internal/core/service/filesystem.go:37-82`, `api/internal/core/service/photo.go:117-175`

**b. ✅ Single-pass file I/O**
`getMetaData` opens the file once and reuses the handle for dimensions, MIME, MD5, and EXIF via `Seek(0,0)`.

*Relevant files:* `api/internal/core/service/filesystem.go:108-181`

**c. ✅ Skip unchanged files during re-index**
Uses in-memory index cache (`file_hash` + `FileModifiedTime`) to skip unmodified files. New `FileModifiedTime` field on `Photo` and `PhotoFile` domain models.

*Relevant files:* `api/internal/core/service/filesystem.go:37-82`, `api/internal/core/domain/photo.go:8-30`

**d. ✅ Async thumbnail generation**
Thumbnails removed from indexing pipeline — generated on-demand in handler via `GenerateThumbnailForPhoto`.

---

## 2. Thumbnail Generation

### Current Issues
- **Fixed 600x600 size**: No responsive sizes. Every thumbnail is 600x600 regardless of use case (grid vs detail view).
- **Video thumbnails spawn ffmpeg per file**: Slow and resource-intensive.
- **No format optimization**: JPEG thumbnails are encoded with default quality. No progressive JPEG or WebP.
- **Memory overhead**: `imaging.Open` decodes the entire image into memory, then creates a second full copy for the thumbnail.

### Recommendations

**a. ✅ Multiple thumbnail sizes**
`GenerateThumbnail` now accepts `width, height`; generates small (200x200), medium (600x600), large (1200x1200); stored in `.thumbnails/{s,m,l}/`.

*Relevant files:* `api/internal/adapter/storage/filesystem/repository/filesystem.go:87-123`

**b. ❌ WebP for thumbnails (deferred) — [#35](https://github.com/r1chjames/photobox/issues/35)**
Deferred indefinitely — requires CGO/libwebp, limited pure-Go encoders available. JPEG thumbnails remain the default.

**c. ❌ Streaming/sequential decode for large images — [#40](https://github.com/r1chjames/photobox/issues/40)**
`imaging.Open` still decodes the full image. For photos >20MP this allocates a large RGBA buffer. Consider using a downscaling decoder.

**d. ❌ Video thumbnail pooling — [#41](https://github.com/r1chjames/photobox/issues/41)**
ffmpeg is still spawned per-call. Consider `github.com/u2takey/ffmpeg-go` or a persistent ffmpeg process pool.

---

## 3. Thumbnail Serving

### Status
**Phase 1 & 2 completed.** Thumbnails are now served from the filesystem with DB fallback.

### Resolved Issues
- ✅ **Thumbnails moved out of PostgreSQL** — stored on disk at `<PhotoDir>/.thumbnails/{safeId}.jpg`
- ✅ **Dedicated thumbnail query** — `GetThumbnailBytes` selects only the `thumbnail` column
- ✅ **Proper HTTP headers** — `Content-Type: image/jpeg`, `Cache-Control: public, max-age=31536000, immutable`, and `ETag` added
- ✅ **Filesystem-first serving** — `GetPhotoThumbnail` uses `ctx.File()` when `thumbnail_path` is set; falls back to DB bytes
- ✅ **Startup migration** — `MigrateThumbnailsToFilesystem` extracts existing DB thumbnails to disk on boot

### Remaining Recommendations
- ✅ **Redis/Valkey caching** — `CacheService` with `go-redis/v9` caches thumbnail paths with 24h TTL; graceful degradation when disabled or unreachable. Config via `CACHE_ENABLED`, `CACHE_HOST`, `CACHE_PORT`, `CACHE_PASSWORD`, `CACHE_DB`.
- ❌ **Replace ZIP batch with individual URLs** — [#42](https://github.com/r1chjames/photobox/issues/42) — The ZIP batch endpoint (`/photo/thumbnail/zip`) still exists. Consider removing it and relying on HTTP/2 multiplexing.
- ❌ **Separate rate limiter for thumbnail endpoints** — [#39](https://github.com/r1chjames/photobox/issues/39) — Only global (100 req/s) and auth (5 req/s) rate limiters exist. Thumbnails are bursty; consider a dedicated limiter.

*Relevant files:* `api/internal/adapter/handler/http/photo.go:243-279`, `webapp/src/utils/ThumbnailUtils.ts:32-64`

---

## 4. Serialization Format: JSON vs Protobuf

### Analysis

| Format | Pros | Cons | Verdict |
|--------|------|------|---------|
| **JSON (current)** | Human-readable, easy debug, native browser support, Gin handles automatically | Larger payload, slower parse/serialize, no schema enforcement | Keep for now |
| **Protobuf** | Compact, fast deserialization, strong typing | Requires `.proto` definitions, extra build step, harder to debug, browser needs library | Not worth it |
| **MessagePack** | Binary JSON, smaller, fast, no schema required | Browser needs decoder library | Moderate benefit |
| **JSON + zstd/gzip** | Transparent, huge size reduction, minimal code change | Slight CPU cost | **Recommended** |

### Recommendation
**Do not migrate to Protobuf.** The API is not a high-throughput microservices mesh. The overhead of maintaining `.proto` files, generating code, and adding a decoder to the React frontend outweighs the benefits for a photo gallery app.

**Instead:**
1. ✅ **Response compression enabled** — Gzip middleware added to Gin router:
   ```go
   import "github.com/gin-contrib/gzip"
   router.Use(gzip.Gzip(gzip.DefaultCompression))
   ```
   This reduces JSON payload size by 60-80% for list endpoints with no frontend changes.

   *File:* `api/internal/adapter/handler/http/router.go:43`

2. ❌ **`goccy/go-json` not explicitly configured** — [#43](https://github.com/r1chjames/photobox/issues/43) — It's in the dependency tree via Gin but not explicitly imported for marshaling. Low priority unless profiling shows serialization as a bottleneck.

3. ❌ **MessagePack not adopted** — Added complexity outweighs benefit for current payload sizes.

*Files:* `api/internal/adapter/handler/http/router.go:30-61`

---

## 5. Database Performance

### Current Issues
- **`GetPhotosWithGeodata`** loads all GPS-tagged photos into memory then maps them in Go.
- **`ListPhotosByTags`** uses `string_to_array(tags, ',')` which is not indexable. Tags are stored as comma-separated strings.
- **`GetPhotoById`** with `includeThumbnail=true` fetches the entire row including the TOAST'd thumbnail column.
- **No connection pool tuning** for the workload. 100 max open might be too high or too low depending on concurrency.
- **`metadata` is `datatypes.JSON`**: Queries like `metadata::jsonb -> 'Exif' ->> 'GPSLatitude'` work but can't use indexes well.

### Recommendations

**a. ✅ Extract key metadata fields to columns**
`latitude` and `longitude` columns added to `Photo` domain model with btree partial index. GPS extracted from EXIF during indexing (DMS → decimal degrees with hemisphere handling).

*Relevant files:* `api/internal/core/domain/photo.go:32-33`

**b. ✅ Normalize tags**
`PhotoTag` junction table created with composite PK (`photo_id`, `tag`) and indexes. `ListPhotosByTags` uses JOINs instead of `string_to_array`. `GetAllTags` queries `PhotoTag` directly. `UpdatePhotoTags` syncs both legacy `tags` column and junction table.

*Relevant files:* `api/internal/core/domain/photo.go:36-40`, `api/internal/adapter/storage/database/repository/photos.go:280-332`

**c. ✅ Add indexed columns for geodata (supersedes JSON path approach)**
Instead of a SQL partial index on JSON, `latitude`/`longitude` columns were added with `idx_photos_lat_lng` partial index. `GetPhotosWithGeodata` now uses `lat/lng BETWEEN` instead of JSON path queries.

**d. ✅ Query-only thumbnail access**
`Omit("thumbnail")` used in all list queries. `GetThumbnailBytes` selects only the `thumbnail` column.

**e. ✅ Enable `pg_stat_statements`**
Enabled in `database.go` at startup via `CREATE EXTENSION IF NOT EXISTS pg_stat_statements`.

---

## 6. Security Improvements

### Current Issues
- **Default admin password**: `DEFAULT_ADMIN_PASSWORD` defaults to `"password"` if not set.
- **Path traversal risk**: `GetPhotoBin` does `strings.ReplaceAll(photoBinary, \'\, \'\')` then `ctx.File(unescapedPath)`. If the database is compromised or a path is crafted, this could read arbitrary files.
- **No request size limits**: Upload endpoints could receive unlimited payloads.
- **No TLS configuration**: Database DSN has no `sslmode` parameter.
- **`IndexPhotos` runs in goroutine with detached context**: Errors during indexing are only logged, not returned to the caller. The `202 Accepted` response is fine, but there's no status tracking beyond the job table.
- **Thumbnail endpoint lacks auth check**: It uses the `photo` group which has `authMiddleware`, but the raw bytes are returned with no content-type validation.

### Recommendations

**a. ✅ Remove insecure defaults**
Done. `DEFAULT_ADMIN_PASSWORD` defaults to `""` and causes a hard error (`os.Exit(1)`) if not set.

*File:* `api/internal/appconfig/appconfig.go:60-62`

**b. ✅ Path traversal hardening**
`PhotoService.PhotoBinary` validates `FilesystemPath` is within `PhotoDir` before serving. Returns `domain.ErrForbidden` on mismatch.

*File:* `api/internal/core/service/photo.go:85-91`

**c. ❌ Add request body limits — [#36](https://github.com/r1chjames/photobox/issues/36)**
No `MaxBytesReader` or `MaxMultipartMemory` limit configured. Upload endpoints can receive unlimited payloads.

**d. ❌ Database SSL — [#37](https://github.com/r1chjames/photobox/issues/37)**
No `sslmode` env var in appconfig. Only test config has `sslmode=disable`. Production DSN may be transmitting credentials in plaintext.

*File:* `api/internal/appconfig/appconfig.go:25-32`

**e. ❌ Security headers — [#38](https://github.com/r1chjames/photobox/issues/38)**
No middleware for `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, or `Content-Security-Policy`.

**f. ❌ Rate limit thumbnail endpoints — [#39](https://github.com/r1chjames/photobox/issues/39)**
Global 100 req/s and auth 5 req/s rate limiters exist, but no dedicated limiter for thumbnail/resource-heavy endpoints.

---

## 7. Implementation Roadmap

### Phase 1: Quick Wins ✅ COMPLETE (`5a9d7b6`)
1. ✅ Add `gzip` middleware to Gin router
2. ✅ Add `Content-Type`, `Cache-Control`, and `ETag` headers to thumbnail endpoint
3. ✅ Add `GetThumbnailBytes` query-only method to photo repository
4. ✅ Add path traversal validation to `GetPhotoBin`
5. ✅ Remove default admin password fallback
6. ✅ Enable `pg_stat_statements`

### Phase 2: Thumbnail Infrastructure ✅ COMPLETE (`d605521`)
1. ✅ Add `ThumbnailPath` field to `Photo` domain model (GORM auto-migrate)
2. ✅ Add `PhotoThumbnailPath` service + repository method
3. ✅ Update `SavePhoto`/`SavePhotos` to write thumbnail bytes to disk and store path
4. ✅ Update `GetPhotoThumbnail` handler to serve from filesystem; fall back to DB
5. ✅ Add startup migration (`MigrateThumbnailsToFilesystem`) for existing DB thumbnails
6. ✅ Update tests and mocks for new interface methods

### Phase 3: Indexing Optimization ✅ COMPLETE (`cada184`)
1. ✅ Refactor `PerformPhotoIndex` to batch saves (100 per flush via `CreatePhotosInfo`)
2. ✅ Single-pass file I/O in `getMetaData` — file opened once, dimensions/MIME/MD5/EXIF read via `Seek(0,0)`
3. ✅ Skip unchanged files using in-memory index cache (`file_hash` + `FileModifiedTime`)
4. ✅ Generate thumbnails asynchronously — removed from indexing pipeline, generated on-demand in handler via `GenerateThumbnailForPhoto`
5. ✅ Add `FileModifiedTime` to `Photo` and `PhotoFile` domain models

### Phase 4: Schema & Query Improvements ✅ COMPLETE (`11c0359`)
1. ✅ Add `latitude`, `longitude` columns to `Photo` with btree partial index
2. ✅ Extract GPS from EXIF during indexing (DMS → decimal degrees with hemisphere handling)
3. ✅ Create `PhotoTag` junction table with composite PK (`photo_id`, `tag`) and indexes
4. ✅ Add `idx_photos_lat_lng` partial index for geospatial queries
5. ✅ Update `GetPhotosWithGeodata` to use indexed `lat/lng BETWEEN` instead of JSON path
6. ✅ Update `ListPhotosByTags` to use junction table JOINs instead of `string_to_array`
7. ✅ Update `GetAllTags` to query `PhotoTag` directly
8. ✅ Update `UpdatePhotoTags` to sync both legacy `tags` column and junction table

### Phase 5: Advanced ⚠️ PARTIAL (`f094907`)
1. ❌ WebP thumbnail generation — [#35](https://github.com/r1chjames/photobox/issues/35); deferred; requires CGO/libwebp, limited pure-Go encoders available
2. ✅ Multiple thumbnail sizes — `GenerateThumbnail` now accepts `width, height`; generates small (200x200), medium (600x600), large (1200x1200); stored in `.thumbnails/{s,m,l}/`
3. ✅ Valkey/Redis caching layer — `CacheService` with `go-redis/v9`; caches thumbnail paths with 24h TTL; graceful degradation when disabled or unreachable; config via `CACHE_ENABLED`, `CACHE_HOST`, `CACHE_PORT`, `CACHE_PASSWORD`, `CACHE_DB`
4. ❌ HTTP/2 server push for thumbnail batches — deferred; modern browsers deprecating server push, HTTP/2 multiplexing + individual requests is sufficient

**New env vars for cache:**
| Variable | Default | Description |
|----------|---------|-------------|
| `CACHE_ENABLED` | `false` | Enable Valkey/Redis caching |
| `CACHE_HOST` | `localhost` | Cache server host |
| `CACHE_PORT` | `6379` | Cache server port |
| `CACHE_PASSWORD` | `""` | Auth password (optional) |
| `CACHE_DB` | `0` | Redis DB number |

**Thumbnail endpoint now supports:**
```
GET /photo/{id}/thumbnail?size=s   # 200x200
GET /photo/{id}/thumbnail?size=m   # 600x600 (default)
GET /photo/{id}/thumbnail?size=l   # 1200x1200
```

---

## 8. Remaining Work

Items not yet implemented, ordered by practical impact:

| Priority | Item | Section Ref | Issue | Effort | Notes |
|----------|------|-------------|-------|--------|-------|
| **P1** | Request body limits (upload size cap) | §6c | [#36](https://github.com/r1chjames/photobox/issues/36) | Low | 2 lines of middleware; prevents OOM from large uploads |
| **P1** | Database SSL (`sslmode`) configuration | §6d | [#37](https://github.com/r1chjames/photobox/issues/37) | Low | Add `DB_SSL_MODE` env var to appconfig, append to DSN |
| **P2** | Security headers middleware | §6e | [#38](https://github.com/r1chjames/photobox/issues/38) | Low | Add `X-Content-Type-Options`, `X-Frame-Options`, `CSP` to router |
| **P2** | Separate rate limiter for thumbnail endpoints | §6f | [#39](https://github.com/r1chjames/photobox/issues/39) | Low | Dedicated limiter for resource-heavy endpoints |
| **P2** | Streaming/sequential decode for large images | §2c | [#40](https://github.com/r1chjames/photobox/issues/40) | Medium | Use downscaling decoder for photos >20MP |
| **P2** | Video thumbnail pooling | §2d | [#41](https://github.com/r1chjames/photobox/issues/41) | Medium | Reuse ffmpeg process or switch to lib-based extraction |
| **P3** | Replace ZIP batch with individual URLs | §3 | [#42](https://github.com/r1chjames/photobox/issues/42) | Low | Remove `/photo/thumbnail/zip` endpoint |
| **P3** | WebP thumbnail generation | §2b | [#35](https://github.com/r1chjames/photobox/issues/35) | High | Deferred until CGO/libwebp or viable pure-Go encoder |
| **P3** | Configure `goccy/go-json` explicitly | §4 | [#43](https://github.com/r1chjames/photobox/issues/43) | Low | Marginal gain unless profiling shows bottleneck |
| **P3** | Remove legacy `tags` column migration | §5b | [#44](https://github.com/r1chjames/photobox/issues/44) | Medium | Once `PhotoTag` junction table is sole source of truth |

---

## Files to Reference for Implementation

| File | Relevance |
|------|-----------|
| `api/internal/adapter/handler/http/photo.go` | Thumbnail handlers, cache headers, path validation |
| `api/internal/adapter/handler/http/router.go` | Compression middleware, rate limiters |
| `api/internal/core/service/filesystem.go` | Indexing pipeline, worker pool, metadata extraction |
| `api/internal/core/service/photo.go` | SavePhoto/SavePhotos, album caching |
| `api/internal/adapter/storage/filesystem/repository/filesystem.go` | Thumbnail generation, ffmpeg |
| `api/internal/adapter/storage/database/repository/photos.go` | DB queries, batch operations |
| `api/internal/core/domain/photo.go` | Domain model, indexes |
| `api/internal/appconfig/appconfig.go` | Security defaults, DB SSL |
| `webapp/src/utils/ThumbnailUtils.ts` | Frontend batch fetching logic |
