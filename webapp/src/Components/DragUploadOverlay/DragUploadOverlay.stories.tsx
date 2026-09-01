import type {Meta, StoryObj} from '@storybook/react-vite';
import {DragUploadOverlay} from './DragUploadOverlay';
import {MockPhotosAdapter} from '../../Adapters/MockPhotosAdapter';

const meta: Meta<typeof DragUploadOverlay> = {
    component: DragUploadOverlay,
    title: 'DragUploadOverlay',
};

export default meta;
type Story = StoryObj<typeof DragUploadOverlay>;

export const Default: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter(),
        albumName: 'General',
    },
};
