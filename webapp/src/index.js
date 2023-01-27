import React from 'react';
import ReactDOM from 'react-dom';
import { QueryClient, QueryClientProvider } from "react-query";
import { MantineProvider } from '@mantine/core';
import './index.css';
import * as serviceWorker from './serviceWorker';
import App from './App.tsx';

const baseApiUrl = "http://127.0.0.1:8080/api";
const queryClient = new QueryClient();

ReactDOM.render(
        <MantineProvider theme={{ loader: 'bars' }}  withGlobalStyles withNormalizeCSS>
            <QueryClientProvider client={queryClient}>
                <App baseApiUrl={baseApiUrl}/>
            </QueryClientProvider>
        </MantineProvider>
,document.getElementById('root')
);

serviceWorker.unregister();
