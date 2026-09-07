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

const mockCoverPhoto = {
    id: 'photo-1',
    name: 'p1',
    albumId: 'a',
    thumbnailUrl: '',
    blurhash: 'LEHV6nWB2yk8pyo0adR*.7kCMdnj',
    dominantColor: '#667788',
};

const useAlbumCardMock = useAlbumCard as ReturnType<typeof vi.fn>;

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

        // Default mock implementation: resolved populated state
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'data:image/png;base64,mockthumb',
            subBlobUrls: [],
            photoCount: 10,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
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

    it('keeps full card shell and title while metadata loads', () => {
        useAlbumCardMock.mockReturnValue({
            photos: [],
            heroBlobUrl: undefined,
            subBlobUrls: [],
            photoCount: 0,
            isMetadataLoading: true,
            isCountLoading: true,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        // The card shell and title are present while metadata loads.
        expect(screen.getByText('Test Album')).toBeInTheDocument();
        expect(container.querySelector('.album-card')).toBeInTheDocument();
        // No spinner: media slots show skeleton frames instead.
        expect(container.querySelector('.mantine-Loader-root')).not.toBeInTheDocument();
        expect(container.querySelector('.album-media-grid .mantine-Skeleton-root')).toBeInTheDocument();
    });

    it('renders blurhash placeholder beneath the hero image in populated state', async () => {
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'data:image/png;base64,mockthumb',
            subBlobUrls: [],
            photoCount: 10,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            // Blurhash decodes to a real canvas placeholder under the img.
            expect(container.querySelector('.album-main-thumb canvas')).toBeInTheDocument();
            const heroImage = screen.getByRole('img');
            expect(heroImage).toHaveAttribute('alt', 'Test Album');
        });
        expect(screen.getByText('10 photos')).toBeInTheDocument();
    });

    it('renders sub thumbnails as images when album has 2+ photos', async () => {
        useAlbumCardMock.mockReturnValue({
            photos: [
                { ...mockCoverPhoto },
                { ...mockCoverPhoto, id: 'photo-2', name: 'p2' },
                { ...mockCoverPhoto, id: 'photo-3', name: 'p3' },
            ],
            heroBlobUrl: 'data:image/png;base64,mockthumb',
            subBlobUrls: [
                'data:image/png;base64,mocksub1',
                'data:image/png;base64,mocksub2',
            ],
            photoCount: 3,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        await waitFor(() => {
            // Hero + 2 sub images, each sub-slot filled.
            expect(container.querySelectorAll('.album-sub-thumb img')).toHaveLength(2);
            expect(container.querySelector('.album-main-thumb img')).toBeInTheDocument();
        });
        const subImages = container.querySelectorAll('.album-sub-thumb img');
        expect(subImages[0]).toHaveAttribute('src', 'data:image/png;base64,mocksub1');
        expect(subImages[1]).toHaveAttribute('src', 'data:image/png;base64,mocksub2');
    });

    it('shows IconPhotoOff and zero photo count for an empty album', () => {
        useAlbumCardMock.mockReturnValue({
            photos: [],
            heroBlobUrl: undefined,
            subBlobUrls: [],
            photoCount: 0,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        expect(screen.getByText('0 photos')).toBeInTheDocument();
        expect(screen.getByText('Test Album')).toBeInTheDocument();
        expect(container.querySelector('.album-main-thumb svg')).toBeInTheDocument();
        expect(container.querySelector('.album-main-thumb img')).not.toBeInTheDocument();
    });

    it('shows hero only with empty sub slots for a one-photo album', () => {
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'data:image/png;base64,mockthumb',
            subBlobUrls: [],
            photoCount: 1,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        expect(container.querySelectorAll('.album-sub-thumb img')).toHaveLength(0);
        expect(container.querySelectorAll('.album-sub-thumb > div')).toHaveLength(2);
        expect(screen.getByText('1 photos')).toBeInTheDocument();
    });

    it('renders shell and title without crashing when the query errors', () => {
        useAlbumCardMock.mockReturnValue({
            photos: [],
            heroBlobUrl: undefined,
            subBlobUrls: [],
            photoCount: 0,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: true,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        expect(screen.getByText('Test Album')).toBeInTheDocument();
        expect(container.querySelector('.album-card')).toBeInTheDocument();
        // Failed metadata: unknown contents — static skeleton, no '0 photos'.
        expect(container.querySelector('.album-main-thumb .mantine-Skeleton-root')).toBeInTheDocument();
        expect(container.querySelector('.album-main-thumb svg')).not.toBeInTheDocument();
        expect(container.querySelector('.album-main-thumb img')).not.toBeInTheDocument();
    });

    it('never flashes a zero count while the count is loading', () => {
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'data:image/png;base64,mockthumb',
            subBlobUrls: [],
            photoCount: 0,
            isMetadataLoading: false,
            isCountLoading: true,
            isError: false,
        });

        const { container } = render(
            <AlbumCard
                photosAdapter={mockPhotosAdapter}
                albumsAdapter={mockAlbumsAdapter}
                source={mockAlbum}
                albumViewCallback={mockAlbumViewCallback}
            />
        );

        // No '0 photos' flash: a skeleton pill occupies the count slot.
        expect(screen.queryByText('0 photos')).not.toBeInTheDocument();
        expect(container.querySelector('.album-info-overlay .mantine-Skeleton-root')).toBeInTheDocument();
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
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'http://example.com/image.jpg',
            subBlobUrls: [],
            photoCount: 5,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
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
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'data:image/jpeg;base64,testdata',
            subBlobUrls: [],
            photoCount: 3,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
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
        useAlbumCardMock.mockReturnValue({
            photos: [{ ...mockCoverPhoto }],
            heroBlobUrl: 'blob:http://localhost/abc123',
            subBlobUrls: [],
            photoCount: 7,
            isMetadataLoading: false,
            isCountLoading: false,
            isError: false,
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
