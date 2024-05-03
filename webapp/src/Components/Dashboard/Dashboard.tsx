import React from 'react';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import {RestApiAdapter} from "../../Adapters/RestApiAdapter";

interface IProps {
    baseApiUrl: string;
    albumsAdapter: AlbumsAdapter;
    photosAdapter: PhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <div>
        <AlbumIndexView
            albumsAdapter={new AlbumsAdapter(new RestApiAdapter(props.baseApiUrl))}
            photosAdapter={new PhotosAdapter(new RestApiAdapter(props.baseApiUrl))}
        />
        <PhotoIndexView
            baseApiUrl={props.baseApiUrl}
            photosAdapter={new PhotosAdapter(new RestApiAdapter(props.baseApiUrl))}
        />
      </div>
  );
};
