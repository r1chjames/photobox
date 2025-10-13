import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";

const usePhotoGrid = (photosAdapter: IPhotosAdapter, albumsAdapter: IAlbumsAdapter, albumId: string | undefined) => {
    const limit = 30;

    const { data: album } = useQuery({
        queryKey: ['album', albumId],
        queryFn: () => albumsAdapter.getAlbumInfoById(albumId!),
        enabled: !!albumId,
    });
    const albumName = albumId ? album?.name : "All Photos";

    async function fetchPhotos({ pageParam = "" }) {
        const fromId = pageParam;
        const retrievedPhotos = albumId !== undefined
            ? await photosAdapter.getPhotosInfoInAlbum(albumId!, fromId, limit, true)
            : await photosAdapter.getAllPhotosInfo(fromId, limit, true);

        const photosData = retrievedPhotos ?? [];

        return {
            data: photosData,
            nextCursor: photosData.length === limit ? photosData[photosData.length - 1].id : undefined,
        };
    };

    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
    } = useInfiniteQuery({
        queryKey: ['albumPhotos', albumId],
        queryFn: fetchPhotos,
        initialPageParam: "",
        getNextPageParam: (lastPage) => lastPage.nextCursor,
        // The select function transforms the paged data into a single, flat, de-duplicated array.
        select: (data) => {
            // 1. Flatten all photos from all pages and filter out any null/undefined items.
            const allPhotos = data.pages.flatMap(page => page.data).filter(Boolean);

            // 2. De-duplicate the flat list using a Map to ensure uniqueness based on photo.id.
            const uniquePhotos = Array.from(new Map(allPhotos.map(photo => [photo.id, photo])).values());

            // 3. Return a simplified structure containing the flat list.
            return {
                pages: data.pages, // Keep original pages for react-query's internal logic
                pageParams: data.pageParams,
                // This new top-level 'photos' array is the clean, de-duplicated data.
                photos: uniquePhotos,
            };
        }
    });

    // The component will consume this flat 'photos' array.
    const photos = data?.photos ?? [];

    return { photos, albumName, allRetrieved: !hasNextPage, fetchNextPage, isFetchingNextPage };
};

export default usePhotoGrid;