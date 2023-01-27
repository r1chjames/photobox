import React from "react";
import {
    BrowserRouter,
    Routes,
    Route
} from "react-router-dom";
import {Dashboard} from "./Components/Dashboard/Dashboard";
import {PhotoIndexView} from "./Components/PhotoIndexView/PhotoIndexView";
import {PhotoDetail} from "./Components/PhotoDetail/PhotoDetail";
import {AlbumIndexView} from "./Components/AlbumIndexView/AlbumIndexView";
import {CreateAlbumView} from "./Components/CreateAlbumView/CreateAlbumView";
import {SettingsView} from "./Components/SettingsView/SettingsView";
import {ViewContainer} from "./Components/ViewContainer/ViewContainer";

interface IProps {
    baseApiUrl: string;
}

export const App: React.FunctionComponent<IProps> = (props) => {
    return (
        <React.StrictMode>
            <BrowserRouter>
                <Routes>
                    <Route
                        path="/"
                        element={
                            <ViewContainer>
                                <Dashboard baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/photos"
                        element={
                            <ViewContainer>
                                <PhotoIndexView baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/photo/:id"
                        element={
                            <ViewContainer>
                                <PhotoDetail baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/albums"
                        element={
                            <ViewContainer>
                                <AlbumIndexView baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/album/:id"
                        element={
                            <ViewContainer>
                                <PhotoIndexView baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/album/new/:name"
                        element={
                            <ViewContainer>
                                <CreateAlbumView baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                    <Route
                        path="/settings"
                        element={
                            <ViewContainer>
                                <SettingsView baseApiUrl={props.baseApiUrl}/>
                            </ViewContainer>
                        }
                    />
                </Routes>
            </BrowserRouter>
        </React.StrictMode>
    );
}

export default App;
