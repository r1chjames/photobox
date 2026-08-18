import React, { useEffect, useState, useCallback } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { MemoryGroup } from '../../Models/MemoryGroup';
import { Card, Image, Text, Group, Title, Skeleton, Stack, Badge } from '@mantine/core';
import { IconClock } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { fetchThumbnailWithAuth, getCachedThumbnail, revokeThumbnail } from '../../utils/ThumbnailUtils';

interface MemoriesProps {
    photosAdapter: IPhotosAdapter;
}

/**
 * "On This Day" — surfaces photos taken on today's date in previous years,
 * grouped by year. Backed by GET /photos/memories.
 */
export const Memories: React.FC<MemoriesProps> = ({ photosAdapter }) => {
    const [groups, setGroups] = useState<MemoryGroup[]>([]);
    const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    const loadMemories = useCallback(async () => {
        try {
            const result = await photosAdapter.getMemories();
            const safe = result || [];
            setGroups(safe);

            // Load thumbnails for all photos across groups
            const urls = new Map<string, string>();
            for (const group of safe) {
                for (const photo of group.photos) {
                    const cached = getCachedThumbnail(photo.id);
                    if (cached) {
                        urls.set(photo.id, cached);
                    } else {
                        const url = await fetchThumbnailWithAuth(photosAdapter, photo.id);
                        urls.set(photo.id, url);
                    }
                }
            }
            setThumbnailUrls(urls);
        } catch {
            setGroups([]);
        } finally {
            setLoading(false);
        }
    }, [photosAdapter]);

    useEffect(() => {
        loadMemories();
        return () => {
            thumbnailUrls.forEach((_, id) => revokeThumbnail(id));
        };
        // NOTE: react-hooks plugin not registered; deps intentionally [loadMemories].
    }, [loadMemories]);

    if (loading) {
        return (
            <Stack gap="md">
                <Skeleton height={24} width={200} />
                <Skeleton height={120} />
                <Skeleton height={120} />
            </Stack>
        );
    }

    if (groups.length === 0) {
        return (
            <Stack align="center" gap="xs" py="xl">
                <IconClock size={32} opacity={0.4} />
                <Text c="dimmed">No memories for today yet.</Text>
            </Stack>
        );
    }

    return (
        <Stack gap="lg">
            <Title order={3}>On This Day</Title>
            {groups.map(group => (
                <div key={group.year}>
                    <Group gap="xs" mb="sm">
                        <Title order={4}>{group.year}</Title>
                        <Badge variant="light" color="teal" size="sm">
                            {group.yearsAgo} year{group.yearsAgo !== 1 ? 's' : ''} ago
                        </Badge>
                    </Group>
                    <Group gap="sm" wrap="wrap">
                        {group.photos.map(photo => (
                            <Card
                                key={photo.id}
                                shadow="sm"
                                radius="md"
                                withBorder
                                w={160}
                                p={0}
                                style={{ cursor: 'pointer' }}
                                onClick={() => navigate(`/photo/${photo.id}`)}
                            >
                                <Card.Section>
                                    {thumbnailUrls.get(photo.id) ? (
                                        <Image src={thumbnailUrls.get(photo.id)} h={110} fit="cover" />
                                    ) : (
                                        <Skeleton height={110} />
                                    )}
                                </Card.Section>
                                <Text size="xs" fw={500} p="xs" lineClamp={1}>{photo.name}</Text>
                            </Card>
                        ))}
                    </Group>
                </div>
            ))}
        </Stack>
    );
};
