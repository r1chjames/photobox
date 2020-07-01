import React from 'react';
import ReactDOM from 'react-dom';
import { BrowserRouter as Router } from 'react-router-dom';
import './index.css';
import * as serviceWorker from './serviceWorker';
import Routes from './Routing/Routes';

const baseApiUrl = "http://192.168.0.123:8080/api";

ReactDOM.render(
    <Router>
        <div className="App">
            {/*<Navigation />*/}
            <Routes baseApiUrl={baseApiUrl} />
        </div>
    </Router>
    ,document.getElementById('root'));

serviceWorker.unregister();
