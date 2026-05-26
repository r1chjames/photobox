import React, { useCallback, useEffect, useState } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center, Text, Group, Button, Skeleton, ActionIcon, Tooltip } from '@mantine/core';
import { IconCopy, IconPhotoOff, IconRefresh } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import { fetchThumbnailsBatch, getCachedThumbnail, fetchThumbnailWithAuth, revokeThumbnail } from '../../utils/ThumbnailUtils';
import { BlurhashCanvas } from '../BlurhashCanvas/BlurhashCanvas';

interface DuplicatesViewProps {
    photosAdapter: IPhotosAdapter;
}

interface DuplicatePhotoItemProps {
    photo: Photo;
    thumbnailUrl: string | undefined;
    onRetry: (photoId: string) => void;
}

const DuplicatePhotoItem: React.FC<DuplicatePhotoItemProps> = ({ photo, thumbnailUrl, onRetry }) => {
    const [isLoaded, setIsLoaded] = useState(false);
    const [hasError, setHasError] = useState(false);

    useEffect(() => {
        setIsLoaded(false);
        setHasError(false);
    }, [thumbnailUrl]);

    return (
        <div style={{ textAlign: 'center' }}>
            <div style={{ position: 'relative', width: 80, height: 80 }}>
                {(!thumbnailUrl || !isLoaded) && !hasError && (
                    photo.blurhash ? (
                        <BlurhashCanvas
                            hash={photo.blurhash}
                            style={{ position: 'absolute', top: 0, left: 0, borderRadius: 4 }}
                        />
                    ) : photo.dominantColor ? (
                        <div
                            style={{
                                position: 'absolute',
                                top: 0,
                                left: 0,
                                width: 80,
                                height: 80,
                                backgroundColor: photo.dominantColor,
                                borderRadius: 4,
                            }}
                        />
                    ) : (
                        <Skeleton
                            height={80}
                            width={80}
                            style={{ position: 'absolute', top: 0, left: 0, borderRadius: 4 }}
                        />
                    )
                )}
                {thumbnailUrl && !hasError && (
                    <img
                        src={thumbnailUrl}
                        alt={photo.name}
                        onLoad={() => setIsLoaded(true)}
                        onError={() => setHasError(true)}
                        style={{ width: 80, height: 80, objectFit: 'cover', borderRadius: 4, position: 'relative', zIndex: 1 }}
                    />
                )}
                {hasError && (
                    <div style={{
                        position: 'absolute',
                        top: 0, left: 0, right: 0, bottom: 0,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        background: 'var(--mantine-color-gray-2)',
                        borderRadius: 4,
                    }}>
                        <Tooltip label="Retry loading thumbnail">
                            <ActionIcon
                                variant="light"
                                size="sm"
                                onClick={() => onRetry(photo.id)}
                            >
                                <IconRefresh size={16} />
                            </ActionIcon>
                        </Tooltip>
                    </div>
                )}
            </div>
            <Text size="xs" mt={4} lineClamp={1} style={{ maxWidth: 80 }}>{photo.name}</Text>
        </div>
    );
};

export const DuplicatesView: React.FC<DuplicatesViewProps> = ({ photosAdapter }) => {
    const [groups, setGroups] = useState<Map<string, Photo[]>>(new Map());
    const [loading, setLoading] = useState(true);
    const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());

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

    // Collect all photo IDs from all groups for thumbnail fetching
    const allPhotoIds = React.useMemo(() => {
        const ids: string[] = [];
        groups.forEach(photos => {
            photos.forEach(p => ids.push(p.id));
        });
        return ids;
    }, [groups]);

    // Merge React state with global cache
    const mergedThumbnailUrls = React.useMemo(() => {
        const merged = new Map(thumbnailUrls);
        for (const photoId of allPhotoIds) {
            if (!merged.has(photoId)) {
                const cached = getCachedThumbnail(photoId);
                if (cached) {
                    merged.set(photoId, cached);
                }
            }
        }
        return merged;
    }, [thumbnailUrls, allPhotoIds]);

    // Batch load thumbnails for all photos in duplicate groups
    useEffect(() => {
        if (allPhotoIds.length === 0) return;

        const loadThumbnails = async () => {
            // Sync any thumbnails already in the global cache
            const cachedOnly = new Map<string, string>();
            for (const photoId of allPhotoIds) {
                const cached = getCachedThumbnail(photoId);
                if (cached) {
                    cachedOnly.set(photoId, cached);
                }
            }
            if (cachedOnly.size > 0) {
                setThumbnailUrls(prev => {
                    const next = new Map(prev);
                    let changed = false;
                    cachedOnly.forEach((url, id) => {
                        if (!next.has(id)) {
                            next.set(id, url);
                            changed = true;
                        }
                    });
                    return changed ? next : prev;
                });
            }

            // Fetch thumbnails that aren't cached yet
            const uncachedIds = allPhotoIds.filter(id => !getCachedThumbnail(id));
            if (uncachedIds.length === 0) return;

            const newThumbnails = await fetchThumbnailsBatch(photosAdapter, uncachedIds);
            setThumbnailUrls(prev => {
                const next = new Map(prev);
                newThumbnails.forEach((url, id) => next.set(id, url));
                return next;
            });
        };

        loadThumbnails();
    }, [allPhotoIds, photosAdapter]);

    const handleRetryThumbnail = useCallback(async (photoId: string) => {
        revokeThumbnail(photoId);
        setThumbnailUrls(prev => {
            const next = new Map(prev);
            next.delete(photoId);
            return next;
        });
        try {
            const url = await fetchThumbnailWithAuth(photosAdapter, photoId);
            setThumbnailUrls(prev => {
                const next = new Map(prev);
                next.set(photoId, url);
                return next;
            });
        } catch {
            // Error state will be handled by the component's onError handler
        }
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
                            <DuplicatePhotoItem
                                key={photo.id}
                                photo={photo}
                                thumbnailUrl={mergedThumbnailUrls.get(photo.id)}
                                onRetry={handleRetryThumbnail}
                            />
                        ))}
                    </Group>
                </div>
            ))}
        </div>
    );
};
