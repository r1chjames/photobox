import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {GridSkeleton} from './GridSkeleton';

const meta: Meta<typeof GridSkeleton> = {
    title: 'PhotoGrid/GridSkeleton',
    component: GridSkeleton,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 900}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof GridSkeleton>;

export const Default: Story = {
    args: {
        columns: 5,
        rows: 3,
    },
};

export const Compact: Story = {
    args: {
        columns: 8,
        rows: 4,
        gap: 4,
    },
};
