import type {Meta, StoryObj} from '@storybook/react-vite';
import {Slideshow} from './Slideshow';
import {SBModelBuilder} from '../../utils/SBModelBuilder';
import {MockPhotosAdapter} from '../../Adapters/MockPhotosAdapter';

const builder = new SBModelBuilder().newAlbumWithPhotos(5);
const photos = builder.getPhotos();
const mockPhotosAdapter = new MockPhotosAdapter().withPhotos(photos);

const meta: Meta<typeof Slideshow> = {
    component: Slideshow,
    parameters: {
        layout: 'fullscreen',
    },
};

export default meta;
type Story = StoryObj<typeof Slideshow>;

export const Default: Story = {
    args: {
        photos,
        startIndex: 0,
        photosAdapter: mockPhotosAdapter,
        onClose: () => console.log('close'),
    },
};
