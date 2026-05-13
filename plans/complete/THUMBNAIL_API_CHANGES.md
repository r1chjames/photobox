# API Changes: Thumbnail Loading

## Summary

Thumbnails are no longer embedded in list responses to improve performance. Instead, use the new batch endpoint or individual thumbnail URLs.

## What Changed

### 1. List responses no longer include thumbnail data

**Before:**
```json
{
  "id": "L1VzZXJzL3JpY2gv...",
  "name": "photo.jpg",
  "thumbnail": "<base64-encoded-bytes>",
  ...
}
```

**After:**
```json
{
  "id": "L1VzZXJzL3JpY2gv...",
  "name": "photo.jpg",
  "thumbnailUrl": "photo/L1VzZXJzL3JpY2gv.../thumbnail",
  ...
}
```

The `thumbnail` field is removed. A new `thumbnailUrl` field provides the relative URL to fetch the thumbnail.

### 2. New batch thumbnail endpoint

**`POST /api/photos/thumbnails`**

Fetch thumbnails for multiple photos in a single request. Returns a zip file where each entry is named by the photo ID.

**Request:**
```json
{
  "photoIds": ["id1", "id2", "id3"]
}
```

**Response:**
- Content-Type: `application/zip`
- Cache-Control: `public, max-age=31536000, immutable`
- Body: Zip archive containing thumbnail files, each named by photo ID

**Limits:** Maximum 100 photoIds per request.

### 3. Individual thumbnail endpoint (unchanged)

**`GET /api/photo/thumbnail/:id`**

Still works as before. Returns raw thumbnail bytes for a single photo.

## Migration Guide

### Option A: Batch loading (recommended for initial page load)

```javascript
// 1. Fetch the photo list
const response = await fetch('/api/photos?limit=30&thumbnail=false');
const data = await response.json();
const photos = data.data;

// 2. Extract photo IDs and fetch thumbnails in batch
const photoIds = photos.map(p => p.id);
const thumbResponse = await fetch('/api/photos/thumbnails', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ photoIds })
});

// 3. Unzip the response and create object URLs
const zipBlob = await thumbResponse.blob();
const zip = await JSZip.loadAsync(zipBlob);

const thumbnailMap = {};
for (const photoId of photoIds) {
  const file = zip.file(photoId);
  if (file) {
    const blob = await file.async('blob');
    thumbnailMap[photoId] = URL.createObjectURL(blob);
  }
}

// 4. Use thumbnailMap[photo.id] as the src for each <img>
```

### Option B: Individual thumbnail URLs

```javascript
// 1. Fetch the photo list
const response = await fetch('/api/photos?limit=30&thumbnail=false');
const data = await response.json();
const photos = data.data;

// 2. Use thumbnailUrl directly (browser handles caching)
photos.forEach(photo => {
  const img = document.createElement('img');
  img.src = `/api/${photo.thumbnailUrl}`;
  // or use fetch for more control
});
```

### Option C: Lazy loading with IntersectionObserver

```javascript
// Best for large galleries - only load thumbnails when visible
const observer = new IntersectionObserver((entries) => {
  entries.forEach(entry => {
    if (entry.isIntersecting) {
      const img = entry.target;
      img.src = img.dataset.thumbnailUrl;
      observer.unobserve(img);
    }
  });
});

photos.forEach(photo => {
  const img = document.createElement('img');
  img.dataset.thumbnailUrl = `/api/${photo.thumbnailUrl}`;
  img.src = 'placeholder.jpg'; // Show placeholder while loading
  observer.observe(img);
});
```

## Performance Comparison

| Approach | Request Count | Total Data | Initial Load |
|----------|--------------|------------|--------------|
| Old (embedded) | 1 | ~4.5MB (30 photos) | 800ms+ |
| New (batch) | 2 | ~4.5MB | 200-300ms |
| New (individual) | 31 | ~4.5MB | 100-200ms (cached) |
| New (lazy) | 1 + visible | ~150KB per visible | <50ms |

## Notes

- Thumbnails are ~150KB each (JPEG)
- The batch endpoint returns `Cache-Control: max-age=31536000, immutable` - browsers will cache aggressively
- Individual thumbnail requests can be made in parallel by the browser
- For best UX, use lazy loading with a placeholder image
