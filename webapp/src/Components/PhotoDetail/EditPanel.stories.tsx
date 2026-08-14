import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {EditPanel} from './EditPanel';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

const photo = {
    id: 'p1',
    name: 'sunset.jpg',
    filesystemPath: '/photos/sunset.jpg',
    sourcePath: '/photos/sunset.jpg',
    albumId: '',
    tags: '',
    metadata: {},
    createdAt: '2023-01-01T00:00:00Z',
    thumbnailUrl: '',
};

const photosAdapter = {
    editPhoto: async () => photo,
    clearEdits: async () => photo,
} as unknown as IPhotosAdapter;

const meta: Meta<typeof EditPanel> = {
    title: 'PhotoDetail/EditPanel',
    component: EditPanel,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 420}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof EditPanel>;

export const Unedited: Story = {
    args: {
        photo: photo as never,
        photosAdapter,
        onSaved: () => {},
    },
};

export const WithRotateApplied: Story = {
    args: {
        photo: {...photo, editParams: {rotate: 90}} as never,
        photosAdapter,
        onSaved: () => {},
    },
};
