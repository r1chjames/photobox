import type {Meta, StoryObj} from '@storybook/react-vite';

import {SettingsView} from "./SettingsView";
import {MockSettingsAdapter} from "../../Adapters/MockSettingsAdapter";

import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof SettingsView> = {
    component: SettingsView,
};

const settings = [
    { key: "thumbnail_width", value: "600", friendlyName: "Thumbnail width", category: "Photo", description: "Thumbnail width used during thumbnail generation" },
    { key: "thumbnail_height", value: "600", friendlyName: "Thumbnail height", category: "Photo", description: "Thumbnail height used during thumbnail generation" },
    { key: "default_new_albums_dir", value: "/photos", friendlyName: "New album storage location", category: "System", description: "Default location on disk to store new albums" },
    { key: "index_frequency_cron", value: "0 1 * * *", friendlyName: "CRON expression for indexing", category: "System", description: "CRON expression used to initiate indexing" }
];


export default meta;
type Story = StoryObj<typeof SettingsView>;

export const Primary: Story = {
    args: {
        settingsAdapter: new MockSettingsAdapter()
            .withSettings(settings),
        photosAdapter: new MockPhotosAdapter()
    },
};
