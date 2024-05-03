import type {Meta, StoryObj} from '@storybook/react';
import {AlbumIndexView} from './AlbumIndexView';
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {MockRestApiAdapter} from "../../Adapters/RestApiAdapter";
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";

const albumResp = [{"id":"a1","name":"album 1","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""},
                   {"id":"a2","name":"album 2","description":"an album","tags":"","metadata":{"a":"","b":""},"createdAt":"","updatedAt":""}];

const photoResp = {"id":"p1","name":"photo 1","filesystemPath":"/tmp/photo1.jpg","albumId":"a1","tags":"","metadata":{"a":"","b":""},"createdAt":""}

const meta: Meta<typeof AlbumIndexView> = {
    component: AlbumIndexView,
};

export default meta;
type Story = StoryObj<typeof AlbumIndexView>;

export const Primary: Story = {
    args: {
        albumsAdapter: new AlbumsAdapter(new MockRestApiAdapter(albumResp)),
        photosAdapter: new PhotosAdapter(new MockRestApiAdapter(photoResp)),
    },
};