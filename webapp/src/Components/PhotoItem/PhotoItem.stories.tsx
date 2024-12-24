import type {Meta, StoryObj} from '@storybook/react';

import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {PhotoItem} from "./PhotoItem";

const photos = new SBModelBuilder().newAlbumWithPhotos(3);

const meta: Meta<typeof PhotoItem> = {
  component: PhotoItem,
};

export default meta;
type Story = StoryObj<typeof PhotoItem>;

export const FirstInAlbum: Story = {
  args: {
    groupKey: 1,
    photoSequence: 0,
    source: photos.getPhotos()[0],
    previousPhoto: () => photos.getPhotos()[0],
    nextPhoto: () => photos.getPhotos()[1],
    firstInAlbum: true,
    lastInAlbum: false
  },
};

export const MiddleOfAlbum: Story = {
  args: {
    groupKey: 1,
    photoSequence: 1,
    source: photos.getPhotos()[1],
    previousPhoto: () => photos.getPhotos()[0],
    nextPhoto: () => photos.getPhotos()[2],
    firstInAlbum: false,
    lastInAlbum: false
  },
};

export const LastInAlbum: Story = {
  args: {
    groupKey: 1,
    photoSequence: 1,
    source: photos.getPhotos()[2],
    previousPhoto: () => photos.getPhotos()[1],
    nextPhoto: () => photos.getPhotos()[2],
    firstInAlbum: false,
    lastInAlbum: true
  },
};