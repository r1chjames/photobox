import type {Meta, StoryObj} from '@storybook/react-vite';
import {BulkActionsToolbar} from './BulkActionsToolbar';

const meta: Meta<typeof BulkActionsToolbar> = {
    component: BulkActionsToolbar,
};

export default meta;
type Story = StoryObj<typeof BulkActionsToolbar>;

export const SomeSelected: Story = {
    args: {
        selectedCount: 5,
        totalCount: 20,
        onSelectAll: () => console.log('select all'),
        onDeselectAll: () => console.log('deselect all'),
        onFavorite: () => console.log('favorite'),
        onDelete: () => console.log('delete'),
        onDownload: () => console.log('download'),
        onCancel: () => console.log('cancel'),
    },
};

export const AllSelected: Story = {
    args: {
        selectedCount: 20,
        totalCount: 20,
        onSelectAll: () => console.log('select all'),
        onDeselectAll: () => console.log('deselect all'),
        onFavorite: () => console.log('favorite'),
        onDelete: () => console.log('delete'),
        onDownload: () => console.log('download'),
        onCancel: () => console.log('cancel'),
    },
};

export const NoneSelected: Story = {
    args: {
        selectedCount: 0,
        totalCount: 20,
        onSelectAll: () => console.log('select all'),
        onDeselectAll: () => console.log('deselect all'),
        onFavorite: () => console.log('favorite'),
        onDelete: () => console.log('delete'),
        onDownload: () => console.log('download'),
        onCancel: () => console.log('cancel'),
    },
};
