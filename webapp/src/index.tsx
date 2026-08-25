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

// One-time migration (issue #164): the legacy hand-written service worker at
// /service-worker.js was cache-first with no update logic, freezing browsers on
// the build that was current when it first installed. It is removed from the
// app; unregister any surviving registration and drop its cache so
// already-poisoned browsers self-heal on the next load. The workbox worker at
// /sw.js (registered by the injected /registerSW.js) is the only SW from now on.
if ('serviceWorker' in navigator) {
    const legacyScript = '/service-worker.js';
    navigator.serviceWorker
        .getRegistrations()
        .then((registrations) =>
            Promise.all(
                registrations
                    .filter(
                        (r) =>
                            r.active?.scriptURL.endsWith(legacyScript) ||
                            r.waiting?.scriptURL.endsWith(legacyScript) ||
                            r.installing?.scriptURL.endsWith(legacyScript),
                    )
                    .map((r) => r.unregister()),
            ),
        )
        .catch(() => {
            // Non-critical cleanup; ignore failures.
        });
}

// Remove the legacy worker's cache ('photobox-v2') left behind by old builds.
if ('caches' in window) {
    caches.delete('photobox-v2').catch(() => {});
}
