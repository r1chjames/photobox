import React from 'react';
import ReactDOM from 'react-dom';
import {
  BrowserRouter,
  Routes,
  Route
} from "react-router-dom";
import './index.css';
import * as serviceWorker from './serviceWorker';
import {Dashboard} from "./Components/Dashboard/Dashboard";
import {PhotoIndexView} from "./Components/PhotoIndexView/PhotoIndexView";
import {PhotoDetail} from "./Components/PhotoDetail/PhotoDetail";
import {AlbumIndexView} from "./Components/AlbumIndexView/AlbumIndexView";
import {CreateAlbumView} from "./Components/CreateAlbumView/CreateAlbumView";
import {SettingsView} from "./Components/SettingsView/SettingsView";

const baseApiUrl = "http://127.0.0.1:8080/api";

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={<Dashboard baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/photos"
          element={<PhotoIndexView baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/photo/:id"
          element={<PhotoDetail baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/albums"
          element={<AlbumIndexView baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/album/:id"
          element={<PhotoIndexView baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/album/new/:name"
          element={<CreateAlbumView baseApiUrl={baseApiUrl}/>}
        />
        <Route
          path="/settings"
          element={<SettingsView baseApiUrl={baseApiUrl}/>}
        />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);

serviceWorker.unregister();
