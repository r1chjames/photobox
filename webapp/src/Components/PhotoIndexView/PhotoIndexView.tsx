import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import {MasonryInfiniteGrid} from '@egjs/react-infinitegrid';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import './PhotoIndexView.css';
import { PhotoItem } from '../PhotoItem/PhotoItem';
//import { ImplAlbumsAdapter } from '../../Adapters/ImplAlbumsAdapter';
import {Loader, MantineProvider} from '@mantine/core';
// import {GroupedItemWrapper} from '../GroupedItemWrapper/GroupedItemWrapper';

interface IProps {
  photosAdapter: PhotosAdapter
}

const getAllPhotos = async(propsPhotosAdapter: PhotosAdapter, albumId: string) => {
  const photosAdapter = propsPhotosAdapter;
  const photos = await photosAdapter.getPhotosInfoInAlbum(albumId, 1, 100);
  return photos;
};

//const getAlbumName = async(baseApiUrl: string, albumId: string) => {
//  const albumsAdapter = new ImplAlbumsAdapter(baseApiUrl);
//  if (albumId !== null) {
//    const album = await albumsAdapter.getAlbumInfoById(albumId);
//    return album.name;
//  }
//};

const getPhotoDate = (imageSource: Photo): string => {
  const exifVal = imageSource.metadata.exif;
  // @ts-ignore
  return exifVal !== null ? exifVal.DateTime : imageSource.createdAt;
};

const updatePhotoCollection = (items: Map<string, JSX.Element[]>, date: string, itemToAdd: JSX.Element) => {
  if (date == null || undefined) {
    console.log(itemToAdd)
  }
  const itemToUpdate = items.has(date) ? items.get(date) : [];
  itemToUpdate!.push(itemToAdd);
  return items.set(date, itemToUpdate!);
};

const parseItems = (items: Map<string, JSX.Element[]>) => {
  const itemsToReturn: JSX.Element[] = [];
  items.forEach((values, key) => {
    itemsToReturn.push(
      <div>
        {key}
        {values}
      </div>
    );
  });

  return itemsToReturn;
};

// const onAppend = async (params: OnAppend) => {
//   // @ts-ignore
//   params.startLoading();
//   const imageSources = await getAllPhotos();
//   if (imageSources.length === this.state.images.length) {
//     return;
//   }
//   const items = this.loadItems(
//     (parseFloat(params.groupKey as string) || 0) + 1,
//     20,
//     imageSources
//   );
//
//   this.setState({ images: this.state.images.concat(items) });
// }

// const onLayoutComplete = (params: OnLayoutComplete) => {
//   // eslint-disable-next-line @typescript-eslint/ban-ts-comment
//   // @ts-ignore
//   return !params.isLayout && params.endLoading();
// };

export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {
  const [photos, setPhotos] = useState<Map<string, JSX.Element[]>>();
//  const [albumName, setAlbumName] = useState('');
  const { id } = useParams();
//  const [albumNameLoaded, setAlbumNameLoaded] = useState(false);
  const [photosLoaded, setPhotosLoaded] = useState(false);

  useEffect(() => {
//    if (id !== undefined) {
//      (async function retrieveAlbumName() {
//        const retrievedAlbumName = await getAlbumName(props.baseApiUrl, id);
//        setAlbumName(retrievedAlbumName);
//      })().then(() => setAlbumNameLoaded(true));
//    } else {
//      setAlbumName('Photos');
//      setAlbumNameLoaded(true);
//    }

    (async function retrievePhotos() {
      const retrievedPhotos = await getAllPhotos(props.photosAdapter, id as string);
      const photoItems = loadItems(0, retrievedPhotos, "TODO");
      setPhotos(photoItems);
    })().then(() => setPhotosLoaded(true));
  },        [setPhotos, id, props.photosAdapter]);

  const apiCallsCompleted = () => (photosLoaded);

  const loadItems = (groupKey: number, imageSources: Photo[], baseApiUrl: string) => {
    let items = new Map<string, JSX.Element[]>();
    const start = 0;

    if (imageSources) {
      for (let i = 0; i < imageSources.length; i += 1) {
        const imageSource = imageSources[start + i];
        if (typeof imageSource !== 'undefined') {
          items = updatePhotoCollection(items,
                                        getPhotoDate(imageSource),
                                        (
                                          <PhotoItem
                                            baseApiUrl={baseApiUrl}
                                            groupKey={groupKey}
                                            num={1 + start + i}
                                            key={start + i}
                                            source={imageSource}
                                            allPhotos={imageSources}
                                          />
                                        )
          );
        }
      }
      return items;
    }
  };

  const content = () => {
    if (apiCallsCompleted() && !photos) {
      return (
        <div>
          This album is empty
        </div>
      );
    }

    if (apiCallsCompleted() && photos) {
      return (
        <MasonryInfiniteGrid
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          // onAppend={onAppend}
          // onLayoutComplete={onLayoutComplete}
        >
          {parseItems(photos)}
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
