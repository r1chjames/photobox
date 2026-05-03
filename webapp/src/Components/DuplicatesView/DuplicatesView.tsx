import React, { useEffect, useState } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center, Text, Group, Button } from '@mantine/core';
import { IconCopy, IconPhotoOff } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';

interface DuplicatesViewProps {
    photosAdapter: IPhotosAdapter;
}

export const DuplicatesView: React.FC<DuplicatesViewProps> = ({ photosAdapter }) => {
    const [groups, setGroups] = useState<Map<string, Photo[]>>(new Map());
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadDuplicates = async () => {
            try {
                const photos = await photosAdapter.getDuplicatePhotos();
                const map = new Map<string, Photo[]>();
                photos.forEach(p => {
                    if (p.fileHash) {
                        const existing = map.get(p.fileHash) || [];
                        existing.push(p);
                        map.set(p.fileHash, existing);
                    }
                });
                setGroups(map);
            } catch (e) {
                setGroups(new Map());
            } finally {
                setLoading(false);
            }
        };
        loadDuplicates();
    }, [photosAdapter]);

    const handleCopyHash = (hash: string) => {
        navigator.clipboard.writeText(hash).then(() => {
            notifications.show({ title: 'Copied', message: 'File hash copied to clipboard', color: 'blue' });
        });
    };

    const handleDeleteAllButOne = (hash: string, photos: Photo[]) => {
        modals.openConfirmModal({
            title: 'Delete duplicates?',
            children: `This will permanently delete ${photos.length - 1} duplicate(s), keeping "${photos[0].name}". This action cannot be undone.`,
            labels: { confirm: 'Delete duplicates', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    const toDelete = photos.slice(1);
                    await Promise.all(toDelete.map(p => photosAdapter.deletePhoto(p.id)));
                    setGroups(prev => {
                        const next = new Map(prev);
                        const remaining = [photos[0]];
                        if (remaining.length > 1) {
                            next.set(hash, remaining);
                        } else {
                            next.delete(hash);
                        }
                        return next;
                    });
                    notifications.show({
                        title: 'Duplicates deleted',
                        message: `${toDelete.length} duplicate(s) removed, kept "${photos[0].name}"`,
                        color: 'green',
                    });
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

    const handleDeleteAll = (hash: string, photos: Photo[]) => {
        modals.openConfirmModal({
            title: 'Delete all photos?',
            children: `This will permanently delete all ${photos.length} photos in this group. This action cannot be undone.`,
            labels: { confirm: 'Delete all', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    await Promise.all(photos.map(p => photosAdapter.deletePhoto(p.id)));
                    setGroups(prev => {
                        const next = new Map(prev);
                        next.delete(hash);
                        return next;
                    });
                    notifications.show({
                        title: 'Photos deleted',
                        message: `${photos.length} photo(s) permanently deleted`,
                        color: 'green',
                    });
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

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    if (groups.size === 0) {
        return (
            <>
                <Title size="h4" mb="md">Duplicates</Title>
                <EmptyState
                    title="No duplicates found"
                    description="Photos with identical file content will appear here for review."
                    icon={<IconPhotoOff size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">Duplicates ({groups.size} groups)</Title>
            {Array.from(groups.entries()).map(([hash, photos]) => (
                <div key={hash} style={{ marginBottom: 24, padding: 16, border: '1px solid var(--mantine-color-default-border)', borderRadius: 8 }}>
                    <Group justify="space-between" mb="sm">
                        <Text size="sm" fw={500}>{photos.length} duplicates</Text>
                        <Group gap="xs">
                            <Button size="compact-xs" variant="light" leftSection={<IconCopy size={12} />} onClick={() => handleCopyHash(hash)}>
                                Copy hash
                            </Button>
                            <Button size="compact-xs" variant="light" color="red" onClick={() => handleDeleteAllButOne(hash, photos)}>
                                Delete all but one
                            </Button>
                            <Button size="compact-xs" variant="light" color="red" onClick={() => handleDeleteAll(hash, photos)}>
                                Delete all
                            </Button>
                        </Group>
                    </Group>
                    <Group gap="sm">
                        {photos.map(photo => (
                            <div key={photo.id} style={{ textAlign: 'center' }}>
                                <img
                                    src={photo.thumbnailUrl || ''}
                                    alt={photo.name}
                                    style={{ width: 80, height: 80, objectFit: 'cover', borderRadius: 4 }}
                                />
                                <Text size="xs" mt={4} lineClamp={1} style={{ maxWidth: 80 }}>{photo.name}</Text>
                            </div>
                        ))}
                    </Group>
                </div>
            ))}
        </div>
    );
};
