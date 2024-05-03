import React, {Component} from 'react';
import {
  BrowserRouter,
  Routes,
  Route
} from "react-router-dom";
import {PhotoIndexView} from '../Components/PhotoIndexView/PhotoIndexView';
import {Dashboard} from '../Components/Dashboard/Dashboard';
import {AlbumIndexView} from '../Components/AlbumIndexView/AlbumIndexView';
import {SettingsView} from '../Components/SettingsView/SettingsView';
import {PhotoDetail} from '../Components/PhotoDetail/PhotoDetail';
import {CreateAlbumView} from '../Components/CreateAlbumView/CreateAlbumView';
import {AlbumsAdapter} from "../Adapters/AlbumsAdapter";
import {RestApiAdapter} from "../Adapters/RestApiAdapter";
import {PhotosAdapter} from "../Adapters/PhotosAdapter";

interface IProps {
  baseApiUrl: string;
}

export default class Router extends Component<IProps> {

  public render = () => {
    return (
      <BrowserRouter>
        <Routes>
          <Route
            path="/"
            element={() => <Dashboard baseApiUrl={this.props.baseApiUrl}/>}
          />
          <Route
            path="/photos"
            element={() => <PhotoIndexView baseApiUrl={this.props.baseApiUrl}/>}
          />
          <Route
            path="/photo/:id"
            element={() => <PhotoDetail baseApiUrl={this.props.baseApiUrl}/>}
          />
          <Route
            path="/albums"
            element={() => <AlbumIndexView
                albumsAdapter={new AlbumsAdapter(new RestApiAdapter(""))}
                photosAdapter={new PhotosAdapter(new RestApiAdapter(""))}
            />}
          />
          <Route
            path="/album/:id"
            element={() => <PhotoIndexView baseApiUrl={this.props.baseApiUrl}/>}
          />
          <Route
            path="/album/new/:name"
            element={() => <CreateAlbumView baseApiUrl={this.props.baseApiUrl}/>}
          />
          <Route
            path="/settings"
            element={() => <SettingsView baseApiUrl={this.props.baseApiUrl}/>}
          />
        </Routes>
      </BrowserRouter>
    );
  };
}
