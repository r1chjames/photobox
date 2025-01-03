import type {Meta, StoryObj} from '@storybook/react';
import { InfoSnackbar } from './InfoSnackbar';

const meta: Meta<typeof InfoSnackbar> = {
    component: InfoSnackbar,
};

export default meta;
type Story = StoryObj<typeof InfoSnackbar>;

export const Primary: Story = {
    args: {
        text: 'Snackbar text',
        show: true,
        handleStopShowing: () => console.log()
    },
};