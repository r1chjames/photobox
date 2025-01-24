import type {Meta, StoryObj} from '@storybook/react';
import {AppBar, Labels} from "./AppBar";
import {AlbumGrid} from "../AlbumGrid/AlbumGrid";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import React from 'react';
import {PhotoGrid} from "../PhotoGrid/PhotoGrid";
import {Dashboard} from "../Dashboard/Dashboard";
import {PhotoDetail} from "../PhotoDetail/PhotoDetail";
import {SettingsView} from "../SettingsView/SettingsView";
import {MockSettingsAdapter} from "../../Adapters/MockSettingsAdapter";
import {Setting} from "../../Models/Setting";

const meta: Meta<typeof AppBar> = {
    component: AppBar,
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
                maxDisplayed={20}
            />
        </AppBar>
    )
};

export const PhotoInfo: Story = {
    render: () => (
        <AppBar activeLink={Labels.Photos}>
            <PhotoDetail
                photo={albumWithPhotos.getPhotos()[0]}
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
    new Setting("thumbnail_width", "600", "Thumbnail width", "Photo", "Thumbnail width used during thumbnail generation"),
    new Setting("thumbnail_height", "600", "Thumbnail height", "Photo", "Thumbnail height used during thumbnail generation"),
    new Setting("default_new_albums_dir", "/photos", "New album storage location", "System", "Default location on disk to store new albums"),
    new Setting("index_frequency_cron", "0 1 * * *", "CRON expression for indexing", "System", "CRON expression used to initiate indexing")
];

export const Settings: Story = {
    render: () => (
        <AppBar activeLink={Labels.Settings}>
            <SettingsView
                settingsAdapter={new MockSettingsAdapter()
            .withSettings(settings)}
            />
        </AppBar>
    )
};