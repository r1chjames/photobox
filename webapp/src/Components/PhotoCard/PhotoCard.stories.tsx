import type {Meta, StoryObj} from '@storybook/react';

import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {PhotoCard} from "./PhotoCard";

const photos = new SBModelBuilder().newAlbumWithPhotos(3);

const meta: Meta<typeof PhotoCard> = {
  component: PhotoCard,
};

export default meta;
type Story = StoryObj<typeof PhotoCard>;

export const FirstInAlbum: Story = {
  args: {
    source: photos.getPhotos()[0],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: true,
    lastInAlbum: false
  },
};

export const MiddleOfAlbum: Story = {
  args: {
    source: photos.getPhotos()[1],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: false,
    lastInAlbum: false
  },
};

export const LastInAlbum: Story = {
  args: {
    source: photos.getPhotos()[2],
    previousPhoto: () => console.log("previous"),
    nextPhoto: () => console.log("next"),
    firstInAlbum: false,
    lastInAlbum: true
  },
};