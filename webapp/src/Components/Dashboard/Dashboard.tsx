import React from 'react';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";

interface IProps {
    albumsAdapter: AlbumsAdapter;
    photosAdapter: PhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <div>
        <AlbumIndexView
            albumsAdapter={props.albumsAdapter}
            photosAdapter={props.photosAdapter}
        />
        <PhotoIndexView
            photosAdapter={props.photosAdapter}
        />
      </div>
  );
};
