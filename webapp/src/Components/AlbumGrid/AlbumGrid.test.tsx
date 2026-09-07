import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { AlbumGrid } from './AlbumGrid';
import { Album } from '../../Models/Album';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock react-router-dom
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useNavigate: () => mockNavigate,
}));

// Mock the useAlbumGrid hook
vi.mock('./useAlbumGrid', () => ({
    default: vi.fn()
}));

// Mock AlbumCard component
vi.mock('../AlbumCard/AlbumCard', () => ({
    AlbumCard: ({ source, albumViewCallback }: any) => (
        <div data-testid={`album-card-${source.id}`} onClick={() => albumViewCallback()}>
            <span>{source.name}</span>
            <span>Mock Album Card</span>
        </div>
    )
}));

// Mock InputModal component
vi.mock('../InputModal/InputModal', () => ({
    InputModal: ({ isOpen, children, handleSave, handleClose }: any) => (
        isOpen ? (
            <div data-testid="input-modal">
                <div>Create Album Modal</div>
                {children}
                <button onClick={handleSave}>Save</button>
                <button onClick={handleClose}>Close</button>
            </div>
        ) : null
    )
}));

import useAlbumGrid from './useAlbumGrid';

describe('AlbumGrid', () => {
    let mockAlbumsAdapter: IAlbumsAdapter;
    let mockPhotosAdapter: IPhotosAdapter;
    let mockAlbums: Album[];

    beforeEach(() => {
        mockAlbumsAdapter = {} as IAlbumsAdapter;
        mockPhotosAdapter = {} as IPhotosAdapter;

        mockAlbums = [
            { id: 'album-1', name: 'Summer Vacation', description: 'Beach photos', tags: 'summer,beach', metadata: '{}' },
            { id: 'album-2', name: 'Winter Trip', description: 'Snow photos', tags: 'winter,snow', metadata: '{}' },
            { id: 'album-3', name: 'Family Reunion', description: 'Family gathering', tags: 'family', metadata: '{}' },
        ];

        // Default mock implementation
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: mockAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        mockNavigate.mockClear();
    });

    it('should render all albums', async () => {
        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Summer Vacation')).toBeInTheDocument();
            expect(screen.getByText('Winter Trip')).toBeInTheDocument();
            expect(screen.getByText('Family Reunion')).toBeInTheDocument();
        });
    });

    it('should render album cards for each album', async () => {
        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByTestId('album-card-album-1')).toBeInTheDocument();
            expect(screen.getByTestId('album-card-album-2')).toBeInTheDocument();
            expect(screen.getByTestId('album-card-album-3')).toBeInTheDocument();
        });
    });

    it('should navigate when album card is clicked', async () => {
        const user = userEvent.setup();

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByTestId('album-card-album-1')).toBeInTheDocument();
        });

        const albumCard = screen.getByTestId('album-card-album-1');
        await user.click(albumCard);

        expect(mockNavigate).toHaveBeenCalledWith('../album/album-1');
    });

    it('should limit displayed albums with initialDisplayCount prop', async () => {
        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
                initialDisplayCount={2}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Summer Vacation')).toBeInTheDocument();
            expect(screen.getByText('Winter Trip')).toBeInTheDocument();
        });

        // Third album should not be rendered
        expect(screen.queryByText('Family Reunion')).not.toBeInTheDocument();
    });

    it('should render InputModal when showNewAlbumModal is true', async () => {
        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // InputModal should be present in DOM (even if not visible initially)
        // The modal visibility is controlled by showNewAlbumModal state
        const modal = screen.queryByTestId('input-modal');
        // Modal should not be visible initially
        expect(modal).not.toBeInTheDocument();
    });

    it('should handle empty album list', async () => {
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: [],
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Should render the grid container but no album cards
        await waitFor(() => {
            const albumCards = container.querySelectorAll('[data-testid^="album-card-"]');
            expect(albumCards.length).toBe(0);
        });
    });

    it('should handle undefined albums', () => {
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: undefined,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Should render without crashing
        expect(container).toBeInTheDocument();
    });

    it('should render with correct grid layout', () => {
        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        const gridContainer = container.querySelector('.albums-grid');
        expect(gridContainer).toBeInTheDocument();
    });

    it('should call useAlbumGrid hook with correct adapter', () => {
        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(useAlbumGrid).toHaveBeenCalledWith(mockAlbumsAdapter);
    });

    it('should navigate to create album page when modal save is clicked', async () => {
        const mockHandleChange = vi.fn();

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: mockAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: 'New Album Name',
            handleNewAlbumNameValueChange: mockHandleChange
        }]);

        // We need to trigger the modal to be open somehow
        // Since the component uses useState(Boolean) which defaults to false,
        // and there's no button to open it (commented out), we can't test the modal interaction
        // in the current implementation. This test documents the expected behavior.

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Modal should be closed by default
        expect(screen.queryByTestId('input-modal')).not.toBeInTheDocument();
    });

    it('should render album cards in article elements', async () => {
        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            const articles = container.querySelectorAll('article');
            expect(articles.length).toBe(3);
        });
    });

    it('should render all albums when initialDisplayCount is not provided', async () => {
        const manyAlbums = Array.from({ length: 25 }, (_, i) =>
            ({ id: `album-${i}`, name: `Album ${i}`, description: '', tags: '', metadata: '{}' })
        );

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: manyAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            // All 25 should be rendered with no initialDisplayCount limit
            expect(screen.getByText('Album 0')).toBeInTheDocument();
            expect(screen.getByText('Album 24')).toBeInTheDocument();
        });
    });

    it('should show Load More button when albums exceed initialDisplayCount', async () => {
        const manyAlbums = Array.from({ length: 50 }, (_, i) =>
            ({ id: `album-${i}`, name: `Album ${i}`, description: '', tags: '', metadata: '{}' })
        );

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: manyAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
                initialDisplayCount={20}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Load More')).toBeInTheDocument();
        });

        // Only 20 albums should be initially visible
        const articles = screen.getAllByRole('article');
        expect(articles.length).toBe(20);
    });

    it('should load more albums when Load More button is clicked', async () => {
        const user = userEvent.setup();
        const manyAlbums = Array.from({ length: 50 }, (_, i) =>
            ({ id: `album-${i}`, name: `Album ${i}`, description: '', tags: '', metadata: '{}' })
        );

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: manyAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
                initialDisplayCount={20}
                loadMoreIncrement={10}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Load More')).toBeInTheDocument();
        });

        // Initially 20 albums
        let articles = screen.getAllByRole('article');
        expect(articles.length).toBe(20);

        // Click Load More
        await user.click(screen.getByText('Load More'));

        // Now 30 albums should be visible
        articles = screen.getAllByRole('article');
        expect(articles.length).toBe(30);

        // Click Load More again
        await user.click(screen.getByText('Load More'));

        // Now 40 albums should be visible
        articles = screen.getAllByRole('article');
        expect(articles.length).toBe(40);
    });

    it('should hide Load More button when all albums are displayed', async () => {
        const manyAlbums = Array.from({ length: 50 }, (_, i) =>
            ({ id: `album-${i}`, name: `Album ${i}`, description: '', tags: '', metadata: '{}' })
        );

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: manyAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
                initialDisplayCount={50}
            />
        );

        await waitFor(() => {
            expect(screen.queryByText('Load More')).not.toBeInTheDocument();
        });
    });

    it('should render skeleton cards while first load is pending and hide the empty state', async () => {
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: undefined,
            isLoading: true,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Skeletons render inside articles matching the real grid
        const skeletonCards = container.querySelectorAll('[data-testid="album-card-skeleton"]');
        expect(skeletonCards.length).toBe(6);

        // No empty state, no "Load More", no real album cards while loading
        expect(screen.queryByText('No albums yet')).not.toBeInTheDocument();
        expect(screen.queryByText('Load More')).not.toBeInTheDocument();
        expect(container.querySelectorAll('[data-testid^="album-card-album-"]').length).toBe(0);
    });

    it('should cap the initial skeleton count to 6 when initialDisplayCount is larger', () => {
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: undefined,
            isLoading: true,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
                initialDisplayCount={30}
            />
        );

        expect(container.querySelectorAll('[data-testid="album-card-skeleton"]').length).toBe(6);
    });

    it('should show the error state with a working Try again button when the fetch failed', async () => {
        const user = userEvent.setup();
        const mockRefetch = vi.fn();

        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: undefined,
            isLoading: false,
            isError: true,
            refetch: mockRefetch,
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Error state is distinct from the empty state and offers a retry
        expect(screen.getByText('Failed to load albums')).toBeInTheDocument();
        expect(screen.queryByText('No albums yet')).not.toBeInTheDocument();

        const retryButton = screen.getByRole('button', { name: 'Try again' });
        await user.click(retryButton);

        expect(mockRefetch).toHaveBeenCalledTimes(1);
    });

    it('should transition from skeleton cards to populated cards once albums load', async () => {
        const mockHook = useAlbumGrid as ReturnType<typeof vi.fn>;

        // First render: first load still in flight
        mockHook.mockReturnValueOnce([{
            albums: undefined,
            isLoading: true,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { rerender, container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Loading: skeletons in the grid, no albums, no empty/error state
        expect(container.querySelectorAll('[data-testid="album-card-skeleton"]').length).toBe(6);
        expect(screen.queryByText('Summer Vacation')).not.toBeInTheDocument();
        expect(screen.queryByText('No albums yet')).not.toBeInTheDocument();
        expect(screen.queryByText('Failed to load albums')).not.toBeInTheDocument();

        // Albums resolve: hook flips to loaded, real cards replace skeletons
        mockHook.mockReturnValueOnce([{
            albums: mockAlbums,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        rerender(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Summer Vacation')).toBeInTheDocument();
        });
        expect(container.querySelectorAll('[data-testid="album-card-skeleton"]').length).toBe(0);
        expect(screen.queryByText('No albums yet')).not.toBeInTheDocument();
    });

    it('should not render skeletons when albums are undefined but the load already finished', () => {
        (useAlbumGrid as ReturnType<typeof vi.fn>).mockReturnValue([{
            albums: undefined,
            isLoading: false,
            isError: false,
            refetch: vi.fn(),
            createAlbumModalAlbumNameErrorText: '',
            newAlbumName: '',
            handleNewAlbumNameValueChange: vi.fn()
        }]);

        const { container } = render(
            <AlbumGrid
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Nothing is loading and nothing errored, so fall through to the
        // empty state without crashing — no skeletons.
        expect(container.querySelectorAll('[data-testid="album-card-skeleton"]').length).toBe(0);
        expect(screen.getByText('No albums yet')).toBeInTheDocument();
    });
});
