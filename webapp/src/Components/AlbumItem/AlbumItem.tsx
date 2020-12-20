import React, { useEffect, useState } from 'react';
import './AlbumItem.css';
import { Album } from '../../Models/Album';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo, PhotoCount } from '../../Models/Photo';
import { LoadingScreen } from '../LoadingScreen/LoadingScreen';
import { Card, CardActionArea, CardContent, CardMedia, Typography } from '@material-ui/core';

interface IProps {
  baseApiUrl: string;
  source: Album;
  albumViewCallback: (albumId: string) => void;
}

const getUrlOfFirstImageInAlbum = async(photosAdapter: PhotosAdapter, baseApiUrl: string, albumId: string) => {
  const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(albumId, 1, 1);
  return `${baseApiUrl}/photo/${photos[0].id}/thumbnail`;
};

const getPhotoCountInAlbum = async(photosAdapter: PhotosAdapter, albumId: string) => {
  return photosAdapter.getPhotoCountInAlbum(albumId);
};

export const AlbumItem: React.FunctionComponent<IProps> = (props) => {

  const [thumbnailUrl, setThumbnailUrl] = useState('');
  const [photoCount, setPhotoCount] = useState(0);

  useEffect(() => {
    const photosAdapter = new PhotosAdapter(props.baseApiUrl);
    (async function retrieveThumbnailUrl() {
      const retrievedThumbnailUrl = await getUrlOfFirstImageInAlbum(photosAdapter, props.baseApiUrl, props.source.id);
      setThumbnailUrl(retrievedThumbnailUrl);
    })();

    (async function retrieveAPhotoCount() {
      const retrievedPhotoCount: PhotoCount = await getPhotoCountInAlbum(photosAdapter, props.source.id);
      setPhotoCount(retrievedPhotoCount.photoCount);
    })();
  },        [setThumbnailUrl, setPhotoCount, props.baseApiUrl, props.source.id]);

  const content = () => {
    if (thumbnailUrl !== '') {
      return (
        <Card
          className="albumItem__cardWrapper"
          onClick={() => props.albumViewCallback(props.source.id)}
          elevation={0}
        >
          <CardActionArea>
            <CardMedia
              className="albumItem__cardImage"
              image={thumbnailUrl}
              title={props.source.name}
            />
            <CardContent>
              <Typography variant="subtitle1" component="div">
                {props.source.name}
              </Typography>
              <Typography variant="caption" component="div">
                {photoCount} photos
              </Typography>
            </CardContent>
          </CardActionArea>
        </Card>
      );
    }
    return (
      <LoadingScreen />
    );
  };

  return content();
};
