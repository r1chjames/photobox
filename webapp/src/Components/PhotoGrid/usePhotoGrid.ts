import { useInfiniteQuery, useQuery } from "@tanstack/react-query";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";

const usePhotoGrid = (photosAdapter: IPhotosAdapter, albumsAdapter: IAlbumsAdapter, albumIdentifier: string | undefined, startDate?: string, endDate?: string, tags?: string, mediaType?: string, searchQuery?: string, favoritesOnly?: boolean) => {
    const limit = 60;

    // When an albumIdentifier is present, it's used to fetch album details.
    // This hook assumes the identifier could be a name (in Storybook) or an ID.
    const { data: album } = useQuery({
        queryKey: ['album', albumIdentifier],
        queryFn: () => albumsAdapter.getAlbumInfoById(albumIdentifier!),
        enabled: !!albumIdentifier, // Only run this query if an albumIdentifier is provided.
    });

    // Use the fetched album's name if available, otherwise default to "All Photos".
    const albumName = albumIdentifier ? album?.name : (searchQuery ? `Search: "${searchQuery}"` : (tags ? `Tag: ${tags}` : (favoritesOnly ? 'Favorites' : (mediaType === 'video' ? 'Videos' : 'All Photos'))));
    const albumId = album?.id; // The actual ID of the album, to be used for fetching photos.

    const {
        data,
        fetchNextPage,
        hasNextPage,
        isFetchingNextPage,
        isFetching,
        refetch,
    } = useInfiniteQuery({
        // The query key for photos is now dependent on the actual albumId and date filters.
        // This ensures that if the albumId or date range changes, the photos are re-fetched.
        queryKey: ['albumPhotos', albumId, startDate, endDate, tags, mediaType, searchQuery, favoritesOnly],
        async queryFn({ pageParam = "" }) {
            const fromId = pageParam;

            // If an album context is intended, but we don't have the albumId yet, we can't fetch photos.
            // This is primarily handled by the `enabled` flag below, but serves as a safeguard.
            if (albumIdentifier && !albumId) {
                return { data: [], nextCursor: undefined };
            }

            let retrievedPhotos;
            if (searchQuery) {
                retrievedPhotos = await photosAdapter.searchPhotos(searchQuery, fromId, limit);
            } else if (tags) {
                retrievedPhotos = await photosAdapter.getPhotosByTag(tags, fromId, limit, false);
            } else if (mediaType === 'video') {
                retrievedPhotos = await photosAdapter.getVideos(fromId, limit);
            } else if (albumId) {
                retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, fromId, limit, false, startDate, endDate);
            } else {
                retrievedPhotos = await photosAdapter.getAllPhotosInfo(fromId, limit, false, startDate, endDate);
            }

            const photosData = retrievedPhotos ?? [];

            return {
                data: photosData,
                nextCursor: photosData.length === limit ? photosData[photosData.length - 1]?.id : undefined,
            };
        },
        // This query is enabled only if:
        // 1. We are not in an album context (albumIdentifier is null/undefined).
        // 2. We ARE in an album context AND we have successfully fetched the album's ID.
        // 3. For search queries, enabled whenever searchQuery is provided.
        enabled: (!albumIdentifier || (!!albumIdentifier && !!albumId)) && (!tags || !!tags) && (!!searchQuery || !searchQuery) && (!mediaType || !!mediaType) && (!favoritesOnly || !!favoritesOnly),
        initialPageParam: "",
        placeholderData: (previousData) => previousData,
        getNextPageParam: (lastPage) => lastPage.nextCursor,
        select: (data) => {
            let allPhotos = data.pages.flatMap(page => page.data).filter(Boolean);
            // Client-side safety filter for mediaType
            if (mediaType) {
                allPhotos = allPhotos.filter(p => p.mediaType === mediaType);
            }
            // Client-side filter for favorites
            if (favoritesOnly) {
                allPhotos = allPhotos.filter(p => p.favorite);
            }
            const uniquePhotos = Array.from(new Map(allPhotos.map(photo => [photo.id, photo])).values());
            return {
                pages: data.pages,
                pageParams: data.pageParams,
                photos: uniquePhotos,
            };
        }
    });

    const photos = data?.photos ?? [];

    return { photos, albumName, allRetrieved: !hasNextPage, fetchNextPage, isFetchingNextPage, isFetching, refetch };
};

export default usePhotoGrid;
