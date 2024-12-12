import {useEffect, useState} from 'react';
import {Photo} from "../../Models/Photo";
import {Album} from "../../Models/Album";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

const useAlbumItem = (photosAdapter: IPhotosAdapter, source: Album) => {

    const [thumbnailUrl, setThumbnailUrl] = useState('');
    const [photoCount, setPhotoCount] = useState(0);
    const [isLoading, setIsLoading] = useState(true);

    const getPhotoCountInAlbum = async () => {
        const count = await photosAdapter.getPhotoCountInAlbum(source.id);
        setPhotoCount(count);
        setIsLoading(false);
    };

    const retrieveThumbnailUrl = async () => {
        const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(source.id, 1, 1, true);
        const firstPhotoInAlbum = (photos && photos.length > 0) ? photos[0].thumbnailPath : 'placeholder';
        setThumbnailUrl(firstPhotoInAlbum);
        setIsLoading(false);
    }

    useEffect(() => {
        retrieveThumbnailUrl();
        getPhotoCountInAlbum();
    },[]);

    return [{thumbnailUrl, photoCount, isLoading}]
};

export default useAlbumItem;