import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Loader, Center, Chip, Group } from '@mantine/core';
import { IconTag } from '@tabler/icons-react';

interface TagsViewProps {
    photosAdapter: IPhotosAdapter;
}

export const TagsView: React.FC<TagsViewProps> = ({ photosAdapter }) => {
    const [tags, setTags] = useState<string[]>([]);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const loadTags = async () => {
            try {
                const allTags = await photosAdapter.getAllTags();
                setTags(allTags.sort());
            } catch {
                setTags([]);
            } finally {
                setLoading(false);
            }
        };
        loadTags();
    }, [photosAdapter]);

    if (loading) {
        return (
            <Center h="50vh">
                <Loader size="lg" />
            </Center>
        );
    }

    if (tags.length === 0) {
        return (
            <>
                <Title size="h4" mb="md">Tags</Title>
                <EmptyState
                    title="No tags yet"
                    description="Tags help you organize and find photos. Add tags to photos from the photo detail view."
                    icon={<IconTag size="2rem" />}
                />
            </>
        );
    }

    return (
        <div>
            <Title size="h4" mb="md">Tags ({tags.length})</Title>
            <Group gap="xs">
                {tags.map(tag => (
                    <Chip
                        key={tag}
                        onClick={() => navigate(`/tags/${encodeURIComponent(tag)}`)}
                        variant="filled"
                        color="blue"
                    >
                        {tag}
                    </Chip>
                ))}
            </Group>
        </div>
    );
};
