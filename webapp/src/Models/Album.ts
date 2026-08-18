export interface Album {
  id: string;
  name: string;
  description: string;
  tags: string;
  metadata: string | Record<string, unknown>;
}

// isSmartAlbum inspects the metadata for the smart marker. The API returns
// metadata as a JSON object, but legacy/mock data may be a JSON string.
export const isSmartAlbum = (album: Album): boolean => {
  if (!album.metadata) return false;
  let meta: Record<string, unknown> | null = null;
  if (typeof album.metadata === 'string') {
    try {
      meta = JSON.parse(album.metadata);
    } catch {
      return false;
    }
  } else {
    meta = album.metadata;
  }
  return meta?.smart === true;
};
