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
                    {
                        urlPattern: /^https:\/\/.*\/.*/i,
                        handler: 'NetworkFirst',
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
                ],
            },
        }),
    ],
    resolve: {
        alias: {},
    },
    build: {
        outDir: 'build',
        sourcemap: true,
    },
    test: {
        setupFiles: ['./src/test/setup.ts'],
        environment: 'jsdom',
        alias: [
            {
                find: '@tabler/icons-react',
                replacement: require.resolve('@tabler/icons-react/dist/cjs/tabler-icons-react.cjs'),
            },
        ],
        server: {
            deps: {
                inline: ['@tabler/icons-react'],
            },
        },
    },
})
