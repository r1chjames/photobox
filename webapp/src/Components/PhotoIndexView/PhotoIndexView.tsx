import React from 'react';
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import './PhotoIndexView.css';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoIndexView from "./usePhotoIndexView";
import {PhotoItem} from "../PhotoItem/PhotoItem";
import {Photo} from "../../Models/Photo";
import {Skeleton, Title} from "@mantine/core";

interface IProps {
    albumId?: string;
    photosAdapter: IPhotosAdapter;
    maxDisplayed? : number;
}

const defaultProps = {
    albumId: "undefined",
    maxDisplayed: 20000000
}

export const PhotoIndexView: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
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
  //             <PhotoItem
  //                 src={photo.sourcePath}
  //                 thumbnail={photo.thumbnailPath}
  //                 groupKey={groupKey}
  //                 source={photo}
  //                 photoSequence={i}
  //                 previousPhoto={() => console.log("previous")}
  //                 nextPhoto={() => console.log("next")}
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

  const GridImageItem = (params: {photo: Photo, groupKey: number, index: number}) => {
      return (
          <div className="item" style={{width: "50px"}}>
              <PhotoItem
                  groupKey={params.groupKey}
                  key={params.index}
                  source={params.photo}
                  photoSequence={params.index}
                  previousPhoto={(currentIndex: number) => (currentIndex !== 0) ? photos[currentIndex-1]:params.photo}
                  nextPhoto={(currentIndex: number) => (currentIndex != photos.length) ? photos[currentIndex+1]:params.photo}
                  firstInAlbum={params.index === 0}
                  lastInAlbum={params.index === photos.length}
              />
        </div>);
  }

const AlbumTitle = () => {
    return props.albumId !== "undefined"
        ?
            <Title size="h4">{props.albumId}</Title>
        :
            <></>
}

  const content = () => {
    return (photos.length === 0)
        ?
          <div>
            This album is empty
          </div>
        :
        <div>
            <AlbumTitle />
            <JustifiedInfiniteGrid
                placeholder={<Skeleton height={7} mt={6} radius="md" />}
                options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
                className="container"
                stretch={true}
                passUnstretchRow={true}
                sizeRange={[228,228]}
                stretchRange={[144,320]}
                onRequestAppend={onRequestAppend}
                >
                {photos
                    .map((photo, index) => <GridImageItem data-grid-groupkey={index} groupKey={index} index={index} photo={photo}/>)}
            </JustifiedInfiniteGrid>
        </div>
  };

  return (
    <div className="photoIndexView__photoIndex">
      {content()}
    </div>
  );
};
