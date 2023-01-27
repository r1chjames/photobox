import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";

import { PhotoItem } from './PhotoItem';
import { Photo } from '../../Models/Photo';

export default {
  component: PhotoItem
} as ComponentMeta<typeof PhotoItem>;

export const primary: ComponentStory<typeof PhotoItem> = (args) => (
        <PhotoItem {...args} />
        );

primary.args = {
  baseApiUrl: 'https://google.com/',
  num: 1,
  groupKey: 1,
  source: new Photo('1', 'test', '', 'a', '', {'': ''}, ''),
  allPhotos: [new Photo('1', 'test', '', 'a', '', {'': ''}, '')],
};
