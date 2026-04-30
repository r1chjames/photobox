import React, { useEffect, useState, useCallback } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { Photo } from '../../Models/Photo';
import { EmptyState } from '../EmptyState/EmptyState';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import { Button, Group, Title, Text, Card, Image, Loader, Center } from '@mantine/core';
import { IconTrash, IconRestore, IconPhotoOff } from '@tabler/icons-react';

interface TrashViewProps {
    photosAdapter: IPhotosAdapter;
}

export const TrashView: React.FC<TrashViewProps> = ({ photosAdapter }) => {
    const [photos, setPhotos] = useState<Photo[]>([]);
    const [loading, setLoading] = useState(true);

    const loadTrash = useCallback(async () => {
        try {
            const trashed = await photosAdapter.getTrashedPhotos('', 100);
            setPhotos(trashed);
        } catch (e) {
            notifications.show({
                title: 'Failed to load trash',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setLoading(false);
        }
    }, [photosAdapter]);

    useEffect(() => {
        loadTrash();
    }, [loadTrash]);

    const handleRestore = async (photoId: string) => {
        try {
            await photosAdapter.restorePhoto(photoId);
            notifications.show({
                title: 'Photo restored',
                message: 'The photo has been restored to your library',
                color: 'green',
            });
            setPhotos(prev => prev.filter(p => p.id !== photoId));
        } catch (e) {
            notifications.show({
                title: 'Restore failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    };

    const handleDeletePermanently = (photo: Photo) => {
        modals.openConfirmModal({
            title: 'Permanently delete photo?',
            children: `This will permanently delete "${photo.name}". This action cannot be undone.`,
            labels: { confirm: 'Delete permanently', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    await photosAdapter.deletePhoto(photo.id);
                    notifications.show({
                        title: 'Photo deleted',
                        message: `"${photo.name}" has been permanently deleted`,
                        color: 'green',
                    });
                    setPhotos(prev => prev.filter(p => p.id !== photo.id));
                } catch (e) {
                    notifications.show({
                        title: 'Delete failed',
                        message: e instanceof Error ? e.message : 'An error occurred',
                        color: 'red',
                    });
                }
            },
        });
    };

    const handleEmptyTrash = () => {
        modals.openConfirmModal({
            title: 'Empty trash?',
            children: `This will permanently delete all ${photos.length} photo${photos.length !== 1 ? 's' : ''} in the trash. This action cannot be undone.`,
            labels: { confirm: 'Empty trash', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    for (const photo of photos) {
                        await photosAdapter.deletePhoto(photo.id);
                    }
                    notifications.show({
                        title: 'Trash emptied',
                        message: 'All photos have been permanently deleted',
                        color: 'green',
                    });
                    setPhotos([]);
                } catch (e) {
                    notifications.show({
                        title: 'Failed to empty trash',
                        message: e instanceof Error ? e.message : 'An error occurred',
                        color: 'red',
                    });
                }
            },
        });
    };

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    if (photos.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Trash</Title>
                <EmptyState
                    title="Trash is empty"
                    description="Deleted photos will appear here for 30 days before being permanently removed."
                    icon={<IconPhotoOff size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Group justify="space-between" mb="md">
                <Title size="h4">Trash ({photos.length})</Title>
                <Button color="red" variant="light" leftSection={<IconTrash size={16} />} onClick={handleEmptyTrash}>
                    Empty trash
                </Button>
            </Group>
            <Group gap="md">
                {photos.map(photo => {
                    const imgSrc = photo.thumbnail && (photo.thumbnail.startsWith('http') || photo.thumbnail.startsWith('data:image'))
                        ? photo.thumbnail
                        : `data:image/png;base64,${photo.thumbnail}`;
                    return (
                        <Card key={photo.id} shadow="sm" radius="md" withBorder w={200}>
                            <Card.Section>
                                <Image src={imgSrc} h={150} fit="cover" />
                            </Card.Section>
                            <Text size="sm" fw={500} mt="sm" lineClamp={1}>{photo.name}</Text>
                            <Group gap="xs" mt="sm">
                                <Button size="xs" variant="light" leftSection={<IconRestore size={14} />} onClick={() => handleRestore(photo.id)}>
                                    Restore
                                </Button>
                                <Button size="xs" color="red" variant="subtle" onClick={() => handleDeletePermanently(photo)}>
                                    Delete
                                </Button>
                            </Group>
                        </Card>
                    );
                })}
            </Group>
        </div>
    );
};
