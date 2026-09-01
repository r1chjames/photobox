import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '../../test/test-utils';
import { CreateAlbumView } from './CreateAlbumView';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock react-router-dom
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useParams: () => ({ name: 'test-album' }),
}));

describe('CreateAlbumView', () => {
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            uploadPhoto: vi.fn().mockResolvedValue(undefined),
        } as unknown as IPhotosAdapter;
    });

    it('should render Add photos button with album name', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.getByText('Add photos to test-album')).toBeInTheDocument();
    });

    it('should render within a container div', () => {
        const { container } = render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const dropzoneContainer = container.querySelector('.createAlbumView__dropzone');
        expect(dropzoneContainer).toBeInTheDocument();
    });

    it('should accept PhotosAdapter as prop', () => {
        const { container } = render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(container).toBeInTheDocument();
    });

    it('should render drag hint text', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.getByText(/or drag photos anywhere on this window to upload/i)).toBeInTheDocument();
    });

    it('should not show error initially', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.queryByText(/failed/i)).not.toBeInTheDocument();
    });
});
