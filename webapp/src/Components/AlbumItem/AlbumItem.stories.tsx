import type {Meta, StoryObj} from '@storybook/react';
import {AlbumItem} from './AlbumItem';
import { getPhotosInfoInAlbum, getPhotoCountInAlbum } from "../../Adapters/PhotosAdapter.mock";
import {Photo} from "../../Models/Photo";
import {newAlbum, newPhoto} from "../../utils/Storybook";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof AlbumItem> = {
    component: AlbumItem,
};

export default meta;
type Story = StoryObj<typeof AlbumItem>;

const photosInAlbums: Photo[] = [
    newPhoto(1), newPhoto(2)
];

export const Primary: Story = {
    // @ts-ignore
    async beforeEach() {
        getPhotosInfoInAlbum.mockReturnValue(photosInAlbums);
        getPhotoCountInAlbum.mockReturnValue(1);
    },
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withGetPhotosInfoInAlbum(photosInAlbums),
        source: newAlbum(1),
        albumViewCallback: () => console.log("Clicked"),
    },
};