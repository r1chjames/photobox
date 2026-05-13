import React, {Suspense} from 'react';
import {Navigate, Outlet, Route, Routes, useLocation} from "react-router-dom";
import {useAuth} from "./AuthContext";
import {useAdapters} from "./AdapterContext";
import {AppBar} from "../Components/AppBar/AppBar";
import {Loader, Center} from "@mantine/core";
import { AnimatePresence, motion } from 'framer-motion';

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
const TagsView = React.lazy(() => import('../Components/TagsView/TagsView').then(m => ({default: m.TagsView})));
const TagPhotosView = React.lazy(() => import('../Components/TagsView/TagPhotosView').then(m => ({default: m.TagPhotosView})));
const SharedView = React.lazy(() => import('../Components/SharedView/SharedView').then(m => ({default: m.SharedView})));
const DuplicatesView = React.lazy(() => import('../Components/DuplicatesView/DuplicatesView').then(m => ({default: m.DuplicatesView})));

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

const AnimatedOutlet = () => {
    const location = useLocation();
    return (
        <AnimatePresence mode="wait">
            <motion.div
                key={location.pathname}
                initial={{ opacity: 0, y: 8 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -8 }}
                transition={{ duration: 0.2 }}
                style={{ height: '100%' }}
            >
                <Outlet />
            </motion.div>
        </AnimatePresence>
    );
};

const AppShellLayout = () => {
    return (
        <AppBar>
            <Suspense fallback={<PageLoader/>}>
                <AnimatedOutlet />
            </Suspense>
        </AppBar>
    );
};

const Router: React.FunctionComponent = () => {
    const {photosAdapter, albumsAdapter, settingsAdapter, usersAdapter, sharesAdapter} = useAdapters();


    return (
        <Routes>
            <Route path="/login" element={<LoginCard usersAdapter={usersAdapter} />} />
            <Route path="/shared/:token" element={<SharedView />} />
            <Route element={<ProtectedRoute><AppShellLayout /></ProtectedRoute>}>
                <Route path="/" element={<Dashboard albumsAdapter={albumsAdapter} photosAdapter={photosAdapter} />} />
                <Route path="/photos" element={<PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />} />
                <Route path="/photo/:id" element={<PhotoDetail photosAdapter={photosAdapter} sharesAdapter={sharesAdapter}/>} />
                <Route path="/albums" element={<AlbumGrid albumsAdapter={albumsAdapter} photosAdapter={photosAdapter} />} />
                <Route path="/album/:id" element={<PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />} />
                <Route path="/album/new/:name" element={<CreateAlbumView photosAdapter={photosAdapter} />} />
                <Route path="/search" element={<SearchView photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />} />
                <Route path="/favorites" element={<FavoritesView photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />} />
                <Route path="/trash" element={<TrashView photosAdapter={photosAdapter} />} />
                <Route path="/shares" element={<ShareManagement sharesAdapter={sharesAdapter} />} />
                <Route path="/users" element={<UserManagement usersAdapter={usersAdapter} />} />
                <Route path="/map" element={<MapView photosAdapter={photosAdapter} />} />
                <Route path="/settings" element={<SettingsView settingsAdapter={settingsAdapter} photosAdapter={photosAdapter} />} />
                <Route path="/tags" element={<TagsView photosAdapter={photosAdapter} />} />
                <Route path="/tags/:tag" element={<TagPhotosView photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} />} />
                <Route path="/duplicates" element={<DuplicatesView photosAdapter={photosAdapter} />} />
                <Route path="/videos" element={<PhotoGrid photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} sharesAdapter={sharesAdapter} mediaType="video" />} />
            </Route>
        </Routes>
    );
};

export default Router;
