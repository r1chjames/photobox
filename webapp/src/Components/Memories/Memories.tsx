import React, { useEffect, useState, useCallback } from 'react';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { Photo } from '../../Models/Photo';
import { Card, Image, Text, Group, Title, Skeleton } from '@mantine/core';
import { IconClock } from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { fetchThumbnailWithAuth, getCachedThumbnail, revokeThumbnail } from '../../utils/ThumbnailUtils';

interface MemoriesProps {
    photosAdapter: IPhotosAdapter;
}

export const Memories: React.FC<MemoriesProps> = ({ photosAdapter }) => {
    const [memories, setMemories] = useState<Photo[]>([]);
    const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    const loadMemories = useCallback(async () => {
        try {
            const today = new Date();
            const month = today.getMonth() + 1;
            const day = today.getDate();

            const allPhotos = await photosAdapter.getAllPhotosInfo('', 500, false);
            const filtered = allPhotos.filter(photo => {
                if (!photo.createdAt) return false;
                const date = new Date(photo.createdAt);
                return date.getMonth() + 1 === month && date.getDate() === day;
            });

            // Sort by year descending
            filtered.sort((a, b) => new Date(b.createdAt).getFullYear() - new Date(a.createdAt).getFullYear());
            const sliced = filtered.slice(0, 6);
            setMemories(sliced);

            // Load thumbnails
            const urls = new Map<string, string>();
            for (const photo of sliced) {
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
            setMemories([]);
        } finally {
            setLoading(false);
        }
    }, [photosAdapter]);

    useEffect(() => {
        loadMemories();
        return () => {
            thumbnailUrls.forEach((_, id) => revokeThumbnail(id));
        };
    }, [loadMemories]);

    if (loading) {
        return (
            <Card withBorder radius="md" p="md" mb="md">
                <Group gap="xs" mb="sm">
                    <IconClock size="1.25rem" />
                    <Title size="h5">On This Day</Title>
                </Group>
                <Group gap="sm">
                    {Array.from({ length: 4 }).map((_, i) => (
                        <Skeleton key={i} height={120} width={160} radius="sm" />
                    ))}
                </Group>
            </Card>
        );
    }

    if (memories.length === 0) {
        return null;
    }

    return (
        <Card withBorder radius="md" p="md" mb="md">
            <Group gap="xs" mb="sm">
                <IconClock size="1.25rem" />
                <Title size="h5">On This Day</Title>
            </Group>
            <Group gap="sm">
                {memories.map(photo => {
                    const year = new Date(photo.createdAt).getFullYear();
                    const url = thumbnailUrls.get(photo.id);
                    return (
                        <Card key={photo.id} p={0} radius="sm" withBorder style={{ cursor: 'pointer' }} onClick={() => navigate(`/photo/${photo.id}`)}>
                            {url ? (
                                <Image src={url} height={120} width={160} fit="cover" radius="sm" />
                            ) : (
                                <Skeleton height={120} width={160} radius="sm" />
                            )}
                            <Text size="xs" ta="center" mt={4} c="dimmed">{year}</Text>
                        </Card>
                    );
                })}
            </Group>
        </Card>
    );
};
