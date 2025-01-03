import type {Meta, StoryObj} from '@storybook/react';
import {Dashboard} from "./Dashboard";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";

const meta: Meta<typeof Dashboard> = {
    component: Dashboard,
};

const albumWithPhotos = new SBModelBuilder()
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4)
    .newAlbumWithPhotos(4);

export default meta;
type Story = StoryObj<typeof Dashboard>;

export const Primary: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos()),
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumWithPhotos.getAlbums())
    },
};