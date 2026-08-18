import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {Memories} from './Memories';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

const photo = (id: string, name: string) => ({
    id,
    name,
    filesystemPath: '',
    sourcePath: '',
    albumId: '',
    tags: '',
    metadata: {},
    createdAt: '2023-01-01T00:00:00Z',
    thumbnailUrl: '',
});

const photosAdapter = {
    getMemories: async () => [
        {year: 2023, yearsAgo: 3, photos: [photo('p1', 'beach_2023.jpg'), photo('p2', 'sunset_2023.jpg')]},
        {year: 2022, yearsAgo: 4, photos: [photo('p3', 'family_2022.jpg')]},
    ],
} as unknown as IPhotosAdapter;

const emptyAdapter = {
    getMemories: async () => [],
} as unknown as IPhotosAdapter;

const meta: Meta<typeof Memories> = {
    title: 'Views/Memories',
    component: Memories,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 800}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof Memories>;

export const Populated: Story = {
    args: {
        photosAdapter,
    },
};

export const Empty: Story = {
    args: {
        photosAdapter: emptyAdapter,
    },
};
