import type {Meta, StoryObj} from '@storybook/react';
import {AlbumItem} from './AlbumItem';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {Album} from "../../Models/Album";
import {anything, instance, mock, when} from "ts-mockito";

const meta: Meta<typeof AlbumItem> = {
    component: AlbumItem,
};

export default meta;
type Story = StoryObj<typeof AlbumItem>;

const album = new Album("1234", "Album 1", "A great album", "","")

const mockPhotosAdapter: PhotosAdapter = mock(PhotosAdapter);
when(mockPhotosAdapter.getPhotosInfoInAlbum(anything(), anything(), anything())).thenResolve("https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg");
when(mockPhotosAdapter.getPhotoCountInAlbum(anything())).thenResolve(1);

export const Primary: Story = {
    args: {
        photosAdapter: instance(mockPhotosAdapter),
        source: album,
        albumViewCallback: () => console.log("Clicked"),
    },
};