import {useEffect, useState} from 'react';
import {Album} from "../../Models/Album";
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";

const useAlbumIndexView = (albumsAdapter: AlbumsAdapter) => {

    const [albums, setAlbums] = useState<Album[]>([]);
    const [createAlbumModalAlbumNameErrorText, setCreateAlbumModalAlbumNameErrorText] = useState('Required');
    const [newAlbumName, setNewAlbumName] = useState('');

    const handleNewAlbumNameValueChange = (fieldValue: string) => {
        if (fieldValue.length > 1) {
            setCreateAlbumModalAlbumNameErrorText('');
            setNewAlbumName(fieldValue);
        } else {
            setCreateAlbumModalAlbumNameErrorText('Invalid length');
        }
    };


    const getAllAlbums = async() => {
        const albumSources: Album[] = await albumsAdapter.getAllAlbumsInfo();
        setAlbums(albumSources);
    };

    useEffect(() => {
        getAllAlbums();
    },[]);

    return [{albums, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}]
};

export default useAlbumIndexView;