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
}

const defaultProps = {
    albumId: "undefined",
}

export const PhotoIndexView: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const [{photos, retrievePhotos}] = usePhotoIndexView(props.photosAdapter, props.albumId);
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
    const nextGroupKey = (+e.groupKey! || 0) + 1;
    console.log("getting some more");
    console.log("current group key: " + e.groupKey)
    console.log("next group key: " + nextGroupKey)
    e.wait();
    e.currentTarget.appendPlaceholders(5, nextGroupKey);
    setTimeout(() => {
      e.ready();
      retrievePhotos(nextGroupKey, 30);
    }, 1000);
  }

  const GridImageItem = (params: {photo: Photo, groupKey: number, index: number}) =>
      <div className="item" style={{ width: "50px" }}>
        <PhotoItem
            src={params.photo.sourcePath}
            thumbnail={params.photo.thumbnailPath}
            groupKey={params.groupKey}
            key={params.index}
            source={params.photo}
            photoSequence={params.index}
            previousPhoto={() => console.log("previous")}
            nextPhoto={() => console.log("next")}
            lastInAlbum={params.index < photos.length}
        />
      </div>;

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
                layoutOptions={{ margin: 5, column: [0, 5] }}
                className="container"
                gap={5}
                stretch={true}
                passUnstretchRow={true}
                sizeRange={[228,228]}
                stretchRange={[144,320]}
                onRequestAppend={onRequestAppend}
                >
                {photos.map((photo, index) => <GridImageItem data-grid-groupkey={index} groupKey={index} index={index} photo={photo}/>)}
            </JustifiedInfiniteGrid>
        </div>
  };

  return (
    <div className="photoIndexView__photoIndex">
      {content()}
    </div>
  );
};
