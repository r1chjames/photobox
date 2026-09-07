import {useState} from 'react';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {useQuery} from "@tanstack/react-query";

const useAlbumGrid = (albumsAdapter: IAlbumsAdapter) => {

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
        return albumsAdapter.getAllAlbumsInfo();
    };

    const {data: albums, isLoading, isError, refetch} = useQuery({
        queryKey: ['getAllAlbums'],
        queryFn: getAllAlbums
    });

    return [{albums, isLoading, isError, refetch, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}]
};

export default useAlbumGrid;