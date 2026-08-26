import { useQuery } from '@tanstack/react-query';
import { Album } from "../../Models/Album";
import { IPhotosAdapter } from "../../Adapters/IPhotosAdapter";
import { fetchThumbnailWithAuth, getCachedThumbnail } from "../../utils/ThumbnailUtils";

// Drives the prototype's 3-image album card: one hero thumbnail plus two
// stacked secondary thumbnails, all fetched through the authenticated
// thumbnail pipeline with its IndexedDB/blob caching.
const useAlbumCard = (photosAdapter: IPhotosAdapter, source: Album) => {
    const { data: photoIds, isLoading: isThumbnailLoading } = useQuery({
        queryKey: ['albumThumbnail', source.id],
        queryFn: async () => {
            const photos = await photosAdapter.getPhotosInfoInAlbum(source.id, "", 3, false);
            return photos ? photos.map(p => p.id) : [];
        },
        staleTime: 60_000,
    });

    const firstPhotoId = photoIds?.[0] ?? null;

    const { data: thumbnailBlobUrl } = useQuery({
        queryKey: ['albumThumbnailBlob', firstPhotoId],
        queryFn: async () => {
            if (!firstPhotoId) return 'no_image.png';
            const cached = getCachedThumbnail(firstPhotoId);
            if (cached) return cached;
            return fetchThumbnailWithAuth(photosAdapter, firstPhotoId);
        },
        enabled: !!firstPhotoId,
        staleTime: 60_000,
    });

    const { data: subThumbnailUrls } = useQuery({
        queryKey: ['albumSubThumbnails', source.id, (photoIds ?? []).slice(1).join(',')],
        queryFn: async () => {
            const ids = (photoIds ?? []).slice(1, 3);
            return Promise.all(ids.map(async (id) => {
                const cached = getCachedThumbnail(id);
                if (cached) return cached;
                try {
                    return await fetchThumbnailWithAuth(photosAdapter, id);
                } catch {
                    return null;
                }
            }));
        },
        enabled: (photoIds?.length ?? 0) > 1,
        staleTime: 60_000,
    });

    const { data: photoCount, isLoading: isCountLoading } = useQuery({
        queryKey: ['albumPhotoCount', source.id],
        queryFn: () => photosAdapter.getPhotoCountInAlbum(source.id),
        staleTime: 60_000,
    });

    const isLoading = isThumbnailLoading || isCountLoading;
    const thumbnailUrl = thumbnailBlobUrl ?? undefined;

    return {
        thumbnailUrl,
        subThumbnailUrls: (subThumbnailUrls ?? []).filter((u): u is string => !!u),
        photoCount: photoCount ?? 0,
        isLoading,
    };
};

export default useAlbumCard;
