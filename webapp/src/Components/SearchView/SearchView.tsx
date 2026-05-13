import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { EmptyState } from '../EmptyState/EmptyState';
import { IconSearch, IconCalendar, IconTag, IconX } from '@tabler/icons-react';
import { Chip, Group, Button, Stack, Text } from '@mantine/core';

interface SearchViewProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const SearchView: React.FC<SearchViewProps> = ({ photosAdapter, albumsAdapter, sharesAdapter }) => {
    const [searchParams, setSearchParams] = useSearchParams();
    const query = searchParams.get('q')?.trim() ?? '';

    const [fromDate, setFromDate] = useState<string>('');
    const [toDate, setToDate] = useState<string>('');
    const [selectedTags, setSelectedTags] = useState<string[]>([]);
    const [availableTags, setAvailableTags] = useState<string[]>([]);
    const [tagsLoading, setTagsLoading] = useState(false);

    useEffect(() => {
        const fetchTags = async () => {
            setTagsLoading(true);
            try {
                const tags = await photosAdapter.getAllTags();
                setAvailableTags(tags);
            } catch {
                // Silently fail — tags are optional
            } finally {
                setTagsLoading(false);
            }
        };
        fetchTags();
    }, [photosAdapter]);

    const hasActiveFilters = query || fromDate || toDate || selectedTags.length > 0;

    const handleClearFilters = () => {
        setSearchParams({});
        setFromDate('');
        setToDate('');
        setSelectedTags([]);
    };

    const handleTagToggle = (tag: string) => {
        setSelectedTags(prev =>
            prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]
        );
    };

    // Build the tags string for PhotoGrid (comma-separated)
    const tagsParam = selectedTags.length > 0 ? selectedTags.join(',') : undefined;

    // Note: PhotoGrid does not currently accept fromDate/toDate props.
    // Date range filtering UI is ready for when backend support is added.
    // Currently, date inputs are displayed but not passed to PhotoGrid.

    if (!hasActiveFilters) {
        return (
            <EmptyState
                title="Search your photos"
                description="Type a keyword in the search box above to find photos by name, album, or metadata."
                icon={<IconSearch size="2rem" />}
            />
        );
    }

    return (
        <Stack gap="md" style={{ height: '100%' }}>
            {/* Filter bar */}
            <div style={{
                padding: '12px 16px',
                background: 'var(--mantine-color-body)',
                borderBottom: '1px solid var(--mantine-color-gray-3)',
            }}>
                <Stack gap="sm">
                    {/* Date range row */}
                    <Group gap="xs" wrap="nowrap">
                        <IconCalendar size={18} style={{ flexShrink: 0, opacity: 0.6 }} />
                        <Text size="sm" fw={500} style={{ flexShrink: 0 }}>From:</Text>
                        <input
                            type="date"
                            value={fromDate}
                            onChange={e => setFromDate(e.target.value)}
                            style={{
                                padding: '4px 8px',
                                border: '1px solid var(--mantine-color-gray-4)',
                                borderRadius: '4px',
                                fontSize: '14px',
                                background: 'var(--mantine-color-body)',
                                color: 'var(--mantine-color-text)',
                                maxWidth: '160px',
                            }}
                        />
                        <Text size="sm" fw={500} style={{ flexShrink: 0 }}>To:</Text>
                        <input
                            type="date"
                            value={toDate}
                            onChange={e => setToDate(e.target.value)}
                            style={{
                                padding: '4px 8px',
                                border: '1px solid var(--mantine-color-gray-4)',
                                borderRadius: '4px',
                                fontSize: '14px',
                                background: 'var(--mantine-color-body)',
                                color: 'var(--mantine-color-text)',
                                maxWidth: '160px',
                            }}
                        />
                    </Group>

                    {/* Tags row */}
                    {availableTags.length > 0 && (
                        <Group gap="xs" wrap="wrap">
                            <IconTag size={18} style={{ flexShrink: 0, opacity: 0.6 }} />
                            <Text size="sm" fw={500} style={{ flexShrink: 0 }}>Tags:</Text>
                            {tagsLoading ? (
                                <Text size="sm" c="dimmed">Loading tags...</Text>
                            ) : (
                                <Group gap="xs" wrap="wrap">
                                    {availableTags.map(tag => (
                                        <Chip
                                            key={tag}
                                            checked={selectedTags.includes(tag)}
                                            onChange={() => handleTagToggle(tag)}
                                            size="sm"
                                            variant="outline"
                                        >
                                            {tag}
                                        </Chip>
                                    ))}
                                </Group>
                            )}
                        </Group>
                    )}

                    {/* Clear filters button */}
                    <Group gap="xs">
                        <Button
                            variant="light"
                            color="red"
                            size="xs"
                            leftSection={<IconX size={14} />}
                            onClick={handleClearFilters}
                        >
                            Clear all filters
                        </Button>
                        {selectedTags.length > 0 && (
                            <Text size="xs" c="dimmed">
                                {selectedTags.length} tag{selectedTags.length !== 1 ? 's' : ''} selected
                            </Text>
                        )}
                    </Group>
                </Stack>
            </div>

            {/* Photo grid */}
            <div style={{ flex: 1, overflow: 'hidden' }}>
                <PhotoGrid
                    photosAdapter={photosAdapter}
                    albumsAdapter={albumsAdapter}
                    sharesAdapter={sharesAdapter}
                    searchQuery={query || undefined}
                    tags={tagsParam}
                />
            </div>
        </Stack>
    );
};
