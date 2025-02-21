import {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";
import {useInfiniteQuery} from "@tanstack/react-query";

const usePhotoGrid = (photosAdapter: IPhotosAdapter, albumId: string) => {

    // const [photos, setPhotos] = useState<Photo[]>([]);
    const [allRetrieved, setAllRetrieved] = useState(false);
    const [page, setPage] = useState(1);
    const limit = 30;

    async function retrievePhotos() {
        let retrievedPhotos: Photo[];
        if (albumId !== 'undefined') {
            retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, page, limit, true);
        } else {
            retrievedPhotos = await photosAdapter.getAllPhotosInfo(page, limit, true);
        }
        if (retrievedPhotos && retrievedPhotos.length > 0) {
            setPhotos([...photos, ...retrievedPhotos])
        } else {
            setAllRetrieved(true);
        }
    }

    // useEffect(() => {
    //     retrievePhotos(1, 30);
    // },[]);

    const {data: photos,
        fetchNextPage,
        hasNextPage,
        isFetching,
        isFetchingNextPage} = useInfiniteQuery({
        queryKey: ['albumPhotos'],
        queryFn: retrievePhotos,
        initialPageParam: 0,
        getNextPageParam: (lastPage, pages) => lastPage.nextCursor,
    });

    return [{photos, allRetrieved, fetchNextPage, setPage}]
};

export default usePhotoGrid;