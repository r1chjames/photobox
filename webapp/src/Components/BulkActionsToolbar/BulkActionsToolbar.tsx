import React from 'react';
import { ActionIcon, Button, Group, Text, Tooltip } from '@mantine/core';
import {
    IconDownload,
    IconHeart,
    IconSelectAll,
    IconSquare,
    IconTag,
    IconTagOff,
    IconTrash,
    IconX,
} from '@tabler/icons-react';

interface BulkActionsToolbarProps {
    selectedCount: number;
    totalCount: number;
    onSelectAll: () => void;
    onDeselectAll: () => void;
    onFavorite: () => void;
    onDelete: () => void;
    onDownload: () => void;
    onAddTag?: () => void;
    onRemoveTag?: () => void;
    onCancel: () => void;
}

export const BulkActionsToolbar: React.FC<BulkActionsToolbarProps> = ({
    selectedCount,
    totalCount,
    onSelectAll,
    onDeselectAll,
    onFavorite,
    onDelete,
    onDownload,
    onAddTag,
    onRemoveTag,
    onCancel,
}) => {
    return (
        <Group justify="space-between" p="sm" style={{ background: 'var(--mantine-color-body)', borderBottom: '1px solid var(--mantine-color-default-border)' }}>
            <Group gap="sm">
                <Text size="sm" fw={500}>
                    {selectedCount} selected
                </Text>
                {selectedCount < totalCount ? (
                    <Button variant="subtle" size="compact-sm" leftSection={<IconSelectAll size={14} />} onClick={onSelectAll}>
                        Select all
                    </Button>
                ) : (
                    <Button variant="subtle" size="compact-sm" leftSection={<IconSquare size={14} />} onClick={onDeselectAll}>
                        Deselect all
                    </Button>
                )}
            </Group>
            <Group gap="xs">
                <Tooltip label="Favorite">
                    <ActionIcon variant="light" color="pink" onClick={onFavorite} disabled={selectedCount === 0}>
                        <IconHeart size="1.25rem" />
                    </ActionIcon>
                </Tooltip>
                {onAddTag && (
                    <Tooltip label="Add tag">
                        <ActionIcon variant="light" color="green" onClick={onAddTag} disabled={selectedCount === 0}>
                            <IconTag size="1.25rem" />
                        </ActionIcon>
                    </Tooltip>
                )}
                {onRemoveTag && (
                    <Tooltip label="Remove tag">
                        <ActionIcon variant="light" color="orange" onClick={onRemoveTag} disabled={selectedCount === 0}>
                            <IconTagOff size="1.25rem" />
                        </ActionIcon>
                    </Tooltip>
                )}
                <Tooltip label="Download">
                    <ActionIcon variant="light" color="blue" onClick={onDownload} disabled={selectedCount === 0}>
                        <IconDownload size="1.25rem" />
                    </ActionIcon>
                </Tooltip>
                <Tooltip label="Delete">
                    <ActionIcon variant="light" color="red" onClick={onDelete} disabled={selectedCount === 0}>
                        <IconTrash size="1.25rem" />
                    </ActionIcon>
                </Tooltip>
                <Tooltip label="Cancel">
                    <ActionIcon variant="default" onClick={onCancel}>
                        <IconX size="1.25rem" />
                    </ActionIcon>
                </Tooltip>
            </Group>
        </Group>
    );
};
