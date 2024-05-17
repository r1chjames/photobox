import type {Meta, StoryObj} from '@storybook/react';
import {AlbumIndexView} from './AlbumIndexView';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";
import {getAllAlbumsInfo} from "../../Adapters/AlbumsAdapter.mock";
import { getPhotosInfoInAlbum, getPhotoCountInAlbum } from "../../Adapters/PhotosAdapter.mock";
import {Album} from "../../Models/Album";
import {Photo} from "../../Models/Photo";

const allAlbums: Album[] = [
    new Album("a1","album 1","an album","", ""),
    new Album("a2","album 2","an album","", "")
];

const photosInAlbums: Photo[] = [
    new Photo("p1", "photo 1", "/tmp/photo1.jpg", "a1", "", {"a":"","b":""}, ""),
    new Photo("p2", "photo 2", "/tmp/photo2.jpg", "a1", "", {"a":"","b":""}, "")
];

const meta: Meta<typeof AlbumIndexView> = {
    component: AlbumIndexView,
};

export default meta;
type Story = StoryObj<typeof AlbumIndexView>;

export const Primary: Story = {
    // @ts-ignore
    async beforeEach() {
        getAllAlbumsInfo.mockReturnValue(allAlbums);
        getPhotosInfoInAlbum.mockReturnValue(photosInAlbums);
        getPhotoCountInAlbum.mockReturnValue(2);
    },
    args: {
        albumsAdapter: new AlbumsAdapter(new MockRestApiAdapter(allAlbums)),
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photosInAlbums)),
    },
};