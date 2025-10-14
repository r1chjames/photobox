import React, {useCallback, useEffect, useRef, useState} from 'react';
import { useParams } from "react-router-dom";
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import usePhotoGrid from "./usePhotoGrid";
import {PhotoCard} from "../PhotoCard/PhotoCard";
import {Flex, Modal, Skeleton, Space, Text, ThemeIcon, Title} from "@mantine/core";
import {useMediaQuery} from "@mantine/hooks";
import {IconPhotoX} from "@tabler/icons-react";
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

// By adding a custom comparison function to React.memo, we prevent re-renders unless the photo's ID changes.
const GridImageItem = React.memo(
    ({photo, onImageClick}: { photo: Photo, onImageClick: (id: string) => void }) => {
        const imgSrc = photo.thumbnail && (photo.thumbnail.startsWith('http') || photo.thumbnail.startsWith('data:image'))
            ? photo.thumbnail
            : `data:image/png;base64,${photo.thumbnail}`;

        return (
            <div className="item" onClick={() => onImageClick(photo.id)}>
                <div className="thumbnail">
                    <img
                        src={imgSrc}
                        alt={photo.name}
                    />
                </div>
            </div>
        );
    },
    (prevProps, nextProps) => prevProps.photo.id === nextProps.photo.id
);
GridImageItem.displayName = 'GridImageItem';

export const PhotoGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const { id } = useParams<{ id: string }>();
    const props = {...defaultProps, ...propsIn};
    const [isImageModalOpen, setImageModalOpen] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const isMobile = useMediaQuery('(max-width: 50em)');

    // The hook now provides a simple, flat, de-duplicated array of photos.
    const {photos, albumName, allRetrieved, fetchNextPage, isFetchingNextPage} = usePhotoGrid(props.photosAdapter, props.albumsAdapter, id);

    useEffect(() => {
        // This effect checks if the page content is enough to be scrollable.
        // If not, it fetches more photos until the viewport is full.
        const checkAndLoadMore = () => {
            const hasOverflow = document.documentElement.scrollHeight > window.innerHeight;
            if (!hasOverflow && !allRetrieved && !isFetchingNextPage) {
                fetchNextPage();
            }
        };

        // Use a timeout to ensure the grid has had time to render.
        const timer = setTimeout(checkAndLoadMore, 300);
        return () => clearTimeout(timer);
    }, [photos, allRetrieved, isFetchingNextPage, fetchNextPage]);


    const photosRef = useRef(photos);
    photosRef.current = photos;

    const appendDebounceRef = useRef<NodeJS.Timeout | null>(null);

    const onRequestAppend = useCallback(() => {
        if (isFetchingNextPage || allRetrieved) {
            return;
        }
        if (appendDebounceRef.current) {
            return;
        }
        appendDebounceRef.current = setTimeout(() => {
            appendDebounceRef.current = null;
        }, 500);
        fetchNextPage();
    }, [isFetchingNextPage, allRetrieved, fetchNextPage]);

    useEffect(() => {
        const handleScroll = () => {
            // Check if we are near the bottom of the page.
            if (window.innerHeight + document.documentElement.scrollTop < document.documentElement.offsetHeight - 300) {
                return;
            }
            onRequestAppend();
        };

        window.addEventListener('scroll', handleScroll);
        return () => window.removeEventListener('scroll', handleScroll);
    }, [onRequestAppend]);

    const onImageClick = useCallback((photoId: string) => {
        const photoIndex = photosRef.current.findIndex(p => p.id === photoId);
        if (photoIndex !== -1) {
            setCurrentIndex(photoIndex);
            setImageModalOpen(true);
        }
    }, []);

    const closeModal = useCallback(() => {
        setImageModalOpen(false);
    }, []);

    const handlePreviousPhoto = () => {
        setCurrentIndex(prevIndex => prevIndex > 0 ? prevIndex - 1 : 0);
    }

    const handleNextPhoto = () => {
        setCurrentIndex(prevIndex => prevIndex < photosRef.current.length - 1 ? prevIndex + 1 : prevIndex);
    }

    const EmptyAlbumContent = () => (
        <>
            <AlbumTitle/>
            <Flex justify="center" align="center" direction="column" wrap="wrap" pos={"absolute"} w={"100%"} h={"100%"}>
                <ThemeIcon radius="md" size="xl" color="orange"><IconPhotoX size="5rem"/></ThemeIcon>
                <Space h="md"/>
                <Text size="h4">This album is empty</Text>
            </Flex>
        </>
    );

    const AlbumTitle = () => <Title size="h4">{albumName}</Title>;

    const GridContent = () => {
        return (
            <div>
                <AlbumTitle/>
                <JustifiedInfiniteGrid
                    placeholder={<Skeleton height={7} mt={6} radius="md"/>}
                    className="container"
                    gap={10}
                    stretch={true}
                    passUnstretchRow={true}
                    onRequestAppend={onRequestAppend}
                    threshold={300}
                >
                    {photos.map((photo: Photo) => (
                        <GridImageItem
                            data-grid-groupkey={photo.id}
                            key={photo.id}
                            photo={photo}
                            onImageClick={onImageClick}
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
                            source={photos[currentIndex]}
                            previousPhoto={handlePreviousPhoto}
                            nextPhoto={handleNextPhoto}
                            firstInAlbum={currentIndex === 0}
                            lastInAlbum={currentIndex === photos.length - 1}
                            closeModal={closeModal}
                        />
                    </Modal>
                )}
            </div>
        )
    }

    return (photos && photos.length > 0) ? <GridContent/> : <EmptyAlbumContent/>;
};
