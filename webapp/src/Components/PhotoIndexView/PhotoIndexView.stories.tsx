import type {Meta, StoryObj} from '@storybook/react';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import {PhotoIndexView} from "./PhotoIndexView";

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}

const meta: Meta<typeof PhotoIndexView> = {
    component: PhotoIndexView,
};

export default meta;
type Story = StoryObj<typeof PhotoIndexView>;

export const Primary: Story = {
    args: {
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photoResp)),
    },
};