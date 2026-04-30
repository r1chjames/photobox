import type {Meta, StoryObj} from '@storybook/react-vite';
import {KeyboardShortcutsHelp} from './KeyboardShortcutsHelp';

const meta: Meta<typeof KeyboardShortcutsHelp> = {
    component: KeyboardShortcutsHelp,
};

export default meta;
type Story = StoryObj<typeof KeyboardShortcutsHelp>;

export const Open: Story = {
    args: {
        opened: true,
        onClose: () => console.log('close'),
    },
};
