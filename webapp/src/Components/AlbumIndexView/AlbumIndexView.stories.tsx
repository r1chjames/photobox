import * as React from "react";
// import { Meta, StoryObj } from "@storybook/react";
import { AlbumIndexView } from './AlbumIndexView';


const albumResp = [{"id":"a1","name":"album 1","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""},
                   {"id":"a2","name":"album 2","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""}];

// const meta: Meta<typeof AlbumIndexView> = {
//     component: AlbumIndexView,
// };
//
// export default meta;
// type Story = StoryObj<typeof AlbumIndexView>;
//
//
// export const Primary: Story = {
//     args: {
//         baseApiUrl: "http://localhost"
//     }
// }


export default {
    title: 'Examples/Fetch',
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
};

const Template = (args: any) => <AlbumIndexView {...args} />;

export const FetchCall = Template.bind({});


//     {
//     component: AlbumIndexView,
//     parameters: {
//         mockData: [
//             {
//                 url: 'http://localhost/albums',
//                 method: 'GET',
//                 status: 200,
//                 response: albumResp
//             },
//         ],
//     },
// };
