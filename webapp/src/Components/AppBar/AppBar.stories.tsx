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