# Photobox API Changes Required

This document collects all backend API changes required to support the frontend feature enhancement plan. Each item includes the endpoint, method, request/response shape, and which frontend feature it enables.

---

## Photos

### `DELETE /api/photos/:id`
**Enables:** Photo deletion, batch delete, trash
- **Auth:** Yes
- **Behavior:** Move photo to a `.trash/` subdirectory on the filesystem (or set a `deleted_at` flag in DB). Do NOT permanently delete.
- **Response:** `{ success: true, message: "Photo moved to trash" }`

### `POST /api/photos/trash/restore/:id`
**Enables:** Trash restore
- **Auth:** Yes
- **Behavior:** Move photo back from `.trash/` to its original album directory.
- **Response:** `{ success: true, message: "Photo restored" }`

### `DELETE /api/photos/trash/empty`
**Enables:** Empty trash
- **Auth:** Yes (admin)
- **Behavior:** Permanently delete all photos in `.trash/` and remove DB records.
- **Response:** `{ success: true, message: "Trash emptied" }`

### `GET /api/photos/trash`
**Enables:** Trash view
- **Auth:** Yes
- **Response:** `PaginatedServerResponse<Photo[]>` — same shape as `/api/photos`

### `PATCH /api/photos/:id/favorite`
**Enables:** Favorites / starring
- **Auth:** Yes
- **Body:** `{ favorite: boolean }`
- **Behavior:** Toggle the `favorite` flag on the Photo record.
- **Response:** `{ success: true, data: Photo }`

### `GET /api/photos?favorites=true`
**Enables:** Favorites filter
- **Auth:** Yes
- **Behavior:** When `favorites=true` query param is present, return only photos where `favorite = true`.

### `POST /api/photos/download`
**Enables:** Batch download as zip
- **Auth:** Yes
- **Body:** `{ photoIds: string[] }`
- **Behavior:** Stream a ZIP archive containing the requested original photos.
- **Response:** `Content-Type: application/zip` with `Content-Disposition: attachment; filename="photobox-download.zip"`

### `POST /api/photos/:id/rotate`
**Enables:** Photo rotation
- **Auth:** Yes
- **Query:** `?direction=cw|ccw`
- **Behavior:** Apply lossless JPEG rotation via EXIF orientation manipulation, or re-encode with imaging library. Regenerate thumbnail.
- **Response:** `{ success: true, data: Photo }`

---

## Albums

### `PATCH /api/albums/:id`
**Enables:** Album editing (rename, description, cover)
- **Auth:** Yes (admin)
- **Body:** `{ name?: string, description?: string, coverPhotoId?: string, tags?: string }`
- **Behavior:** Update album fields. If `name` changes, rename the underlying filesystem directory.
- **Response:** `{ success: true, data: Album }`

### `DELETE /api/albums/:id`
**Enables:** Album deletion
- **Auth:** Yes (admin)
- **Query:** `?deletePhotos=false|true` (default: false)
- **Behavior:** If `deletePhotos=false`, move photos to a default "Uncategorized" album. If `true`, move photos to trash. Remove album record.
- **Response:** `{ success: true, message: "Album deleted" }`

### `POST /api/albums`
**Enables:** Manual empty album creation
- **Auth:** Yes
- **Body:** `{ name: string, description?: string }`
- **Behavior:** Create a new empty album record and corresponding filesystem directory under `PHOTO_DIR`.
- **Response:** `{ success: true, data: Album }`

---

## Search

### `GET /api/search?q=...&limit=...&fromId=...`
**Enables:** Full-text search
- **Auth:** Yes
- **Query:** `q` (search string), `limit`, `fromId` (cursor)
- **Behavior:** Full-text search across photo filenames, album names, and EXIF metadata (camera, lens, date). Use PostgreSQL `tsvector`.
- **Response:** `PaginatedServerResponse<(Photo | Album)[]>` with a `type` field on each result (`photo` | `album`)

### `GET /api/photos/timeline`
**Enables:** Timeline scrubber
- **Auth:** Yes
- **Response:** `{ data: [{ year: number, month: number, count: number }] }`

---

## Sharing

### `POST /api/share`
**Enables:** Shareable links
- **Auth:** Yes
- **Body:** `{ resourceType: 'photo' | 'album', resourceId: string, expiry?: string, password?: string }`
- **Behavior:** Create a cryptographically random token, store in `shared_links` table with expiry and optional password hash.
- **Response:** `{ success: true, data: { token: string, url: string, expiry?: string } }`

### `GET /api/shared/:token`
**Enables:** Public shared view
- **Auth:** No
- **Behavior:** Look up shared link by token. If expired or password-protected, return 410/401. Otherwise return the photo or album data.
- **Response:** `ServerResponse<Photo | Album>`

### `GET /api/shares`
**Enables:** Manage shares page
- **Auth:** Yes (admin)
- **Response:** `ServerResponse<SharedLink[]>`

### `DELETE /api/shares/:token`
**Enables:** Revoke share
- **Auth:** Yes (admin)
- **Behavior:** Delete shared link record.

---

## Geodata

### `GET /api/photos/geodata`
**Enables:** Map view
- **Auth:** Yes
- **Query:** `?north=...&south=...&east=...&west=...` (optional bounding box)
- **Behavior:** Return photos that have GPS coordinates in their EXIF metadata.
- **Response:** `ServerResponse<{ id: string, lat: number, lng: number, thumbnail: string, dateTaken: string }[]>`

---

## Users

### `GET /api/users`
**Enables:** User management UI
- **Auth:** Yes (admin)
- **Response:** `ServerResponse<User[]>`

### `PATCH /api/users/:id`
**Enables:** User role/approval changes
- **Auth:** Yes (admin)
- **Body:** `{ role?: 'administrator' | 'viewer' | 'contributor', approved?: boolean }`
- **Response:** `ServerResponse<User>`

### `DELETE /api/users/:id`
**Enables:** User deletion
- **Auth:** Yes (admin)
- **Response:** `{ success: true }`

---

## Schema Additions

The following fields/columns should be added to existing models:

| Model | Field | Type | Default | Notes |
|-------|-------|------|---------|-------|
| `Photo` | `favorite` | `boolean` | `false` | Indexed for fast favorites queries |
| `Photo` | `deleted_at` | `timestamp` | `null` | Soft-delete / trash flag |
| `Photo` | `blurhash` | `string` | `null` | Computed during indexing |
| `Photo` | `dominant_color` | `string` | `null` | Hex color for placeholders |
| `Album` | `cover_photo_id` | `string` (FK) | `null` | User-selected cover |
| `Album` | `description` | `string` | `''` | Already exists but unused |

New table: `shared_links`
```sql
CREATE TABLE photobox.shared_links (
    token VARCHAR(64) PRIMARY KEY,
    resource_type VARCHAR(20) NOT NULL, -- 'photo' | 'album'
    resource_id VARCHAR(255) NOT NULL,
    created_by UUID REFERENCES photobox.users(id),
    expiry TIMESTAMP,
    password_hash VARCHAR(255),
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Implementation Priority

1. **P0 ( unblock Phase 2-3 ):** `DELETE /api/photos/:id`, `PATCH /api/photos/:id/favorite`, `PATCH /api/albums/:id`, `POST /api/albums`
2. **P1 ( unblock Phase 4 ):** `GET /api/search`, `GET /api/photos/timeline`, `POST /api/photos/download`
3. **P2 ( unblock Phase 5 ):** `POST /api/share`, `GET /api/shared/:token`, `GET /api/photos/geodata`
4. **P3 ( nice-to-have ):** `POST /api/photos/:id/rotate`, `GET /api/users`, user management endpoints
