import React from 'react';
import ReactDOM from "react-dom/client";
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './App.tsx';
import {BrowserRouter} from "react-router-dom";

const baseApiUrl = "https://photodev.r1chjames.co.uk/api";

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
    <React.StrictMode>
        <BrowserRouter>
            <App baseApiUrl={baseApiUrl}/>
        </BrowserRouter>
    </React.StrictMode>
);

serviceWorker.unregister();