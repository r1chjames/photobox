import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'

export default defineConfig({
    base: '/',
    plugins: [react()],
    resolve: {
        alias: {
            // /esm/icons/index.mjs only exports the icons statically, so no separate chunks are created
            '@tabler/icons-react': '@tabler/icons-react/dist/esm/icons/index.mjs',
        },
    },
    build: {
        outDir: 'build',
        sourcemap: true,
    }
})