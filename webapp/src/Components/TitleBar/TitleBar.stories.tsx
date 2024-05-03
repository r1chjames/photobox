import type {Meta, StoryObj} from '@storybook/react';
import { TitleBar } from './TitleBar';

const meta: Meta<typeof TitleBar> = {
    component: TitleBar,
};

export default meta;
type Story = StoryObj<typeof TitleBar>;

export const Primary: Story = {
    args: {},
};