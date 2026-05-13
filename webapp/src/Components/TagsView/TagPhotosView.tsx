import React from 'react';
import { useParams } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { Title } from '@mantine/core';

interface TagPhotosViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const TagPhotosView: React.FC<TagPhotosViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    const { tag } = useParams<{ tag: string }>();
    const decodedTag = decodeURIComponent(tag || '');

    return (
        <div style={{ height: '100%' }}>
            <Title size="h4" mb="md">Tag: {decodedTag}</Title>
            <PhotoGrid
                photosAdapter={photosAdapter}
                albumsAdapter={albumsAdapter}
                sharesAdapter={sharesAdapter}
                tags={decodedTag}
            />
        </div>
    );
};
