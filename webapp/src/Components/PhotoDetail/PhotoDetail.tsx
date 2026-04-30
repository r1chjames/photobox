import React, {useCallback, useEffect, useState} from 'react';
import {useDisclosure} from '@mantine/hooks';
import {Button, Dialog, Drawer, Group, Image, Loader, ScrollArea, Table} from '@mantine/core';
import {useMediaQuery} from '@mantine/hooks';
import {IconDownload, IconHeart, IconHeartFilled, IconRotateClockwise, IconShare2} from '@tabler/icons-react';
import {valueType} from "../../utils/TypeUtils";
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {ISharesAdapter} from '../../Adapters/ISharesAdapter';
import {ShareModal} from '../ShareModal/ShareModal';
import {useParams} from "react-router-dom";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";
import {useQuery} from "@tanstack/react-query";
import {notifications} from '@mantine/notifications';

interface IProps {
    photosAdapter: IPhotosAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const {id} = useParams();
    const [opened, {toggle, close}] = useDisclosure(true);
    const [isFavorite, setIsFavorite] = useState(false);
    const [showShareModal, setShowShareModal] = useState(false);
    const isMobile = useMediaQuery('(max-width: 50em)');

    const fetchPhoto = async () => {
        return props.photosAdapter.getPhotoInfoById(id!);
    }

    const fetchPhotoBin = async () => {
        return fetchPhotoBinWithAuth(props.photosAdapter, id!);
    }

    const {data: photo} = useQuery({
        queryKey: ['fetchPhoto', id],
        queryFn: fetchPhoto,
        enabled: !!id,
    });

    const {data: photoUrl} = useQuery({
        queryKey: ['fetchPhotoBin', id],
        queryFn: fetchPhotoBin,
        enabled: !!id,
    });

    useEffect(() => {
        if (photo) {
            setIsFavorite(photo.favorite ?? false);
        }
    }, [photo]);

    const handleDownload = useCallback(async () => {
        if (!photo) return;
        try {
            await props.photosAdapter.downloadPhoto(photo.id, photo.name);
            notifications.show({
                title: 'Download started',
                message: `Downloading ${photo.name}`,
                color: 'blue',
            });
        } catch (e) {
            notifications.show({
                title: 'Download failed',
                message: e instanceof Error ? e.message : 'Failed to download photo',
                color: 'red',
            });
        }
    }, [props.photosAdapter, photo]);

    const handleFavorite = useCallback(async () => {
        if (!photo) return;
        const newFavorite = !isFavorite;
        try {
            await props.photosAdapter.favoritePhoto(photo.id, newFavorite);
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
    }, [props.photosAdapter, photo, isFavorite]);

    const handleRotate = useCallback(async (direction: 'cw' | 'ccw') => {
        if (!photo) return;
        try {
            await props.photosAdapter.rotatePhoto(photo.id, direction);
            notifications.show({
                title: 'Photo rotated',
                message: `Rotated ${direction === 'cw' ? 'clockwise' : 'counter-clockwise'}`,
                color: 'green',
            });
        } catch (e) {
            notifications.show({
                title: 'Rotation failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, photo]);

    const blobUrlRef = React.useRef<string | undefined>(undefined);

    useEffect(() => {
        if (photoUrl && typeof photoUrl === 'string' && photoUrl.startsWith('blob:')) {
            blobUrlRef.current = photoUrl;
        }
        return () => {
            revokeBlobUrl(blobUrlRef.current);
            blobUrlRef.current = undefined;
        };
    }, [photoUrl]);

    const tableRow = (key: string, value: string) =>
        <Table.Tr key={key}>
            <Table.Td>{key}</Table.Td>
            <Table.Td>{value}</Table.Td>
        </Table.Tr>

    const buildRows = (metadata: Record<string, any>): React.ReactNode => {
        if (!metadata) {
            return null;
        }
        return Object.entries(metadata)
            .filter(([key, value]) => (value !== undefined) && (![...Array(100).keys()].map(v => v.toString()).includes(key)) && (value !== ""))
            .map(([key, value]): React.ReactNode => {
                switch (valueType(value)) {
                    case ('json'):
                        return buildRows(JSON.parse(value));
                    case ('object'):
                        return buildRows(value);
                    case ('string'):
                        return tableRow(key, value);
                    default:
                        return null;
                }
            })
    };

    return (
        photo && photoUrl ?
            <>
                <Image
                    radius={"md"}
                    mah="600px"
                    fit="scale-down"
                    src={photoUrl}
                />
                <Group justify="center" mt="md" gap="md">
                    <Button onClick={handleFavorite} leftSection={isFavorite ? <IconHeartFilled size={16} /> : <IconHeart size={16} />} color={isFavorite ? 'pink' : undefined}>
                        {isFavorite ? 'Favorited' : 'Favorite'}
                    </Button>
                    {props.sharesAdapter && (
                        <Button onClick={() => setShowShareModal(true)} leftSection={<IconShare2 size={16} />} variant="light">
                            Share
                        </Button>
                    )}
                    <Button onClick={() => handleRotate('cw')} leftSection={<IconRotateClockwise size={16} />} variant="light">
                        Rotate
                    </Button>
                    <Button onClick={handleDownload} leftSection={<IconDownload size={16} />}>Download</Button>
                    <Button onClick={toggle}>Metadata</Button>
                </Group>
                {isMobile ? (
                    <Drawer opened={opened} onClose={close} title="Metadata" position="bottom" size="md">
                        <ScrollArea>
                            <Table>
                                <Table.Thead>
                                    <Table.Tr>
                                        <Table.Th>Parameter</Table.Th>
                                        <Table.Th>Value</Table.Th>
                                    </Table.Tr>
                                </Table.Thead>
                                <Table.Tbody>
                                    {photo.metadata && buildRows(photo.metadata)}
                                </Table.Tbody>
                            </Table>
                        </ScrollArea>
                    </Drawer>
                ) : (
                    <Dialog opened={opened} withCloseButton onClose={close} size="lg" radius="md" mah="50%"
                            position={{top: "30%", right: 50, bottom: 50}}>
                        <ScrollArea h={400}>
                            <Table>
                                <Table.Thead>
                                    <Table.Tr>
                                        <Table.Th>Parameter</Table.Th>
                                        <Table.Th>Value</Table.Th>
                                    </Table.Tr>
                                </Table.Thead>
                                <Table.Tbody>
                                    {photo.metadata && buildRows(photo.metadata)}
                                </Table.Tbody>
                            </Table>
                        </ScrollArea>
                    </Dialog>
                )}
                {props.sharesAdapter && photo && (
                    <ShareModal
                        opened={showShareModal}
                        onClose={() => setShowShareModal(false)}
                        resourceType="photo"
                        resourceId={photo.id}
                        resourceName={photo.name}
                        sharesAdapter={props.sharesAdapter}
                    />
                )}
            </>
            : <Loader size={"md"}/>
    );
};
