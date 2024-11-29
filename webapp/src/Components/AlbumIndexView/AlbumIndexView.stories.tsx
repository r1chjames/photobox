import type {Meta, StoryObj} from '@storybook/react';
import {AlbumIndexView} from './AlbumIndexView';
import {Album} from "../../Models/Album";
import {Photo} from "../../Models/Photo";
import {newAlbum, newPhoto} from "../../utils/Storybook";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const allAlbums: Album[] = [
    newAlbum(1), newAlbum(2)
];

const photosInAlbums: Photo[] = [
    newPhoto(1), newPhoto(2)
];

const meta: Meta<typeof AlbumIndexView> = {
    component: AlbumIndexView,
};

export default meta;
type Story = StoryObj<typeof AlbumIndexView>;

export const Primary: Story = {
    args: {
        albumsAdapter: new MockAlbumsAdapter()
            .withGetAllAlbumsResponse(allAlbums)
            .withGetCountOfPhotosInAlbum(2)
            .withGetAlbumInfoById(""),
        photosAdapter: new MockPhotosAdapter()
            .withGetAllPhotosInfo(photosInAlbums),
    },
};