import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import { AlbumItem } from './AlbumItem';
import { Album } from "../../Models/Album";
import {Photo, PhotoCount} from "../../Models/Photo";

const photoResp = JSON.stringify(new Photo("p1", "photo 1", "/tmp/photo1.jpg", "a1", "", {a: "", b: ""}, ""))
const countResp = JSON.stringify(new PhotoCount(1))

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
