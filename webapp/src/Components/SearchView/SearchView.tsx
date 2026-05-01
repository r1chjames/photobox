import React from 'react';
import { useSearchParams } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { EmptyState } from '../EmptyState/EmptyState';
import { IconSearch } from '@tabler/icons-react';

interface SearchViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const SearchView: React.FC<SearchViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    const [searchParams] = useSearchParams();
    const query = searchParams.get('q')?.trim() ?? '';

    if (!query) {
        return (
            <EmptyState
                title="Search your photos"
                description="Type a keyword in the search box above to find photos by name, album, or metadata."
                icon={<IconSearch size="2rem" />}
            />
        );
    }

    return (
        <PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} searchQuery={query} />
    );
};
