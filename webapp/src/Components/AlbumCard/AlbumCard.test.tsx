import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { AlbumCard } from './AlbumCard';
import { Album } from '../../Models/Album';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';

// Mock the useAlbumCard hook
vi.mock('./useAlbumCard', () => ({
    default: vi.fn()
}));

import useAlbumCard from './useAlbumCard';

describe('AlbumCard', () => {
    let mockPhotosAdapter: IPhotosAdapter;
    let mockAlbumsAdapter: IAlbumsAdapter;
    let mockAlbum: Album;
    let mockAlbumViewCallback: (albumId: string) => void;

    beforeEach(() => {
        mockPhotosAdapter = {} as IPhotosAdapter;
        mockAlbumsAdapter = {
            updateAlbum: vi.fn().mockResolvedValue(undefined),
            deleteAlbum: vi.fn().mockResolvedValue(undefined),
        } as unknown as IAlbumsAdapter;
        mockAlbum = {
            id: 'album-1',
            name: 'Test Album',
            description: 'Test Description',
            tags: 'test,tags',
            metadata: '{}'
        };
        mockAlbumViewCallback = vi.fn();

        // Default mock implementation
        (useAlbumCard as ReturnType<typeof vi.fn>).mockReturnValue({
            thumbnailUrl: 'data:image/png;base64,mockthumb',
            photoCount: 10,
            isLoading: false
        });
    });

    it('should render album card with name and photo count', async () => {
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
            expect(screen.getByText('10 photos')).toBeInTheDocument();
        });
    });

    it('should display thumbnail image', async () => {
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            const image = screen.getByRole('img');
            expect(image).toBeInTheDocument();
            expect(image).toHaveAttribute('src', 'data:image/png;base64,mockthumb');
        });
    });

    it('should show loader when loading', () => {
        (useAlbumCard as ReturnType<typeof vi.fn>).mockReturnValue({
            thumbnailUrl: '',
            photoCount: 0,
            isLoading: true
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        const loader = container.querySelector('.mantine-Loader-root');
        expect(loader).toBeInTheDocument();
    });

    it('should call albumViewCallback when clicked', async () => {
        const user = userEvent.setup();

        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
        });

        const nameText = screen.getByText('Test Album');
        await user.click(nameText);

        expect(mockAlbumViewCallback).toHaveBeenCalledWith('album-1');
    });

    it('should handle thumbnail URL with http prefix', async () => {
        (useAlbumCard as ReturnType<typeof vi.fn>).mockReturnValue({
            thumbnailUrl: 'http://example.com/image.jpg',
            photoCount: 5,
            isLoading: false
        });

        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            const image = screen.getByRole('img');
            expect(image).toHaveAttribute('src', 'http://example.com/image.jpg');
        });
    });

    it('should handle thumbnail URL with data:image prefix', async () => {
        (useAlbumCard as ReturnType<typeof vi.fn>).mockReturnValue({
            thumbnailUrl: 'data:image/jpeg;base64,testdata',
            photoCount: 3,
            isLoading: false
        });

        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            const image = screen.getByRole('img');
            expect(image).toHaveAttribute('src', 'data:image/jpeg;base64,testdata');
        });
    });

    it('should use thumbnail URL from hook directly', async () => {
        (useAlbumCard as ReturnType<typeof vi.fn>).mockReturnValue({
            thumbnailUrl: 'blob:http://localhost/abc123',
            photoCount: 7,
            isLoading: false
        });

        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            const image = screen.getByRole('img');
            expect(image).toHaveAttribute('src', 'blob:http://localhost/abc123');
        });
    });

    it('should display photo count badge', async () => {
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            const badge = screen.getByText('10 photos');
            expect(badge).toBeInTheDocument();
            // Badge should be rendered within a badge component
            expect(badge.closest('.mantine-Badge-root')).toBeInTheDocument();
        });
    });

    it('should render with different album name', async () => {
        const differentAlbum = {
            id: 'album-2',
            name: 'Vacation Photos',
            description: 'Summer trip',
            tags: 'vacation,summer',
            metadata: '{}'
        };

        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={differentAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Vacation Photos')).toBeInTheDocument();
        });
    });

    it('should call useAlbumCard hook with correct parameters', () => {
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        expect(useAlbumCard).toHaveBeenCalledWith(mockPhotosAdapter, mockAlbum);
    });
    it('should render Smart badge for smart albums', async () => {
        mockAlbum.metadata = JSON.stringify({ smart: true, rules: { camera: 'Canon' } });
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
            expect(screen.getByText('Smart')).toBeInTheDocument();
        });
    });

    it('should not render Smart badge for regular albums', async () => {
        mockAlbum.metadata = '{}';
        render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
        });
        expect(screen.queryByText('Smart')).not.toBeInTheDocument();
    });
});


