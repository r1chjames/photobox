import type { StorybookConfig } from "@storybook/react-vite";
import { mergeConfig } from 'vite';

const config: StorybookConfig = {
  stories: ["../src/**/*.mdx", "../src/**/*.stories.@(js|jsx|mjs|ts|tsx)"],
  addons: ["@storybook/addon-links", "@storybook/addon-docs"],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
  async viteFinal(config) {
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