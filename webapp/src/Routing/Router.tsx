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

interface IProps {
    baseApiUrl: string;
}

interface IPRProps {
    children: React.ReactNode;
}

const ProtectedRoute = (props: IPRProps) => {
    const token = localStorage.getItem('token');
    if (token == null) {
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
                            <Dashboard
                                albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
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
                            <PhotoGrid
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/albums"
                    element={
                        <ProtectedRoute>
                            <AlbumGrid
                                albumsAdapter={new AlbumsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/album/:albumid"
                    element={
                        <ProtectedRoute>
                            <PhotoGrid
                                photosAdapter={new PhotosAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
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
                            <SettingsView
                                settingsAdapter={new SettingsAdapter(new RestApiAdapter(this.props.baseApiUrl))}
                            />
                        </ProtectedRoute>
                    }
                />
            </Routes>
        );
    };
}
