import React from 'react';
import ReactDOM from 'react-dom';
import {
    BrowserRouter,
    Routes,
    Route
} from "react-router-dom";import './index.css';
import * as serviceWorker from './serviceWorker';
import App from "./App";
import {Dashboard} from "./Components/Dashboard/Dashboard";
import {PhotoIndexView} from "./Components/PhotoIndexView/PhotoIndexView";
import {PhotoDetail} from "./Components/PhotoDetail/PhotoDetail";
import {AlbumIndexView} from "./Components/AlbumIndexView/AlbumIndexView";
import {CreateAlbumView} from "./Components/CreateAlbumView/CreateAlbumView";
import {SettingsView} from "./Components/SettingsView/SettingsView";
// import Router from './Routing/Router';

const baseApiUrl = "https://127.0.0.1:8080/api";

ReactDOM.render(
    <React.StrictMode>
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
                    element={() => <AlbumIndexView baseApiUrl={this.props.baseApiUrl}/>}
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
        <App />
    </React.StrictMode>,
    document.getElementById('root')
);
// <BrowserRouter>
    //     {/*<div className="App">*/}
    //         <Routes>
    //             <Route>
    //                 baseApiUrl={baseApiUrl}
    //             </Route>
    //         </Routes>
    //     {/*</div>*/}
    // </BrowserRouter>
    // ,document.getElementById('root'));

serviceWorker.unregister();
