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
}));

describe('PhotoDetail', () => {
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    Camera: 'Canon EOS R5',
                    ISO: '400',
                    FocalLength: '50mm',
                    ExifVersion: '0232',
                },
            }),
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
        } as unknown as IPhotosAdapter;

        const { container } = render(<PhotoDetail photosAdapter={slowAdapter} />);

        // Mantine Loader should be present while loading (it's a span with mantine-Loader-root class)
        const loader = container.querySelector('.mantine-Loader-root');
        expect(loader).toBeInTheDocument();

        // Wait for component to finish loading
        await waitFor(() => {
            expect(screen.getByText('Metadata')).toBeInTheDocument();
        });
    });

    it('should render metadata button', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(screen.getByText('Metadata')).toBeInTheDocument();
        });
    });

    it('should display metadata in dialog', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            // Metadata should be visible since dialog is open by default
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('Canon EOS R5')).toBeInTheDocument();
        });
    });

    it('should filter numeric keys from metadata', async () => {
        const adapterWithNumericKeys = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    Camera: 'Canon EOS R5',
                    '0': 'should be filtered',
                    '1': 'should be filtered',
                    '99': 'should be filtered',
                    ISO: 'ISO 400',  // Changed from '400' to 'ISO 400' to avoid JSON parsing
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithNumericKeys} />);

        await waitFor(() => {
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('ISO')).toBeInTheDocument();
        }, { timeout: 3000 });

        expect(screen.queryByText('should be filtered')).not.toBeInTheDocument();
    });

    it('should handle nested metadata objects', async () => {
        const adapterWithNestedMetadata = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    Camera: {
                        Make: 'Canon',
                        Model: 'EOS R5',
                    },
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithNestedMetadata} />);

        await waitFor(() => {
            expect(screen.getByText('Make')).toBeInTheDocument();
            expect(screen.getByText('Canon')).toBeInTheDocument();
            expect(screen.getByText('Model')).toBeInTheDocument();
            expect(screen.getByText('EOS R5')).toBeInTheDocument();
        });
    });

    it('should handle metadata with empty values', async () => {
        const adapterWithEmptyValues = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    Camera: 'Canon EOS R5',
                    EmptyField: '',
                    UndefinedField: undefined,
                    ISO: 'ISO 400',  // Changed from '400' to 'ISO 400' to avoid JSON parsing
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithEmptyValues} />);

        await waitFor(() => {
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('ISO')).toBeInTheDocument();
        }, { timeout: 3000 });

        // Empty and undefined fields should be filtered out
        expect(screen.queryByText('EmptyField')).not.toBeInTheDocument();
        expect(screen.queryByText('UndefinedField')).not.toBeInTheDocument();
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
        } as unknown as IPhotosAdapter;

        // Component should render without crashing even when adapter fails
        const { container } = render(<PhotoDetail photosAdapter={adapterWithError} />);

        // Should show loader when error occurs (since component renders loader when no photo data)
        const loader = container.querySelector('.mantine-Loader-root');
        expect(loader).toBeInTheDocument();
    });

    it('should render metadata table with headers', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(screen.getByText('Parameter')).toBeInTheDocument();
            expect(screen.getByText('Value')).toBeInTheDocument();
        });
    });

    it('should handle JSON string metadata', async () => {
        const adapterWithJSONMetadata = {
            getPhotoInfoById: vi.fn().mockResolvedValue({
                id: 'test-photo-id',
                name: 'test-photo.jpg',
                metadata: {
                    EXIF: '{"Make":"Canon","Model":"EOS R5"}',
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithJSONMetadata} />);

        await waitFor(() => {
            expect(screen.getByText('Make')).toBeInTheDocument();
            expect(screen.getByText('Canon')).toBeInTheDocument();
        });
    });
});