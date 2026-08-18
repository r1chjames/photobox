import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';
import { IAlbumsAdapter } from '../../Adapters/IAlbumsAdapter';
import { ISharesAdapter } from '../../Adapters/ISharesAdapter';
import { PhotoGrid } from '../PhotoGrid/PhotoGrid';
import { EmptyState } from '../EmptyState/EmptyState';
import { IconSearch, IconCalendar, IconTag, IconCamera, IconMapPin, IconAspectRatio, IconX } from '@tabler/icons-react';
import { Chip, Group, Button, Stack, Text, Select, TextInput } from '@mantine/core';
import { notifications } from '@mantine/notifications';
import { SmartAlbumRules } from '../../Models/SmartAlbumRules';

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
    const [camera, setCamera] = useState<string>('');
    const [hasGps, setHasGps] = useState<boolean>(false);
    const [orientation, setOrientation] = useState<string>('');
    const [selectedTags, setSelectedTags] = useState<string[]>([]);
    const [availableTags, setAvailableTags] = useState<string[]>([]);
    const [tagsLoading, setTagsLoading] = useState(false);
    const [savingSmartAlbum, setSavingSmartAlbum] = useState(false);

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

    const hasActiveFilters = query || fromDate || toDate || camera || hasGps || orientation || selectedTags.length > 0;

    const handleClearFilters = () => {
        setSearchParams({});
        setFromDate('');
        setToDate('');
        setCamera('');
        setHasGps(false);
        setOrientation('');
        setSelectedTags([]);
    };

    const handleSaveAsSmartAlbum = async () => {
        if (!albumsAdapter.createSmartAlbum) return;
        setSavingSmartAlbum(true);
        try {
            const name = window.prompt('Name this smart album', query ? `Search: ${query}` : 'Smart album');
            if (!name || !name.trim()) return;
            const rules: SmartAlbumRules = {
                startDate: startDateParam,
                endDate: endDateParam,
                camera: camera || undefined,
                hasGps,
                orientation: orientation || undefined,
                tags: selectedTags.length ? selectedTags : undefined,
                favorite: undefined,
                lowQuality: undefined,
            };
            await albumsAdapter.createSmartAlbum(name.trim(), rules);
            notifications.show({
                title: 'Smart album created',
                message: `"${name.trim()}" will stay up to date automatically`,
                color: 'teal',
            });
        } catch (e) {
            notifications.show({
                title: 'Failed to create smart album',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setSavingSmartAlbum(false);
        }
    };

    const handleTagToggle = (tag: string) => {
        setSelectedTags(prev =>
            prev.includes(tag) ? prev.filter(t => t !== tag) : [...prev, tag]
        );
    };

    // Build the tags string for PhotoGrid (comma-separated)
    const tagsParam = selectedTags.length > 0 ? selectedTags.join(',') : undefined;

    // Convert date inputs to the backend's expected ISO format. Date-only
    // inputs are interpreted as local midnight and sent as UTC.
    const startDateParam = fromDate ? new Date(`${fromDate}T00:00:00`).toISOString() : undefined;
    const endDateParam = toDate ? new Date(`${toDate}T23:59:59`).toISOString() : undefined;

    if (!hasActiveFilters) {
        return (
            <EmptyState
                title="Search your photos"
                description="Type a keyword in the search box above, or combine filters below to narrow your library."
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

                    {/* Advanced filter row */}
                    <Group gap="xs" wrap="wrap">
                        <TextInput
                            size="sm"
                            placeholder="Camera (e.g. Canon, iPhone)"
                            value={camera}
                            onChange={e => setCamera(e.currentTarget.value)}
                            leftSection={<IconCamera size={14} />}
                            style={{ maxWidth: 220 }}
                        />
                        <Chip checked={hasGps} onChange={() => setHasGps(prev => !prev)} size="sm" variant="outline">
                            <Group gap={4}>
                                <IconMapPin size={14} />
                                Has GPS
                            </Group>
                        </Chip>
                        <Select
                            size="sm"
                            placeholder="Orientation"
                            clearable
                            data={[
                                { value: 'landscape', label: 'Landscape' },
                                { value: 'portrait', label: 'Portrait' },
                                { value: 'square', label: 'Square' },
                            ]}
                            value={orientation || null}
                            onChange={(v) => setOrientation(v ?? '')}
                            leftSection={<IconAspectRatio size={14} />}
                            style={{ maxWidth: 180 }}
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
                        {hasActiveFilters && albumsAdapter.createSmartAlbum && (
                            <Button
                                variant="light"
                                color="teal"
                                size="xs"
                                onClick={handleSaveAsSmartAlbum}
                                loading={savingSmartAlbum}
                            >
                                Save as smart album
                            </Button>
                        )}
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
                    startDate={startDateParam}
                    endDate={endDateParam}
                    camera={camera || undefined}
                    hasGps={hasGps}
                    orientation={orientation || undefined}
                />
            </div>
        </Stack>
    );
};
