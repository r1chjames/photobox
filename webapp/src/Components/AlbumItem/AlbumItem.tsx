import React, { useEffect, useState } from 'react';
import './AlbumItem.css';
import { Album } from '../../Models/Album';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import { LoadingScreen } from '../LoadingScreen/LoadingScreen';
import { Card, CardActionArea, CardContent, CardMedia, Typography } from '@material-ui/core';

interface IProps {
  baseApiUrl: string;
  source: Album;
  albumViewCallback: (albumId: string) => void;
}

const getUrlOfFirstImageInAlbum = async(baseApiUrl: string, albumId: string) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(albumId);
  return `${baseApiUrl}/photo/bin?photoId=${photos[0].id}`;
};

export const AlbumItem: React.FunctionComponent<IProps> = (props) => {

  const [thumbnailUrl, setThumbnailUrl] = useState('');

  useEffect(() => {
    (async function retrieveThumbnailUrl() {
      const retrievedThumbnailUrl = await getUrlOfFirstImageInAlbum(props.baseApiUrl, props.source.id);
      setThumbnailUrl(retrievedThumbnailUrl);
    })();
  });

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
              <Typography variant="subtitle1" component="body">
                {props.source.name}
              </Typography>
              <Typography variant="caption" component="body">
                50 Photos
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
