import type {Meta, StoryObj} from '@storybook/react';
import {AlbumIndexView} from './AlbumIndexView';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";


const meta: Meta<typeof AlbumIndexView> = {
    component: AlbumIndexView,
};

const albumWithPhotos = new SBModelBuilder().newAlbumWithPhotos(5, 4);

export default meta;
type Story = StoryObj<typeof AlbumIndexView>;

export const Primary: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withGetAllAlbumsResponse(albumWithPhotos.getAlbums())
            .withGetCountOfPhotosInAlbum(albumWithPhotos.getPhotos().length),
        photosAdapter: new MockPhotosAdapter()
            .withGetAllPhotosInfo(albumWithPhotos.getPhotos()),
    },
};