import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, fireEvent } from '../../test/test-utils';
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

// Mock @tanstack/react-virtual — in jsdom there is no real layout, so the
// virtualizer reports every row as visible. This lets the tests assert on the
// rendered photos without a real scroll container.
vi.mock('@tanstack/react-virtual', () => ({
    useVirtualizer: ({ count }: any) => ({
        getVirtualItems: () => Array.from({ length: count }, (_, index) => ({ key: index, index, start: index * 100, size: 100 })),
        getTotalSize: () => count * 100,
        scrollToIndex: vi.fn(),
        measure: vi.fn(),
    }),
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
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
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

    it('should render virtualized grid', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByTestId('virtual-grid')).toBeInTheDocument();
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
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
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
                60, // limit
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
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
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
            expect(screen.getByTestId('virtual-grid')).toBeInTheDocument();
        });
    });

    it('should display skeleton loaders for images while loading', async () => {
        render(
            <PhotoGrid
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            const grid = screen.getByTestId('virtual-grid');
            expect(grid).toBeInTheDocument();
        });
    });

    it('should render all photos in the virtualized grid', async () => {
        const manyPhotos = Array.from({ length: 65 }, (_, i) => ({
            id: `photo-${i}`,
            name: `photo${i}.jpg`,
            thumbnailUrl: `/api/photo/thumbnail/photo-${i}`,
            createdAt: '2024-01-01',
        }));

        const manyPhotosAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockResolvedValue(manyPhotos),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
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
    });

    it('shows photo-shaped skeleton tiles, not a spinning loader, while the next page is fetching', async () => {
        // First call resolves a full page (60 photos => nextCursor set);
        // the second page hangs on a deferred promise so we can assert the
        // mid-fetch state.
        const pageOne = Array.from({ length: 60 }, (_, i) => ({
            id: `photo-${i}`,
            name: `photo${i}.jpg`,
            thumbnailUrl: `/api/photo/thumbnail/photo-${i}`,
            createdAt: '2024-01-01',
        }));
        // Second page hangs forever so we can assert the mid-fetch state.
        const secondPagePromise = new Promise(() => {});
        const pagingAdapter = {
            getPhotosInfoInAlbum: vi.fn()
                .mockResolvedValueOnce(pageOne)
                .mockImplementationOnce(() => secondPagePromise),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
        } as unknown as IPhotosAdapter;

        const { container } = render(
            <PhotoGrid
                photosAdapter={pagingAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        // First page rendered, so gridColumns > 0 (no initial skeleton).
        await waitFor(() => {
            expect(screen.getAllByRole('img').length).toBe(60);
        });

        const scrollContainer = container.querySelector('[class^="density-"]') as HTMLElement;
        expect(scrollContainer).not.toBeNull();
        fireEvent.scroll(scrollContainer);

        await waitFor(() => {
            expect(screen.getByTestId('pagination-skeleton')).toBeInTheDocument();
        });

        // Photo-shaped skeleton tiles inside the masonry grid — never the
        // spinning pagination Loader (issue #173).
        const grid = screen.getByTestId('virtual-grid');
        expect(grid.querySelector('[data-testid="pagination-skeleton"]')).not.toBeNull();
        expect(grid.querySelectorAll('.mantine-Skeleton-root').length).toBeGreaterThan(0);
        expect(container.querySelector('.mantine-Loader-root')).toBeNull();
    });

    it('replaces pagination skeleton tiles with new photos when the next page resolves', async () => {
        const pageOne = Array.from({ length: 60 }, (_, i) => ({
            id: `photo-${i}`,
            name: `photo${i}.jpg`,
            thumbnailUrl: `/api/photo/thumbnail/photo-${i}`,
            createdAt: '2024-01-01',
        }));
        const pageTwo = Array.from({ length: 10 }, (_, i) => ({
            id: `photo-${i + 60}`,
            name: `photo${i + 60}.jpg`,
            thumbnailUrl: `/api/photo/thumbnail/photo-${i + 60}`,
            createdAt: '2024-01-02',
        }));
        const pagingAdapter = {
            getPhotosInfoInAlbum: vi.fn()
                .mockResolvedValueOnce(pageOne)
                .mockResolvedValueOnce(pageTwo),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
        } as unknown as IPhotosAdapter;

        const { container } = render(
            <PhotoGrid
                photosAdapter={pagingAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getAllByRole('img').length).toBe(60);
        });

        const scrollContainer = container.querySelector('[class^="density-"]') as HTMLElement;
        fireEvent.scroll(scrollContainer);

        // Page 2 lands: the skeleton batch unmounts atomically and the new
        // photo tiles render in its place (issue #173).
        await waitFor(() => {
            expect(screen.getAllByRole('img').length).toBe(70);
        });
        expect(screen.queryByTestId('pagination-skeleton')).toBeNull();
        expect(screen.getByAltText('photo60.jpg')).toBeInTheDocument();
    });

    it('renders no canvas and falls back to a skeleton for a tile with an invalid blurhash', async () => {
        const invalidBlurhashAdapter = {
            getPhotosInfoInAlbum: vi.fn().mockResolvedValue([
                {
                    id: 'photo-bad-hash',
                    name: 'badhash.jpg',
                    thumbnailUrl: '/api/photo/thumbnail/photo-bad-hash',
                    createdAt: '2024-01-01',
                    blurhash: 'invalid!!!',
                },
            ]),
            getThumbnailUrl: vi.fn((photoId: string) => `/api/photo/thumbnail/${photoId}`),
            getPhotoThumbnailBlob: vi.fn().mockResolvedValue('data:image/png;base64,mockthumb'),
        } as unknown as IPhotosAdapter;

        render(
            <PhotoGrid
                photosAdapter={invalidBlurhashAdapter}
                albumsAdapter={mockAlbumsAdapter}
            />
        );

        await waitFor(() => {
            expect(screen.getByAltText('badhash.jpg')).toBeInTheDocument();
        });

        // The invalid hash must never produce an (empty) canvas tile; the
        // placeholder falls back to the Mantine skeleton (issue #173).
        const grid = screen.getByTestId('virtual-grid');
        expect(grid.querySelector('canvas')).toBeNull();
        expect(grid.querySelector('.mantine-Skeleton-root')).not.toBeNull();
    });
});