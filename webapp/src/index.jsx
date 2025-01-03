import React from 'react';
import ReactDOM from "react-dom/client";
import './index.css';
import App from './App.tsx';
import '@mantine/core/styles.css';
import {BrowserRouter} from "react-router-dom";

const baseApiUrl = "http://photobox:8080/api";

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
    <React.StrictMode>
        <BrowserRouter>
            <App baseApiUrl={baseApiUrl}/>
        </BrowserRouter>
    </React.StrictMode>
);