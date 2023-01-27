import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";

import { PhotoDetail } from './PhotoDetail';

const photosResp = [{"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}];
const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""};

export default {
    component: PhotoDetail,
    parameters: {
        mockData: [
            {
                url: 'http://localhost/photo/1',
                method: 'GET',
                status: 200,
                response: photoResp
            },
            {
                url: 'http://localhost/photos?page=1&limit=1',
                method: 'GET',
                status: 200,
                response: photosResp
            }
        ],
    },
} as ComponentMeta<typeof PhotoDetail>;

export const primary: ComponentStory<typeof PhotoDetail> = (args) => (
        <PhotoDetail {...args} />
        );

primary.args = {
  baseApiUrl: 'https://google.com/'
};
