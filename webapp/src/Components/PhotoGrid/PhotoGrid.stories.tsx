import type {Meta, StoryObj} from '@storybook/react';
import {PhotoGrid} from "./PhotoGrid";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import React from "react";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";

const queryClient = new QueryClient();

const meta: Meta<typeof PhotoGrid> = {
    component: PhotoGrid,
    decorators: [
        (Story) => (
            <QueryClientProvider client={queryClient}>
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
