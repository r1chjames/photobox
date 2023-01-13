import { BrowserRouter } from "react-router-dom";
import {addDecorator} from "@storybook/react";

addDecorator(story => <BrowserRouter initialEntries={['/']}>{story()}</BrowserRouter>);


export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
}

