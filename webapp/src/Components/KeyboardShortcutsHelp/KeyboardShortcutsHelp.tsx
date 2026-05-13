import React from 'react';
import {Modal, Table, Kbd, Text} from '@mantine/core';

interface KeyboardShortcutsHelpProps {
    opened: boolean;
    onClose: () => void;
}

export const KeyboardShortcutsHelp: React.FC<KeyboardShortcutsHelpProps> = ({ opened, onClose }) => {
    const shortcuts = [
        { keys: ['←', '→'], action: 'Navigate between photos in lightbox' },
        { keys: ['Esc'], action: 'Close lightbox / go back' },
        { keys: ['D'], action: 'Download current photo' },
        { keys: ['F'], action: 'Toggle favorite on current photo' },
        { keys: ['E'], action: 'Toggle EXIF / info panel on photo detail' },
        { keys: ['Enter'], action: 'Open focused photo' },
        { keys: ['?'], action: 'Show this help dialog' },
    ];

    return (
        <Modal opened={opened} onClose={onClose} title="Keyboard Shortcuts" centered size="md">
            <Text size="sm" c="dimmed" mb="md">
                Press these keys anywhere in the app to quickly navigate and interact with your photos.
            </Text>
            <Table>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Key</Table.Th>
                        <Table.Th>Action</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {shortcuts.map((shortcut, index) => (
                        <Table.Tr key={index}>
                            <Table.Td>
                                <span style={{ display: 'flex', gap: 4 }}>
                                    {shortcut.keys.map((key, i) => (
                                        <Kbd key={i} size="sm">{key}</Kbd>
                                    ))}
                                </span>
                            </Table.Td>
                            <Table.Td>{shortcut.action}</Table.Td>
                        </Table.Tr>
                    ))}
                </Table.Tbody>
            </Table>
        </Modal>
    );
};
