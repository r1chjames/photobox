import {useEffect, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Photo} from "../../Models/Photo";

const usePhotoIndexView = (photosAdapter: IPhotosAdapter, albumId: string) => {

    const [photos, setPhotos] = useState<Photo[]>([]);
    const [allRetrieved, setAllRetrieved] = useState(false);

    async function retrievePhotos(nextGroupKey: number, count: number) {
        let retrievedPhotos: Photo[];
        console.log(`ngk: ${nextGroupKey} + cnt: ${count}`)
        if (albumId !== 'undefined') {
            retrievedPhotos = await photosAdapter.getPhotosInfoInAlbum(albumId, nextGroupKey, count, true);
        } else {
            retrievedPhotos = await photosAdapter.getAllPhotosInfo(nextGroupKey, count, true);
        }
        console.log(`retrieved: ${retrievedPhotos.length}`)
        if (retrievedPhotos.length > 0) {
            setPhotos([...photos, ...retrievedPhotos])
        } else {
            setAllRetrieved(true);
        }
    }

    useEffect(() => {
        retrievePhotos(1, 30);
    },[]);

    return [{photos, allRetrieved, retrievePhotos}]
};

export default usePhotoIndexView;