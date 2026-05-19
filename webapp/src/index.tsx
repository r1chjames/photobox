import React from 'react';
import ReactDOM from "react-dom/client";
import './index.css';
import App from './App';
import '@mantine/core/styles.css';
import {BrowserRouter} from "react-router-dom";
import {QueryClient, QueryClientProvider} from "@tanstack/react-query";
import {AuthProvider} from "./Routing/AuthContext";

const root = ReactDOM.createRoot(document.getElementById("root")!);
const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            staleTime: 30_000,
            retry: 3,
            retryDelay: (attempt) => Math.min(1000 * 2 ** attempt, 30000),
        },
    },
});

const baseApiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080/api";

root.render(
    <React.StrictMode>
        <BrowserRouter>
            <AuthProvider>
                <QueryClientProvider client={queryClient}>
                    <App baseApiUrl={baseApiUrl}/>
                </QueryClientProvider>
            </AuthProvider>
        </BrowserRouter>
    </React.StrictMode>
);

// Register service worker for PWA support
if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/service-worker.js')
            .then(() => {
                console.log('Service Worker registered');
            })
            .catch((err) => {
                console.log('Service Worker registration failed', err);
            });
    });
}
