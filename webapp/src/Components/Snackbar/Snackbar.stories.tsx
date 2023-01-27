import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";

import { InfoSnackbar } from './InfoSnackbar';

export default {
  component: InfoSnackbar
} as ComponentMeta<typeof InfoSnackbar>;

export const primary: ComponentStory<typeof InfoSnackbar> = (args) => (
        <InfoSnackbar {...args} />
        );

primary.args = {
  text: 'Snackbar text',
  show: true,
  handleStopShowing: () => console.log()
};
