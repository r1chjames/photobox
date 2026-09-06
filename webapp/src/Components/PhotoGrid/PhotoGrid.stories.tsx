import type {Meta, StoryObj} from '@storybook/react-vite';
import {PhotoGrid} from "./PhotoGrid";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";

const meta: Meta<typeof PhotoGrid> = {
    component: PhotoGrid,
    decorators: [
        (Story) => (
            <QueryClientProvider client={new QueryClient()}>
                <Story />
            </QueryClientProvider>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof PhotoGrid>;

const albumsAndPhotos = new SBModelBuilder().newAlbumWithPhotos(150);


export const AllPhotos: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumsAndPhotos.getPhotos()),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumsAndPhotos.getAlbums())
    },
    parameters: {
        reactRouter: {
            initialEntries: ['/'],
            routePath: '/',
        }
    }
};

export const AlbumPhotos: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumsAndPhotos.getPhotos()),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumsAndPhotos.getAlbums())
    },
    parameters: {
        reactRouter: {
            initialEntries: ['/album/Album 1'],
            routePath: '/album/:id',
        }
    }
};

// Mobile viewport variant (issue #97) — verifies the 2-column grid + bottom nav
// layout at phone width.
export const AllPhotosMobile: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumsAndPhotos.getPhotos()),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumsAndPhotos.getAlbums())
    },
    parameters: {
        viewport: {
            defaultViewport: 'mobile2',
        },
        reactRouter: {
            initialEntries: ['/'],
            routePath: '/',
        }
    }
};

// Adversarial composition (issue #174): extreme aspect ratios — panoramas,
// very tall portraits, squares and landscapes in a fixed cycle — stress the
// planner's featured/tall budgets and consecutive-variant limits. The grid
// must stay aligned to the 16px rhythm with no runaway repeats.
const adversarialDims: Array<[number, number]> = [
    [6000, 2000],   // panorama -> featured (2-col)
    [4032, 3024],   // landscape -> standard
    [1500, 4000],   // very tall portrait -> tall
    [3024, 4032],   // portrait
    [3000, 3000],   // square
    [4032, 3024],   // landscape -> standard
];
const adversarialBuilder = new SBModelBuilder().newAlbumWithPhotos(90);
const adversarialPhotos = adversarialBuilder
    .getPhotos()
    .map((photo, index) => {
        const [width, height] = adversarialDims[index % adversarialDims.length];
        return {...photo, width, height};
    });

export const AdversarialMixedDimensions: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(adversarialPhotos),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(adversarialBuilder.getAlbums()),
    },
    parameters: {
        reactRouter: {
            initialEntries: ['/'],
            routePath: '/',
        }
    }
};

const emptyAlbum = new SBModelBuilder().newEmptyAlbum();
export const EmptyAlbumPhotos: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter(),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(emptyAlbum.getAlbums())
    },
    parameters: {
        reactRouter: {
            initialEntries: ['/album/Album 1'],
            routePath: '/album/:id',
        }
    }
};

// Favorites view (issue #126) — favourites are fetched server-side via
// `favorites=true`; only favourited photos should appear.
const favouritesPhotos = new SBModelBuilder().newAlbumWithPhotos(12).getPhotos()
    .map((photo, index) => ({ ...photo, favorite: index % 2 === 0 }));
export const Favorites: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(favouritesPhotos),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(new SBModelBuilder().newEmptyAlbum().getAlbums()),
        favoritesOnly: true,
    },
    parameters: {
        reactRouter: {
            initialEntries: ['/favorites'],
            routePath: '/favorites',
        }
    }
};
