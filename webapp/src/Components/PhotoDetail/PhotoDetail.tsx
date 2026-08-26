import React, {useCallback, useEffect, useState} from 'react';
import {useHotkeys} from '@mantine/hooks';
import {Button, Loader, Modal, TextInput} from '@mantine/core';
import {IconChevronLeft, IconChevronRight, IconDownload, IconHeart, IconHeartFilled, IconPhotoEdit, IconRotateClockwise, IconShare2, IconX} from '@tabler/icons-react';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {ISharesAdapter} from '../../Adapters/ISharesAdapter';
import {ShareModal} from '../ShareModal/ShareModal';
import {useNavigate, useParams} from "react-router-dom";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";
import {fetchThumbnailWithAuth, getCachedThumbnail, revokeThumbnail} from "../../utils/ThumbnailUtils";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import {notifications} from '@mantine/notifications';
import {optimisticallyUpdatePhoto} from '../../utils/queryClientHelpers';
import {MetadataPanel} from './MetadataPanel';
import {EditPanel} from './EditPanel';
import {Photo} from '../../Models/Photo';

interface IProps {
    photosAdapter: IPhotosAdapter;
    sharesAdapter?: ISharesAdapter;
}

const FilmstripThumb: React.FC<{ photo: Photo; active: boolean; photosAdapter: IPhotosAdapter; onClick: () => void }> = ({photo, active, photosAdapter, onClick}) => {
    const [url, setUrl] = useState<string | undefined>(photo.thumbnailUrl ?? getCachedThumbnail(photo.id) ?? undefined);

    useEffect(() => {
        let alive = true;
        const cached = getCachedThumbnail(photo.id);
        if (cached) {
            setUrl(cached);
            return;
        }
        fetchThumbnailWithAuth(photosAdapter, photo.id)
            .then(u => { if (alive && u) setUrl(u); })
            .catch(() => { /* thumbnail may be missing; leave blank */ });
        return () => { alive = false; };
    }, [photo.id, photosAdapter]);

    return (
        <div
            className={active ? 'filmstrip-thumb active' : 'filmstrip-thumb'}
            onClick={onClick}
            role="button"
            tabIndex={0}
            aria-label={photo.name}
            onKeyDown={(e) => e.key === 'Enter' && onClick()}
        >
            {url && <img src={url} alt={photo.name} loading="lazy" />}
        </div>
    );
};

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const {id} = useParams();
    const navigate = useNavigate();
    const [isFavorite, setIsFavorite] = useState(false);
    const [showShareModal, setShowShareModal] = useState(false);
    const [showTagModal, setShowTagModal] = useState(false);
    const [showEditPanel, setShowEditPanel] = useState(false);
    const [tagInput, setTagInput] = useState('');
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
    const handleClose = useCallback(() => navigate(-1), [navigate]);

    useHotkeys([
        ['ArrowLeft', handlePrevious],
        ['ArrowRight', handleNext],
        ['Escape', handleClose],
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

    // Poster frame for videos — the generated thumbnail shown before playback.
    const [posterUrl, setPosterUrl] = useState<string | undefined>(undefined);
    useEffect(() => {
        let active = true;
        if (photo?.mediaType === 'video' && id) {
            fetchThumbnailWithAuth(props.photosAdapter, id).then(url => {
                if (active && url) setPosterUrl(url);
            }).catch(() => { /* non-fatal */ });
        }
        return () => {
            active = false;
            if (posterUrl) revokeThumbnail(id!);
        };
        // NOTE: react-hooks plugin not registered; deps intentionally [photo?.mediaType, id, props.photosAdapter].
    }, [photo?.mediaType, id, props.photosAdapter]);

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
            queryClient.invalidateQueries({queryKey: ['fetchPhoto', photo.id]});
            queryClient.invalidateQueries({queryKey: ['fetchPhotoBin', photo.id]});
        } catch (e) {
            notifications.show({
                title: 'Rotation failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        }
    }, [props.photosAdapter, photo, queryClient]);

    const handleEditSaved = useCallback(() => {
        if (!photo) return;
        queryClient.invalidateQueries({queryKey: ['fetchPhoto', photo.id]});
        queryClient.invalidateQueries({queryKey: ['fetchPhotoBin', photo.id]});
    }, [photo]);

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

    // Keep the active filmstrip thumbnail centered as navigation changes.
    useEffect(() => {
        document.querySelector('.filmstrip-thumb.active')?.scrollIntoView({
            behavior: 'smooth',
            inline: 'center',
            block: 'nearest',
        });
    }, [id]);

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

    const tags = photo?.tags ? photo.tags.split(',').map(t => t.trim()).filter(Boolean) : [];

    return (
        photo && photoUrl ?
            <div className="lightbox-overlay active">
                <div className="lightbox-header">
                    <button type="button" className="icon-btn" title="Edit" aria-label="Edit photo" onClick={() => setShowEditPanel(prev => !prev)}>
                        <IconPhotoEdit size={18} stroke={1.8} />
                    </button>
                    <button type="button" className="icon-btn" title={isFavorite ? 'Remove from favorites' : 'Add to favorites'} aria-label="Toggle favorite" onClick={handleFavorite}>
                        {isFavorite ? <IconHeartFilled size={18} color="var(--pb-favorite)" /> : <IconHeart size={18} stroke={1.8} />}
                    </button>
                    {props.sharesAdapter && (
                        <button type="button" className="icon-btn" title="Share" aria-label="Share photo" onClick={() => setShowShareModal(true)}>
                            <IconShare2 size={18} stroke={1.8} />
                        </button>
                    )}
                    <button type="button" className="icon-btn" title="Rotate" aria-label="Rotate photo" onClick={() => handleRotate('cw')}>
                        <IconRotateClockwise size={18} stroke={1.8} />
                    </button>
                    <button type="button" className="icon-btn" title="Download" aria-label="Download photo" onClick={handleDownload}>
                        <IconDownload size={18} stroke={1.8} />
                    </button>
                    <button type="button" className="icon-btn" title="Close" aria-label="Close viewer" onClick={handleClose}>
                        <IconX size={18} stroke={1.8} />
                    </button>
                </div>

                <div className="lightbox-body">
                    {!isFirstPhoto && (
                        <button type="button" className="lightbox-nav-btn" onClick={handlePrevious} aria-label="Previous photo">
                            <IconChevronLeft size={22} stroke={2} />
                        </button>
                    )}

                    <div className="lightbox-main-image-container">
                        {photo.mediaType === 'video' ? (
                            <video
                                src={photoUrl}
                                controls
                                poster={posterUrl}
                                preload="metadata"
                                className="lightbox-main-image"
                            />
                        ) : (
                            <img className="lightbox-main-image" src={photoUrl} alt={photo.name} />
                        )}
                    </div>

                    {!isLastPhoto && (
                        <button type="button" className="lightbox-nav-btn" onClick={handleNext} aria-label="Next photo">
                            <IconChevronRight size={22} stroke={2} />
                        </button>
                    )}

                    <div className="metadata-drawer">
                        <div className="drawer-section">
                            <div className="drawer-title-row">Details</div>
                            <MetadataPanel photo={photo} photosAdapter={props.photosAdapter} />
                            {showEditPanel && (
                                <EditPanel photo={photo} photosAdapter={props.photosAdapter} onSaved={handleEditSaved} />
                            )}
                        </div>
                        {tags.length > 0 && (
                            <div className="drawer-section">
                                <div className="drawer-title-row">Tags</div>
                                <div className="tags-container">
                                    {tags.map(tag => (
                                        <button type="button" key={tag} className="tag-chip" onClick={() => setShowTagModal(true)}>
                                            {tag}
                                        </button>
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                </div>

                <div className="lightbox-filmstrip">
                    {(allPhotos ?? []).map(p => (
                        <FilmstripThumb
                            key={p.id}
                            photo={p}
                            active={p.id === id}
                            photosAdapter={props.photosAdapter}
                            onClick={() => navigateToPhoto(p.id)}
                        />
                    ))}
                </div>

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
            </div>
            : <Loader size={"md"}/>
    );
};
