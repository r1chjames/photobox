import type {Meta, StoryObj} from '@storybook/react-vite';
import {AlbumCard} from './AlbumCard';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {Photo} from "../../Models/Photo";

const meta: Meta<typeof AlbumCard> = {
    component: AlbumCard,
    decorators: [
        (Story) => (
            <QueryClientProvider client={new QueryClient({ defaultOptions: { queries: { retry: false } } })}>
                <Story />
            </QueryClientProvider>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof AlbumCard>;

// Valid blurhash + dominant colour so placeholder-first covers decode to a
// canvas frame (or colour block) before the fetched thumbnails crossfade in.
const COVER_BLURHASH = 'LEHV6nWB2yk8pyo0adR*.7kCMdnj';
const COVER_COLOR = '#5a6b7c';

const withPlaceholderMeta = (photo: Photo): Photo => ({
    ...photo,
    blurhash: COVER_BLURHASH,
    dominantColor: COVER_COLOR,
});

// Adapter whose metadata/count queries never resolve — pins the card in its
// loading shell (skeleton media + skeleton count pill) for the story.
class PendingPhotosAdapter extends MockPhotosAdapter {
    public getPhotosInfoInAlbum = async (): Promise<Photo[]> => new Promise(() => undefined);
    public getPhotoCountInAlbum = async (): Promise<number> => new Promise(() => undefined);
}

// Adapter whose metadata/count queries reject — the card must keep its full
// shell and fall back to the static empty/placeholder treatment.
class ErroringPhotosAdapter extends MockPhotosAdapter {
    public getPhotosInfoInAlbum = async (): Promise<Photo[]> => {
        throw new Error('Failed to load album covers');
    };
    public getPhotoCountInAlbum = async (): Promise<number> => {
        throw new Error('Failed to load photo count');
    };
}

const populatedAlbum = new SBModelBuilder().newAlbumWithPhotos(3);
const clickCallback = () => console.log('Album clicked');

// Valid blurhash + dominant colour under every cover: placeholders paint
// immediately, then the fetched thumbnails crossfade in above them.
export const Populated: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter().withPhotos(populatedAlbum.getPhotos().map(withPlaceholderMeta)),
        albumsAdapter: new MockAlbumsAdapter(),
        source: populatedAlbum.getAlbums().at(0),
        albumViewCallback: clickCallback,
    },
};

// Metadata/count pending: full-size card shell with skeleton media slots and
// a skeleton count pill — never a bare spinner or a collapsed card.
export const Loading: Story = {
    args: {
        photosAdapter: new PendingPhotosAdapter(),
        albumsAdapter: new MockAlbumsAdapter(),
        source: populatedAlbum.getAlbums().at(0),
        albumViewCallback: clickCallback,
    },
};

// Photos carry only a dominant colour (no blurhash): the colour block acts
// as the static placeholder beneath the loaded thumbnail.
export const NoBlurhash: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(populatedAlbum.getPhotos().map((photo) => ({ ...photo, dominantColor: COVER_COLOR }))),
        albumsAdapter: new MockAlbumsAdapter(),
        source: populatedAlbum.getAlbums().at(0),
        albumViewCallback: clickCallback,
    },
};

// No photos in the album at all: distinct IconPhotoOff empty state with the
// real '0 photos' badge — visually different from the loading shell.
export const EmptyAlbum: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter(),
        albumsAdapter: new MockAlbumsAdapter(),
        source: new SBModelBuilder().newEmptyAlbum().getAlbums().at(0),
        albumViewCallback: clickCallback,
    },
};

// Queries reject: the shell, title and static placeholder treatment stay up
// with no crash (albums have no thumbnail retry affordance).
export const ErrorState: Story = {
    args: {
        photosAdapter: new ErroringPhotosAdapter(),
        albumsAdapter: new MockAlbumsAdapter(),
        source: populatedAlbum.getAlbums().at(0),
        albumViewCallback: clickCallback,
    },
};
