import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import { AlbumIndexView } from './AlbumIndexView';


const albumResp = [{"id":"a1","name":"album 1","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""},
                   {"id":"a2","name":"album 2","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""}];

export default {
    component: AlbumIndexView,
    parameters: {
        mockData: [
            {
                url: 'http://localhost/albums',
                method: 'GET',
                status: 200,
                response: albumResp
            },
        ],
    },
} as ComponentMeta<typeof AlbumIndexView>;

export const Primary: ComponentStory<typeof AlbumIndexView> = (args) => (
        <AlbumIndexView {...args} />
        );
Primary.args = {
    baseApiUrl: "http://localhost"
};
