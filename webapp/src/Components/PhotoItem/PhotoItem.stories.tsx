import type {Meta, StoryObj} from '@storybook/react';

import { PhotoItem } from './PhotoItem';
import {newPhoto} from "../../utils/SBModelBuilder";


const meta: Meta<typeof PhotoItem> = {
  component: PhotoItem,
};

export default meta;
type Story = StoryObj<typeof PhotoItem>;

export const FirstInAlbum: Story = {
  args: {
    src: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    thumbnail: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    groupKey: 1,
    photoSequence: 0,
    source: newPhoto("1", "1"),
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    lastInAlbum: false
  },
};

export const MiddleOfAlbum: Story = {
  args: {
    src: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    thumbnail: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    groupKey: 1,
    photoSequence: 1,
    source: newPhoto("1", "1"),
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    lastInAlbum: false
  },
};

export const LastInAlbum: Story = {
  args: {
    src: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    thumbnail: 'https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg',
    groupKey: 1,
    photoSequence: 1,
    source: newPhoto("1", "1"),
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    lastInAlbum: true
  },
};