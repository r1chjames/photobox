import type {Meta, StoryObj} from '@storybook/react';
import {AlbumIndexView} from './AlbumIndexView';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";


const meta: Meta<typeof AlbumIndexView> = {
    component: AlbumIndexView,
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
type Story = StoryObj<typeof AlbumIndexView>;

export const Primary: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withAlbums(albumWithPhotos.getAlbums()),
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos())
    },
};