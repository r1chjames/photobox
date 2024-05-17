import React, {StrictMode} from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './App.tsx';
import Router from "./Routing/Router";

const baseApiUrl = "http://127.0.0.1:8080/api";

ReactDOM.render(
    <StrictMode>
        <Router baseApiUrl={baseApiUrl}>
            <App />
        </Router>
    </StrictMode>
,document.getElementById('root')
);

serviceWorker.unregister();
