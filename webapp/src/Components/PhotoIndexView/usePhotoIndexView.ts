import {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, albumId: string | undefined) => {

    const [photos, setPhotos] = useState<Photo[]>();
    const [isApiCallsRunning, sApiCallsRunning] = useState(false);

    async function retrievePhotos() {
        let retrievedPhotos;
        if (typeof albumId !== 'undefined') {
            retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, 1, 100, true);
        } else {
            retrievedPhotos = await photosAdapter.getAllPhotosInfo(true);
        }
        setPhotos(retrievedPhotos);
        sApiCallsRunning(false);
    }

    useEffect(() => {
        retrievePhotos();
    },[]);

    return [{photos, isApiCallsRunning}]
};

export default usePhotoIndexView;