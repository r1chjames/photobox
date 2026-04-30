import React, {useCallback, useEffect, useRef, useState} from 'react';
import {Photo} from '../../Models/Photo';
import {ActionIcon, Badge, Button, Card, Flex, Group, Image, Overlay, Stack, Tooltip} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {IconArrowLeftDashed, IconArrowRightDashed, IconCalendar, IconDownload, IconFolder, IconHeart, IconHeartFilled, IconPlayerPlay, IconX} from "@tabler/icons-react";
import {useHotkeys} from "@mantine/hooks";
import {notifications} from '@mantine/notifications';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";

interface IProps {
    photosAdapter: IPhotosAdapter;
    albumsAdapter: IAlbumsAdapter;
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
    const [albumName, setAlbumName] = useState<string | undefined>();
    const [isFavorite, setIsFavorite] = useState(props.source.favorite ?? false);
    const img: React.Ref<HTMLImageElement> = React.createRef();
    const blobUrlRef = useRef<string | undefined>(undefined);

    const fetchImage = useCallback(async () => {
        const imageUrl = await fetchPhotoBinWithAuth(props.photosAdapter, props.source.id);
        blobUrlRef.current = imageUrl;
        setFetchedImage(imageUrl);
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
        void fetchImage();
        void fetchAlbumName();
        return () => {
            revokeBlobUrl(blobUrlRef.current);
            blobUrlRef.current = undefined;
        };
    }, [props.source, fetchImage, fetchAlbumName]);

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

    useHotkeys([
        ['ArrowLeft', () => props.previousPhoto()],
        ['ArrowRight', () => props.nextPhoto()],
        ['Escape', () => props.closeModal()],
        ['d', () => handleDownload()],
        ['f', () => handleFavorite()],
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

    return (
        <Card shadow="sm" radius="md" padding={0}
              onClick={() => navigate(`/photo/${props.source.id}`)}>
            <Card.Section inheritPadding={false} withBorder={false}>
                {fetchedImage ?
                    <Image
                        h={"500px"}
                        fit={"cover"}
                        w={"auto"}
                        radius={0}
                        ref={img}
                        src={fetchedImage}
                        alt={props.source.name}
                    /> :
                    'Loading...'}
                <Overlay color="#000" backgroundOpacity={0} opacity={0.5}>
                    <Flex direction="row" style={{width: "100%", justifyContent: "right"}} gap="xs">
                        {props.onSlideshow && (
                            <Tooltip label="Slideshow">
                                <ActionIcon color="dark" size="l" opacity={1} onClick={(e) => { e.stopPropagation(); props.onSlideshow!(); }} aria-label="Slideshow">
                                    <IconPlayerPlay size="1.75rem"/>
                                </ActionIcon>
                            </Tooltip>
                        )}
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
    );
};
