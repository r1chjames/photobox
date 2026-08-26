import { createTheme, MantineColorsTuple } from '@mantine/core';

// Photobox brand blue (#3b82f6) ramp — matches the Open Design prototype's
// --accent token at index 5 (Mantine's "filled" shade).
const photoboxBlue: MantineColorsTuple = [
    '#eaf2ff',
    '#d3e3ff',
    '#a9c8ff',
    '#7cabff',
    '#5492fb',
    '#3b82f6',
    '#2f6fe0',
    '#255cc0',
    '#1f4fa6',
    '#1a418b',
];

// Favorite heart rose (#f43f5e) ramp for favorite states.
const photoboxRose: MantineColorsTuple = [
    '#ffe9ec',
    '#ffd3d9',
    '#ffa5b1',
    '#fd7587',
    '#fa4d63',
    '#f43f5e',
    '#d92b4c',
    '#bf1f41',
    '#a5163a',
    '#8c0f33',
];

export const theme = createTheme({
    fontFamily: "'DM Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    fontFamilyMonospace: "'JetBrains Mono', Monaco, 'Courier New', monospace",
    headings: {
        fontFamily: "'Space Grotesk', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
        fontWeight: '600',
    },
    defaultRadius: 'md',
    primaryColor: 'photobox-blue',
    colors: {
        'photobox-blue': photoboxBlue,
        'photobox-rose': photoboxRose,
    },
    radius: {
        sm: '8px',
        md: '12px',
        lg: '18px',
        xl: '22px',
    },
});
