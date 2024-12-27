import type {Meta, StoryObj} from '@storybook/react';
import {PhotoGrid} from "./PhotoGrid";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof PhotoGrid> = {
    component: PhotoGrid,
};

export default meta;
type Story = StoryObj<typeof PhotoGrid>;

const photos = new SBModelBuilder().newPhotoCollection(150, "Album 1");

export const AllPhotos: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(photos.getPhotos())
    },
};

export const AlbumPhotos: Story = {
    args: {
        albumId: "Album 1",
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(photos.getPhotos())
    },
};

export const EmptyAlbumPhotos: Story = {
    args: {
        albumId: "Album 1",
        photosAdapter: new MockPhotosAdapter()
    },
};
