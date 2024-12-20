import React, { useEffect, useState } from 'react';
import {JustifiedInfiniteGrid} from '@egjs/react-infinitegrid';
import './PhotoIndexView.css';
import {Loader} from '@mantine/core';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoIndexView from "./usePhotoIndexView";
import {PhotoItem} from "../PhotoItem/PhotoItem";
import {useParams} from "react-router-dom";

interface IProps {
  photosAdapter: IPhotosAdapter;
}

export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {
  const { albumId } = useParams();
  const [{photos, isApiCallsRunning, retrievePhotos}] = usePhotoIndexView(props.photosAdapter, albumId);
  // const [photoItemsByDate, setPhotoItemsByDate] = useState<Map<string, JSX.Element[]>>(new Map<string, JSX.Element[]>);
  const [photoItems, setPhotoItems] = useState<JSX.Element[]>([]);
  const [isLoading, setIsLoading] = useState(false);

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

  const loadItems = (groupKey: number) => {
    setIsLoading(true);
    if (photos) {
      photos.forEach(function (photo, i) {
        if (typeof photo !== 'undefined') {
          photoItems.push(
              <PhotoItem
                src={photo.sourcePath}
                thumbnail={photo.thumbnailPath}
                groupKey={groupKey}
                source={photo}
                photoSequence={i}
                previousPhoto={() => console.log("previous")}
                nextPhoto={() => console.log("next")}
                lastInAlbum={i < photos.length}
            />
          );
        }
      });
    }
    setPhotoItems(photoItems);
    setIsLoading(false);
  };

  useEffect(() => {
    if (!isApiCallsRunning) loadItems(1)
  }, [photos]);

  const onRequestAppend = async (e: any) => {
    const nextGroupKey = (+e.groupKey! || 0) + 1;
    e.wait();
    e.currentTarget.appendPlaceholders(5, nextGroupKey);
    setTimeout(() => {
      e.ready();
      retrievePhotos(nextGroupKey, 30);
    }, 1000);
  }

  const content = () => {
    if (isLoading && !photoItems) {
      return (
        <div>
          This album is empty
        </div>
      );
    }

    const Item = ({ num }: any) => <div className="item" style={{
      width: "250px",
    }}>
      <div className="thumbnail">
        <img
            src={`https://naver.github.io/egjs-infinitegrid/assets/image/${(num % 33) + 1}.jpg`}
            alt="egjs"
        />
      </div>
      <div className="info">{`egjs ${num}`}</div>
    </div>;

    if (!isLoading && photoItems) {
      return (
        <JustifiedInfiniteGrid
          placeholder={<div className="placeholder"></div>}
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          className="container"
          gap={5}
          onRequestAppend={onRequestAppend}
          // onAppend={onAppend}
          // onLayoutComplete={onLayoutComplete}
        >
          {photoItems.map((photo, index) => <Item data-grid-groupkey={index} key={photo.key} num={index} />)}

        </JustifiedInfiniteGrid>
      );
    }

    return (
      <Loader size={"md"} />
    );
  };

  return (
    <div className="photoIndexView__photoIndex">
      {content()}
    </div>
  );
};
