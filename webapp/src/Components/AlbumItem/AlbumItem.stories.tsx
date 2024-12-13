import type {Meta, StoryObj} from '@storybook/react';
import {AlbumItem} from './AlbumItem';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof AlbumItem> = {
    component: AlbumItem,
};

export default meta;
type Story = StoryObj<typeof AlbumItem>;

const albumWithPhotos = new SBModelBuilder().newAlbumWithPhotos(4);

export const Primary: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos()),
        source: albumWithPhotos.getAlbums().at(0),
        albumViewCallback: () => console.log("Clicked"),
    },
};