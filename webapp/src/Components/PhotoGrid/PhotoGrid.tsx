import React, {useCallback, useEffect, useRef, useState} from 'react';
import {useDropzone} from 'react-dropzone';
import { useParams, useNavigate } from "react-router-dom";
import { useQueryClient } from '@tanstack/react-query';
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {ISharesAdapter} from "../../Adapters/ISharesAdapter";
import usePhotoGrid from "./usePhotoGrid";
import {PhotoCard} from "../PhotoCard/PhotoCard";
import {Slideshow} from "../Slideshow/Slideshow";
import {EmptyState} from "../EmptyState/EmptyState";
import {BulkActionsToolbar} from "../BulkActionsToolbar/BulkActionsToolbar";
import {KeyboardShortcutsHelp} from "../KeyboardShortcutsHelp/KeyboardShortcutsHelp";
import {TimelineScrubber} from "../TimelineScrubber/TimelineScrubber";
import {ActionIcon, Button, Card, Checkbox, Group, Loader, Modal, Paper, Progress, SegmentedControl, Skeleton, Stack, Table, Text, TextInput, Title, Tooltip} from "@mantine/core";
import {useHotkeys, useMediaQuery} from "@mantine/hooks";
import {IconLayoutGrid, IconList, IconPhotoOff, IconPlayerPlay, IconRefresh, IconSelect} from "@tabler/icons-react";
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import './PhotoGrid.css';
import {Photo} from "../../Models/Photo";
import {fetchThumbnailsBatch, getCachedThumbnail, fetchThumbnailWithAuth, revokeThumbnail} from "../../utils/ThumbnailUtils";

const getPhotoDisplayDate = (photo: Photo): string => {
    // Try top-level dateTaken first
    if (photo.dateTaken) {
        try { return new Date(photo.dateTaken).toLocaleDateString(); } catch { /* ignore */ }
    }

    // Recursively search metadata for common date fields
    const findDate = (obj: unknown): string | undefined => {
        if (!obj || typeof obj !== 'object') return undefined;
        const record = obj as Record<string, unknown>;
        for (const key of ['DateTimeOriginal', 'DateTime', 'dateTaken', 'CreateDate', 'ModifyDate']) {
            if (key in record && record[key]) {
                const val = String(record[key]);
                // EXIF dates are often in format "2006:01:02 15:04:05"
                if (/^\d{4}:\d{2}:\d{2}/.test(val)) {
                    return val.replace(/^\d{4}:\d{2}:\d{2}/, (m) => m.replace(/:/g, '-'));
                }
                return val;
            }
        }
        // Search nested objects
        for (const val of Object.values(record)) {
            if (val && typeof val === 'object') {
                const found = findDate(val);
                if (found) return found;
            }
        }
        return undefined;
    };

    const raw = findDate(photo.metadata) || photo.createdAt;
    if (!raw) return '';
    try {
        return new Date(raw).toLocaleDateString();
    } catch {
        return raw;
    }
};

const readUploadedFileAsText = (inputFile: File) => {
    const temporaryFileReader = new FileReader();

    return new Promise<string>((resolve, reject) => {
        temporaryFileReader.onerror = () => {
            temporaryFileReader.abort();
            reject(new DOMException('Problem parsing input file.'));
        };

        temporaryFileReader.onload = () => {
            resolve(temporaryFileReader.result as string);
        };
        temporaryFileReader.readAsDataURL(inputFile);
    });
};

