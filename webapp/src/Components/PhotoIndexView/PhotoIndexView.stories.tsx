import type {Meta, StoryObj} from '@storybook/react';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import {PhotoIndexView} from "./PhotoIndexView";
import {getPhotosInfoInAlbum} from "../../Adapters/PhotosAdapter.mock";
import {Photo} from "../../Models/Photo";

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}

const meta: Meta<typeof PhotoIndexView> = {
    component: PhotoIndexView,
};

export default meta;
type Story = StoryObj<typeof PhotoIndexView>;

const photosInAlbums: Photo[] = [
    new Photo("p1", "photo 1", "/tmp/photo1.jpg", "a1", "", {"a":"","b":""}, ""),
    new Photo("p2", "photo 2", "/tmp/photo2.jpg", "a1", "", {"a":"","b":""}, "")
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