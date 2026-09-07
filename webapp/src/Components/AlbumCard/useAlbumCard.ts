import { useQuery } from '@tanstack/react-query';
import { Album } from "../../Models/Album";
import { IPhotosAdapter } from "../../Adapters/IPhotosAdapter";
import { fetchThumbnailWithAuth, getCachedThumbnail } from "../../utils/ThumbnailUtils";

// Drives the prototype's 3-image album card: one hero thumbnail plus two
// stacked secondary thumbnails, all fetched through the authenticated
// thumbnail pipeline with its IndexedDB/blob caching.
//
// The metadata query returns the raw cover photos (blurhash/dominantColor
// included) so the card can paint placeholder-first frames immediately and
// crossfade the fetched blobs in on top (issue #173). Blob failures keep the
// placeholder visible per-slot instead of leaving a blank frame.
const useAlbumCard = (photosAdapter: IPhotosAdapter, source: Album) => {
    const { data: photos, isLoading: isMetadataLoading, isError: isMetadataError } = useQuery({
        queryKey: ['albumThumbnail', source.id],
        queryFn: async () => {
            const result = await photosAdapter.getPhotosInfoInAlbum(source.id, "", 3, false);
            return result ?? [];
        },
        staleTime: 60_000,
    });

    const heroPhotoId = photos?.[0]?.id ?? null;

    const { data: heroBlobUrl, isError: isHeroBlobError } = useQuery({
        queryKey: ['albumThumbnailBlob', heroPhotoId],
        queryFn: async () => {
            if (!heroPhotoId) return undefined;
            const cached = getCachedThumbnail(heroPhotoId);
            if (cached) return cached;
            return fetchThumbnailWithAuth(photosAdapter, heroPhotoId);
        },
        enabled: !!heroPhotoId,
        staleTime: 60_000,
    });

    const { data: subBlobUrls, isError: isSubBlobsError } = useQuery({
        queryKey: ['albumSubThumbnails', source.id, (photos ?? []).slice(1, 3).map(p => p.id).join(',')],
        queryFn: async () => {
            // Keep the result aligned to photos[1..2]: a failed id resolves to
            // null so that slot keeps its placeholder instead of shifting.
            const subs = (photos ?? []).slice(1, 3);
            return Promise.all(subs.map(async (photo) => {
                const cached = getCachedThumbnail(photo.id);
                if (cached) return cached;
                try {
                    return await fetchThumbnailWithAuth(photosAdapter, photo.id);
                } catch {
                    return null;
                }
            }));
        },
        enabled: (photos?.length ?? 0) > 1,
        staleTime: 60_000,
    });

    const { data: photoCount, isLoading: isCountLoading, isError: isCountError } = useQuery({
        queryKey: ['albumPhotoCount', source.id],
        queryFn: () => photosAdapter.getPhotoCountInAlbum(source.id),
        staleTime: 60_000,
    });

    const isLoading = isMetadataLoading || isCountLoading;

    return {
        photos: photos ?? [],
        heroBlobUrl,
        subBlobUrls: subBlobUrls ?? [],
        photoCount: photoCount ?? 0,
        isMetadataLoading,
        isCountLoading,
        isError: isMetadataError || isCountError || isHeroBlobError || isSubBlobsError,
        isLoading,
    };
};

export default useAlbumCard;