const PhotoGridUploadDropzone: React.FunctionComponent<{ photosAdapter: IPhotosAdapter }> = ({ photosAdapter }) => {
    const [isUploading, setIsUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState({ current: 0, total: 0 });

    const handleFileUpload = useCallback(async (acceptedFiles: File[]) => {
        setIsUploading(true);
        setUploadProgress({ current: 0, total: acceptedFiles.length });
        let uploadedCount = 0;
        try {
            for (const file of acceptedFiles) {
                const fileContent = await readUploadedFileAsText(file);
                const photoContent = {
                    name: file.name,
                    albumName: 'General',
                    binaryContent: fileContent,
                };
                await photosAdapter.uploadPhoto(photoContent);
                uploadedCount++;
                setUploadProgress({ current: uploadedCount, total: acceptedFiles.length });
            }
            notifications.show({
                title: 'Upload complete',
                message: `${uploadedCount} photo${uploadedCount !== 1 ? 's' : ''} uploaded successfully`,
                color: 'green',
            });
        } catch (e) {
            const message = e instanceof Error ? e.message : 'Upload failed';
            notifications.show({
                title: 'Upload failed',
                message,
                color: 'red',
            });
        } finally {
            setIsUploading(false);
            setUploadProgress({ current: 0, total: 0 });
        }
    }, [photosAdapter]);

    const onDrop = useCallback((acceptedFiles: File[]) => {
        handleFileUpload(acceptedFiles);
    }, [handleFileUpload]);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { 'image/*': [] },
        disabled: isUploading,
    });

    return (
        <Card mb="md" p="md" withBorder>
            <Paper
                {...getRootProps()}
                p="xl"
                withBorder
                style={{
                    border: isDragActive ? '2px dashed var(--mantine-color-blue-6)' : '2px dashed var(--mantine-color-gray-4)',
                    borderRadius: 'var(--mantine-radius-md)',
                    textAlign: 'center',
                    cursor: isUploading ? 'default' : 'pointer',
                    background: isDragActive ? 'var(--mantine-color-blue-light)' : 'transparent',
                    transition: 'all 0.2s ease',
                }}
            >
                <input {...getInputProps()} />
                <Text size="lg" c={isDragActive ? 'blue' : 'dimmed'}>
                    {isDragActive ? 'Drop photos here...' : 'Drag photos here or click to upload'}
                </Text>
            </Paper>
            {isUploading && uploadProgress.total > 0 && (
                <Stack gap="xs" mt="md">
                    <Text size="sm">Uploading {uploadProgress.current} of {uploadProgress.total} photos...</Text>
                    <Progress value={(uploadProgress.current / uploadProgress.total) * 100} size="lg" />
                </Stack>
            )}
        </Card>
    );
};

interface IProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
    maxDisplayed?: number;
    tags?: string;
    mediaType?: string;
    searchQuery?: string;
    favoritesOnly?: boolean;
}

const defaultProps = {
    maxDisplayed: 20000000
}

interface GridImageItemProps {
    photo: Photo;
    isSelectionMode: boolean;
    isSelected: boolean;
    onImageClick: (id: string) => void;
    onToggleSelect: (id: string) => void;
    onRetry?: (id: string) => void;
}

// By adding a custom comparison function to React.memo, we prevent re-renders unless the photo's ID changes.
const GridImageItem = React.memo(
    ({photo, isSelectionMode, isSelected, onImageClick, onToggleSelect, onRetry, thumbnailUrl}: GridImageItemProps & { thumbnailUrl: string | undefined }) => {
        const [isLoaded, setIsLoaded] = useState(false);
        const [hasError, setHasError] = useState(false);
        const effectiveThumbnailUrl = thumbnailUrl || photo.thumbnailUrl;

        useEffect(() => {
            setIsLoaded(false);
            setHasError(false);
        }, [effectiveThumbnailUrl]);

        const handleClick = () => {
            if (isSelectionMode) {
                onToggleSelect(photo.id);
            } else {
                onImageClick(photo.id);
            }
        };

        const cameraModel = photo.metadata && typeof photo.metadata === 'object' && 'Model' in photo.metadata ? String(photo.metadata.Model) : undefined;

        return (
            <div className="item" onClick={handleClick} style={{ position: 'relative' }}>
                {isSelectionMode && (
                    <div style={{ position: 'absolute', top: 4, left: 4, zIndex: 2 }} onClick={(e) => { e.stopPropagation(); onToggleSelect(photo.id); }}>
                        <Checkbox checked={isSelected} onChange={() => {}} size="md" />
                    </div>
                )}
                <div className="thumbnail" style={{ aspectRatio: '4 / 3' }}>
                    {(!effectiveThumbnailUrl || !isLoaded) && !hasError && (
                        <Skeleton
                            height="100%"
                            width="100%"
                            style={{position: 'absolute', top: 0, left: 0}}
                        />
                    )}
                    {effectiveThumbnailUrl && !hasError && (
                        <img
                            src={effectiveThumbnailUrl}
                            alt={photo.name}
                            onLoad={() => setIsLoaded(true)}
                            onError={() => setHasError(true)}
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
                            borderRadius: 8,
                        }}>
                            <Tooltip label="Retry loading thumbnail">
                                <ActionIcon
                                    variant="light"
                                    size="lg"
                                    onClick={(e) => { e.stopPropagation(); onRetry?.(photo.id); }}
                                >
                                    <IconRefresh size={20} />
                                </ActionIcon>
                            </Tooltip>
                        </div>
                    )}
                    {photo.mediaType === 'video' && effectiveThumbnailUrl && !hasError && (
                        <>
                            <div style={{
                                position: 'absolute',
                                top: '50%',
                                left: '50%',
                                transform: 'translate(-50%, -50%)',
                                background: 'rgba(0,0,0,0.5)',
                                borderRadius: '50%',
                                width: 40,
                                height: 40,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                pointerEvents: 'none',
                            }}>
                                <IconPlayerPlay size={20} color="white" />
                            </div>
                            {photo.duration !== undefined && photo.duration > 0 && (
                                <div style={{
                                    position: 'absolute',
                                    bottom: 6,
                                    right: 6,
                                    background: 'rgba(0,0,0,0.7)',
                                    color: 'white',
                                    fontSize: 11,
                                    padding: '2px 6px',
                                    borderRadius: 4,
                                    pointerEvents: 'none',
                                }}>
                                    {Math.floor(photo.duration / 60)}:{String(photo.duration % 60).padStart(2, '0')}
                                </div>
                            )}
                        </>
                    )}
                    {!isSelectionMode && effectiveThumbnailUrl && !hasError && (
                        <div className="photo-hover-overlay" style={{
                            position: 'absolute',
                            bottom: 0,
                            left: 0,
                            right: 0,
                            padding: '8px 12px',
                            background: 'linear-gradient(to top, rgba(0,0,0,0.7), transparent)',
                            color: 'white',
                            fontSize: 12,
                            opacity: 0,
                            transition: 'opacity 0.2s',
                            pointerEvents: 'none',
                        }}>
                            <div style={{ fontWeight: 500, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                                {photo.name}
                            </div>
                            <div style={{ opacity: 0.8, fontSize: 11 }}>
                                {getPhotoDisplayDate(photo)}
                                {cameraModel ? ` · ${cameraModel}` : ''}
                            </div>
                        </div>
                    )}
                </div>
                <style>{`
                    .item:hover .photo-hover-overlay {
                        opacity: 1 !important;
                    }
                `}</style>
            </div>
        );
    },
    (prevProps, nextProps) =>
        prevProps.photo.id === nextProps.photo.id &&
        prevProps.isSelectionMode === nextProps.isSelectionMode &&
        prevProps.isSelected === nextProps.isSelected &&
        prevProps.thumbnailUrl === nextProps.thumbnailUrl
);
GridImageItem.displayName = 'GridImageItem';

