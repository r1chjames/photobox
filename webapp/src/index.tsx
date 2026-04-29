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
