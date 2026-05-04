# Backend Improvements Plan: Performance & Security

## Priority Summary

| Priority | Area | Impact | Effort |
|----------|------|--------|--------|
| **P0** | Thumbnail storage & serving | High perf gain, reduces DB load | Medium |
| **P0** | Photo indexing batching & I/O | Dramatically faster re-index | Medium |
| **P1** | Database query & index optimization | Faster list views, less memory | Low |
| **P1** | Security hardening | Auth, path traversal, headers | Low-Medium |
| **P2** | Serialization format | Marginal gain unless huge payloads | Medium-High |
| **P2** | Thumbnail generation pipeline | Better video support, progressive JPEG | Medium |

---

## 1. Photo Indexing (Highest Impact)

### Current Issues
- **Individual saves**: `PerformPhotoIndex` calls `save(photo)` per file. The callback `SavePhoto` does a single `INSERT` per photo with album lookup. No batching during the actual index run.
- **Redundant file I/O**: `getMetaData` opens the file once for dimensions, then again for MD5/MIME/EXIF/thumbnail (`utils.OpenFile` at line 132).
- **Album lookup N+1**: Every single photo triggers an album lookup or creation, even though `SavePhotos` exists with caching.
- **Thumbnail always regenerated**: Even if photo hasn't changed, the full thumbnail pipeline runs.

### Recommendations

**a. Batch save during indexing**
Change `PerformPhotoIndex` to collect photo metadata into chunks (e.g., 100-500) and call `SavePhotos` instead of `SavePhoto`. This requires changing the `save func(domain.PhotoFile) error` signature to accept batches or buffering in the service layer.

*Files:* `api/internal/core/service/filesystem.go:37-82`, `api/internal/core/service/photo.go:117-175`

**b. Single-pass file I/O**
Refactor `getMetaData` to open the file once and reuse the handle for dimensions, hash, MIME, and EXIF extraction.

*Files:* `api/internal/core/service/filesystem.go:108-181`

**c. Skip unchanged files during re-index**
Use `file_hash` + `filesystem_path` + `modified time` to detect unchanged files. Only process files that are new or have changed.

*Files:* `api/internal/core/service/filesystem.go:37-82`, `api/internal/core/domain/photo.go:8-30`

**d. Parallelize thumbnail generation separately**
Separate the thumbnail generation from metadata extraction. Thumbnails can be generated in a second pass or lazily, allowing the index to complete much faster.

---

## 2. Thumbnail Generation

### Current Issues
- **Fixed 600x600 size**: No responsive sizes. Every thumbnail is 600x600 regardless of use case (grid vs detail view).
- **Video thumbnails spawn ffmpeg per file**: Slow and resource-intensive.
- **No format optimization**: JPEG thumbnails are encoded with default quality. No progressive JPEG or WebP.
- **Memory overhead**: `imaging.Open` decodes the entire image into memory, then creates a second full copy for the thumbnail.

### Recommendations

**a. Multiple thumbnail sizes**
Generate `small` (200x200), `medium` (600x600), and `large` (1200x1200) variants. Store them in separate columns or migrate to a dedicated table.

*Files:* `api/internal/adapter/storage/filesystem/repository/filesystem.go:87-123`

**b. WebP for thumbnails**
WebP typically produces 25-35% smaller files than JPEG at equivalent quality. Go's `golang.org/x/image/webp` or `github.com/chai2010/webp` can encode.

**c. Streaming/sequential decode for large images**
For very large photos (>20MP), `imaging.Open` allocates a huge RGBA buffer. Consider using `imaging.Decode` with `image.Config` to check dimensions first and use a downscaling decoder if available.

**d. Video thumbnail pooling**
Reuse ffmpeg processes or use a thumbnail extraction library like `github.com/u2takey/ffmpeg-go` with better process management.

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
- **Add Redis/in-memory caching** — Cache hot thumbnails in an LRU cache. Thumbnails are immutable once generated, making them perfect for caching.
- **Replace ZIP batch with individual URLs** — The frontend already fetches thumbnails individually with concurrency limiting (6 at a time). The ZIP batch endpoint adds complexity. Consider removing it and relying on HTTP/2 multiplexing.

