import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '../../test/test-utils';
import { Dashboard } from './Dashboard';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// Mock AlbumGrid component
vi.mock('../AlbumGrid/AlbumGrid', () => ({
    AlbumGrid: ({ initialDisplayCount }: { initialDisplayCount?: number }) => (
        <div data-testid="album-grid" data-max-displayed={initialDisplayCount}>
            Mock Album Grid
        </div>
    )
}));

// Mock PhotoGrid component
vi.mock('../PhotoGrid/PhotoGrid', () => ({
    PhotoGrid: ({ maxDisplayed }: { maxDisplayed?: number }) => (
        <div data-testid="photo-grid" data-max-displayed={maxDisplayed}>
            Mock Photo Grid
        </div>
    )
}));

describe('Dashboard', () => {
    let mockAlbumsAdapter: IAlbumsAdapter;
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockAlbumsAdapter = {} as IAlbumsAdapter;
        mockPhotosAdapter = {} as IPhotosAdapter;
    });

    it('should render Albums section with title', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.getByText('Albums')).toBeInTheDocument();
    });

    it('should render Photos section with title', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.getByText('Photos')).toBeInTheDocument();
    });

    it('should render AlbumGrid component', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.getByTestId('album-grid')).toBeInTheDocument();
    });

    it('should render PhotoGrid component', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.getByTestId('photo-grid')).toBeInTheDocument();
    });

    it('should pass initialDisplayCount=30 to AlbumGrid', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        const albumGrid = screen.getByTestId('album-grid');
        expect(albumGrid).toHaveAttribute('data-max-displayed', '30');
    });

    it('should pass maxDisplayed=50 to PhotoGrid', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        const photoGrid = screen.getByTestId('photo-grid');
        expect(photoGrid).toHaveAttribute('data-max-displayed', '50');
    });

    it('should render both grids on the same page', () => {
        render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        expect(screen.getByTestId('album-grid')).toBeInTheDocument();
        expect(screen.getByTestId('photo-grid')).toBeInTheDocument();
    });

    it('should render section titles as h4', () => {
        const { container } = render(
            <Dashboard
                albumsAdapter={mockAlbumsAdapter}
                photosAdapter={mockPhotosAdapter}
            />
        );

        // Check for h4 elements with specific text
        const titles = container.querySelectorAll('[data-size="h4"]');
        expect(titles.length).toBeGreaterThanOrEqual(2);
    });
});
