import React, { ComponentProps } from 'react';
import { Meta, Story } from '@storybook/react';

import { InfoSnackbar } from '../Components/Snackbar/InfoSnackbar';

export default {
  title: 'Snackbar',
  component: InfoSnackbar
} as Meta;

const template: Story<ComponentProps<typeof InfoSnackbar>> = args => <InfoSnackbar {...args} />;

export const primary = template.bind({});
primary.args = {
  text: 'Snackbar text',
  show: true,
  handleStopShowing: () => console.log()
};
