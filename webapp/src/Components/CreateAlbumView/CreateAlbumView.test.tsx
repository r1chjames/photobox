import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '../../test/test-utils';
import { CreateAlbumView } from './CreateAlbumView';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock react-router-dom
vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useParams: () => ({ name: 'test-album' }),
}));

// Mock react-dropzone
vi.mock('react-dropzone', () => ({
    default: ({ children, onDrop }: any) => {
        const mockFile = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
        return (
            <div data-testid="dropzone" onClick={() => onDrop([mockFile])}>
                {children({
                    getRootProps: () => ({ 'data-testid': 'dropzone-root' }),
                    getInputProps: () => ({ 'data-testid': 'dropzone-input' }),
                })}
            </div>
        );
    }
}));

// Mock InfoSnackbar
vi.mock('../Snackbar/InfoSnackbar', () => ({
    InfoSnackbar: ({ text, show }: any) => (
        show ? <div data-testid="snackbar">{text}</div> : null
    )
}));

describe('CreateAlbumView', () => {
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            uploadPhoto: vi.fn().mockResolvedValue(undefined),
        } as unknown as IPhotosAdapter;
    });

    it('should render dropzone area', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.getByTestId('dropzone')).toBeInTheDocument();
    });

    it('should display drag and drop message', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.getByText('Drag photos here to upload')).toBeInTheDocument();
    });

    it('should not show snackbar initially', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        expect(screen.queryByTestId('snackbar')).not.toBeInTheDocument();
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

    it('should render Dropzone component', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const dropzone = screen.getByTestId('dropzone');
        expect(dropzone).toBeInTheDocument();
    });

    it('should render input element for file selection', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const input = screen.getByTestId('dropzone-input');
        expect(input).toBeInTheDocument();
    });

    it('should render root props container', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const root = screen.getByTestId('dropzone-root');
        expect(root).toBeInTheDocument();
    });

    it('should render section element', () => {
        const { container } = render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const section = container.querySelector('section');
        expect(section).toBeInTheDocument();
    });

    it('should render paragraph with instructions', () => {
        render(<CreateAlbumView photosAdapter={mockPhotosAdapter} />);

        const paragraph = screen.getByText(/Drag photos here to upload/i);
        expect(paragraph).toBeInTheDocument();
    });
});