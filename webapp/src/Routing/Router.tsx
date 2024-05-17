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
import {SettingsAdapter} from "../Adapters/SettingsAdapter";
import {ViewContainer} from "../Components/ViewContainer/ViewContainer";

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
                    element={
                        <ViewContainer>
                            <Dashboard
                                albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/photos"
                    element={
                        <ViewContainer>
                            <PhotoIndexView
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/photo/:id"
                    element={
                        <ViewContainer>
                            <PhotoDetail
                                baseApiUrl={this.props.baseApiUrl}
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/albums"
                    element={
                        <ViewContainer>
                            <AlbumIndexView
                                albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/album/:id"
                    element={
                        <ViewContainer>
                            <PhotoIndexView
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/album/new/:name"
                    element={
                        <ViewContainer>
                            <CreateAlbumView
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
                <Route
                    path="/settings"
                    element={
                        <ViewContainer>
                            <SettingsView
                                settingsAdapter={new SettingsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ViewContainer>
                    }
                />
            </Routes>
        </BrowserRouter>
    );
  };
}
