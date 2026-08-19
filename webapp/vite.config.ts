import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { VitePWA } from 'vite-plugin-pwa'

// https://vitejs.dev/config/
export default defineConfig({
    base: '/',
    plugins: [
        react(),
        VitePWA({
            registerType: 'autoUpdate',
            manifest: {
                name: 'Photobox',
                short_name: 'Photobox',
                description: 'Self-hosted photo library',
                theme_color: '#228be6',
                background_color: '#ffffff',
                display: 'standalone',
                start_url: '/',
                icons: [
                    {
                        src: '/logo192.png',
                        sizes: '192x192',
                        type: 'image/png',
                    },
                    {
                        src: '/logo512.png',
                        sizes: '512x512',
                        type: 'image/png',
                    },
                ],
            },
            workbox: {
                runtimeCaching: [
                    // Photo thumbnails: cache-first so the grid renders offline
                    // (serve from cache, update in background).
                    {
                        urlPattern: /\/api\/photo\/thumbnail\/.*/i,
                        handler: 'CacheFirst',
                        options: {
                            cacheName: 'photobox-thumbnails',
                            expiration: {
                                maxEntries: 500,
                                maxAgeSeconds: 60 * 60 * 24 * 30, // 30 days
                            },
                            cacheableResponse: {
                                statuses: [0, 200],
                            },
                        },
                    },
                    // Photo originals (viewed): cache-on-view with LRU eviction
                    // so recently-viewed photos are available offline.
                    {
                        urlPattern: /\/api\/photo\/bin\/.*/i,
                        handler: 'CacheFirst',
                        options: {
                            cacheName: 'photobox-originals',
                            expiration: {
                                maxEntries: 100,
                                maxAgeSeconds: 60 * 60 * 24 * 7, // 1 week
                            },
                            cacheableResponse: {
                                statuses: [0, 200],
                            },
                        },
                    },
                    // API responses: stale-while-revalidate so lists/albums
                    // render from cache while refreshing in the background.
                    {
                        urlPattern: /\/api\/.*/i,
                        handler: 'StaleWhileRevalidate',
                        options: {
                            cacheName: 'photobox-api-cache',
                            expiration: {
                                maxEntries: 100,
                                maxAgeSeconds: 60 * 60 * 24 * 7, // 1 week
                            },
                            cacheableResponse: {
                                statuses: [0, 200],
                            },
                        },
                    },
                    // Other cross-origin (tiles, etc.): network-first fallback.
                    {
                        urlPattern: /^https:\/\/.*\/.*/i,
                        handler: 'NetworkFirst',
                        options: {
                            cacheName: 'photobox-network-cache',
                            expiration: {
                                maxEntries: 100,
                                maxAgeSeconds: 60 * 60 * 24 * 7,
                            },
                            cacheableResponse: {
                                statuses: [0, 200],
                            },
                        },
                    },
                ],
            },
        }),
    ],
    build: {
        outDir: 'build',
        sourcemap: true,
    },
    test: {
        setupFiles: ['./src/test/setup.ts'],
        environment: 'jsdom',
        poolOptions: {
            threads: {
                maxThreads: 2,
            },
        },
        alias: [
            {
                find: '@tabler/icons-react',
                replacement: '@tabler/icons-react/dist/cjs/tabler-icons-react.cjs',
            },
        ],
        server: {
            deps: {
                inline: ['@tabler/icons-react'],
            },
        },
        deps: {
            optimizer: {
                web: {
                    exclude: ['@tabler/icons-react'],
                },
            },
        },
    },
})
