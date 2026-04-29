import type {Meta, StoryObj} from '@storybook/react-vite';
import {AppBar, Labels} from "./AppBar";
import {AlbumGrid} from "../AlbumGrid/AlbumGrid";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {PhotoGrid} from "../PhotoGrid/PhotoGrid";
import {Dashboard} from "../Dashboard/Dashboard";
import {PhotoDetail} from "../PhotoDetail/PhotoDetail";
import {SettingsView} from "../SettingsView/SettingsView";
import {MockSettingsAdapter} from "../../Adapters/MockSettingsAdapter";

import {QueryClient, QueryClientProvider} from "@tanstack/react-query";

const meta: Meta<typeof AppBar> = {
    component: AppBar,
    decorators: [
        (Story) => (
            <QueryClientProvider client={new QueryClient()}>
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
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(14)
    .newAlbumWithPhotos(5);

export default meta;
type Story = StoryObj<typeof AppBar>;

export const Home: Story = {
    render: () => (
        <AppBar activeLink={Labels.Dashboard}>
            <Dashboard
                albumsAdapter={new MockAlbumsAdapter()
                    .withAlbums(albumWithPhotos.getAlbums())}
                photosAdapter={new MockPhotosAdapter()
                    .withPhotos(albumWithPhotos.getPhotos())}
            />
        </AppBar>
    )
};

export const Photos: Story = {
    render: () => (
        <AppBar activeLink={Labels.Photos}>
            <PhotoGrid
                photosAdapter={new MockPhotosAdapter()
                    .withPhotos(albumWithPhotos.getPhotos())}
                albumsAdapter={new MockAlbumsAdapter()
                    .withAlbums(albumWithPhotos.getAlbums())}
                maxDisplayed={20}
            />
        </AppBar>
    )
};

export const PhotoInfo: Story = {
    render: () => (
        <AppBar activeLink={Labels.Photos}>
            <PhotoDetail
                photosAdapter={new MockPhotosAdapter()
                    .withPhotos(albumWithPhotos.getPhotos())}
            />
        </AppBar>
    )
};

export const Albums: Story = {
    render: () => (
        <AppBar activeLink={Labels.Albums}>
            <AlbumGrid
                albumsAdapter={new MockAlbumsAdapter()
                    .withAlbums(albumWithPhotos.getAlbums())}
                photosAdapter={new MockPhotosAdapter()
                                .withPhotos(albumWithPhotos.getPhotos())}
                maxDisplayed={20}
            />
        </AppBar>
    )
};

const settings = [
    { key: "thumbnail_width", value: "600", friendlyName: "Thumbnail width", category: "Photo", description: "Thumbnail width used during thumbnail generation" },
    { key: "thumbnail_height", value: "600", friendlyName: "Thumbnail height", category: "Photo", description: "Thumbnail height used during thumbnail generation" },
    { key: "default_new_albums_dir", value: "/photos", friendlyName: "New album storage location", category: "System", description: "Default location on disk to store new albums" },
    { key: "index_frequency_cron", value: "0 1 * * *", friendlyName: "CRON expression for indexing", category: "System", description: "CRON expression used to initiate indexing" }
];

export const Settings: Story = {
    render: () => (
        <AppBar activeLink={Labels.Settings}>
            <SettingsView
                settingsAdapter={new MockSettingsAdapter()
                    .withSettings(settings)}
                photosAdapter={new MockPhotosAdapter()
                    .withPhotos(albumWithPhotos.getPhotos())}
            />
        </AppBar>
    )
};