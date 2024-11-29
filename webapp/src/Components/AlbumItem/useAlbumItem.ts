import {useEffect, useState} from 'react';
import {Photo} from "../../Models/Photo";
import {Album} from "../../Models/Album";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

const useAlbumItem = (photosAdapter: IPhotosAdapter, source: Album) => {

    const [thumbnailUrl, setThumbnailUrl] = useState('');
    const [photoCount, setPhotoCount] = useState(0);
    const [isLoading, setIsLoading] = useState(true);

    const getUrlOfFirstImageInAlbum = async () => {
        const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(source.id, 1, 1);
        if (photos && photos.length > 0) {
            return `photo/${photos[0].id}/thumbnail`;
        }
        return '';
    };

    const getPhotoCountInAlbum = async () => {
        const count = await photosAdapter.getPhotoCountInAlbum(source.id);
        setPhotoCount(count);
        setIsLoading(false);
    };

    const retrieveThumbnailUrl = async () => {
        const retrievedThumbnailUrl = await getUrlOfFirstImageInAlbum();
        setThumbnailUrl(retrievedThumbnailUrl);
        setIsLoading(false);
    }

    useEffect(() => {
        retrieveThumbnailUrl();
        getPhotoCountInAlbum();
    },[]);

    return [{thumbnailUrl, photoCount, loading: isLoading}]
};

export default useAlbumItem;