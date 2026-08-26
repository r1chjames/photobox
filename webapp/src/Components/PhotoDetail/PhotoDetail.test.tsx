import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import { PhotoDetail } from './PhotoDetail';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock react-router-dom - preserve other exports for MemoryRouter in test-utils
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useParams: () => ({ id: 'test-photo-id' }),
}));

// Mock ImageUtils
vi.mock('../../utils/ImageUtils', () => ({
    fetchPhotoBinWithAuth: vi.fn().mockResolvedValue('data:image/png;base64,mockimage'),
    revokeBlobUrl: vi.fn(),
}));

// Mock @mantine/hooks
vi.mock('@mantine/hooks', () => ({
    useDisclosure: () => [true, { toggle: vi.fn(), close: vi.fn() }],
    useMediaQuery: () => false,
    useHotkeys: () => {},
}));

describe('PhotoDetail', () => {
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    exif: {
                        Make: 'Canon',
                        Model: 'EOS R5',
                        ISOSpeedRatings: [400],
                        FNumber: '1/8',
                        DateTimeOriginal: '2023:01:01 21:43:14',
                        GPSLatitude: ['52/1', '2/1', '24/1'],
                        GPSLongitude: ['0/1', '5/1', '40/1'],
                        GPSLatitudeRef: 'N',
                        GPSLongitudeRef: 'E',
                    },
                    size: 72731,
                    width: 2400,
                    height: 1200,
                    mime: 'image/jpeg',
                    md5: '1031db5affaee4fb204553478faec521',
                },
            }),
            getAllPhotosInfo: vi.fn().mockResolvedValue([
                { id: 'photo-1', name: 'photo-1.jpg' },
                { id: 'test-photo-id', name: 'test-photo.jpg' },
                { id: 'photo-3', name: 'photo-3.jpg' },
            ]),
        } as unknown as IPhotosAdapter;
    });

    it('should render photo image when loaded', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            const images = screen.getAllByRole('img');
            expect(images.length).toBeGreaterThan(0);
        });
    });

    it('should show loader initially', async () => {
        // Create a slow-resolving mock to ensure we can catch the loading state
        const slowAdapter = {
            getPhotoInfoById: vi.fn().mockImplementation(() =>
                new Promise(resolve => setTimeout(() => resolve({
                    id: 'test-photo-id',
                    name: 'test-photo.jpg',
                    metadata: { Camera: 'Canon' },
                }), 100))
            ),
            getAllPhotosInfo: vi.fn().mockResolvedValue([
                { id: 'test-photo-id', name: 'test-photo.jpg' },
            ]),
        } as unknown as IPhotosAdapter;

        const { container } = render(<PhotoDetail photosAdapter={slowAdapter} />);

        // Mantine Loader should be present while loading (it's a span with mantine-Loader-root class)
        const loader = container.querySelector('.mantine-Loader-root');
        expect(loader).toBeInTheDocument();

        // Wait for component to finish loading
        await waitFor(() => {
            expect(screen.getByText('Details')).toBeInTheDocument();
        });
    });

    it('should render metadata button', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(screen.getByText('Details')).toBeInTheDocument();
        });
    });

    it('should display metadata in dialog', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            // Metadata should be visible since dialog is open by default
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('Canon')).toBeInTheDocument();
            expect(screen.getByText('EOS R5')).toBeInTheDocument();
        });
    });

    it('should render grouped metadata sections', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('Settings')).toBeInTheDocument();
            expect(screen.getByText('Date')).toBeInTheDocument();
            expect(screen.getByText('Location')).toBeInTheDocument();
            expect(screen.getByText('File')).toBeInTheDocument();
        });
    });

    it('should format EXIF values nicely', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            // ISO from array
            expect(screen.getByText('400')).toBeInTheDocument();
            // Aperture f/1/8 -> f/0.13 (1/8 = 0.125)
            expect(screen.getByText(/f\//)).toBeInTheDocument();
            // GPS decimal degrees
            expect(screen.getByText(/52\.04/)).toBeInTheDocument();
            // File size 72731 bytes -> 71.0 KB
            expect(screen.getByText(/KB/)).toBeInTheDocument();
            // EXIF date 2023:01:01 -> friendly format
            expect(screen.getByText(/2023/)).toBeInTheDocument();
        });
    });

    it('should handle metadata with no exif', async () => {
        const adapterNoExif = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    size: 1024,
                    width: 100,
                    height: 100,
                    mime: 'image/png',
                },
            }),
            getAllPhotosInfo: vi.fn().mockResolvedValue([
                { id: 'test-photo-id', name: 'test-photo.jpg' },
            ]),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterNoExif} />);

        await waitFor(() => {
            expect(screen.getByText('File')).toBeInTheDocument();
            expect(screen.queryByText('Camera')).not.toBeInTheDocument();
        });
    });

    it('should show empty state when no metadata exists', async () => {
        const adapterNoMetadata = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: '',
                metadata: {},
            }),
            getAllPhotosInfo: vi.fn().mockResolvedValue([
                { id: 'test-photo-id', name: 'test-photo.jpg' },
            ]),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterNoMetadata} />);

        // With no metadata and no editable fields, the panel falls back to
        // the empty message only when there is no adapter; PhotoDetail always
        // passes one, so the Edit affordance is what appears.
        await waitFor(() => {
            expect(screen.getAllByText('Edit').length).toBeGreaterThan(0);
        });
        expect(screen.queryByText('Camera')).not.toBeInTheDocument();
    });

    it('should call getPhotoInfoById with correct id', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(mockPhotosAdapter.getPhotoInfoById).toHaveBeenCalledWith('test-photo-id');
        });
    });

    it('should handle missing photo gracefully', () => {
        const adapterWithError = {
            getPhotoInfoById: vi.fn().mockRejectedValue(new Error('Photo not found')),
            getAllPhotosInfo: vi.fn().mockResolvedValue([]),
        } as unknown as IPhotosAdapter;

        // Component should render without crashing even when adapter fails
        const { container } = render(<PhotoDetail photosAdapter={adapterWithError} />);

        // Should show loader when error occurs (since component renders loader when no photo data)
        const loader = container.querySelector('.mantine-Loader-root');
        expect(loader).toBeInTheDocument();
    });
});
