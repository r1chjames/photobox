import React, {useCallback, useEffect, useRef, useState} from 'react';
import { useParams, useNavigate } from "react-router-dom";
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import usePhotoGrid from "./usePhotoGrid";
import {PhotoCard} from "../PhotoCard/PhotoCard";
import {Slideshow} from "../Slideshow/Slideshow";
import {EmptyState} from "../EmptyState/EmptyState";
import {BulkActionsToolbar} from "../BulkActionsToolbar/BulkActionsToolbar";
import {KeyboardShortcutsHelp} from "../KeyboardShortcutsHelp/KeyboardShortcutsHelp";
import {ActionIcon, Checkbox, Modal, Skeleton, Title, Tooltip} from "@mantine/core";
import {useHotkeys, useMediaQuery} from "@mantine/hooks";
import {IconPhotoOff, IconSelect} from "@tabler/icons-react";
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import './PhotoGrid.css';
import {Photo} from "../../Models/Photo";

interface IProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    maxDisplayed?: number;
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
}

// By adding a custom comparison function to React.memo, we prevent re-renders unless the photo's ID changes.
const GridImageItem = React.memo(
    ({photo, isSelectionMode, isSelected, onImageClick, onToggleSelect}: GridImageItemProps) => {
        const [isLoaded, setIsLoaded] = useState(false);
        const imgSrc = photo.thumbnail && (photo.thumbnail.startsWith('http') || photo.thumbnail.startsWith('data:image'))
            ? photo.thumbnail
            : `data:image/png;base64,${photo.thumbnail}`;

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
                <div className="thumbnail">
                    {!isLoaded && (
                        <Skeleton
                            height="100%"
                            width="100%"
                            style={{position: 'absolute', top: 0, left: 0}}
                        />
                    )}
                    <img
                        src={imgSrc}
                        alt={photo.name}
                        onLoad={() => setIsLoaded(true)}
                        style={{opacity: isLoaded ? 1 : 0, transition: 'opacity 0.2s'}}
                    />
                    {!isSelectionMode && isLoaded && (
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
                                {photo.createdAt ? new Date(photo.createdAt).toLocaleDateString() : ''}
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
        prevProps.isSelected === nextProps.isSelected
);
GridImageItem.displayName = 'GridImageItem';

export const PhotoGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const props = {...defaultProps, ...propsIn};
    const [isImageModalOpen, setImageModalOpen] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const [showShortcutsHelp, setShowShortcutsHelp] = useState(false);
    const [isSelectionMode, setIsSelectionMode] = useState(false);
    const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
    const [isSlideshowOpen, setIsSlideshowOpen] = useState(false);
    const isMobile = useMediaQuery('(max-width: 50em)');

    // The hook now provides a simple, flat, de-duplicated array of photos.
    const {photos, albumName, allRetrieved, fetchNextPage, isFetchingNextPage} = usePhotoGrid(props.photosAdapter, props.albumsAdapter, id);

    const photosRef = useRef(photos);
    photosRef.current = photos;

    const appendDebounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

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
    }, [props.photosAdapter, selectedIds]);

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
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Title size="h4">{albumName}</Title>
            {photos.length > 0 && (
                <Tooltip label="Select photos">
                    <ActionIcon variant="light" onClick={() => setIsSelectionMode(prev => !prev)} aria-label="Select photos">
                        <IconSelect size="1.25rem" />
                    </ActionIcon>
                </Tooltip>
            )}
        </div>
    );

    if (photos.length === 0) {
        return (
            <>
                <AlbumTitle/>
                <EmptyState
                    title={id ? "This album is empty" : "No photos yet"}
                    description={id ? "Upload photos to see them here." : "Your photo library is empty. Upload photos or configure your photo directory."}
                    icon={<IconPhotoOff size="2rem" />}
                    action={{
                        label: "Upload photos",
                        onClick: () => navigate('/album/new/General'),
                    }}
                />
            </>
        );
    }

    return (
        <div>
            <AlbumTitle/>
            {isSelectionMode && (
                <BulkActionsToolbar
                    selectedCount={selectedIds.size}
                    totalCount={photos.length}
                    onSelectAll={handleSelectAll}
                    onDeselectAll={handleDeselectAll}
                    onFavorite={handleBulkFavorite}
                    onDelete={handleBulkDelete}
                    onDownload={handleBulkDownload}
                    onCancel={() => { setIsSelectionMode(false); setSelectedIds(new Set()); }}
                />
            )}
            <JustifiedInfiniteGrid
                placeholder={<Skeleton height={7} mt={6} radius="md"/>}
                className="container"
                gap={10}
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
                        isSelectionMode={isSelectionMode}
                        isSelected={selectedIds.has(photo.id)}
                        onImageClick={onImageClick}
                        onToggleSelect={onToggleSelect}
                    />
                ))}
            </JustifiedInfiniteGrid>
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
            <KeyboardShortcutsHelp opened={showShortcutsHelp} onClose={() => setShowShortcutsHelp(false)} />
            {isSlideshowOpen && (
                <Slideshow
                    photos={photos}
                    startIndex={currentIndex}
                    photosAdapter={props.photosAdapter}
                    onClose={() => setIsSlideshowOpen(false)}
                />
            )}
        </div>
    );
};
