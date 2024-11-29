import type {Meta, StoryObj} from '@storybook/react';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import {PhotoIndexView} from "./PhotoIndexView";
import {getPhotosInfoInAlbum} from "../../Adapters/PhotosAdapter.mock";
import {Photo} from "../../Models/Photo";
import {newPhoto} from "../../utils/Storybook";

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}

const meta: Meta<typeof PhotoIndexView> = {
    component: PhotoIndexView,
};

export default meta;
type Story = StoryObj<typeof PhotoIndexView>;

const photosInAlbums: Photo[] = [
    newPhoto(1), newPhoto(2)
];

export const Primary: Story = {
    // @ts-ignore
    async beforeEach() {
        getPhotosInfoInAlbum.mockReturnValue(photosInAlbums);
    },
    args: {
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photoResp)),
    },
};