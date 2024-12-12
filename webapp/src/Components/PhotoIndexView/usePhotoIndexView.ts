import React, {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, albumId: string) => {

    const [photos, setPhotos] = useState<Photo[]>();
    const [isLoading, setIsLoading] = useState(false);

    async function retrievePhotos() {
        const retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, 1, 100, true);
        setPhotos(retrievedPhotos);
        setIsLoading(false);
    }

    useEffect(() => {
        retrievePhotos();
    },[]);

    return [{photos, isLoading}]
};

export default usePhotoIndexView;