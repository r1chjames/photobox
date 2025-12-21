import type {Meta, StoryObj} from '@storybook/react';

import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {PhotoCard} from "./PhotoCard";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const photos = new SBModelBuilder().newAlbumWithPhotos(3);
const mockPhotosAdapter = new MockPhotosAdapter().withPhotos(photos.getPhotos());

const meta: Meta<typeof PhotoCard> = {
  component: PhotoCard,
};

export default meta;
type Story = StoryObj<typeof PhotoCard>;

export const FirstInAlbum: Story = {
  args: {
    photosAdapter: mockPhotosAdapter,
    source: photos.getPhotos()[0],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: true,
    lastInAlbum: false,
    closeModal: () => console.log("close"),
  },
};

export const MiddleOfAlbum: Story = {
  args: {
    photosAdapter: mockPhotosAdapter,
    source: photos.getPhotos()[1],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: false,
    lastInAlbum: false,
    closeModal: () => console.log("close"),
  },
};

export const LastInAlbum: Story = {
  args: {
    photosAdapter: mockPhotosAdapter,
    source: photos.getPhotos()[2],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: false,
    lastInAlbum: true,
    closeModal: () => console.log("close"),
  },
};