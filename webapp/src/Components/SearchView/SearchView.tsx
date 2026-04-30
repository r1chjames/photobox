import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center } from '@mantine/core';
import { IconSearch } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';

interface SearchViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const SearchView: React.FC<SearchViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    const [searchParams] = useSearchParams();
    const query = searchParams.get('q')?.trim() ?? '';
    const [results, setResults] = useState<Photo[]>([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (!query) {
            setResults([]);
            return;
        }
        setLoading(true);
        photosAdapter.searchPhotos(query, '', 100)
            .then(setResults)
            .catch(() => setResults([]))
            .finally(() => setLoading(false));
    }, [query, photosAdapter]);

    if (!query) {
        return (
            <EmptyState
                title="Search your photos"
                description="Type a keyword in the search box above to find photos by name, album, or metadata."
                icon={<IconSearch size="2rem" />}
            />
        );
    }

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    if (results.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Search: &quot;{query}&quot;</Title>
                <EmptyState
                    title="No results found"
                    description={`No photos match "${query}". Try a different keyword.`}
                    icon={<IconSearch size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">Search: &quot;{query}&quot; ({results.length} results)</Title>
            <PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />
        </div>
    );
};
