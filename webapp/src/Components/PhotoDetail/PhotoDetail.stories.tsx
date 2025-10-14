import type {Meta, StoryObj} from '@storybook/react';

import { PhotoDetail } from './PhotoDetail';
import {newPhoto} from "../../utils/SBModelBuilder";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import React from "react";

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

export const Primary: Story = {
    args: {
        photo: newPhoto("1", "1")
    },
};
