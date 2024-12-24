import React, {useState} from 'react';
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import './PhotoIndexView.css';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoIndexView from "./usePhotoIndexView";
import {PhotoCard} from "../PhotoCard/PhotoCard";
import {Image, Modal, Skeleton, Title} from "@mantine/core";
import {useMediaQuery} from "@mantine/hooks";

interface IProps {
    albumId?: string;
    photosAdapter: IPhotosAdapter;
    maxDisplayed?: number;
}

const defaultProps = {
    albumId: "undefined",
    maxDisplayed: 20000000
}

export const PhotoIndexView: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const [isImageModalOpen, setImageModalOpen] = useState(false);
    const [currentIndex, setCurrentIndex] = useState(0);
    const isMobile = useMediaQuery('(max-width: 50em)');
    const [{photos, allRetrieved, retrievePhotos}] = usePhotoIndexView(props.photosAdapter, props.albumId);
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

    const handleModalOpen = (idx: number) => {
        console.log(idx);
        console.log(JSON.stringify(photos[idx]))
        setCurrentIndex(idx);
        setImageModalOpen(true);
    };

    const handleModalClose = () => {
        setCurrentIndex(0);
        setImageModalOpen(false);
    };

    const handlePreviousPhoto = () => {
        setCurrentIndex(currentIndex - 1);
    }

    const handleNextPhoto = () => {
        setCurrentIndex(currentIndex + 1);
    }

    const GridImageItem = ({photo, index}: any) =>
        <div className="item" style={{width: "50px"}}>
            <div className="thumbnail">
                <Image
                    radius="md"
                    fit="contain"
                    p={5}
                    src={photo.thumbnailPath}
                    alt={photo.name}
                    onClick={() => handleModalOpen(index)}
                    onError={(e: any) => e.target.src = '/no_image.png'}
                />
            </div>
        </div>;


    const AlbumTitle = () => props.albumId !== "undefined" ? <Title size="h4">{props.albumId}</Title> : <></>

    return (photos.length === 0)
        ?
        <>
            This album is empty
        </>
        :
        <>
            <AlbumTitle/>
            <JustifiedInfiniteGrid
                placeholder={<Skeleton height={7} mt={6} radius="md"/>}
                options={{isConstantSize: false, transitionDuration: 0.2, useFit: true}}
                className="container"
                stretch={true}
                passUnstretchRow={true}
                sizeRange={[228, 228]}
                stretchRange={[144, 320]}
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
                onClose={() => handleModalClose()}>
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
};
