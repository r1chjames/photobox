import type {Meta, StoryObj} from '@storybook/react';

import { PhotoDetail } from './PhotoDetail';
import {newPhoto, SBModelBuilder} from "../../utils/SBModelBuilder";

const meta: Meta<typeof PhotoDetail> = {
    component: PhotoDetail,
};

export default meta;
type Story = StoryObj<typeof PhotoDetail>;

export const Primary: Story = {
    args: {
        photo: newPhoto("1", "1")
    },
};
