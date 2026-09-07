import type {Meta, StoryObj} from '@storybook/react-vite';
import {AlbumGrid} from './AlbumGrid';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {Album} from "../../Models/Album";
import {SmartAlbumRules} from "../../Models/SmartAlbumRules";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";

const meta: Meta<typeof AlbumGrid> = {
    component: AlbumGrid,
    decorators: [
        (Story) => (
            <QueryClientProvider client={new QueryClient({defaultOptions: {queries: {retry: false}}})}>
                <Story />
            </QueryClientProvider>
        ),
    ],
};

const albumWithPhotos = new SBModelBuilder()
    .newAlbumWithPhotos(2)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(7)
    .newAlbumWithPhotos(1)
    .newAlbumWithPhotos(9)
    .newAlbumWithPhotos(88)
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(5);

const emptyAlbums = new SBModelBuilder();

export default meta;
type Story = StoryObj<typeof AlbumGrid>;

/**
 * Minimal IAlbumsAdapter whose getAllAlbumsInfo is a promise the story
 * controls (resolve with albums or reject), letting stories exercise the
 * loading skeleton and error states. Other adapter methods are only reached
 * through the create-album modal flow and stay unimplemented for these
 * stories.
 */
class DeferredAlbumsAdapter implements IAlbumsAdapter {
    private _resolve!: (albums: Album[]) => void;
    private _reject!: (reason?: unknown) => void;
    public readonly getAllAlbumsPromise: Promise<Album[]>;

    constructor() {
        this.getAllAlbumsPromise = new Promise<Album[]>((resolve, reject) => {
            this._resolve = resolve;
            this._reject = reject;
        });
    }

    public getAllAlbumsInfo = (): Promise<Album[]> => this.getAllAlbumsPromise;
    public resolve(albums: Album[]): void { this._resolve(albums); }
    public reject(reason?: unknown): void { this._reject(reason); }

    public getCountOfPhotosInAlbum(): Promise<any> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public getAlbumInfoById(_albumId: string): Promise<Album> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public createAlbum(_name: string, _description?: string): Promise<Album> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public createSmartAlbum(_name: string, _rules: SmartAlbumRules): Promise<Album> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public updateSmartAlbum(_albumId: string, _rules: SmartAlbumRules): Promise<Album> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public updateAlbum(_albumId: string, _updates: { name?: string; description?: string; coverPhotoId?: string }): Promise<Album> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public deleteAlbum(_albumId: string, _deletePhotos?: boolean): Promise<void> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
    public addPhotosToAlbum(_albumId: string, _photoIds: string[]): Promise<void> {
        throw new Error('Not implemented in DeferredAlbumsAdapter');
    }
}

export const Primary: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumWithPhotos.getAlbums()),
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos())
    },
};

export const Empty: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(emptyAlbums.getAlbums()),
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(emptyAlbums.getPhotos())
    },
};

export const Loading: Story = {
    render: () => (
        <AlbumGrid
            albumsAdapter={new DeferredAlbumsAdapter()}
            photosAdapter={new MockPhotosAdapter()}
        />
    ),
};

export const Error: Story = {
    render: () => {
        const adapter = new DeferredAlbumsAdapter();
        // Reject on the next microtask so React Query observes the failure.
        setTimeout(() => adapter.reject(new Error('Network error')), 0);
        return (
            <AlbumGrid
                albumsAdapter={adapter}
                photosAdapter={new MockPhotosAdapter()}
            />
        );
    },
};
