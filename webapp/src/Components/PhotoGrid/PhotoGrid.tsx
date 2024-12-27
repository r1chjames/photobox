import React, {useState} from 'react';
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoGrid from "./usePhotoGrid";
import {PhotoCard} from "../PhotoCard/PhotoCard";
import {Flex, Modal, Skeleton, Space, Text, ThemeIcon, Title} from "@mantine/core";
import {useMediaQuery} from "@mantine/hooks";
import {IconPhotoX} from "@tabler/icons-react";
import './PhotoGrid.css';

interface IProps {
    albumId?: string;
    photosAdapter: IPhotosAdapter;
    maxDisplayed?: number;
}

const defaultProps = {
    albumId: "undefined",
    maxDisplayed: 20000000
}

export const PhotoGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const [isImageModalOpen, setImageModalOpen] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const isMobile = useMediaQuery('(max-width: 50em)');
    const [{photos, allRetrieved, retrievePhotos}] = usePhotoGrid(props.photosAdapter, props.albumId);
    // const [photoItemsByDate, setPhotoItemsByDate] = useState<Map<string, JSX.Element[]>>(new Map<string, JSX.Element[]>);

    // const loadItemsByDate = (groupKey: number) => {
    //   setIsLoading(true);
    //   const photoItems = new Map<string, JSX.Element[]>();
    //   if (photos) {
    //     photos.forEach(function (photo, i) {
    //       if (typeof photo !== 'undefined') {
    //         const currentMapValue = photoItems.get(photo.getPhotoDate()) || [];
    //         currentMapValue.push(
    //             <PhotoCard
    //                 src={photo.sourcePath}
    //                 thumbnail={photo.thumbnailPath}
    //                 groupKey={groupKey}
    //                 source={photo}
    //                 photoSequence={i}
    //                 handlePreviousPhoto={() => console.log("previous")}
    //                 handleNextPhoto={() => console.log("next")}
    //                 lastInAlbum={i < photos.length}
    //             />
    //           );
    //         photoItems.set(photo.getPhotoDate(), currentMapValue);
    //       }
    //     });
    //   }
    //   setPhotoItemsByDate(photoItems);
    //   setIsLoading(false);
    // };

    const onRequestAppend = async (e: any) => {
        if ((photos.length < props.maxDisplayed) && !allRetrieved) {
            const nextGroupKey = (+e.groupKey! || 0) + 1;
            e.wait();
            e.currentTarget.appendPlaceholders(5, nextGroupKey);
            await retrievePhotos(nextGroupKey, 30);
            e.ready();
        }
    }

    const toggleModal = (idx: number) => {
        setCurrentIndex(idx);
        setImageModalOpen(!isImageModalOpen);
    };

    const handlePreviousPhoto = () => {
        setCurrentIndex(currentIndex - 1);
    }

    const handleNextPhoto = () => {
        setCurrentIndex(currentIndex + 1);
    }

    const GridImageItem = ({photo, index}: any) =>
        <div className="item">
            <div className="thumbnail">
                <img
                    src={photo.thumbnailPath}
                    alt={photo.name}
                    onClick={() => toggleModal(index)}
                    // data-grid-maintained-target="true"
                />
            </div>
        </div>;

    const EmptyAlbumContent = () =>
        <>
            <AlbumTitle/>
            <Flex justify="center"
                align="center"
                direction="column"
                wrap="wrap"
                pos={"absolute"}
                w={"100%"}
                h={"100%"}>
                <ThemeIcon radius="md" size="xl" color="orange">
                    <IconPhotoX size="5rem"/>
                </ThemeIcon>
                <Space h="md" />
                <Text size="h4">This album is empty</Text>
            </Flex>
        </>

    const AlbumTitle = () => props.albumId !== "undefined" ? <Title size="h4">{props.albumId}</Title> : <></>

    const GridContent = () => {
        return (
            <>
                <AlbumTitle/>
                <JustifiedInfiniteGrid
                    placeholder={<Skeleton height={7} mt={6} radius="md"/>}
                    className="container"
                    gap={10}
                    stretch={true}
                    passUnstretchRow={true}
                    onRequestAppend={onRequestAppend}>
                    {photos.map((photo, index) =>
                        <GridImageItem
                            data-grid-groupkey={index}
                            key={index}
                            photo={photo}
                            index={index}/>)}
                </JustifiedInfiniteGrid>
                <Modal
                    opened={isImageModalOpen}
                    withCloseButton={false}
                    aria-labelledby="customized-dialog-title"
                    size="auto"
                    padding={"0"}
                    m={"0"}
                    overlayProps={{
                        backgroundOpacity: 0.55,
                    }}
                    fullScreen={isMobile}
                    transitionProps={{transition: 'fade', duration: 200}}
                    onClose={() => toggleModal(0)}>
                    <PhotoCard
                        source={photos[currentIndex]}
                        previousPhoto={() => handlePreviousPhoto()}
                        nextPhoto={() => handleNextPhoto()}
                        firstInAlbum={currentIndex === 0}
                        lastInAlbum={currentIndex === photos.length}
                        closeModal={() => setImageModalOpen(false)}
                    />
                </Modal>
            </>
        )
    }

    return (photos.length === 0) ? EmptyAlbumContent() : GridContent();
};
