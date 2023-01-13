import * as React from "react";
import { ComponentMeta, ComponentStory } from "@storybook/react";
import { TitleBar } from './TitleBar';

export default {
    component: TitleBar,
} as ComponentMeta<typeof TitleBar>;

export const Primary: ComponentStory<typeof TitleBar> = () => (
        <TitleBar />
        );
