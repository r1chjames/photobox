import type {Meta, StoryObj} from '@storybook/react';
import {AppBar} from "./AppBar";
import {AlbumIndexView} from "../AlbumIndexView/AlbumIndexView";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import React from 'react';

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
    .newAlbumWithPhotos(5);

export default meta;
type Story = StoryObj<typeof AppBar>;

export const Primary: Story = {
};

export const Dashboard: Story = {
};

export const Photos: Story = {
};

export const Albums: Story = {
    render: (args) => (
        <AppBar>
            <AlbumIndexView
                albumsAdapter={new MockAlbumsAdapter()
                    .withAlbums(albumWithPhotos.getAlbums())}
                photosAdapter={new MockPhotosAdapter()
                                .withPhotos(albumWithPhotos.getPhotos())}
            />
        </AppBar>
    )
};