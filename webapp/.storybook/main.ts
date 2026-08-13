import type { StorybookConfig } from "@storybook/react-vite";
import { mergeConfig } from 'vite';

const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  addons: ["@storybook/addon-links", "@storybook/addon-docs", "@storybook/addon-viewport"],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
  async viteFinal(config) {
    // Exclude the PWA/service-worker plugin from Storybook builds. Its
    // workbox precache step fails on Storybook's large sb-manager bundle
    // (exceeds the 2 MiB maximumFileSizeToCacheInBytes default), and
    // Storybook output doesn't need a service worker at all.
    const flatten = (plugins: any): any[] => (Array.isArray(plugins) ? plugins.flat(Infinity) : []);
    config.plugins = flatten(config.plugins).filter(
      (plugin) => !(plugin && 'name' in plugin && typeof plugin.name === 'string' && plugin.name.startsWith('vite-plugin-pwa'))
    );

    // Merge custom configuration into the final config.
    return mergeConfig(config, {
      resolve: {
        alias: {
          '@tabler/icons-react': '@tabler/icons-react/dist/esm/icons/index.mjs',
        },
      },
    });
  },
};
export default config;