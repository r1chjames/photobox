import React, { useEffect, useState } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center, Text, Group, Button } from '@mantine/core';
import { IconCopy, IconPhotoOff } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';
import { notifications } from '@mantine/notifications';

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
                        <Button size="compact-xs" variant="light" leftSection={<IconCopy size={12} />} onClick={() => handleCopyHash(hash)}>
                            Copy hash
                        </Button>
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
