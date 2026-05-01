import { useQuery } from '@tanstack/react-query';
import { Album } from "../../Models/Album";
import { IPhotosAdapter } from "../../Adapters/IPhotosAdapter";
import { fetchThumbnailWithAuth, getCachedThumbnail } from "../../utils/ThumbnailUtils";

const useAlbumCard = (photosAdapter: IPhotosAdapter, source: Album) => {
    const { data: firstPhotoId, isLoading: isThumbnailLoading } = useQuery({
        queryKey: ['albumThumbnail', source.id],
        queryFn: async () => {
            const photos = await photosAdapter.getPhotosInfoInAlbum(source.id, "", 1, false);
            return (photos && photos.length > 0) ? photos[0].id : null;
        },
        staleTime: 60_000,
    });

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

    const { data: photoCount, isLoading: isCountLoading } = useQuery({
        queryKey: ['albumPhotoCount', source.id],
        queryFn: () => photosAdapter.getPhotoCountInAlbum(source.id),
        staleTime: 60_000,
    });

    const isLoading = isThumbnailLoading || isCountLoading;
    const thumbnailUrl = thumbnailBlobUrl ?? undefined;

    return { thumbnailUrl, photoCount: photoCount ?? 0, isLoading };
};

export default useAlbumCard;
