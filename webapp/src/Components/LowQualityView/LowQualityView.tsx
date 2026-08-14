import React from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';

interface LowQualityViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

/**
 * Surfaces photos flagged as low quality (blurry, badly exposed, mostly
 * solid) by the deterministic quality analyzer. Read-only review surface —
 * users restore/hide photos from the grid, nothing is auto-deleted.
 */
export const LowQualityView: React.FC<LowQualityViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    return (
        <PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} lowQualityOnly />
    );
};
