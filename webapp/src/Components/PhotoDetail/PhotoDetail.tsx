import React, {useCallback, useEffect, useState} from 'react';
import {useDisclosure, useHotkeys} from '@mantine/hooks';
import {ActionIcon, Button, Chip, Dialog, Drawer, Group, Image, Loader, Modal, ScrollArea, TextInput} from '@mantine/core';
import {useMediaQuery} from '@mantine/hooks';
import {IconArrowLeftDashed, IconArrowRightDashed, IconDownload, IconHeart, IconHeartFilled, IconListDetails, IconRotateClockwise, IconShare2, IconTag} from '@tabler/icons-react';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {ISharesAdapter} from '../../Adapters/ISharesAdapter';
import {ShareModal} from '../ShareModal/ShareModal';
import {useNavigate, useParams} from "react-router-dom";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import {notifications} from '@mantine/notifications';
import {optimisticallyUpdatePhoto} from '../../utils/queryClientHelpers';
import {MetadataPanel} from './MetadataPanel';

interface IProps {
    photosAdapter: IPhotosAdapter;
    sharesAdapter?: ISharesAdapter;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const {id} = useParams();
    const navigate = useNavigate();
    const [opened, {toggle, close}] = useDisclosure(true);
    const [isFavorite, setIsFavorite] = useState(false);
    const [showShareModal, setShowShareModal] = useState(false);
    const [showTagModal, setShowTagModal] = useState(false);
    const [tagInput, setTagInput] = useState('');
    const isMobile = useMediaQuery('(max-width: 50em)');
    const queryClient = useQueryClient();
    // Fetch surrounding photos to determine prev/next navigation
    const {data: allPhotos} = useQuery({
        queryKey: ['allPhotosForNavigation'],
        queryFn: () => props.photosAdapter.getAllPhotosInfo('', 1000, false),
        staleTime: 1000 * 60 * 5,
    });

    const currentIndex = allPhotos?.findIndex(p => p.id === id) ?? -1;
    const prevPhotoId = currentIndex > 0 ? allPhotos![currentIndex - 1].id : null;
    const nextPhotoId = currentIndex >= 0 && currentIndex < (allPhotos?.length ?? 0) - 1 ? allPhotos![currentIndex + 1].id : null;
    const isFirstPhoto = currentIndex <= 0;
    const isLastPhoto = currentIndex >= (allPhotos?.length ?? 0) - 1 || currentIndex === -1;

    const navigateToPhoto = useCallback((photoId: string | null) => {
        if (photoId) {
            navigate(`/photo/${photoId}`);
        }
    }, [navigate]);

    const handlePrevious = useCallback(() => navigateToPhoto(prevPhotoId), [navigateToPhoto, prevPhotoId]);
    const handleNext = useCallback(() => navigateToPhoto(nextPhotoId), [navigateToPhoto, nextPhotoId]);

    useHotkeys([
        ['ArrowLeft', handlePrevious],
        ['ArrowRight', handleNext],
    ]);

    const fetchPhoto = async () => {
        return props.photosAdapter.getPhotoInfoById(id!);
    }

    const {data: photo} = useQuery({
        queryKey: ['fetchPhoto', id],
        queryFn: fetchPhoto,
        enabled: !!id,
    });

    const fetchPhotoBin = async () => {
        return fetchPhotoBinWithAuth(props.photosAdapter, id!, photo?.mediaType);
    }

    const {data: photoUrl} = useQuery({
        queryKey: ['fetchPhotoBin', id, photo?.mediaType],
        queryFn: fetchPhotoBin,
        enabled: !!id && !!photo,
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
            optimisticallyUpdatePhoto(queryClient, photo.id, { favorite: newFavorite });
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
    }, [props.photosAdapter, photo, isFavorite, queryClient]);

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

