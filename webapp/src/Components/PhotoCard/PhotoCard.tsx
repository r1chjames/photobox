import React, {useCallback, useEffect, useRef, useState} from 'react';
import {Photo} from '../../Models/Photo';
import {ActionIcon, Badge, Button, Card, Flex, Group, Image, Overlay, Stack, Tooltip} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {IconArrowLeftDashed, IconArrowRightDashed, IconCalendar, IconDownload, IconFolder, IconHeart, IconHeartFilled, IconMaximize, IconMinimize, IconPlayerPlay, IconRotateClockwise, IconShare2, IconX} from "@tabler/icons-react";
import {useHotkeys} from "@mantine/hooks";
import {notifications} from '@mantine/notifications';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {ISharesAdapter} from "../../Adapters/ISharesAdapter";
import {ShareModal} from "../ShareModal/ShareModal";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";
import {fetchThumbnailWithAuth, revokeThumbnail} from "../../utils/ThumbnailUtils";

interface IProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
    sharesAdapter?: ISharesAdapter;
    source: Photo;
    previousPhoto: () => void;
    nextPhoto: () => void;
    firstInAlbum: boolean;
    lastInAlbum: boolean;
    closeModal: () => void;
    onSlideshow?: () => void;
}

export const PhotoCard: React.FunctionComponent<IProps> = (props) => {
    const navigate = useNavigate()
    const [fetchedImage, setFetchedImage] = useState<string | undefined>();
    const [thumbnailUrl, setThumbnailUrl] = useState<string | undefined>();
    const [albumName, setAlbumName] = useState<string | undefined>();
    const [isFavorite, setIsFavorite] = useState(props.source.favorite ?? false);
    const [showShareModal, setShowShareModal] = useState(false);
    const [isFullscreen, setIsFullscreen] = useState(false);
    const img: React.Ref<HTMLImageElement> = React.createRef();
    const cardRef = useRef<HTMLDivElement>(null);
    const blobUrlRef = useRef<string | undefined>(undefined);
    const touchStartRef = useRef<{ x: number; y: number } | null>(null);
    const isSwipingRef = useRef(false);

    const fetchImage = useCallback(async () => {
        const imageUrl = await fetchPhotoBinWithAuth(props.photosAdapter, props.source.id);
        blobUrlRef.current = imageUrl;
        setFetchedImage(imageUrl);
    }, [props.photosAdapter, props.source.id]);

    const fetchThumbnail = useCallback(async () => {
        const url = await fetchThumbnailWithAuth(props.photosAdapter, props.source.id);
        setThumbnailUrl(url);
    }, [props.photosAdapter, props.source.id]);

    const fetchAlbumName = useCallback(async () => {
        if (props.source.albumId) {
            try {
                const album = await props.albumsAdapter.getAlbumInfoById(props.source.albumId);
                setAlbumName(album?.name);
            } catch (error) {
                console.error('Failed to fetch album name:', error);
            }
        }
    }, [props.albumsAdapter, props.source.albumId]);

    useEffect(() => {
        setFetchedImage(undefined);
        setThumbnailUrl(undefined);
        void fetchImage();
        void fetchThumbnail();
        void fetchAlbumName();
        return () => {
            revokeBlobUrl(blobUrlRef.current);
            blobUrlRef.current = undefined;
            revokeThumbnail(props.source.id);
        };
    }, [props.source, fetchImage, fetchThumbnail, fetchAlbumName]);

    const handleDownload = useCallback(async () => {
        try {
            await props.photosAdapter.downloadPhoto(props.source.id, props.source.name);
            notifications.show({
                title: 'Download started',
                message: `Downloading ${props.source.name}`,
                color: 'blue',
            });
        } catch (e) {
            notifications.show({
                title: 'Download failed',
                message: e instanceof Error ? e.message : 'Failed to download photo',
                color: 'red',
            });
        }
    }, [props.photosAdapter, props.source.id, props.source.name]);

    const handleFavorite = useCallback(async () => {
        const newFavorite = !isFavorite;
        try {
            await props.photosAdapter.favoritePhoto(props.source.id, newFavorite);
            setIsFavorite(newFavorite);
            notifications.show({
                title: newFavorite ? 'Added to favorites' : 'Removed from favorites',
                message: newFavorite ? 'Photo added to your favorites' : 'Photo removed from your favorites',
                color: 'pink',
            });
        } catch (e) {
            notifications.show({
                title: 'Failed to update favorite',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, props.source.id, isFavorite]);

    const handleRotate = useCallback(async (direction: 'cw' | 'ccw') => {
        try {
            await props.photosAdapter.rotatePhoto(props.source.id, direction);
            notifications.show({
                title: 'Photo rotated',
                message: `Rotated ${direction === 'cw' ? 'clockwise' : 'counter-clockwise'}`,
                color: 'green',
            });
            // Refresh image by re-fetching
            setFetchedImage(undefined);
            void fetchImage();
        } catch (e) {
            notifications.show({
                title: 'Rotation failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, props.source.id, fetchImage]);

    const toggleFullscreen = useCallback(async () => {
        try {
            if (!document.fullscreenElement) {
                await cardRef.current?.requestFullscreen();
                setIsFullscreen(true);
            } else {
                await document.exitFullscreen();
                setIsFullscreen(false);
            }
        } catch (e) {
            console.error('Fullscreen error:', e);
        }
    }, []);

    const handleNativeShare = useCallback(async () => {
        const shareUrl = `${window.location.origin}/photo/${props.source.id}`;
        const shareData = {
            title: props.source.name,
            text: `Check out this photo: ${props.source.name}`,
            url: shareUrl,
        };
        if (navigator.share) {
            try {
                await navigator.share(shareData);
                notifications.show({ title: 'Shared', message: 'Photo shared successfully', color: 'green' });
            } catch (e) {
                if ((e as Error).name !== 'AbortError') {
                    notifications.show({ title: 'Share failed', message: (e as Error).message, color: 'red' });
                }
            }
        } else {
            try {
                await navigator.clipboard.writeText(shareUrl);
                notifications.show({ title: 'Link copied', message: 'Photo link copied to clipboard', color: 'blue' });
            } catch (e) {
                notifications.show({ title: 'Copy failed', message: 'Could not copy link to clipboard', color: 'red' });
            }
        }
    }, [props.source.id, props.source.name]);

    useEffect(() => {
        const handler = () => setIsFullscreen(!!document.fullscreenElement);
        document.addEventListener('fullscreenchange', handler);
        return () => document.removeEventListener('fullscreenchange', handler);
    }, []);

    useHotkeys([
        ['ArrowLeft', () => props.previousPhoto()],
        ['ArrowRight', () => props.nextPhoto()],
        ['Escape', () => props.closeModal()],
        ['d', () => handleDownload()],
        ['f', () => handleFavorite()],
        ['r', () => handleRotate('cw')],
        ['s', () => handleNativeShare()],
        ['F11', (e) => { e.preventDefault(); toggleFullscreen(); }],
    ]);

    const handleClose = useCallback((e: React.MouseEvent) => {
        e.stopPropagation();
        props.closeModal();
    }, [props.closeModal]);

    const handlePrevious = useCallback((e: React.MouseEvent) => {
        e.stopPropagation();
        props.previousPhoto();
    }, [props.previousPhoto]);

    const handleNext = useCallback((e: React.MouseEvent) => {
        e.stopPropagation();
        props.nextPhoto();
    }, [props.nextPhoto]);

    const previousButton = useCallback(() => {
        if (!props.firstInAlbum) {
            return (
                <ActionIcon color="dark" size="xl" onClick={handlePrevious} aria-label="Previous photo">
                    <IconArrowLeftDashed size="2.125rem"/>
                </ActionIcon>
            );
        }
        return null;
    }, [props.firstInAlbum, handlePrevious]);

    const nextButton = useCallback(() => {
        if (!props.lastInAlbum) {
            return (
                <ActionIcon color="dark" size="xl" onClick={handleNext} aria-label="Next photo">
                    <IconArrowRightDashed size="2.125rem"/>
                </ActionIcon>
            );
        }
        return null;
    }, [props.lastInAlbum, handleNext]);

    const handleTouchStart = useCallback((e: React.TouchEvent) => {
        touchStartRef.current = { x: e.touches[0].clientX, y: e.touches[0].clientY };
        isSwipingRef.current = false;
    }, []);

    const handleTouchMove = useCallback((e: React.TouchEvent) => {
        if (!touchStartRef.current) return;
        const dx = Math.abs(e.touches[0].clientX - touchStartRef.current.x);
        if (dx > 10) {
            isSwipingRef.current = true;
        }
    }, []);

    const handleTouchEnd = useCallback((e: React.TouchEvent) => {
        if (!touchStartRef.current) return;
        const dx = e.changedTouches[0].clientX - touchStartRef.current.x;
        const dy = e.changedTouches[0].clientY - touchStartRef.current.y;
        const absDx = Math.abs(dx);
        const absDy = Math.abs(dy);
        if (absDx > absDy && absDx > 50) {
            if (dx > 0 && !props.firstInAlbum) {
                props.previousPhoto();
            } else if (dx < 0 && !props.lastInAlbum) {
                props.nextPhoto();
            }
        } else if (absDy > absDx && dy > 80) {
            props.closeModal();
        }
        touchStartRef.current = null;
    }, [props.firstInAlbum, props.lastInAlbum, props.previousPhoto, props.nextPhoto, props.closeModal]);

    const handleCardClick = useCallback(() => {
        if (isSwipingRef.current) {
            isSwipingRef.current = false;
            return;
        }
        navigate(`/photo/${props.source.id}`);
    }, [navigate, props.source.id]);

    return (
        <>
            <Card shadow="sm" radius="md" padding={0}
                  ref={cardRef}
                  onClick={handleCardClick}
                  onTouchStart={handleTouchStart}
                  onTouchMove={handleTouchMove}
                  onTouchEnd={handleTouchEnd}>
                <Card.Section inheritPadding={false} withBorder={false}>
                    {fetchedImage ?
                        (props.source.mediaType === 'video' ? (
                            <video
                                ref={img as React.Ref<HTMLVideoElement>}
                                src={fetchedImage}
                                controls
                                style={{ height: '500px', width: 'auto', objectFit: 'cover' }}
                            />
                        ) : (
                            <Image
                                h={"500px"}
                                fit={"cover"}
                                w={"auto"}
                                radius={0}
                                ref={img}
                                src={fetchedImage}
                                alt={props.source.name}
                            />
                        )) :
                        <Image
                            h="500px"
                            fit="cover"
                            w="auto"
                            radius={0}
                            src={thumbnailUrl}
                            alt={props.source.name}
                            style={{ filter: 'blur(10px) brightness(0.8)', transition: 'filter 0.3s ease' }}
                        />}
                    <Overlay color="#000" backgroundOpacity={0} opacity={0.5}>
                        <Flex direction="row" style={{width: "100%", justifyContent: "right"}} gap="xs">
                            {props.onSlideshow && (
                                <Tooltip label="Slideshow">
                                    <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); props.onSlideshow!(); }} aria-label="Slideshow">
                                        <IconPlayerPlay size="1.75rem"/>
                                    </ActionIcon>
                                </Tooltip>
                            )}
                            <Tooltip label="Native share (S)">
                                <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); handleNativeShare(); }} aria-label="Native share">
                                    <IconShare2 size="1.75rem"/>
                                </ActionIcon>
                            </Tooltip>
                            {props.sharesAdapter && (
                                <Tooltip label="Create share link">
                                    <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); setShowShareModal(true); }} aria-label="Create share link">
                                        <IconShare2 size="1.75rem"/>
                                    </ActionIcon>
                                </Tooltip>
                            )}
                            <Tooltip label={isFullscreen ? 'Exit fullscreen (F11)' : 'Enter fullscreen (F11)'}>
                                <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); toggleFullscreen(); }} aria-label="Toggle fullscreen">
                                    {isFullscreen ? <IconMinimize size="1.75rem" /> : <IconMaximize size="1.75rem" />}
                                </ActionIcon>
                            </Tooltip>
                            <Tooltip label="Rotate clockwise (R)">
                                <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); handleRotate('cw'); }} aria-label="Rotate clockwise">
                                    <IconRotateClockwise size="1.75rem"/>
                                </ActionIcon>
                            </Tooltip>
                            <Tooltip label={isFavorite ? 'Remove from favorites (F)' : 'Add to favorites (F)'}>
                                <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); handleFavorite(); }} aria-label="Toggle favorite">
                                    {isFavorite ? <IconHeartFilled size="1.75rem" color="var(--mantine-color-pink-filled)" /> : <IconHeart size="1.75rem" />}
                                </ActionIcon>
                            </Tooltip>
                            <Tooltip label="Download (D)">
                                <ActionIcon color="dark" size="l" opacity={1} onClick={handleDownload} aria-label="Download">
                                    <IconDownload size="1.75rem"/>
                                </ActionIcon>
                            </Tooltip>
                            <Tooltip label="Close (Esc)">
                                <ActionIcon color="dark" size="l" opacity={1} onClick={handleClose} aria-label="Close">
                                    <IconX size="1.75rem"/>
                                </ActionIcon>
                            </Tooltip>
                        </Flex>
                        <Flex direction="row" style={{
                            width: "100%",
                            height: "100%",
                            justifyContent: "space-between",
                            alignItems: "center"
                        }}>
                            {previousButton()}
                            {nextButton()}
                        </Flex>
                    </Overlay>
                </Card.Section>
                <Stack gap="xs" mt="md" mb="xs" px="xs" style={{width: '100%', minWidth: 0}}>
                    <Group gap="xs" style={{flexWrap: 'wrap'}}>
                        {props.source.createdAt && (
                            <Badge leftSection={<IconCalendar size={14} />} variant="light" color="blue">
                                {new Date(props.source.createdAt).toLocaleString()}
                            </Badge>
                        )}
                        {albumName && (
                            <Badge leftSection={<IconFolder size={14} />} variant="light" color="teal">
                                {albumName}
                            </Badge>
                        )}
                    </Group>
                    <Button onClick={(e: React.MouseEvent) => {
                        e.stopPropagation();
                        navigate(`/photo/${props.source.id}`);
                    }} fullWidth>
                        View Details
                    </Button>
                </Stack>
            </Card>
            {props.sharesAdapter && (
                <ShareModal
                    opened={showShareModal}
                    onClose={() => setShowShareModal(false)}
                    resourceType="photo"
                    resourceId={props.source.id}
                    resourceName={props.source.name}
                    sharesAdapter={props.sharesAdapter}
                />
            )}
        </>
    );
};
