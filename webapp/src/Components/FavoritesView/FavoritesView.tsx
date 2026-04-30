import React, { useEffect, useState, useCallback } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center } from '@mantine/core';
import { IconHeart } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';

interface FavoritesViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
}

export const FavoritesView: React.FC<FavoritesViewProps> = ({ photosAdapter, albumsAdapter }) => {
    const [photos, setPhotos] = useState<Photo[]>([]);
    const [loading, setLoading] = useState(true);

    const loadFavorites = useCallback(async () => {
        try {
            // Client-side filter until backend supports ?favorites=true
            const all = await photosAdapter.getAllPhotosInfo('', 1000, true);
            setPhotos(all.filter(p => p.favorite));
        } catch (e) {
            setPhotos([]);
        } finally {
            setLoading(false);
        }
    }, [photosAdapter]);

    useEffect(() => {
        loadFavorites();
    }, [loadFavorites]);

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    if (photos.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Favorites</Title>
                <EmptyState
                    title="No favorites yet"
                    description="Photos you mark as favorite will appear here. Click the heart icon on any photo to add it."
                    icon={<IconHeart size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">Favorites ({photos.length})</Title>
            <PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} />
        </div>
    );
};
