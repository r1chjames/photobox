import type {Meta, StoryObj} from '@storybook/react-vite';
import {EmptyState} from './EmptyState';
import {IconSearch, IconAlbumOff} from '@tabler/icons-react';

const meta: Meta<typeof EmptyState> = {
    component: EmptyState,
};

export default meta;
type Story = StoryObj<typeof EmptyState>;

export const Default: Story = {
    args: {
        title: 'No photos yet',
        description: 'Your photo library is empty. Upload photos or configure your photo directory.',
    },
};

export const WithAction: Story = {
    args: {
        title: 'No albums yet',
        description: 'Albums appear automatically from your photo folders, or create one manually.',
        icon: <IconAlbumOff size="2rem" />,
        action: {
            label: 'Create album',
            onClick: () => console.log('Create album clicked'),
        },
    },
};

export const SearchEmpty: Story = {
    args: {
        title: 'No results found',
        description: 'No photos match your search. Try a different keyword.',
        icon: <IconSearch size="2rem" />,
    },
};

export const Minimal: Story = {
    args: {
        title: 'This album is empty',
    },
};
