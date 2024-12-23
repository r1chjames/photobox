import type {Meta, StoryObj} from '@storybook/react';
import {PhotoIndexView} from "./PhotoIndexView";
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof PhotoIndexView> = {
    component: PhotoIndexView,
};

export default meta;
type Story = StoryObj<typeof PhotoIndexView>;

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
