import {useEffect, useState} from 'react';
import {Photo} from "../../Models/Photo";
import {Album} from "../../Models/Album";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

const useAlbumCard = (photosAdapter: IPhotosAdapter, source: Album) => {

    const [thumbnailUrl, setThumbnailUrl] = useState('');
    const [photoCount, setPhotoCount] = useState();
    const [isLoading, setIsLoading] = useState(true);


    const retrieveThumbnailUrl = async () => {
        const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(source.id, "", 1, true);
        const firstPhotoInAlbum = (photos && photos.length > 0) ? photos[0].thumbnail : 'no_image.png';
        setThumbnailUrl(firstPhotoInAlbum);
    }

    const getPhotoCountInAlbum = async () => {
        const count = await photosAdapter.getPhotoCountInAlbum(source.id);
        setPhotoCount(count);
    };

    useEffect(() => {
        retrieveThumbnailUrl();
        getPhotoCountInAlbum();
        setIsLoading(false);
    },[]);

    return [{thumbnailUrl, photoCount, isLoading}]
};

export default useAlbumCard;