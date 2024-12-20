import {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, albumId: string | undefined) => {

    const [photos, setPhotos] = useState<Photo[]>([]);
    const [isApiCallsRunning, setIsApiCallsRunning] = useState(false);

    async function retrievePhotos(nextGroupKey: number, count: number) {
        let retrievedPhotos;
        if (typeof albumId !== 'undefined') {
            retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, nextGroupKey, count, true);
        } else {
            retrievedPhotos = await photosAdapter.getAllPhotosInfo(true);
        }
        setPhotos(retrievedPhotos);
        setIsApiCallsRunning(false);
    }

    useEffect(() => {
        retrievePhotos(1, 30);
    },[]);

    return [{photos, isApiCallsRunning, retrievePhotos}]
};

export default usePhotoIndexView;