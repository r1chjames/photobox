import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {TrashView} from './TrashView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {ISettingsAdapter} from '../../Adapters/ISettingsAdapter';

const daysAgo = (days: number) => new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString();

const basePhoto = {
    id: 'p1',
    name: 'vacation_beach.jpg',
    filesystemPath: '/photos/album1/vacation_beach.jpg',
    sourcePath: '/photos/album1/vacation_beach.jpg',
    albumId: 'album1',
    tags: 'holiday,beach',
    metadata: {Make: 'Canon', Model: 'EOS R6'},
    createdAt: '2023-06-01T10:00:00Z',
    thumbnailUrl: '',
    deletedAt: daysAgo(10),
    favorite: false,
};

const photosAdapter = {
    getTrashedPhotos: async () => [
        {...basePhoto, id: 'p1', name: 'vacation_beach.jpg'},
        {...basePhoto, id: 'p2', name: 'family_dinner.jpg', deletedAt: daysAgo(50)},
        {...basePhoto, id: 'p3', name: 'old_memories.png', deletedAt: daysAgo(80)},
    ],
    restorePhoto: async () => undefined,
    deletePhoto: async () => undefined,
} as unknown as IPhotosAdapter;

const settingsAdapter = {
    getAllSettings: async () => [{ key: 'trash_retention_days', value: '60', friendlyName: 'Trash retention (days)', category: 'System', description: 'Days to keep photos in trash' }],
} as unknown as ISettingsAdapter;

const emptyPhotosAdapter = {
    getTrashedPhotos: async () => [],
    restorePhoto: async () => undefined,
    deletePhoto: async () => undefined,
} as unknown as IPhotosAdapter;

const meta: Meta<typeof TrashView> = {
    title: 'Views/TrashView',
    component: TrashView,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 900}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof TrashView>;

export const Populated: Story = {
    args: {
        photosAdapter,
        settingsAdapter,
    },
};

export const Empty: Story = {
    args: {
        photosAdapter: emptyPhotosAdapter,
        settingsAdapter,
    },
};

export const RetentionDisabled: Story = {
    args: {
        photosAdapter,
        settingsAdapter: {
            getAllSettings: async () => [{ key: 'trash_retention_days', value: '0', friendlyName: 'Trash retention (days)', category: 'System', description: 'Days to keep photos in trash' }],
        } as unknown as ISettingsAdapter,
    },
};
