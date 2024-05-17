import type {Meta, StoryObj} from '@storybook/react';
import {AlbumItem} from './AlbumItem';
import {Album} from "../../Models/Album";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import { getPhotosInfoInAlbum, getPhotoCountInAlbum } from "../../Adapters/PhotosAdapter.mock";
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {Photo} from "../../Models/Photo";

const meta: Meta<typeof AlbumItem> = {
    component: AlbumItem,
};

export default meta;
type Story = StoryObj<typeof AlbumItem>;

const album = new Album("a1", "Album 12222", "A great album", "","")

const photosInAlbums: Photo[] = [
    new Photo("p1", "photo 1", "/tmp/photo1.jpg", "a1", "", {"a":"","b":""}, ""),
    new Photo("p2", "photo 2", "/tmp/photo2.jpg", "a1", "", {"a":"","b":""}, "")
];

export const Primary: Story = {
    // @ts-ignore
    async beforeEach() {
        getPhotosInfoInAlbum.mockReturnValue(photosInAlbums);
        getPhotoCountInAlbum.mockReturnValue(1);
    },
    args: {
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter("")),
        source: album,
        albumViewCallback: () => console.log("Clicked"),
    },
};