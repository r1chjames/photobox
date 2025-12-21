import '@mantine/core/styles.css';

import {MantineProvider} from '@mantine/core';
import {MemoryRouter} from "react-router";
import {theme} from "../src/theme";

export const decorators = [
    (Story) => (
        <MemoryRouter initialEntries={['/']}>
            <Story/>
        </MemoryRouter>
    ),
    (renderStory: any) => (
        <MantineProvider theme={theme}>{renderStory()}</MantineProvider>
    ),
];
export const tags = ['autodocs'];