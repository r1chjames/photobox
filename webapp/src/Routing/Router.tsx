import React, {Suspense} from 'react';
import {Navigate, Route, Routes} from "react-router-dom";
import {useAuth} from "./AuthContext";
import {useAdapters} from "./AdapterContext";
import {AppBar, Labels} from "../Components/AppBar/AppBar";
import {Loader, Center} from "@mantine/core";

const PhotoGrid = React.lazy(() => import('../Components/PhotoGrid/PhotoGrid').then(m => ({default: m.PhotoGrid})));
const Dashboard = React.lazy(() => import('../Components/Dashboard/Dashboard').then(m => ({default: m.Dashboard})));
const AlbumGrid = React.lazy(() => import('../Components/AlbumGrid/AlbumGrid').then(m => ({default: m.AlbumGrid})));
const SettingsView = React.lazy(() => import('../Components/SettingsView/SettingsView').then(m => ({default: m.SettingsView})));
const CreateAlbumView = React.lazy(() => import('../Components/CreateAlbumView/CreateAlbumView').then(m => ({default: m.CreateAlbumView})));
const LoginCard = React.lazy(() => import('../Components/LoginCard/LoginCard').then(m => ({default: m.LoginCard})));
const PhotoDetail = React.lazy(() => import('../Components/PhotoDetail/PhotoDetail').then(m => ({default: m.PhotoDetail})));
const SearchView = React.lazy(() => import('../Components/SearchView/SearchView').then(m => ({default: m.SearchView})));
const FavoritesView = React.lazy(() => import('../Components/FavoritesView/FavoritesView').then(m => ({default: m.FavoritesView})));
const TrashView = React.lazy(() => import('../Components/TrashView/TrashView').then(m => ({default: m.TrashView})));
const ShareManagement = React.lazy(() => import('../Components/ShareManagement/ShareManagement').then(m => ({default: m.ShareManagement})));
const UserManagement = React.lazy(() => import('../Components/UserManagement/UserManagement').then(m => ({default: m.UserManagement})));
const MapView = React.lazy(() => import('../Components/MapView/MapView').then(m => ({default: m.MapView})));

const PageLoader = () => (
    <Center h="100vh">
        <Loader size="lg"/>
    </Center>
);

const ProtectedRoute = (props: { children: React.ReactNode }) => {
    const {isAuthenticated} = useAuth();
    if (!isAuthenticated) {
        return <Navigate to={"/login"}/>;
    }

    return <>{props.children}</>;
};

const Router: React.FunctionComponent = () => {
    const {photosAdapter, albumsAdapter, settingsAdapter, usersAdapter, sharesAdapter} = useAdapters();


    return (
        <Suspense fallback={<PageLoader/>}>
            <Routes>
                <Route
                    path="/"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Dashboard}>
                                <Dashboard
                                    albumsAdapter={albumsAdapter}
                                    photosAdapter={photosAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/login"
                    element={
                        <LoginCard
                            usersAdapter={usersAdapter}
                        />
                    }
                />
                <Route
                    path="/photos"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Photos}>
                                <PhotoGrid
                                    photosAdapter={photosAdapter}
                                    albumsAdapter={albumsAdapter}
                                    sharesAdapter={sharesAdapter}
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
                                <PhotoDetail photosAdapter={photosAdapter} sharesAdapter={sharesAdapter}/>
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
                                    albumsAdapter={albumsAdapter}
                                    photosAdapter={photosAdapter}
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
                                    photosAdapter={photosAdapter}
                                    albumsAdapter={albumsAdapter}
                                    sharesAdapter={sharesAdapter}
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
                                photosAdapter={photosAdapter}
                            />
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/search"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Photos}>
                                <SearchView
                                    photosAdapter={photosAdapter}
                                    albumsAdapter={albumsAdapter}
                                    sharesAdapter={sharesAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/favorites"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Favorites}>
                                <FavoritesView
                                    photosAdapter={photosAdapter}
                                    albumsAdapter={albumsAdapter}
                                    sharesAdapter={sharesAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/trash"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Trash}>
                                <TrashView
                                    photosAdapter={photosAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/shares"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Settings}>
                                <ShareManagement
                                    sharesAdapter={sharesAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/users"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Settings}>
                                <UserManagement
                                    usersAdapter={usersAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/map"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Map}>
                                <MapView
                                    photosAdapter={photosAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
                <Route
                    path="/settings"
                    element={
                        <ProtectedRoute>
                            <AppBar activeLink={Labels.Settings}>
                                <SettingsView
                                    settingsAdapter={settingsAdapter}
                                    photosAdapter={photosAdapter}
                                />
                            </AppBar>
                        </ProtectedRoute>
                    }
                />
            </Routes>
        </Suspense>
    );
};

export default Router;
