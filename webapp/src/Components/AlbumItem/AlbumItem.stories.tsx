import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import { AlbumItem } from './AlbumItem';
import { Album } from "../../Models/Album";

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}
const countResp = {"photoCount": 1};

export default {
    component: AlbumItem,
    parameters: {
        mockData: [
            {
                url: 'http://localhost/photos?page=1&limit=1&albumId=a1',
                method: 'GET',
                status: 200,
                response: photoResp
            },
            {
                url: 'http://localhost/photos/count?albumId=a1&albumId=a1',
                method: 'GET',
                status: 200,
                response: countResp
            },
        ],
    },
} as ComponentMeta<typeof AlbumItem>;

export const Primary: ComponentStory<typeof AlbumItem> = (args) => (
        <AlbumItem {...args} />
        );
Primary.args = {
    baseApiUrl: "http://localhost",
    source: new Album("a1", "Album 1", "A great day out","", ""),
    albumViewCallback: () => null,
};
