import { IPhotosAdapter } from '../Adapters/IPhotosAdapter';
import { getCachedBlobUrl, cacheBlobUrl, releaseBlobUrl, clearBlobCache } from '../hooks/useBlobUrl';

export async function fetchThumbnailWithAuth(
  photosAdapter: IPhotosAdapter,
  photoId: string,
  retries = 2
): Promise<string> {
  const cached = getCachedBlobUrl(photoId);
  if (cached) {
    return cached;
  }

  try {
    const blob = await photosAdapter.getPhotoThumbnailBlob(photoId);

    return cacheBlobUrl(photoId, blob);
  } catch (error) {
    if (retries > 0) {
      const delay = Math.pow(2, (2 - retries)) * 1000;
      await new Promise(r => setTimeout(r, delay));
      return fetchThumbnailWithAuth(photosAdapter, photoId, retries - 1);
    }
    throw error;
  }
}

export async function fetchThumbnailsBatch(
  photosAdapter: IPhotosAdapter,
  photoIds: string[]
): Promise<Map<string, string>> {
  const thumbnailMap = new Map<string, string>();

  if (photoIds.length === 0) return thumbnailMap;

  // Process in chunks of 6 to stay within browser concurrent connection limits
  const CHUNK_SIZE = 6;

  for (let i = 0; i < photoIds.length; i += CHUNK_SIZE) {
    const chunk = photoIds.slice(i, i + CHUNK_SIZE);
    const results = await Promise.allSettled(
      chunk.map(async (id) => {
        const cached = getCachedBlobUrl(id);
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
  }

  return thumbnailMap;
}

export function getCachedThumbnail(photoId: string): string | undefined {
  return getCachedBlobUrl(photoId);
}

export function revokeThumbnail(photoId: string) {
  releaseBlobUrl(photoId);
}

export function clearThumbnailCache() {
  clearBlobCache();
}
