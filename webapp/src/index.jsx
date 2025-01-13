import React from 'react';
import ReactDOM from "react-dom/client";
import './index.css';
import App from './App.tsx';
import '@mantine/core/styles.css';
import {BrowserRouter} from "react-router-dom";
import './appGlobals';

const root = ReactDOM.createRoot(document.getElementById("root"));

root.render(
    <React.StrictMode>
        <BrowserRouter>
            <App baseApiUrl={globalThis.app.baseApiUrl}/>
        </BrowserRouter>
    </React.StrictMode>
);