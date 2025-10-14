import type {Meta, StoryObj} from '@storybook/react';
import {AlbumGrid} from './AlbumGrid';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import React from "react";


const meta: Meta<typeof AlbumGrid> = {
    component: AlbumGrid,
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
    .newAlbumWithPhotos(5);

export default meta;
type Story = StoryObj<typeof AlbumGrid>;

export const Primary: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumWithPhotos.getAlbums()),
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos())
    },
};