export const PhotoGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const props = {...defaultProps, ...propsIn};
    const [isImageModalOpen, setImageModalOpen] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [showShortcutsHelp, setShowShortcutsHelp] = useState(false);
    const [isSelectionMode, setIsSelectionMode] = useState(false);
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
    const [isSlideshowOpen, setIsSlideshowOpen] = useState(false);
    const [density, setDensity] = useState<'compact' | 'comfortable' | 'large'>(() => {
        try {
            const saved = localStorage.getItem('photobox-grid-density');
            if (saved === 'compact' || saved === 'comfortable' || saved === 'large') return saved;
        } catch { /* ignore */ }
        return 'comfortable';
    });
    const [viewMode, setViewMode] = useState<'grid' | 'list'>(() => {
        try {
            const saved = localStorage.getItem('photobox-view-mode');
            if (saved === 'grid' || saved === 'list') return saved;
        } catch { /* ignore */ }
        return 'grid';
    });
    const [dateFilter, setDateFilter] = useState<{ year: number; month: number } | null>(null);
    const [tagModalOpen, setTagModalOpen] = useState(false);
    const [tagModalMode, setTagModalMode] = useState<'add' | 'remove'>('add');
    const [tagInput, setTagInput] = useState('');
    const [albumModalOpen, setAlbumModalOpen] = useState(false);
    const [availableAlbums, setAvailableAlbums] = useState<{ id: string; name: string }[]>([]);
    const [albumModalLoading, setAlbumModalLoading] = useState(false);
    const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());
    const isMobile = useMediaQuery('(max-width: 50em)');

    useEffect(() => {
        try {
            localStorage.setItem('photobox-grid-density', density);
        } catch { /* ignore */ }
    }, [density]);

    useEffect(() => {
        try {
            localStorage.setItem('photobox-view-mode', viewMode);
        } catch { /* ignore */ }
    }, [viewMode]);

    const startDate = dateFilter ? `${dateFilter.year}-${String(dateFilter.month).padStart(2, '0')}-01T00:00:00` : undefined;
    const endDate = dateFilter ? (() => {
        const lastDay = new Date(dateFilter.year, dateFilter.month, 0).getDate();
        return `${dateFilter.year}-${String(dateFilter.month).padStart(2, '0')}-${String(lastDay).padStart(2, '0')}T23:59:59`;
    })() : undefined;

    // The hook now provides a simple, flat, de-duplicated array of photos.
    const {photos, albumName, allRetrieved, fetchNextPage, isFetchingNextPage, isFetching, refetch} = usePhotoGrid(props.photosAdapter, props.albumsAdapter, id, startDate, endDate, props.tags, props.mediaType, props.searchQuery, props.favoritesOnly);

    const photoIdsKey = React.useMemo(() => photos.map(p => p.id).join(','), [photos]);

    // Merge React state with global cache so grid items always see cached thumbnails
    // even if the async batch fetch hasn't updated state yet (e.g. browser-cache hit).
    const mergedThumbnailUrls = React.useMemo(() => {
        const merged = new Map(thumbnailUrls);
        photos.forEach(photo => {
            if (!merged.has(photo.id)) {
                const cached = getCachedThumbnail(photo.id);
                if (cached) {
                    merged.set(photo.id, cached);
                }
            }
        });
        return merged;
    }, [thumbnailUrls, photoIdsKey]);

    // Batch load thumbnails for new photos
    useEffect(() => {
        if (photos.length === 0) return;

        const loadThumbnails = async () => {
            // 1. Sync any thumbnails that are already in the global cache but missing from React state
            const cachedOnly = new Map<string, string>();
            for (const photo of photos) {
                const cached = getCachedThumbnail(photo.id);
                if (cached) {
                    cachedOnly.set(photo.id, cached);
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

            // 2. Fetch thumbnails that aren't cached yet
            const uncachedIds = photos
                .map(p => p.id)
                .filter(id => !getCachedThumbnail(id));

            console.log('[PhotoGrid] Loading thumbnails for', uncachedIds.length, 'photos');

            if (uncachedIds.length === 0) return;

            const newThumbnails = await fetchThumbnailsBatch(props.photosAdapter, uncachedIds);
            console.log('[PhotoGrid] Received', newThumbnails.size, 'thumbnails');
            setThumbnailUrls(prev => {
                const next = new Map(prev);
                newThumbnails.forEach((url, id) => next.set(id, url));
                console.log('[PhotoGrid] Setting thumbnailUrls state, total:', next.size);
                return next;
            });
        };

        loadThumbnails();
    }, [photoIdsKey, props.photosAdapter]);

    const photosRef = useRef(photos);
    photosRef.current = photos;

    const appendDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    // Pull-to-refresh state
    const [pullStartY, setPullStartY] = useState(0);
    const [pullDistance, setPullDistance] = useState(0);
    const [isRefreshing, setIsRefreshing] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);
    const PULL_THRESHOLD = 80;

    useEffect(() => {
        return () => {
            if (appendDebounceRef.current) {
                clearTimeout(appendDebounceRef.current);
            }
        };
    }, []);

    useEffect(() => {
        // Clear selection when photos change or album changes
        setSelectedIds(new Set());
        setIsSelectionMode(false);
    }, [id]);

    const onRequestAppend = useCallback(() => {
        if (isFetchingNextPage || allRetrieved) {
            return;
        }
        if (appendDebounceRef.current) {
            return;
        }
        appendDebounceRef.current = setTimeout(() => {
            appendDebounceRef.current = null;
        }, 100);
        fetchNextPage();
    }, [isFetchingNextPage, allRetrieved, fetchNextPage]);

    const onImageClick = useCallback((photoId: string) => {
        const photoIndex = photosRef.current.findIndex(p => p.id === photoId);
        if (photoIndex !== -1) {
            setCurrentIndex(photoIndex);
            setImageModalOpen(true);
        }
    }, []);

    const handleTouchStart = useCallback((e: React.TouchEvent) => {
        if (containerRef.current && containerRef.current.scrollTop === 0) {
            setPullStartY(e.touches[0].clientY);
        }
    }, []);

    const handleTouchMove = useCallback((e: React.TouchEvent) => {
        if (pullStartY > 0 && containerRef.current && containerRef.current.scrollTop === 0) {
            const delta = e.touches[0].clientY - pullStartY;
            if (delta > 0) {
                setPullDistance(Math.min(delta * 0.5, 120));
            }
        }
    }, [pullStartY]);

    const handleTouchEnd = useCallback(() => {
        if (pullDistance > PULL_THRESHOLD && !isRefreshing) {
            setIsRefreshing(true);
            setPullDistance(0);
            refetch().then(() => {
                notifications.show({
                    title: 'Refreshed',
                    message: 'Photos refreshed successfully',
                    color: 'green',
                });
            }).catch(() => {
                notifications.show({
                    title: 'Refresh failed',
                    message: 'Failed to refresh photos',
                    color: 'red',
                });
            }).finally(() => {
                setIsRefreshing(false);
            });
        } else {
            setPullDistance(0);
        }
        setPullStartY(0);
    }, [pullDistance, isRefreshing, refetch]);

    const onToggleSelect = useCallback((photoId: string) => {
        setSelectedIds(prev => {
            const next = new Set(prev);
            if (next.has(photoId)) {
                next.delete(photoId);
            } else {
                next.add(photoId);
            }
            return next;
        });
    }, []);

    const handleRetryThumbnail = useCallback(async (photoId: string) => {
        revokeThumbnail(photoId);
        setThumbnailUrls(prev => {
            const next = new Map(prev);
            next.delete(photoId);
            return next;
        });
        try {
            const url = await fetchThumbnailWithAuth(props.photosAdapter, photoId);
            setThumbnailUrls(prev => {
                const next = new Map(prev);
                next.set(photoId, url);
                return next;
            });
        } catch {
            // Error state will be handled by the component's onError handler
        }
    }, [props.photosAdapter]);

    const handleSelectAll = useCallback(() => {
        setSelectedIds(new Set(photosRef.current.map(p => p.id)));
    }, []);

    const handleDeselectAll = useCallback(() => {
        setSelectedIds(new Set());
    }, []);

    const handleBulkFavorite = useCallback(async () => {
        try {
            for (const photoId of selectedIds) {
                await props.photosAdapter.favoritePhoto(photoId, true);
            }
            const selectedIdsSet = new Set(selectedIds);
            queryClient.setQueriesData(
                { queryKey: ['albumPhotos'] },
                (oldData: any) => {
                    if (!oldData) return oldData;
                    return {
                        ...oldData,
                        pages: oldData.pages.map((page: any) => ({
                            ...page,
                            data: page.data.map((p: any) =>
                                selectedIdsSet.has(p.id)
                                    ? { ...p, favorite: true }
                                    : p
                            ),
                        })),
                    };
                }
            );
            notifications.show({
                title: 'Favorited',
                message: `${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''} added to favorites`,
                color: 'pink',
            });
            setIsSelectionMode(false);
            setSelectedIds(new Set());
        } catch (e) {
            notifications.show({
                title: 'Failed to favorite',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, selectedIds, queryClient]);

    const handleBulkDelete = useCallback(() => {
        modals.openConfirmModal({
            title: 'Delete selected photos?',
            children: `This will move ${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''} to the trash.`,
            labels: { confirm: 'Delete', cancel: 'Cancel' },
            confirmProps: { color: 'red' },
            onConfirm: async () => {
                try {
                    for (const photoId of selectedIds) {
                        await props.photosAdapter.deletePhoto(photoId);
                    }
                    notifications.show({
                        title: 'Deleted',
                        message: `${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''} moved to trash`,
                        color: 'green',
                    });
                    setIsSelectionMode(false);
                    setSelectedIds(new Set());
                } catch (e) {
                    notifications.show({
                        title: 'Failed to delete',
                        message: e instanceof Error ? e.message : 'An error occurred',
                        color: 'red',
                    });
                }
            },
        });
    }, [props.photosAdapter, selectedIds]);

    const handleBulkDownload = useCallback(async () => {
        try {
            await props.photosAdapter.downloadPhotosAsZip(Array.from(selectedIds));
            notifications.show({
                title: 'Download started',
                message: `Preparing zip of ${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''}`,
                color: 'blue',
            });
            setIsSelectionMode(false);
            setSelectedIds(new Set());
        } catch (e) {
            notifications.show({
                title: 'Download failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, selectedIds]);

    const handleBulkTag = useCallback(async () => {
        const tags = tagInput.split(',').map(t => t.trim()).filter(Boolean);
        if (tags.length === 0) return;
        try {
            await props.photosAdapter.batchUpdatePhotoTags(
                Array.from(selectedIds),
                tags,
                tagModalMode
            );
            notifications.show({
                title: tagModalMode === 'add' ? 'Tags added' : 'Tags removed',
                message: `${tagModalMode === 'add' ? 'Added' : 'Removed'} tags for ${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''}`,
                color: 'green',
            });
            setTagModalOpen(false);
            setTagInput('');
            setIsSelectionMode(false);
            setSelectedIds(new Set());
            refetch();
        } catch (e) {
            notifications.show({
                title: 'Failed to update tags',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, selectedIds, tagInput, tagModalMode, refetch]);

    const handleBulkAddToAlbum = useCallback(async (albumId: string) => {
        try {
            await props.albumsAdapter.addPhotosToAlbum(albumId, Array.from(selectedIds));
            notifications.show({
                title: 'Added to album',
                message: `${selectedIds.size} photo${selectedIds.size !== 1 ? 's' : ''} added to album`,
                color: 'green',
            });
            setAlbumModalOpen(false);
            setIsSelectionMode(false);
            setSelectedIds(new Set());
            refetch();
        } catch (e) {
            notifications.show({
                title: 'Failed to add to album',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.albumsAdapter, selectedIds, refetch]);

    const handleOpenAlbumModal = useCallback(async () => {
        setAlbumModalLoading(true);
        try {
            const albums = await props.albumsAdapter.getAllAlbumsInfo();
            const albumList = Array.isArray(albums) ? albums : (albums?.data ?? []);
            setAvailableAlbums(albumList.map((a: { id: string; name: string }) => ({ id: a.id, name: a.name })));
            setAlbumModalOpen(true);
        } catch (e) {
            notifications.show({
                title: 'Failed to load albums',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setAlbumModalLoading(false);
        }
    }, [props.albumsAdapter]);

    const closeModal = useCallback(() => {
        setImageModalOpen(false);
    }, []);

    const handlePreviousPhoto = useCallback(() => {
        setCurrentIndex(prevIndex => prevIndex > 0 ? prevIndex - 1 : 0);
    }, []);

    const handleNextPhoto = useCallback(() => {
        setCurrentIndex(prevIndex => prevIndex < photosRef.current.length - 1 ? prevIndex + 1 : prevIndex);
    }, []);

    useHotkeys([
        ['?', () => setShowShortcutsHelp(prev => !prev)],
        ['Enter', () => {
            if (!isImageModalOpen && !isSelectionMode && photosRef.current.length > 0) {
                setCurrentIndex(0);
                setImageModalOpen(true);
            }
        }],
    ]);

    const AlbumTitle = () => (
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: 8 }}>
            <Title size="h4">{albumName}</Title>
            <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                {photos.length > 0 && !isMobile && viewMode === 'grid' && (
                    <SegmentedControl
                        size="xs"
                        value={density}
                        onChange={(value) => setDensity(value as 'compact' | 'comfortable' | 'large')}
                        data={[
                            { label: 'Compact', value: 'compact' },
                            { label: 'Comfortable', value: 'comfortable' },
                            { label: 'Large', value: 'large' },
                        ]}
                    />
                )}
                {photos.length > 0 && !isMobile && (
                    <ActionIcon
                        variant={viewMode === 'grid' ? 'filled' : 'light'}
                        size="sm"
                        onClick={() => setViewMode('grid')}
                        aria-label="Grid view"
                    >
                        <IconLayoutGrid size="1rem" />
                    </ActionIcon>
                )}
                {photos.length > 0 && !isMobile && (
                    <ActionIcon
                        variant={viewMode === 'list' ? 'filled' : 'light'}
                        size="sm"
                        onClick={() => setViewMode('list')}
                        aria-label="List view"
                    >
                        <IconList size="1rem" />
                    </ActionIcon>
                )}
                {photos.length > 0 && (
                    <Tooltip label="Select photos">
                        <ActionIcon variant="light" onClick={() => setIsSelectionMode(prev => !prev)} aria-label="Select photos">
                            <IconSelect size="1.25rem" />
                        </ActionIcon>
                    </Tooltip>
                )}
            </div>
        </div>
    );

    if (photos.length === 0 && !isFetching) {
        const isVideos = props.mediaType === 'video';
        return (
            <>
                <AlbumTitle/>
                <EmptyState
                    title={id ? "This album is empty" : isVideos ? "No videos yet" : "No photos yet"}
                    description={id ? "Upload photos to see them here." : isVideos ? "Your video library is empty. Upload videos or configure your photo directory." : "Your photo library is empty. Upload photos or configure your photo directory."}
                    icon={<IconPhotoOff size="2rem" />}
                    action={{
                        label: isVideos ? "Upload videos" : "Upload photos",
                        onClick: () => navigate('/album/new/General'),
                    }}
                />
            </>
        );
    }

    const showPullIndicator = pullDistance > 10;

    return (
        <div
            ref={containerRef}
            className={`density-${density}`}
            onTouchStart={handleTouchStart}
            onTouchMove={handleTouchMove}
            onTouchEnd={handleTouchEnd}
            style={{ overflowY: 'auto', height: '100%', position: 'relative' }}
        >
            {showPullIndicator && (
                <div style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    right: 0,
                    height: pullDistance,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    zIndex: 10,
                    transition: isRefreshing ? 'height 0.3s ease' : undefined,
                }}>
                    {isRefreshing ? (
                        <Loader size="sm" />
                    ) : (
                        <div style={{
                            transform: `rotate(${Math.min(pullDistance / PULL_THRESHOLD, 1) * 360}deg)`,
                            transition: 'transform 0.1s',
                        }}>
                            ↓
                        </div>
                    )}
                </div>
            )}
            <div style={{ transform: `translateY(${showPullIndicator ? pullDistance : 0}px)`, transition: isRefreshing ? 'transform 0.3s ease' : undefined }}>
                <AlbumTitle/>
                <PhotoGridUploadDropzone photosAdapter={props.photosAdapter} />
                {isSelectionMode && (
                    <BulkActionsToolbar
                        selectedCount={selectedIds.size}
                        totalCount={photos.length}
                        onSelectAll={handleSelectAll}
                        onDeselectAll={handleDeselectAll}
                        onFavorite={handleBulkFavorite}
                        onDelete={handleBulkDelete}
                        onDownload={handleBulkDownload}
                        onAddTag={() => { setTagModalMode('add'); setTagModalOpen(true); }}
                        onRemoveTag={() => { setTagModalMode('remove'); setTagModalOpen(true); }}
                        onAddToAlbum={handleOpenAlbumModal}
                        onCancel={() => { setIsSelectionMode(false); setSelectedIds(new Set()); }}
                    />
                )}
                {viewMode === 'grid' ? (
                    <JustifiedInfiniteGrid
                        key={`grid-${density}`}
                        placeholder={<Skeleton height={7} mt={6} radius="md"/>}
                        className="container"
                        gap={density === 'compact' ? 4 : density === 'large' ? 20 : 10}
                        stretch={true}
                        passUnstretchRow={true}
                        onRequestAppend={onRequestAppend}
                        threshold={800}
                        useRecycle={false}
                        preserveUIOnDestroy={true}
                    >
                        {photos.map((photo: Photo, index: number) => (
                            <GridImageItem
                                data-grid-groupkey={Math.floor(index / 30)}
                                key={photo.id}
                                photo={photo}
                                thumbnailUrl={mergedThumbnailUrls.get(photo.id)}
                                isSelectionMode={isSelectionMode}
                                isSelected={selectedIds.has(photo.id)}
                                onImageClick={onImageClick}
                                onToggleSelect={onToggleSelect}
                                onRetry={handleRetryThumbnail}
                            />
                        ))}
                    </JustifiedInfiniteGrid>
                ) : (
                    <div style={{ padding: '0 8px' }}>
                        <Table striped highlightOnHover>
                            <Table.Thead>
                                <Table.Tr>
                                    {isSelectionMode && <Table.Th style={{ width: 40 }}></Table.Th>}
                                    <Table.Th style={{ width: 60 }}>Preview</Table.Th>
                                    <Table.Th>Name</Table.Th>
                                    <Table.Th>Date</Table.Th>
                                    <Table.Th>Tags</Table.Th>
                                    <Table.Th style={{ width: 80 }}>Favorite</Table.Th>
                                </Table.Tr>
                            </Table.Thead>
                            <Table.Tbody>
                                {photos.map((photo: Photo) => (
                                    <Table.Tr
                                        key={photo.id}
                                        onClick={() => {
                                            if (isSelectionMode) {
                                                onToggleSelect(photo.id);
                                            } else {
                                                onImageClick(photo.id);
                                            }
                                        }}
                                        style={{ cursor: 'pointer' }}
                                    >
                                        {isSelectionMode && (
                                            <Table.Td onClick={(e) => { e.stopPropagation(); onToggleSelect(photo.id); }}>
                                                <Checkbox checked={selectedIds.has(photo.id)} onChange={() => {}} size="sm" />
                                            </Table.Td>
                                        )}
                                        <Table.Td>
                                            <img
                                                src={mergedThumbnailUrls.get(photo.id) || photo.thumbnailUrl || ''}
                                                alt={photo.name}
                                                style={{ width: 40, height: 40, objectFit: 'cover', borderRadius: 4 }}
                                            />
                                        </Table.Td>
                                        <Table.Td>{photo.name}</Table.Td>
                                        <Table.Td>{getPhotoDisplayDate(photo)}</Table.Td>
                                        <Table.Td>
                                            <Group gap={4}>
                                                {photo.tags && photo.tags.split(',').map(t => t.trim()).filter(Boolean).map(tag => (
                                                    <span key={tag} style={{ fontSize: 11, padding: '2px 6px', background: 'var(--mantine-color-blue-light)', borderRadius: 4, color: 'var(--mantine-color-blue-light-color)' }}>{tag}</span>
                                                ))}
                                            </Group>
                                        </Table.Td>
                                        <Table.Td>{photo.favorite ? '★' : ''}</Table.Td>
                                    </Table.Tr>
                                ))}
                            </Table.Tbody>
                        </Table>
                        {!allRetrieved && (
                            <div style={{ textAlign: 'center', padding: 16 }}>
                                <Button variant="light" onClick={onRequestAppend} loading={isFetchingNextPage}>
                                    Load more
                                </Button>
                            </div>
                        )}
                    </div>
                )}
            </div>
            {isImageModalOpen && (
                <Modal
                    opened={isImageModalOpen}
                    withCloseButton={false}
                    aria-labelledby="customized-dialog-title"
                    size="auto"
                    padding={"0"}
                    m={"0"}
                    overlayProps={{backgroundOpacity: 0.55}}
                    fullScreen={isMobile}
                    transitionProps={{transition: 'fade', duration: 200}}
                    onClose={closeModal}>
                    <PhotoCard
                        photosAdapter={props.photosAdapter}
                        albumsAdapter={props.albumsAdapter}
                        sharesAdapter={props.sharesAdapter}
                        source={photos[currentIndex]}
                        previousPhoto={handlePreviousPhoto}
                        nextPhoto={handleNextPhoto}
                        firstInAlbum={currentIndex === 0}
                        lastInAlbum={currentIndex === photos.length - 1}
                        closeModal={closeModal}
                        onSlideshow={() => { setImageModalOpen(false); setIsSlideshowOpen(true); }}
                    />
                </Modal>
            )}
            {!id && (
                <TimelineScrubber
                    photosAdapter={props.photosAdapter}
                    onSelectMonth={(year, month) => setDateFilter({ year, month })}
                    onClear={() => setDateFilter(null)}
                    activeYear={dateFilter?.year}
                    activeMonth={dateFilter?.month}
                />
            )}
            <KeyboardShortcutsHelp opened={showShortcutsHelp} onClose={() => setShowShortcutsHelp(false)} />
            {isSlideshowOpen && (
                <Slideshow
                    photos={photos}
                    startIndex={currentIndex}
                    photosAdapter={props.photosAdapter}
                    onClose={() => setIsSlideshowOpen(false)}
                />
            )}
            <Modal
                opened={tagModalOpen}
                onClose={() => setTagModalOpen(false)}
                title={tagModalMode === 'add' ? 'Add tags to selected photos' : 'Remove tags from selected photos'}
                size="sm"
            >
                <TextInput
                    placeholder="Enter tags separated by commas"
                    value={tagInput}
                    onChange={(e) => setTagInput(e.currentTarget.value)}
                    onKeyDown={(e) => { if (e.key === 'Enter') handleBulkTag(); }}
                    autoFocus
                />
                <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, marginTop: 16 }}>
                    <Button variant="default" onClick={() => setTagModalOpen(false)}>Cancel</Button>
                    <Button onClick={handleBulkTag} color={tagModalMode === 'add' ? 'green' : 'orange'}>
                        {tagModalMode === 'add' ? 'Add tags' : 'Remove tags'}
                    </Button>
                </div>
            </Modal>
            <Modal
                opened={albumModalOpen}
                onClose={() => setAlbumModalOpen(false)}
                title="Add to album"
                size="sm"
            >
                {albumModalLoading ? (
                    <div style={{ textAlign: 'center', padding: 16 }}>
                        <Loader size="sm" />
                    </div>
                ) : availableAlbums.length === 0 ? (
                    <Text size="sm" c="dimmed">No albums available. Create an album first.</Text>
                ) : (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                        {availableAlbums.map(album => (
                            <Button
                                key={album.id}
                                variant="light"
                                fullWidth
                                onClick={() => handleBulkAddToAlbum(album.id)}
                            >
                                {album.name}
                            </Button>
                        ))}
                    </div>
                )}
            </Modal>
        </div>
    );
};
