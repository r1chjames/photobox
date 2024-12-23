import {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, albumId: string) => {

    const [photos, setPhotos] = useState<Photo[]>([]);

    async function retrievePhotos(nextGroupKey: number, count: number) {
        let retrievedPhotos: Photo[];
        if (albumId !== 'undefined') {
            retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, nextGroupKey, count, true);
        } else {
            retrievedPhotos = await photosAdapter.getAllPhotosInfo(nextGroupKey, count, true);
        }
        console.log(retrievedPhotos.length)
        setPhotos(retrievedPhotos.length > 0 ? retrievedPhotos : photos);
    }

    useEffect(() => {
        retrievePhotos(1, 30);
    },[]);

    return [{photos, retrievePhotos}]
};

export default usePhotoIndexView;