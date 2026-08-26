import React, { useState } from 'react';
import {Album, isSmartAlbum} from '../../Models/Album';
import {ActionIcon, Badge, Loader, Menu, TextInput} from '@mantine/core';
import {IconDotsVertical, IconPencil, IconTrash, IconPhotoOff} from '@tabler/icons-react';
import { modals } from '@mantine/modals';
import { notifications } from '@mantine/notifications';
import useAlbumCard from "./useAlbumCard";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";

interface IProps {
  photosAdapter: IPhotosAdapter;
  albumsAdapter: IAlbumsAdapter;
  source: Album;
  albumViewCallback: (albumId: string) => void;
  onAlbumUpdated?: () => void;
}

export const AlbumCard: React.FunctionComponent<IProps> = (props) => {
  const { thumbnailUrl, subThumbnailUrls = [], photoCount, isLoading } = useAlbumCard(props.photosAdapter, props.source);
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameValue, setRenameValue] = useState(props.source.name);

  const handleRename = async () => {
    if (!renameValue.trim() || renameValue === props.source.name) {
      setIsRenaming(false);
      return;
    }
    try {
      await props.albumsAdapter.updateAlbum(props.source.id, { name: renameValue.trim() });
      notifications.show({ title: 'Album renamed', message: `Renamed to "${renameValue.trim()}"`, color: 'green' });
      setIsRenaming(false);
      props.onAlbumUpdated?.();
    } catch (e) {
      notifications.show({ title: 'Rename failed', message: e instanceof Error ? e.message : 'Error', color: 'red' });
    }
  };

  const handleDelete = () => {
    modals.openConfirmModal({
      title: 'Delete album?',
      children: `Are you sure you want to delete "${props.source.name}"? Photos will be moved to Uncategorized.`,
      labels: { confirm: 'Delete', cancel: 'Cancel' },
      confirmProps: { color: 'red' },
      onConfirm: async () => {
        try {
          await props.albumsAdapter.deleteAlbum(props.source.id, false);
          notifications.show({ title: 'Album deleted', message: `"${props.source.name}" has been deleted`, color: 'green' });
          props.onAlbumUpdated?.();
        } catch (e) {
          notifications.show({ title: 'Delete failed', message: e instanceof Error ? e.message : 'Error', color: 'red' });
        }
      },
    });
  };

  if (isLoading) {
    return <Loader size="md" />;
  }

  const subThumb = (url: string | undefined, key: string) => (
    <div
      key={key}
      className="album-sub-thumb"
      style={url ? { backgroundImage: `url(${url})` } : undefined}
      aria-hidden="true"
    />
  );

  return (
    <div className="album-card" onClick={() => !isRenaming && props.albumViewCallback(props.source.id)}>
      <div style={{ position: 'absolute', top: 8, right: 8, zIndex: 3 }} onClick={(e) => e.stopPropagation()}>
        <Menu withinPortal position="bottom-end">
          <Menu.Target>
            <ActionIcon variant="transparent" size="sm" radius="xl" color="gray" aria-label="Album actions">
              <IconDotsVertical size="1rem" />
            </ActionIcon>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Item leftSection={<IconPencil size={14} />} onClick={() => { setIsRenaming(true); }}>
              Rename
            </Menu.Item>
            <Menu.Item leftSection={<IconTrash size={14} />} color="red" onClick={handleDelete}>
              Delete
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </div>

      <div className="album-media-grid">
        <div className="album-main-thumb">
          {thumbnailUrl ? (
            <img src={thumbnailUrl} alt={props.source.name} loading="lazy" />
          ) : (
            <div
              style={{
                width: '100%',
                height: '100%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'var(--pb-surface-hover)',
              }}
            >
              <IconPhotoOff size={32} color="var(--pb-muted)" />
            </div>
          )}
        </div>
        <div className="album-sub-thumbs">
          {subThumb(subThumbnailUrls[0], 'sub-0')}
          {subThumb(subThumbnailUrls[1], 'sub-1')}
        </div>
      </div>

      <div className="album-info-overlay">
        {isRenaming ? (
          <TextInput
            size="xs"
            value={renameValue}
            onChange={(e) => setRenameValue(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') handleRename(); if (e.key === 'Escape') setIsRenaming(false); }}
            onBlur={handleRename}
            autoFocus
            style={{ flex: 1 }}
            onClick={(e) => e.stopPropagation()}
          />
        ) : (
          <span className="album-title">{props.source.name}</span>
        )}
        <span style={{ display: 'inline-flex', gap: 6 }}>
          {isSmartAlbum(props.source) && (
            <Badge color="teal" variant="light" size="xs" radius="xl">Smart</Badge>
          )}
          <Badge className="album-count-badge" variant="light" color="gray" size="sm" radius="xl">
            {photoCount} photos
          </Badge>
        </span>
      </div>
    </div>
  );
};
