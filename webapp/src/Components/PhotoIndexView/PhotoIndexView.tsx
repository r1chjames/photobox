import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { JustifiedLayout, OnLayoutComplete } from '@egjs/react-infinitegrid';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import './PhotoIndexView.css';
import { PhotoItem } from '../PhotoItem/PhotoItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { MainContent } from '../MainContent/MainContent';
import {LoadingScreen} from '../LoadingScreen/LoadingScreen';

interface IProps {
  baseApiUrl: string;
}

const getAllPhotos = async(baseApiUrl: string, albumId: string) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  let photos: Photo[];
  photos = await photosAdapter.getPhotosInfoInAlbum(albumId, 1, 100);
  return photos;
};

const getAlbumName = async(baseApiUrl: string, albumId: string) => {
  const albumsAdapter = new AlbumsAdapter(baseApiUrl);
  if (albumId !== null) {
    const album = await albumsAdapter.getAlbumInfoById(albumId);
    return album.name;
  }
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

const onLayoutComplete = (params: OnLayoutComplete) => {
  // @ts-ignore
  return !params.isLayout && params.endLoading();
};

export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {
  const [photos, setPhotos] = useState<JSX.Element[]>();
  const [albumName, setAlbumName] = useState('');
  const { id } = useParams();
  const [albumNameLoaded, setAlbumNameLoaded] = useState(false);
  const [photosLoaded, setPhotosLoaded] = useState(false);

  useEffect(() => {
    if (id !== undefined) {
      (async function retrieveAlbumName() {
        const retrievedAlbumName = await getAlbumName(props.baseApiUrl, id);
        setAlbumName(retrievedAlbumName);
      })().then(() => setAlbumNameLoaded(true));
    } else {
      setAlbumName('Photos');
      setPhotosLoaded(true);
    }

    (async function retrievePhotos() {
      const retrievedPhotos = await getAllPhotos(props.baseApiUrl, id);
      const photoItems = loadItems(0, retrievedPhotos, props.baseApiUrl);
      setPhotos(photoItems);
    })().then(() => setPhotosLoaded(true));
  },        [setPhotos, id, props.baseApiUrl]);

  const apiCallsCompleted = () => (albumNameLoaded && photosLoaded);

  const loadItems = (groupKey: number, imageSources: Photo[], baseApiUrl: string) => {
    const items = [];
    const start = 0;

    if (imageSources) {
      for (let i = 0; i < imageSources.length; i += 1) {
        const imageSource = imageSources[start + i];
        if (typeof imageSource !== 'undefined') {
          items.push(
            <PhotoItem
              baseApiUrl={baseApiUrl}
              groupKey={groupKey}
              num={1 + start + i}
              key={start + i}
              source={imageSource}
              allPhotos={imageSources}
            />
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
        <JustifiedLayout
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          // onAppend={onAppend}
          onLayoutComplete={onLayoutComplete}
        >
          {photos}
        </JustifiedLayout>
      );
    }

    return (
      <LoadingScreen />
    );
  };

  return (
    <MuiThemeProvider>
      <MainContent title={albumName}>
        <div className="photoIndexView__photoIndex">
          {content()}
        </div>
      </MainContent>
    </MuiThemeProvider>
  );
};
