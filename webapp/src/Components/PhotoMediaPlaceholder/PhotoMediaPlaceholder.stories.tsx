import React from 'react';
import type {Meta, StoryObj} from '@storybook/react-vite';
import {PhotoMediaPlaceholder} from './PhotoMediaPlaceholder';

// Valid 3x4-component blurhash (also used by the issue #173 photo-grid
// fixtures); decodes to 32x24 pixels and paints a blurred placeholder.
const VALID_BLURHASH = 'LEHV6nWB2yk8pyo0adR*.7kCMdnj';

const meta: Meta<typeof PhotoMediaPlaceholder> = {
    component: PhotoMediaPlaceholder,
    decorators: [
        (Story) => (
            <div style={{position: 'relative', width: 320, height: 220, borderRadius: 8, overflow: 'hidden', background: '#16181d'}}>
                <Story />
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof PhotoMediaPlaceholder>;

export const WithBlurhash: Story = {
    args: {
        blurhash: VALID_BLURHASH,
    },
};

export const DominantColorOnly: Story = {
    args: {
        dominantColor: '#667788',
    },
};

export const NoPlaceholderData: Story = {
    args: {},
};

export const InvalidBlurhashFallsBackToSkeleton: Story = {
    args: {
        // Not a valid blurhash: decode throws / validation fails, so the
        // placeholder must fall through to the skeleton — never an empty
        // canvas (issue #173).
        blurhash: 'invalid!!!',
    },
};
