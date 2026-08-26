import React from 'react';
import { Menu } from '@mantine/core';
import {
    IconDownload,
    IconHeart,
    IconPhotoPlus,
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
    onAddToAlbum?: () => void;
    onCancel: () => void;
}

export const BulkActionsToolbar: React.FC<BulkActionsToolbarProps> = ({
    selectedCount,
    onSelectAll,
    onDeselectAll,
    onFavorite,
    onDelete,
    onDownload,
    onAddTag,
    onRemoveTag,
    onAddToAlbum,
    onCancel,
}) => {
    return (
        <div className="selection-bar">
            <span className="selection-bar-count">
                {selectedCount} photo{selectedCount !== 1 ? 's' : ''} selected
            </span>
            <span className="selection-bar-divider" />
            <button type="button" className="selection-bar-btn" onClick={onDownload} title="Download">
                <IconDownload size={16} stroke={1.8} /> Download
            </button>
            <button type="button" className="selection-bar-btn" onClick={onFavorite} title="Favorite">
                <IconHeart size={16} stroke={1.8} /> Favorite
            </button>
            {onAddToAlbum && (
                <button type="button" className="selection-bar-btn" onClick={onAddToAlbum} title="Add to album">
                    <IconPhotoPlus size={16} stroke={1.8} /> Add to album
                </button>
            )}
            {(onAddTag || onRemoveTag) && (
                <Menu withinPortal position="top" offset={8}>
                    <Menu.Target>
                        <button type="button" className="selection-bar-btn" title="Tags">
                            <IconTag size={16} stroke={1.8} /> Tags
                        </button>
                    </Menu.Target>
                    <Menu.Dropdown>
                        {onAddTag && (
                            <Menu.Item leftSection={<IconTag size={14} />} onClick={onAddTag}>Add tags</Menu.Item>
                        )}
                        {onRemoveTag && (
                            <Menu.Item leftSection={<IconTagOff size={14} />} onClick={onRemoveTag}>Remove tags</Menu.Item>
                        )}
                    </Menu.Dropdown>
                </Menu>
            )}
            <button type="button" className="selection-bar-btn" onClick={onSelectAll} title="Select all">
                <IconSelectAll size={16} stroke={1.8} /> All
            </button>
            <button type="button" className="selection-bar-btn" onClick={onDeselectAll} title="Deselect all">
                <IconSquare size={16} stroke={1.8} /> None
            </button>
            <button type="button" className="selection-bar-btn selection-bar-btn-danger" onClick={onDelete} title="Delete">
                <IconTrash size={16} stroke={1.8} /> Delete
            </button>
            <span className="selection-bar-divider" />
            <button type="button" className="selection-bar-close" onClick={onCancel} aria-label="Cancel selection">
                <IconX size={16} stroke={2} />
            </button>
        </div>
    );
};
