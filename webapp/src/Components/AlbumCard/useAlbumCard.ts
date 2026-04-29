import { useQuery } from '@tanstack/react-query';
import { Album } from "../../Models/Album";
import { IPhotosAdapter } from "../../Adapters/IPhotosAdapter";

const useAlbumCard = (photosAdapter: IPhotosAdapter, source: Album) => {
    const { data: thumbnailUrl, isLoading: isThumbnailLoading } = useQuery({
        queryKey: ['albumThumbnail', source.id],
        queryFn: async () => {
            const photos = await photosAdapter.getPhotosInfoInAlbum(source.id, "", 1, true);
            return (photos && photos.length > 0) ? photos[0].thumbnail : 'no_image.png';
        },
        staleTime: 60_000,
    });

    const { data: photoCount, isLoading: isCountLoading } = useQuery({
        queryKey: ['albumPhotoCount', source.id],
        queryFn: () => photosAdapter.getPhotoCountInAlbum(source.id),
        staleTime: 60_000,
    });

    const isLoading = isThumbnailLoading || isCountLoading;

    return { thumbnailUrl: thumbnailUrl ?? '', photoCount: photoCount ?? 0, isLoading };
};

export default useAlbumCard;
