import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {MasonryInfiniteGrid} from '@egjs/react-infinitegrid';
import { Photo } from '../../Models/Photo';
import './PhotoIndexView.css';
import {Loader, MantineProvider} from '@mantine/core';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import useAlbumIndexView from "../AlbumIndexView/useAlbumIndexView";
import usePhotoIndexView from "./usePhotoIndexView";
// import {GroupedItemWrapper} from '../GroupedItemWrapper/GroupedItemWrapper';

interface IProps {
  photosAdapter: IPhotosAdapter;
  albumId: string;
}



export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {
  const [{photos, isLoading}] = usePhotoIndexView(props.photosAdapter, props.albumId);
  const [photoItemsByDate, setPhotoItemsByDate] = useState();

  const updatePhotoCollectionByDate = (items: Map<string, JSX.Element>) => {
    // let itemDateMap = new Map<string, JSX.Element[]>;
    // items.forEach((jsxItem, date) => {
    //   const currentMapValue = itemDateMap.get(date);
    //   itemDateMap.set(
    //       date,
    //       currentMapValue.push(jsxItem)
    //   )
    // })
    return Object.groupBy(items, (entries, index) => {
      return entries.project;
    });
  };


  const loadItems = (groupKey: number) => {
    let items= new Map<string, JSX.Element>;
    const start = 0;

    if (photos) {
      for (let i = 0; i < photos.length; i += 1) {
        const photo = photos[start + i];
        if (typeof photo !== 'undefined') {
          items.set(
              photo.getPhotoDate(),
                  <PhotoItem
                      src={photo.sourcePath}
                      thumbnail={photo.thumbnailPath}
                      groupKey={groupKey}
                      num={1 + start + i}
                      key={start + i}
                      source={photo}
                      allPhotos={photos}
                  />
              )
            }
        }
      }
    setPhotoItemsByDate(updatePhotoCollectionByDate(items));
    };

  useEffect(() =>
      loadItems(1),
      [photoItemsByDate]
  );

  {/*const parseItems = (items: Map<string, JSX.Element[]>) => {*/}
  {/*  const itemsToReturn: JSX.Element[] = [];*/}
  {/*  items.forEach((values, key) => {*/}
  {/*    itemsToReturn.push(*/}
  {/*        <div>*/}
  {/*          {key}*/}
  {/*          {values}*/}
  {/*        </div>*/}
  {/*    );*/}
  {/*  });*/}

  {/*  return itemsToReturn;*/}
  {/*};*/}

  const content = () => {
    if (isLoading && !photos) {
      return (
        <div>
          This album is empty
        </div>
      );
    }

    if (isLoading && photos) {
      return (
        <MasonryInfiniteGrid
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          // onAppend={onAppend}
          // onLayoutComplete={onLayoutComplete}
        >
          {photoItems}
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
