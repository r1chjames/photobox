import React, { useEffect, useState } from 'react';
import {MasonryInfiniteGrid} from '@egjs/react-infinitegrid';
import './PhotoIndexView.css';
import {Loader, MantineProvider} from '@mantine/core';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoIndexView from "./usePhotoIndexView";
import {PhotoItem} from "../PhotoItem/PhotoItem";

interface IProps {
  photosAdapter: IPhotosAdapter;
  albumId: string;
}

export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {
  const [{photos, isApiCallsRunning}] = usePhotoIndexView(props.photosAdapter, props.albumId);
  const [photoItemsByDate, setPhotoItemsByDate] = useState<Map<string, JSX.Element[]>>(new Map<string, JSX.Element[]>);
  const [isLoading, setIsLoading] = useState(false);

  const loadItems = (groupKey: number) => {
    setIsLoading(true);
    const photoItems = new Map<string, JSX.Element[]>();
    if (photos) {
      photos.forEach(function (photo, i) {
        if (typeof photo !== 'undefined') {
          const currentMapValue = photoItems.get(photo.getPhotoDate()) || [];
          currentMapValue.push(
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
          photoItems.set(photo.getPhotoDate(), currentMapValue);
        }
      });
    }
    setPhotoItemsByDate(photoItems);
    setIsLoading(false);
  };

  useEffect(() => {
    if (!isApiCallsRunning) loadItems(1)
  }, [photos]);

  const content = () => {
    if (isLoading && !photoItemsByDate) {
      return (
        <div>
          This album is empty
        </div>
      );
    }

    if (!isLoading && photoItemsByDate) {
      return (
        <MasonryInfiniteGrid
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          // onAppend={onAppend}
          // onLayoutComplete={onLayoutComplete}
        >
          {photoItemsByDate}
        </MasonryInfiniteGrid>
      );
    }

    return (
      <Loader size={"md"} />
    );
  };

  return (
    <MantineProvider>
        <div className="photoIndexView__photoIndex">
          {content()}
        </div>
    </MantineProvider>
  );
};
