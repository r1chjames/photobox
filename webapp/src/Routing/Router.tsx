import React, {Component} from 'react';
import {Route, Routes} from "react-router-dom";
import {PhotoIndexView} from '../Components/PhotoIndexView/PhotoIndexView';
import {Dashboard} from '../Components/Dashboard/Dashboard';
import {AlbumIndexView} from '../Components/AlbumIndexView/AlbumIndexView';
import {SettingsView} from '../Components/SettingsView/SettingsView';
import {CreateAlbumView} from '../Components/CreateAlbumView/CreateAlbumView';
import {AlbumsAdapter} from "../Adapters/AlbumsAdapter";
import {RestApiAdapter} from "../Adapters/RestApiAdapter";
import {PhotosAdapter} from "../Adapters/PhotosAdapter";
import {SettingsAdapter} from "../Adapters/SettingsAdapter";

interface IProps {
    baseApiUrl: string;
}

export default class Router extends Component<IProps> {

    public render = () => {
        return (
            <Routes>
                <Route
                    path="/"
                    element={
                        <Dashboard
                            albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                <Route
                    path="/photos"
                    element={
                        <PhotoIndexView
                            photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                {/*<Route*/}
                {/*    path="/photo/:id"*/}
                {/*    element={*/}
                {/*        <PhotoDetail*/}
                {/*            // baseApiUrl={this.props.baseApiUrl}*/}
                {/*            photo={new Photo("p1", "photo 1", "/tmp/photo1.jpg", "https://4.img-dpreview.com/files/p/E~TS1180x0~articles/3925134721/0266554465.jpeg", "a1", "", {"a":"","b":""}, "")}*/}
                {/*            // photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}*/}
                {/*        />*/}
                {/*    }*/}
                {/*/>*/}
                <Route
                    path="/albums"
                    element={
                        <AlbumIndexView
                            albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                <Route
                    path="/album/:albumid"
                    element={
                        <PhotoIndexView
                            photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                <Route
                    path="/album/new/:name"
                    element={
                        <CreateAlbumView
                            photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                <Route
                    path="/settings"
                    element={
                        <SettingsView
                            settingsAdapter={new SettingsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
            </Routes>
        );
    };
}
