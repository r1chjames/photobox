import React, { ComponentProps } from 'react';
import { Meta, Story } from '@storybook/react';

import { PhotoItem } from '../Components/PhotoItem/PhotoItem';
import { Photo } from '../Models/Photo';

export default {
  title: 'PhotoItem',
  component: PhotoItem
} as Meta;

const template: Story<ComponentProps<typeof PhotoItem>> = args => <PhotoItem {...args} />;

export const primary = template.bind({});
primary.args = {
  baseApiUrl: 'https://google.com/',
  num: 1,
  groupKey: 1,
  source: new Photo('1', 'test', '', 'a', '', {'': ''}, ''),
  allPhotos: [new Photo('1', 'test', '', 'a', '', {'': ''}, '')],
};
