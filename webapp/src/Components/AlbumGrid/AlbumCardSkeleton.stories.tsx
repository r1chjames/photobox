import type {Meta, StoryObj} from '@storybook/react-vite';
import {AlbumCardSkeleton} from './AlbumCardSkeleton';

const meta: Meta<typeof AlbumCardSkeleton> = {
    component: AlbumCardSkeleton,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 1000, display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 24}}>
                {Array.from({length: 3}).map((_, i) => <Story key={i} />)}
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof AlbumCardSkeleton>;

export const GridOfThree: Story = {};
