# Frontend Thumbnail Loading Changes

## Summary

Thumbnails are no longer embedded in list API responses. Instead, they are loaded via dedicated thumbnail endpoints. This reduces list response sizes from ~4.5MB to ~50KB for 30 photos.

## What Changed

### 1. Photo Model
- **Removed**: `thumbnail: string` field
- **Added**: `thumbnailUrl: string` field (relative path like `photo/L1VzZXJz.../thumbnail`)

### 2. API Calls
- All list endpoints now use `thumbnail=false` (no thumbnail data in response)
- New batch endpoint: `POST /api/photos/thumbnails`
- New helper method: `photosAdapter.getThumbnailUrl(photoId)` returns full URL

### 3. Components Updated

| Component | Change |
|-----------|--------|
| **PhotoGrid** | Uses `photosAdapter.getThumbnailUrl(photo.id)` for `<img src>` |
| **PhotoCard** | Uses `photosAdapter.getThumbnailUrl()` for blurred placeholder |
| **AlbumCard** | `useAlbumCard` hook now fetches photo ID only, constructs URL |
| **Memories** | Uses `photosAdapter.getThumbnailUrl()` directly |
| **MapView** | Uses `photosAdapter.getThumbnailUrl()` for map popup thumbnails |
| **TrashView** | Uses `photosAdapter.getThumbnailUrl()` for trash photo thumbnails |

### 4. How Thumbnail URLs Work

```typescript
// Get full thumbnail URL for a photo
const url = photosAdapter.getThumbnailUrl(photo.id);
// Returns: "http://localhost:8080/api/photo/thumbnail/L1VzZXJz..."

// Use directly in img src
<img src={url} alt={photo.name} />
```

The browser handles:
- Parallel loading of multiple thumbnails
- Caching (service worker caches `/thumbnail/` paths)
- Connection pooling

### 5. Batch Thumbnail Loading (Available but not used by default)

For cases where you need to preload many thumbnails at once:

```typescript
const thumbnails = await photosAdapter.getPhotoThumbnails(['id1', 'id2', 'id3']);
// Returns Map<photoId, blobUrl>
```

This fetches thumbnails as a zip file and creates blob URLs. Useful for:
- Preloading thumbnails for offline viewing
- Slideshow preloading
- Custom caching strategies

**Note**: Blob URLs must be cleaned up with `URL.revokeObjectURL()` when no longer needed.

## Performance Impact

| Metric | Before | After |
|--------|--------|-------|
| List response size (30 photos) | ~4.5MB | ~50KB |
| Initial page load | 800ms+ | <100ms |
| Thumbnail loading | Embedded (blocking) | Parallel (non-blocking) |
| Browser caching | None | Automatic + Service Worker |

## Migration Notes for Frontend

1. **No more base64 conversion**: Remove the pattern:
   ```typescript
   // OLD - no longer needed
   const imgSrc = photo.thumbnail.startsWith('http')
       ? photo.thumbnail
       : `data:image/png;base64,${photo.thumbnail}`;
   ```

2. **Use adapter method**:
   ```typescript
   // NEW
   const imgSrc = photosAdapter.getThumbnailUrl(photo.id);
   ```

3. **Mock adapters**: Update mock data to use `thumbnailUrl` instead of `thumbnail`
