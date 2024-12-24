import React from 'react';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
    albumsAdapter: IAlbumsAdapter;
    photosAdapter: IPhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <>
        <AlbumIndexView
            albumsAdapter={props.albumsAdapter}
            photosAdapter={props.photosAdapter}
            maxDisplayed={20}
        />
        <PhotoIndexView
            photosAdapter={props.photosAdapter}
            maxDisplayed={50}
        />
      </>
  );
};
