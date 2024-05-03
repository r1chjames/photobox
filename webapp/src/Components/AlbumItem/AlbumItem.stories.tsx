import type {Meta, StoryObj} from '@storybook/react';
import {AlbumItem} from './AlbumItem';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {Album} from "../../Models/Album";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}
const album = new Album("1234", "Album 1", "A great album", "","")
// const countResp = {"photoCount": 1};

const meta: Meta<typeof AlbumItem> = {
    component: AlbumItem,
};

export default meta;
type Story = StoryObj<typeof AlbumItem>;

export const Primary: Story = {
    args: {
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photoResp)),
        source: album,
        albumViewCallback: () => console.log("Clicked"),
    },
};