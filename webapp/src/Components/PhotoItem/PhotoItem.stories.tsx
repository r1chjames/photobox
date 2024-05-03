import type {Meta, StoryObj} from '@storybook/react';

import { PhotoItem } from './PhotoItem';
import { Photo } from '../../Models/Photo';


const meta: Meta<typeof PhotoItem> = {
  component: PhotoItem,
};

export default meta;
type Story = StoryObj<typeof PhotoItem>;

export const Primary: Story = {
  args: {
    baseApiUrl: 'https://google.com/',
    num: 1,
    groupKey: 1,
    source: new Photo('1', 'test', '', 'a', '', {'': ''}, ''),
    allPhotos: [new Photo('1', 'test', '', 'a', '', {'': ''}, '')],
  },
};