import type { Preview } from "@storybook/react";
import React from "react";
import { MemoryRouter } from "react-router";

export const decorators = [
  (Story) => (
      <MemoryRouter initialEntries={['/']}>
          <Story />
          </MemoryRouter>
  ),
];

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;
