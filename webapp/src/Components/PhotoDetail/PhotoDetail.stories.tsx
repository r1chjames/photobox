import type {Meta, StoryObj} from '@storybook/react';

import { PhotoDetail } from './PhotoDetail';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";

// const photosResp = [{"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}];
const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""};

const meta: Meta<typeof PhotoDetail> = {
    component: PhotoDetail,
};

export default meta;
type Story = StoryObj<typeof PhotoDetail>;

export const Primary: Story = {
    args: {
    baseApiUrl: 'https://google.com/',
    photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photoResp))
    },
};
