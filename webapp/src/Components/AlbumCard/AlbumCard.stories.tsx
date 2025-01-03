import type {Meta, StoryObj} from '@storybook/react';
import {AlbumCard} from './AlbumCard';
import {SBModelBuilder} from "../../utils/SBModelBuilder";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof AlbumCard> = {
    component: AlbumCard,
};

export default meta;
type Story = StoryObj<typeof AlbumCard>;

const albumWithPhotos = new SBModelBuilder().newAlbumWithPhotos(4);

export const Primary: Story = {
    args: {
        photosAdapter: new MockPhotosAdapter()
            .withPhotos(albumWithPhotos.getPhotos()),
        source: albumWithPhotos.getAlbums().at(0),
        albumViewCallback: () => console.log("Clicked"),
    },
};