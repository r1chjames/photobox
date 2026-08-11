import type {Meta, StoryObj} from '@storybook/react-vite';
import {NetworkStatusBannerAlert} from './NetworkStatusBanner';

const meta: Meta<typeof NetworkStatusBannerAlert> = {
    component: NetworkStatusBannerAlert,
};

export default meta;
type Story = StoryObj<typeof NetworkStatusBannerAlert>;

export const BackendUnreachable: Story = {
    args: {
        status: 'backend-unreachable',
        onDismiss: () => console.log('dismiss'),
    },
};

export const Offline: Story = {
    args: {
        status: 'offline',
        onDismiss: () => console.log('dismiss'),
    },
};