    const handleSaveTags = useCallback(async () => {
        if (!photo) return;
        const tags = tagInput.split(',').map(t => t.trim()).filter(Boolean);
        try {
            await props.photosAdapter.updatePhotoTags(photo.id, tags);
            notifications.show({
                title: 'Tags updated',
                message: `Updated tags for ${photo.name}`,
                color: 'green',
            });
            setShowTagModal(false);
        } catch (e) {
            notifications.show({
                title: 'Failed to update tags',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, photo, tagInput]);

    useEffect(() => {
        if (photo) {
            setTagInput(photo.tags || '');
        }
    }, [photo]);

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

    return (
        photo && photoUrl ?
            <>
                <div style={{ position: 'relative' }}>
                    {photo.mediaType === 'video' ? (
                        <video
                            src={photoUrl}
                            controls
                            style={{ maxHeight: '600px', width: '100%', borderRadius: '8px' }}
                        />
                    ) : (
                        <Image
                            radius={"md"}
                            mah="600px"
                            fit="scale-down"
                            src={photoUrl}
                        />
                    )}
                    {!isFirstPhoto && (
                        <ActionIcon
                            variant="light"
                            size="xl"
                            onClick={handlePrevious}
                            aria-label="Previous photo"
                            style={{
                                position: 'absolute',
                                left: 8,
                                top: '50%',
                                transform: 'translateY(-50%)',
                                zIndex: 10,
                            }}
                        >
                            <IconArrowLeftDashed size="2.125rem" />
                        </ActionIcon>
                    )}
                    {!isLastPhoto && (
                        <ActionIcon
                            variant="light"
                            size="xl"
                            onClick={handleNext}
                            aria-label="Next photo"
                            style={{
                                position: 'absolute',
                                right: 8,
                                top: '50%',
                                transform: 'translateY(-50%)',
                                zIndex: 10,
                            }}
                        >
                            <IconArrowRightDashed size="2.125rem" />
                        </ActionIcon>
                    )}
                </div>
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
                    <Button onClick={toggle} leftSection={<IconListDetails size={16} />} variant="light">Metadata</Button>
                </Group>
                {isMobile ? (
                    <Drawer opened={opened} onClose={close} title="Metadata" position="bottom" size="md">
                        <ScrollArea>
                            <Group gap="xs" mb="sm">
                                <Button size="compact-sm" variant="light" leftSection={<IconTag size={14} />} onClick={() => setShowTagModal(true)}>Edit tags</Button>
                            </Group>
                            {photo.tags && (
                                <Group gap="xs" mb="sm">
                                    {photo.tags.split(',').map(t => t.trim()).filter(Boolean).map(tag => (
                                        <Chip key={tag} size="xs" checked={false} onClick={() => {}}>{tag}</Chip>
                                    ))}
                                </Group>
                            )}
                            <MetadataPanel photo={photo} />
                        </ScrollArea>
                    </Drawer>
                ) : (
                    <Dialog opened={opened} withCloseButton onClose={close} size="lg" radius="md" mah="50%"
                            position={{top: "30%", right: 50, bottom: 50}}>
                        <ScrollArea h={400}>
                            <Group gap="xs" mb="sm">
                                <Button size="compact-sm" variant="light" leftSection={<IconTag size={14} />} onClick={() => setShowTagModal(true)}>Edit tags</Button>
                            </Group>
                            {photo.tags && (
                                <Group gap="xs" mb="sm">
                                    {photo.tags.split(',').map(t => t.trim()).filter(Boolean).map(tag => (
                                        <Chip key={tag} size="xs" checked={false} onClick={() => {}}>{tag}</Chip>
                                    ))}
                                </Group>
                            )}
                            <MetadataPanel photo={photo} />
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
                <Modal
                    opened={showTagModal}
                    onClose={() => setShowTagModal(false)}
                    title="Edit tags"
                    size="sm"
                >
                    <TextInput
                        placeholder="Enter tags separated by commas"
                        value={tagInput}
                        onChange={(e) => setTagInput(e.currentTarget.value)}
                        onKeyDown={(e) => { if (e.key === 'Enter') handleSaveTags(); }}
                        autoFocus
                    />
                    <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8, marginTop: 16 }}>
                        <Button variant="default" onClick={() => setShowTagModal(false)}>Cancel</Button>
                        <Button onClick={handleSaveTags}>Save tags</Button>
                    </div>
                </Modal>
            </>
            : <Loader size={"md"}/>
    );
};
