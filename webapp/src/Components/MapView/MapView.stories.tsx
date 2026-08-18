import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {MapView} from './MapView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

const photosAdapter = {
    getGeodata: async () => [
        {id: 'p1', lat: 52.04, lng: 0.094, dateTaken: '2023-01-01T00:00:00Z'},
        {id: 'p2', lat: 51.58, lng: -1.01, dateTaken: '2023-01-02T00:00:00Z'},
        {id: 'p3', lat: 49.95, lng: 0.61, dateTaken: '2023-01-03T00:00:00Z'},
        {id: 'p4', lat: 52.04, lng: 0.095, dateTaken: '2023-01-04T00:00:00Z'},
        {id: 'p5', lat: 52.041, lng: 0.093, dateTaken: '2023-01-05T00:00:00Z'},
    ],
} as unknown as IPhotosAdapter;

const emptyAdapter = {
    getGeodata: async () => [],
} as unknown as IPhotosAdapter;

const meta: Meta<typeof MapView> = {
    title: 'Views/MapView',
    component: MapView,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', height: 500}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof MapView>;

export const Populated: Story = {
    args: {
        photosAdapter,
    },
};

export const Empty: Story = {
    args: {
        photosAdapter: emptyAdapter,
    },
};
