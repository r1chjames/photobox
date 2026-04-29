import React, {useCallback, useEffect, useRef, useState} from 'react';
import {Photo} from '../../Models/Photo';
import {ActionIcon, Badge, Button, Card, Flex, Group, Image, Overlay, Stack} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {IconArrowLeftDashed, IconArrowRightDashed, IconCalendar, IconFolder, IconX} from "@tabler/icons-react";
import {useHotkeys} from "@mantine/hooks";
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
}

export const PhotoCard: React.FunctionComponent<IProps> = (props) => {
    const navigate = useNavigate()
    const [fetchedImage, setFetchedImage] = useState<string | undefined>();
    const [albumName, setAlbumName] = useState<string | undefined>();
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

    useHotkeys([
        ['ArrowLeft', () => props.previousPhoto()],
        ['ArrowRight', () => props.nextPhoto()],
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
                <ActionIcon color="dark" size="xl" onClick={handlePrevious}>
                    <IconArrowLeftDashed size="2.125rem"/>
                </ActionIcon>
            );
        }
        return null;
    }, [props.firstInAlbum, handlePrevious]);

    const nextButton = useCallback(() => {
        if (!props.lastInAlbum) {
            return (
                <ActionIcon color="dark" size="xl" onClick={handleNext}>
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
                    <Flex direction="row" style={{width: "100%", justifyContent: "right"}}>
                        <ActionIcon color="dark" size="l" opacity={1} onClick={handleClose}>
                            <IconX size="1.75rem"/>
                        </ActionIcon>
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
