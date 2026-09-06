import type {Meta, StoryObj} from '@storybook/react-vite';
import {TimelineScrubber} from './TimelineScrubber';
import {MockPhotosAdapter} from '../../Adapters/MockPhotosAdapter';
import {SBModelBuilder} from '../../utils/SBModelBuilder';

// The desktop rail (issue #174) is absolutely positioned against its nearest
// positioned ancestor, exactly as PhotoGrid mounts it. These stories wrap the
// scrubber in a relative container so the overlay renders in isolation.
const meta: Meta<typeof TimelineScrubber> = {
    component: TimelineScrubber,
    decorators: [
        (Story) => (
            <div style={{position: 'relative', height: 480, background: '#141518'}}>
                <Story />
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof TimelineScrubber>;

const adapter = new MockPhotosAdapter()
    .withPhotos(new SBModelBuilder().newAlbumWithPhotos(120).getPhotos());

// Desktop: the rail overlays the right gutter — a 48px strip with the track
// line, month ticks and hover/select chips.
export const DesktopRail: Story = {
    args: {
        photosAdapter: adapter,
        onSelectMonth: () => {},
        onClear: () => {},
    },
};

// Active filter: a month is selected (chip filled) and the clear button shows.
export const DesktopRailActiveFilter: Story = {
    args: {
        photosAdapter: adapter,
        onSelectMonth: () => {},
        onClear: () => {},
        activeYear: 2024,
        activeMonth: 1,
    },
};

// Mobile: the scrubber falls back to the floating action button + drawer.
export const MobileFab: Story = {
    args: {
        photosAdapter: adapter,
        onSelectMonth: () => {},
        onClear: () => {},
    },
    parameters: {
        viewport: {
            defaultViewport: 'mobile2',
        },
    },
};
