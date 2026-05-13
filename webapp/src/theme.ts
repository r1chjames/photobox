import { createTheme, MantineColorsTuple } from '@mantine/core';

const indigoDark: MantineColorsTuple = [
    '#edf2ff',
    '#dbe4ff',
    '#bac8ff',
    '#91a7ff',
    '#748ffc',
    '#5c7cfa',
    '#4c6ef5',
    '#4263eb',
    '#3b5bdb',
    '#364fc7',
];

export const theme = createTheme({
    fontFamily: 'BlinkMacSystemFont, Segoe UI, Monaco, sans-serif',
    fontFamilyMonospace: 'Monaco, Courier, monospace',
    headings: { fontFamily: 'Greycliff CF, sans-serif' },
    defaultRadius: 'md',
    primaryColor: 'indigo',
    colors: {
        indigo: indigoDark,
    },
});
