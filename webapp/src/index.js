import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './App.tsx';
import {BrowserRouter} from "react-router-dom";

const baseApiUrl = "http://127.0.0.1:8080/api";

ReactDOM.render(
    <React.StrictMode>
        <BrowserRouter>
            <App baseApiUrl={""}/>
        </BrowserRouter>
    </React.StrictMode>
    , document.getElementById('root')
);

serviceWorker.unregister();
