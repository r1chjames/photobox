import React, {useEffect, useState} from 'react';
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
  if (photos && photos.length > 0) {
    return `${baseApiUrl}/photo/${photos[0].id}/thumbnail`;
  }
  return '';
};

const getPhotoCountInAlbum = async(photosAdapter: PhotosAdapter, albumId: string) => {
  return photosAdapter.getPhotoCountInAlbum(albumId);
};

export const AlbumItem: React.FunctionComponent<IProps> = (props) => {

  const [thumbnailUrl, setThumbnailUrl] = useState('');
  const [photoCount, setPhotoCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const photosAdapter = new PhotosAdapter(props.baseApiUrl);
    (async function retrieveThumbnailUrl() {
      const retrievedThumbnailUrl = await getUrlOfFirstImageInAlbum(photosAdapter, props.baseApiUrl, props.source.id);
      setThumbnailUrl(retrievedThumbnailUrl);
      setLoading(false);
    })();

    (async function retrieveAPhotoCount() {
      const retrievedPhotoCount: PhotoCount = await getPhotoCountInAlbum(photosAdapter, props.source.id);
      setPhotoCount(retrievedPhotoCount.photoCount);
      setLoading(false);
    })();
  },        [setThumbnailUrl, setPhotoCount, props.baseApiUrl, props.source.id]);

  const content = () => {
    if (loading) {
      return (
        <LoadingScreen />
      );
    }
    return (
      <Card
        className="albumItem__cardWrapper"
        onClick={() => props.albumViewCallback(props.source.id)}
        elevation={0}
      >
        <CardActionArea>
          <CardMedia
            className="albumItem__cardImage"
            component="img"
            src={thumbnailUrl}
            alt="No thumbnail"
            onError={(e: any) => e.target.src = '/no_image.png'}
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
  };

  return content();
};
