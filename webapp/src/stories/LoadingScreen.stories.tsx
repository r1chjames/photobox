import React, { ComponentProps } from 'react';
import { Meta, Story } from '@storybook/react';
import { LoadingScreen } from './../Components/LoadingScreen/LoadingScreen';

export default {
  title: 'LoadingScreen',
  component: LoadingScreen
} as Meta;

const Template: Story<ComponentProps<typeof LoadingScreen>> = args => <LoadingScreen {...args} />;

export const Primary = Template.bind({});
