import React from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';

interface FavoritesViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const FavoritesView: React.FC<FavoritesViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    return (
        <PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} favoritesOnly />
    );
};
