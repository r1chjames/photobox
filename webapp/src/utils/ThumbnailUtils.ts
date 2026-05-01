import { IPhotosAdapter } from '../Adapters/IPhotosAdapter';

const thumbnailCache = new Map<string, string>();

export async function fetchThumbnailWithAuth(
  photosAdapter: IPhotosAdapter,
  photoId: string
): Promise<string> {
  const cached = thumbnailCache.get(photoId);
  if (cached) {
    console.log('[ThumbnailUtils] Cache hit for', photoId);
    return cached;
  }

  console.log('[ThumbnailUtils] Fetching thumbnail for', photoId);
  const blob = await photosAdapter.getPhotoThumbnailBlob(photoId);
  console.log('[ThumbnailUtils] Got response for', photoId, typeof blob, blob instanceof Blob ? blob.size : 'N/A');
  
  if (typeof blob === 'string') {
    thumbnailCache.set(photoId, blob);
    return blob;
  }

  const url = URL.createObjectURL(blob);
  thumbnailCache.set(photoId, url);
  return url;
}

export async function fetchThumbnailsBatch(
  photosAdapter: IPhotosAdapter,
  photoIds: string[]
): Promise<Map<string, string>> {
  const thumbnailMap = new Map<string, string>();

  const results = await Promise.allSettled(
    photoIds.map(async (id) => {
      const cached = thumbnailCache.get(id);
      if (cached) {
        return { id, url: cached };
      }
      const url = await fetchThumbnailWithAuth(photosAdapter, id);
      return { id, url };
    })
  );

  for (const result of results) {
    if (result.status === 'fulfilled') {
      thumbnailMap.set(result.value.id, result.value.url);
    }
  }

  return thumbnailMap;
}

export function getCachedThumbnail(photoId: string): string | undefined {
  return thumbnailCache.get(photoId);
}

export function revokeThumbnail(photoId: string) {
  const url = thumbnailCache.get(photoId);
  if (url && url.startsWith('blob:')) {
    URL.revokeObjectURL(url);
  }
  thumbnailCache.delete(photoId);
}

export function clearThumbnailCache() {
  thumbnailCache.forEach((url) => {
    if (url.startsWith('blob:')) {
      URL.revokeObjectURL(url);
    }
  });
  thumbnailCache.clear();
}
