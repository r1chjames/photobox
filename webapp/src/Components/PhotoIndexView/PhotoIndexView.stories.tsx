import type {Meta, StoryObj} from '@storybook/react';
import {PhotoIndexView} from "./PhotoIndexView";
import {Photo} from "../../Models/Photo";
import {newPhoto} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";


const meta: Meta<typeof PhotoIndexView> = {
    component: PhotoIndexView,
};

export default meta;
type Story = StoryObj<typeof PhotoIndexView>;

const photosInAlbums: Photo[] = [
    newPhoto("1", "1"), newPhoto("2", "1")
];

export const Primary: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withGetPhotosInfoInAlbum(photosInAlbums)
            .withGetPhotoCountInAlbum(photosInAlbums.length)
    },
};