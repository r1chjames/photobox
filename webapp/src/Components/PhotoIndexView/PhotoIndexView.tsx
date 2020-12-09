import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { JustifiedLayout, OnLayoutComplete } from '@egjs/react-infinitegrid';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import './PhotoIndexView.css';
import { ImageItem } from '../ImageItem/ImageItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { TitleBar } from '../TitleBar/TitleBar';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { MainContent } from '../MainContent/MainContent';
import { Typography } from '@material-ui/core';
// import { Album } from '../../Models/Album';
// import { useParams } from 'react-router-dom';

interface IProps {
  baseApiUrl: string;
  // albumId: string | null;
}

// interface IState {
//   images: ReactElement[];
//   albumTitle: string;
//   albumId: string;
// }

const getAllPhotos = async(baseApiUrl: string, albumId: string) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  let photos: Photo[];
  // if (albumId != null) {
  photos = await photosAdapter.getPhotosInfoInAlbum(albumId);
  // } else {
  //   photos = await photosAdapter.getAllPhotosInfo();
  // }
  return photos;
};

const loadItems = (groupKey: number, imageSources: Photo[], baseApiUrl: string) => {
  const items = [];
  const start = 0;

  for (let i = 0; i < imageSources.length; i += 1) {
    const imageSource = imageSources[start + i];
    if (typeof imageSource !== 'undefined') {
      items.push(
        <ImageItem
          baseApiUrl={baseApiUrl}
          groupKey={groupKey}
          num={1 + start + i}
          key={start + i}
          source={imageSource}
        />
      );
    }
  }
  return items;
};

const getAlbumName = async(baseApiUrl: string, albumId: string) => {
  const albumsAdapter = new AlbumsAdapter(baseApiUrl);
  if (albumId === null) return '';
  const album = await albumsAdapter.getAlbumInfoById(albumId);
  return album.name;
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

  useEffect(() => {
    (async function retrieveAlbumName() {
      const retrievedAlbumName = await getAlbumName(props.baseApiUrl, id);
      setAlbumName(retrievedAlbumName);
    })();

    (async function retrievePhotos() {
      const retrievedPhotos = await getAllPhotos(props.baseApiUrl, id);
      const photoItems = loadItems(0, retrievedPhotos, props.baseApiUrl);
      setPhotos(photoItems);
    })();
  },        [setPhotos, id, props.baseApiUrl]);

  return (
    <MuiThemeProvider>
      <MainContent>
        <TitleBar />
        <Typography variant="h4" component="h1">
          {albumName}
        </Typography>
        <div className="photoIndexView__photoIndex">
          <JustifiedLayout
            options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
            layoutOptions={{ margin: 5, column: [0, 5] }}
            // onAppend={onAppend}
            onLayoutComplete={onLayoutComplete}
          >
            {photos}
          </JustifiedLayout>
        </div>
      </MainContent>
    </MuiThemeProvider>
  );
};
