import React, {Component} from 'react';
import {Navigate, Route, Routes} from "react-router-dom";
import {PhotoGrid} from '../Components/PhotoGrid/PhotoGrid';
import {Dashboard} from '../Components/Dashboard/Dashboard';
import {AlbumGrid} from '../Components/AlbumGrid/AlbumGrid';
import {SettingsView} from '../Components/SettingsView/SettingsView';
import {CreateAlbumView} from '../Components/CreateAlbumView/CreateAlbumView';
import {AlbumsAdapter} from "../Adapters/AlbumsAdapter";
import {RestApiAdapter} from "../Adapters/RestApiAdapter";
import {PhotosAdapter} from "../Adapters/PhotosAdapter";
import {SettingsAdapter} from "../Adapters/SettingsAdapter";
import {LoginCard} from "../Components/LoginCard/LoginCard";
import {UsersAdapter} from "../Adapters/UsersAdapter";
import {AppBar, Labels} from "../Components/AppBar/AppBar";
import {PhotoDetail} from "../Components/PhotoDetail/PhotoDetail";

interface IProps {
    baseApiUrl: string;
}

const ProtectedRoute = (props: { children: React.ReactNode }) => {
    const token = localStorage.getItem('token');
    if (token === null) {
        return <Navigate to={"/login"}/>;
    }

    return <>{props.children}</>;
};

export default class Router extends Component<IProps> {

    public render = () => {
        return (
            <Routes>
                <Route
                    path="/"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Dashboard}>
                                <Dashboard
                                    albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                    photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/login"
                    element={
                        <LoginCard
                            usersAdapter={new UsersAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                        />
                    }
                />
                <Route
                    path="/photos"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Photos}>
                                <PhotoGrid
                                    photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                    albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/photo/:id"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Photos}>
                                <PhotoDetail photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}/>
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/albums"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Albums}>
                                <AlbumGrid
                                    albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                    photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/album/:id"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Albums}>
                                <PhotoGrid
                                    photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                    albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/album/new/:name"
                    element={
                        <ProtectedRoute>
                            <CreateAlbumView
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/settings"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Settings}>
                                <SettingsView
                                    settingsAdapter={new SettingsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                    photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
            </Routes>
        );
    };
}
