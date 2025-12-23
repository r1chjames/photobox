import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { PhotoDetail } from './PhotoDetail';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock react-router-dom
vi.mock('react-router-dom', () => ({
    useParams: () => ({ id: 'test-photo-id' }),
}));

// Mock ImageUtils
vi.mock('../../utils/ImageUtils', () => ({
    fetchPhotoBinWithAuth: vi.fn().mockResolvedValue('data:image/png;base64,mockimage'),
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

    it('should show loader initially', () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);
        // Loader should be present initially
        const loader = screen.queryByRole('progressbar') || screen.queryByRole('status');
        expect(loader).toBeTruthy();
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
                    ISO: '400',
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithNumericKeys} />);

        await waitFor(() => {
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('ISO')).toBeInTheDocument();
            expect(screen.queryByText('should be filtered')).not.toBeInTheDocument();
        });
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
                    ISO: '400',
                },
            }),
        } as unknown as IPhotosAdapter;

        render(<PhotoDetail photosAdapter={adapterWithEmptyValues} />);

        await waitFor(() => {
            expect(screen.getByText('Camera')).toBeInTheDocument();
            expect(screen.getByText('ISO')).toBeInTheDocument();
            // Empty and undefined fields should be filtered out
            expect(screen.queryByText('EmptyField')).not.toBeInTheDocument();
            expect(screen.queryByText('UndefinedField')).not.toBeInTheDocument();
        });
    });

    it('should call getPhotoInfoById with correct id', async () => {
        render(<PhotoDetail photosAdapter={mockPhotosAdapter} />);

        await waitFor(() => {
            expect(mockPhotosAdapter.getPhotoInfoById).toHaveBeenCalledWith('test-photo-id');
        });
    });

    it('should handle missing photo gracefully', async () => {
        const adapterWithError = {
            getPhotoInfoById: vi.fn().mockRejectedValue(new Error('Photo not found')),
        } as unknown as IPhotosAdapter;

        const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

        render(<PhotoDetail photosAdapter={adapterWithError} />);

        // Should show loader when error occurs
        await waitFor(() => {
            expect(consoleErrorSpy).toHaveBeenCalled();
        });

        consoleErrorSpy.mockRestore();
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