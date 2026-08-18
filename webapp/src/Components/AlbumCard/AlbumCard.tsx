import React, { useState } from 'react';
import {Album, isSmartAlbum} from '../../Models/Album';
import {ActionIcon, Badge, Text, Card, Group, Image, Loader, Menu, TextInput} from '@mantine/core';
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
  const { thumbnailUrl, photoCount, isLoading } = useAlbumCard(props.photosAdapter, props.source);
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

  return (
    <Card shadow="sm" radius="md" withBorder mb="10px" mr="10px" w="200px" style={{ position: 'relative' }}>
      <div style={{ position: 'absolute', top: 4, right: 4, zIndex: 2 }}>
        <Menu withinPortal position="bottom-end">
          <Menu.Target>
            <ActionIcon variant="white" size="sm" opacity={0.8} onClick={(e) => e.stopPropagation()}>
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
      <Card.Section onClick={() => props.albumViewCallback(props.source.id)} style={{ cursor: 'pointer' }}>
        {thumbnailUrl ? (
          <Image src={thumbnailUrl} w="200px" h="150px" />
        ) : (
          <div style={{
            width: '200px',
            height: '150px',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: 'var(--mantine-color-gray-2)',
          }}>
            <IconPhotoOff size={32} color="var(--mantine-color-gray-5)" />
          </div>
        )}
      </Card.Section>
      <Group justify="space-between" mt="md" mb="xs" onClick={() => !isRenaming && props.albumViewCallback(props.source.id)} style={{ cursor: 'pointer' }}>
        {isRenaming ? (
          <TextInput
            size="xs"
            value={renameValue}
            onChange={(e) => setRenameValue(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter') handleRename(); if (e.key === 'Escape') setIsRenaming(false); }}
            onBlur={handleRename}
            autoFocus
            style={{ flex: 1 }}
          />
        ) : (
          <Text fw={500} lineClamp={1}>{props.source.name}</Text>
        )}
        <Badge color="blue" variant="light" size="xs">
          {photoCount} photos
        </Badge>
        {isSmartAlbum(props.source) && (
          <Badge color="teal" variant="light" size="xs">Smart</Badge>
        )}
      </Group>
    </Card>
  );
};