*Files:* `api/internal/adapter/handler/http/photo.go:243-279`, `webapp/src/utils/ThumbnailUtils.ts:32-64`

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
1. **Enable response compression** in Gin. Add `gzip` middleware:
   ```go
   import "github.com/gin-contrib/gzip"
   router.Use(gzip.Gzip(gzip.DefaultCompression))
   ```
   This reduces JSON payload size by 60-80% for list endpoints with no frontend changes.

2. **Use `goccy/go-json` explicitly** for large marshal/unmarshal operations. It's already in your dependency tree (via Gin) and is ~2x faster than `encoding/json`.

3. **If you still want binary:** Consider **MessagePack** with `github.com/vmihailenco/msgpack/v5`. It requires only a small frontend library and no schema definitions.

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

**a. Extract key metadata fields to columns**
Move frequently queried metadata (GPS lat/lng, camera model, orientation) into proper indexed columns. This eliminates JSON path queries.

*Files:* `api/internal/core/domain/photo.go:8-30`

**b. Normalize tags**
Move from comma-separated string to a `photo_tags` junction table. Enables proper indexing and efficient queries.

*Files:* `api/internal/core/domain/photo.go:14`, `api/internal/adapter/storage/database/repository/photos.go:284-299`

**c. Add partial index for geodata**
```sql
CREATE INDEX idx_photos_gps ON photobox.photos (id) 
WHERE metadata::jsonb -> 'Exif' ->> 'GPSLatitude' IS NOT NULL;
```
Or better, add `latitude`/`longitude` columns and a GiST index.

**d. Query-only thumbnail access**
As mentioned in section 3e, never fetch the thumbnail column when loading photo lists or metadata.

**e. Add `pg_stat_statements`**
Enable this PostgreSQL extension to identify slow queries in production.

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

**a. Remove insecure defaults**
```go
adminPassword := utils.GetEnv("DEFAULT_ADMIN_PASSWORD", "")
if adminPassword == "" {
    slog.Error("DEFAULT_ADMIN_PASSWORD must be set")
    os.Exit(1)
}
```

*Files:* `api/internal/appconfig/appconfig.go:41`

**b. Path traversal hardening**
Validate that `photoInfo.FilesystemPath` is within `appConfig.PhotoDir` before serving:

```go
func validatePhotoPath(baseDir, requestedPath string) error {
    absBase, _ := filepath.Abs(baseDir)
    absReq, _ := filepath.Abs(requestedPath)
    if !strings.HasPrefix(absReq, absBase) {
        return domain.ErrForbidden
    }
    return nil
}
```

*Files:* `api/internal/adapter/handler/http/photo.go:127-143`

**c. Add request body limits**
```go
router.MaxMultipartMemory = 32 << 20 // 32 MB
// Or for general body size:
router.Use(func(c *gin.Context) {
    c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 50*1024*1024)
})
```

**d. Database SSL**
Add `sslmode` to DSN based on env var:
```go
sslMode := utils.GetEnv("DB_SSL_MODE", "require")
dbURL := fmt.Sprintf("host=%s ... sslmode=%s", ..., sslMode)
```

*Files:* `api/internal/appconfig/appconfig.go:25-32`

**e. Security headers**
Add middleware for `X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Content-Security-Policy`.

**f. Rate limit thumbnail endpoints**
The global 100 req/s limiter applies, but thumbnails are bursty. Consider a separate, stricter limiter for resource-heavy endpoints.

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

### Phase 5: Advanced (Optional / Future)
1. WebP thumbnail generation
2. Multiple thumbnail sizes (small/medium/large)
3. Redis caching layer for hot thumbnails
4. HTTP/2 server push for thumbnail batches

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
