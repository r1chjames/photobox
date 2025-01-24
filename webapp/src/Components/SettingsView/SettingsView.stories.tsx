import type {Meta, StoryObj} from '@storybook/react';

import {SettingsView} from "./SettingsView";
import {MockSettingsAdapter} from "../../Adapters/MockSettingsAdapter";
import {Setting} from "../../Models/Setting";
import {MockPhotosAdapter} from "../../Adapters/MockPhotosAdapter";

const meta: Meta<typeof SettingsView> = {
    component: SettingsView,
};

const settings = [
    new Setting("thumbnail_width", "600", "Thumbnail width", "Photo", "Thumbnail width used during thumbnail generation"),
    new Setting("thumbnail_height", "600", "Thumbnail height", "Photo", "Thumbnail height used during thumbnail generation"),
    new Setting("default_new_albums_dir", "/photos", "New album storage location", "System", "Default location on disk to store new albums"),
    new Setting("index_frequency_cron", "0 1 * * *", "CRON expression for indexing", "System", "CRON expression used to initiate indexing")
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
