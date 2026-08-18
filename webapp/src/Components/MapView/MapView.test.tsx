import React from 'react';
import {render, screen, waitFor, fireEvent} from '../../test/test-utils';
import {vi} from 'vitest';
import {MapView} from './MapView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

// Mock react-leaflet — leaflet needs a real DOM/Canvas that happy-dom lacks.
vi.mock('react-leaflet', () => ({
    MapContainer: ({children}: any) => <div data-testid="map-container">{children}</div>,
    TileLayer: () => <div data-testid="tile-layer" />,
    Marker: ({children}: any) => <div data-testid="marker">{children}</div>,
    Popup: ({children}: any) => <div data-testid="popup">{children}</div>,
    useMap: () => ({fitBounds: vi.fn(), removeLayer: vi.fn(), addLayer: vi.fn()}),
}));

vi.mock('react-leaflet-cluster', () => ({
    default: ({children}: any) => <div data-testid="cluster">{children}</div>,
}));

const photosAdapter = {
    getGeodata: vi.fn().mockResolvedValue([
        {id: 'p1', lat: 52.04, lng: 0.094, dateTaken: '2023-01-01T00:00:00Z'},
        {id: 'p2', lat: 51.58, lng: -1.01, dateTaken: '2023-01-02T00:00:00Z'},
    ]),
} as unknown as IPhotosAdapter;

describe('MapView', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders markers in cluster mode by default', async () => {
        render(<MapView photosAdapter={photosAdapter} />);

        await waitFor(() => {
            expect(screen.getByTestId('map-container')).toBeInTheDocument();
        });
        expect(screen.getByTestId('cluster')).toBeInTheDocument();
        expect(screen.getAllByTestId('marker').length).toBe(2);
        expect(screen.getByText(/2 photos/i)).toBeInTheDocument();
    });

    it('shows empty state when no geotagged photos', async () => {
        const emptyAdapter = {
            getGeodata: vi.fn().mockResolvedValue([]),
        } as unknown as IPhotosAdapter;

        render(<MapView photosAdapter={emptyAdapter} />);

        await waitFor(() => {
            expect(screen.getByText(/no geotagged photos/i)).toBeInTheDocument();
        });
    });

    it('toggles between markers and heatmap modes', async () => {
        render(<MapView photosAdapter={photosAdapter} />);

        await waitFor(() => {
            expect(screen.getByTestId('map-container')).toBeInTheDocument();
        });

        // Switch to heatmap
        fireEvent.click(screen.getByText('Heatmap'));
        // Heatmap mode should not render markers/cluster
        expect(screen.queryByTestId('cluster')).not.toBeInTheDocument();

        // Switch back to markers
        fireEvent.click(screen.getByText('Markers'));
        expect(screen.getByTestId('cluster')).toBeInTheDocument();
    });
});
