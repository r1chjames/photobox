import type {Meta, StoryObj} from '@storybook/react';

import { PhotoDetail } from './PhotoDetail';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import React from "react";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof PhotoDetail> = {
    component: PhotoDetail,
    decorators: [
        (Story) => (
            <QueryClientProvider client={new QueryClient()}>
                <Story />
            </QueryClientProvider>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof PhotoDetail>;

const albumsAndPhotos = new SBModelBuilder().newAlbumWithPhotos(1);

export const Primary: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumsAndPhotos.getPhotos()),
    },
};
