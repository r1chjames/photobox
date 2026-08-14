import React, { useEffect, useState, useCallback, useRef } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { ISettingsAdapter } from '../../Adapters/ISettingsAdapter';
import { Photo } from '../../Models/Photo';
import { EmptyState } from '../EmptyState/EmptyState';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import { Badge, Button, Group, Title, Text, Card, Image, Loader, Center, Skeleton, Tooltip } from '@mantine/core';
import { IconTrash, IconRestore, IconPhotoOff, IconClockHour3 } from '@tabler/icons-react';
import { fetchThumbnailWithAuth, getCachedThumbnail, revokeThumbnail } from '../../utils/ThumbnailUtils';

interface TrashViewProps {
    photosAdapter: IPhotosAdapter;
    settingsAdapter?: ISettingsAdapter;
}

export const TrashView: React.FC<TrashViewProps> = ({ photosAdapter, settingsAdapter }) => {
    const [photos, setPhotos] = useState<Photo[]>([]);
    const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());
    const thumbnailUrlsRef = useRef(thumbnailUrls);
    const [loading, setLoading] = useState(true);
    const [retentionDays, setRetentionDays] = useState<number | null>(null);

    const loadSettings = useCallback(async () => {
        if (!settingsAdapter) return;
        try {
            const settings = await settingsAdapter.getAllSettings();
            const retention = settings?.find(s => s.key === 'trash_retention_days');
            if (retention && retention.value !== '') {
                const parsed = Number(retention.value);
                if (!Number.isNaN(parsed)) {
                    setRetentionDays(parsed);
                }
            }
        } catch {
            // Non-fatal: countdown just won't show if settings can't be read.
        }
    }, [settingsAdapter]);

    const loadTrash = useCallback(async () => {
        try {
            const trashed = await photosAdapter.getTrashedPhotos('', 100);
            const safeTrashed = trashed || [];
            setPhotos(safeTrashed);

            const urls = new Map<string, string>();
            for (const photo of safeTrashed) {
                const cached = getCachedThumbnail(photo.id);
                if (cached) {
                    urls.set(photo.id, cached);
                } else {
                    const url = await fetchThumbnailWithAuth(photosAdapter, photo.id);
                    urls.set(photo.id, url);
                }
            }
            setThumbnailUrls(urls);
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
      thumbnailUrlsRef.current = thumbnailUrls;
    }, [thumbnailUrls]);

    useEffect(() => {
      loadSettings();
      loadTrash();
      return () => {
        thumbnailUrlsRef.current.forEach((_, id) => revokeThumbnail(id));
      };
    }, [loadTrash, loadSettings]);

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

    // Days until the photo is auto-deleted by the retention job (null when
    // retention is disabled, unknown, or the photo has no recorded delete time).
    const autoDeleteDaysLeft = (photo: Photo): number | null => {
        if (!photo.deletedAt || retentionDays === null || retentionDays <= 0) return null;
        const deletedMs = Date.parse(photo.deletedAt);
        if (Number.isNaN(deletedMs)) return null;
        const autoDeleteAt = deletedMs + retentionDays * 24 * 60 * 60 * 1000;
        const daysLeft = Math.ceil((autoDeleteAt - Date.now()) / (24 * 60 * 60 * 1000));
        return Math.max(daysLeft, 0);
    };

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    const retentionDescription = retentionDays === null
        ? 'Deleted photos will appear here for 30 days before being permanently removed.'
        : retentionDays <= 0
            ? 'Automatic cleanup is disabled. Trashed photos are kept until you remove them manually.'
            : `Trashed photos are automatically deleted after ${retentionDays} day${retentionDays !== 1 ? 's' : ''}.`;

    if (photos.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Trash</Title>
                <EmptyState
                    title="Trash is empty"
                    description={retentionDescription}
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
            <Text size="sm" c="dimmed" mb="md">{retentionDescription}</Text>
            <Group gap="md">
                {photos.map(photo => {
                    const daysLeft = autoDeleteDaysLeft(photo);
                    return (
                        <Card key={photo.id} shadow="sm" radius="md" withBorder w={200}>
                            <Card.Section>
                                {thumbnailUrls.get(photo.id) ? (
                                    <Image src={thumbnailUrls.get(photo.id)} h={150} fit="cover" />
                                ) : (
                                    <Skeleton height={150} />
                                )}
                            </Card.Section>
                            <Text size="sm" fw={500} mt="sm" lineClamp={1}>{photo.name}</Text>
                            {daysLeft !== null && (
                                <Group gap="xs" mt="xs">
                                    <Tooltip label={`Auto-deleted in ~${daysLeft} day${daysLeft !== 1 ? 's' : ''}`}>
                                        <Badge variant="light" color={daysLeft <= 7 ? 'red' : 'gray'} leftSection={<IconClockHour3 size={12} />}>
                                            {daysLeft === 0 ? 'due' : `${daysLeft}d`}
                                        </Badge>
                                    </Tooltip>
                                </Group>
                            )}
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
