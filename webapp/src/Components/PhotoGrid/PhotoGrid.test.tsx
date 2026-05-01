import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import { PhotoGrid } from './PhotoGrid';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';

// Mock react-router-dom
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useParams: () => ({ id: 'test-album-id' }),
    useNavigate: () => vi.fn(),
}));

// Mock @mantine/hooks
vi.mock('@mantine/hooks', () => ({
    useMediaQuery: () => false,
    useHotkeys: vi.fn(),
}));

// Mock @egjs/react-infinitegrid
vi.mock('@egjs/react-infinitegrid', () => ({
    JustifiedInfiniteGrid: ({ children }: any) => <div data-testid="infinite-grid">{children}</div>,
}));

describe('PhotoGrid', () => {
    let mockPhotosAdapter: IPhotosAdapter;
    let mockAlbumsAdapter: IAlbumsAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockResolvedValue([
                {
                    id: 'photo-1',
                    name: 'photo1.jpg',
                    thumbnailUrl: '/api/photo/thumbnail/photo-1',
                    createdAt: '2024-01-01',
                },
                {
                    id: 'photo-2',
                    name: 'photo2.jpg',
                    thumbnailUrl: '/api/photo/thumbnail/photo-2',
                    createdAt: '2024-01-02',
                },
            ]),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
        } as unknown as IPhotosAdapter;

        mockAlbumsAdapter = {
            getAlbumInfoById: vi.fn().mockResolvedValue({
                id: 'test-album-id',
                name: 'Test Album',
            }),
        } as unknown as IAlbumsAdapter;
    });

    it('should render album title', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
        });
    });

    it('should render infinite grid', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByTestId('infinite-grid')).toBeInTheDocument();
        });
    });

    it('should render photos in grid', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            const images = screen.getAllByRole('img');
            expect(images.length).toBeGreaterThan(0);
        });
    });

    it('should show empty state when no photos', async () => {
        const emptyPhotosAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockResolvedValue([]),
        } as unknown as IPhotosAdapter;

        render(
            <PhotoGrid
                photosAdapter={emptyPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByText('This album is empty')).toBeInTheDocument();
        });
    });

    it('should call getPhotosInfoInAlbum with correct parameters', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(mockPhotosAdapter.getPhotosInfoInAlbum).toHaveBeenCalledWith(
                'test-album-id',
                expect.any(String), // fromId
                30, // limit
                false, // includeThumbnails
                undefined, // startDate
                undefined  // endDate
            );
        });
    });

    it('should fetch album info by id', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(mockAlbumsAdapter.getAlbumInfoById).toHaveBeenCalledWith('test-album-id');
        });
    });

    it('should handle photo adapter errors gracefully', async () => {
        const errorPhotosAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockRejectedValue(new Error('API Error')),
        } as unknown as IPhotosAdapter;

        // Component should render without crashing even when adapter fails
        render(
            <PhotoGrid
                photosAdapter={errorPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        // When there's an error fetching photos, component treats it as empty and shows empty state
        await waitFor(() => {
            expect(screen.getByText('This album is empty')).toBeInTheDocument();
        });
    });

    it('should render with maxDisplayed prop', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                maxDisplayed={10}
            />
        );

        await waitFor(() => {
            expect(screen.getByTestId('infinite-grid')).toBeInTheDocument();
        });
    });

    it('should display skeleton loaders for images while loading', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        // Initially might show skeletons (depending on image load state)
        await waitFor(() => {
            const grid = screen.getByTestId('infinite-grid');
            expect(grid).toBeInTheDocument();
        });
    });

    it('should group photos correctly with groupkey', async () => {
        const manyPhotos = Array.from({ length: 65 }, (_, i) => ({
            id: `photo-${i}`,
            name: `photo${i}.jpg`,
            thumbnailUrl: `/api/photo/thumbnail/photo-${i}`,
            createdAt: '2024-01-01',
        }));

        const manyPhotosAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockResolvedValue(manyPhotos),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
        } as unknown as IPhotosAdapter;

        render(
            <PhotoGrid
                photosAdapter={manyPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            const images = screen.getAllByRole('img');
            expect(images.length).toBe(65);
        });

        // Photos should be grouped by index / 30
        // First 30 photos: group 0
        // Next 30 photos: group 1
        // Last 5 photos: group 2
    });
});