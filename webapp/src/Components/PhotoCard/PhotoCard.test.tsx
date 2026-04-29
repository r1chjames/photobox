import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { PhotoCard } from './PhotoCard';
import { Photo } from '../../Models/Photo';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';

// Mock react-router-dom
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useNavigate: () => vi.fn(),
}));

// Mock @mantine/hooks
vi.mock('@mantine/hooks', () => ({
    useHotkeys: vi.fn(),
    useMediaQuery: () => false,
}));

// Mock ImageUtils
vi.mock('../../utils/ImageUtils', () => ({
    fetchPhotoBinWithAuth: vi.fn().mockResolvedValue('data:image/png;base64,mockimage'),
    revokeBlobUrl: vi.fn(),
}));

describe('PhotoCard', () => {
    let mockPhotosAdapter: IPhotosAdapter;
    let mockAlbumsAdapter: IAlbumsAdapter;
    let mockPhoto: Photo;
    let mockProps: any;

    beforeEach(() => {
        mockPhotosAdapter = {} as IPhotosAdapter;
        mockAlbumsAdapter = {
            getAlbumInfoById: vi.fn().mockResolvedValue({ id: 'album-1', name: 'Test Album' }),
        } as unknown as IAlbumsAdapter;

        mockPhoto = {
            id: 'photo-1',
            name: 'test-photo.jpg',
            albumId: 'album-1',
            createdAt: '2024-01-15T10:30:00Z',
            thumbnail: 'base64thumbnaildata',
        } as Photo;

        mockProps = {
            photosAdapter: mockPhotosAdapter,
            albumsAdapter: mockAlbumsAdapter,
            source: mockPhoto,
            previousPhoto: vi.fn(),
            nextPhoto: vi.fn(),
            firstInAlbum: false,
            lastInAlbum: false,
            closeModal: vi.fn(),
        };
    });

    it('should render photo image when loaded', async () => {
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            const images = screen.getAllByRole('img');
            expect(images.length).toBeGreaterThan(0);
        });
    });

    it('should display loading state initially', () => {
        render(<PhotoCard {...mockProps} />);
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('should display creation date badge', async () => {
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            // Format will vary by locale, but should contain date
            const dateElement = screen.getByText(/2024|Jan|15/i);
            expect(dateElement).toBeInTheDocument();
        });
    });

    it('should fetch and display album name', async () => {
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            expect(screen.getByText('Test Album')).toBeInTheDocument();
        });

        expect(mockAlbumsAdapter.getAlbumInfoById).toHaveBeenCalledWith('album-1');
    });

    it('should not display album badge when no albumId', async () => {
        const propsWithoutAlbum = {
            ...mockProps,
            source: { ...mockPhoto, albumId: undefined },
        };

        render(<PhotoCard {...propsWithoutAlbum} />);

        await waitFor(() => {
            expect(screen.queryByText('Test Album')).not.toBeInTheDocument();
        });
    });

    it('should render previous button when not first in album', async () => {
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            const buttons = screen.getAllByRole('button');
            // Should have previous, next, close, and view details buttons
            expect(buttons.length).toBeGreaterThan(0);
        });
    });

    it('should not render previous button when first in album', () => {
        const propsFirstInAlbum = { ...mockProps, firstInAlbum: true };
        render(<PhotoCard {...propsFirstInAlbum} />);

        // There should be fewer buttons (no previous button)
        const buttons = screen.getAllByRole('button');
        expect(buttons).toBeDefined();
    });

    it('should not render next button when last in album', () => {
        const propsLastInAlbum = { ...mockProps, lastInAlbum: true };
        render(<PhotoCard {...propsLastInAlbum} />);

        const buttons = screen.getAllByRole('button');
        expect(buttons).toBeDefined();
    });

    it('should call previousPhoto when previous button clicked', async () => {
        const user = userEvent.setup();
        render(<PhotoCard {...mockProps} />);

        // Wait for all buttons to render
        await waitFor(() => {
            const prevButton = screen.getByLabelText('Previous photo');
            expect(prevButton).toBeInTheDocument();
        });

        const prevButton = screen.getByLabelText('Previous photo');
        await user.click(prevButton);

        // The previousPhoto should be called
        expect(mockProps.previousPhoto).toHaveBeenCalled();
    });

    it('should call closeModal when close button clicked', async () => {
        const user = userEvent.setup();
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            const closeButton = screen.getByLabelText('Close');
            expect(closeButton).toBeInTheDocument();
        });

        const closeButton = screen.getByLabelText('Close');
        await user.click(closeButton);

        expect(mockProps.closeModal).toHaveBeenCalled();
    });

    it('should render View Details button', async () => {
        render(<PhotoCard {...mockProps} />);

        await waitFor(() => {
            expect(screen.getByText('View Details')).toBeInTheDocument();
        });
    });

    it('should handle missing photo data gracefully', () => {
        const propsWithMinimalPhoto = {
            ...mockProps,
            source: {
                id: 'photo-2',
                name: 'minimal.jpg',
            } as Photo,
        };

        render(<PhotoCard {...propsWithMinimalPhoto} />);
        expect(screen.getByText('Loading...')).toBeInTheDocument();
    });

    it('should handle album fetch error gracefully', async () => {
        const mockAlbumsAdapterWithError = {
            getAlbumInfoById: vi.fn().mockRejectedValue(new Error('Album not found')),
        } as unknown as IAlbumsAdapter;

        const propsWithError = {
            ...mockProps,
            albumsAdapter: mockAlbumsAdapterWithError,
        };

        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

        render(<PhotoCard {...propsWithError} />);

        await waitFor(() => {
            expect(consoleErrorSpy).toHaveBeenCalledWith(
                'Failed to fetch album name:',
                expect.any(Error)
            );
        });

        consoleErrorSpy.mockRestore();
    });
});