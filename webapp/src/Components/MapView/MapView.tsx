import React, { useEffect, useState, useCallback } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import MarkerClusterGroup from 'react-leaflet-cluster';
import { LatLngBounds } from 'leaflet';
import 'leaflet/dist/leaflet.css';
import 'leaflet.markercluster/dist/MarkerCluster.css';
import 'leaflet.markercluster/dist/MarkerCluster.Default.css';
import 'leaflet.heat';
import L from 'leaflet';
import { IPhotosAdapter, PhotoGeoData } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Title, Skeleton, Text, SegmentedControl } from '@mantine/core';
import { IconMap } from '@tabler/icons-react';

// Leaflet's default marker icon URLs don't resolve under bundlers (they point
// at /marker-icon.png which 404s). Point them at the bundled images.
import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

L.Icon.Default.mergeOptions({
    iconRetinaUrl: markerIcon2x,
    iconUrl: markerIcon,
    shadowUrl: markerShadow,
});

interface MapViewProps {
  photosAdapter: IPhotosAdapter;
}

const BoundsSetter: React.FC<{ bounds: LatLngBounds | null }> = ({ bounds }) => {
  const map = useMap();
  useEffect(() => {
    if (bounds && bounds.isValid()) {
      map.fitBounds(bounds, { padding: [50, 50] });
    }
  }, [bounds, map]);
  return null;
};

// HeatmapLayer renders a leaflet.heat layer for the given points.
const HeatmapLayer: React.FC<{ points: [number, number][] }> = ({ points }) => {
  const map = useMap();
  useEffect(() => {
    const layer = (L as any).heatLayer(points, {
      radius: 25,
      blur: 15,
      maxZoom: 17,
    });
    layer.addTo(map);
    return () => {
      map.removeLayer(layer);
    };
  }, [map, points]);
  return null;
};

export const MapView: React.FC<MapViewProps> = ({ photosAdapter }) => {
  const [photos, setPhotos] = useState<PhotoGeoData[]>([]);
  const [loading, setLoading] = useState(true);
  const [mode, setMode] = useState<'markers' | 'heatmap'>('markers');

  const loadGeodata = useCallback(async () => {
    try {
      const data = await photosAdapter.getGeodata(90, -90, 180, -180);
      setPhotos(data || []);
    } catch {
      setPhotos([]);
    } finally {
      setLoading(false);
    }
  }, [photosAdapter]);

  useEffect(() => {
    loadGeodata();
  }, [loadGeodata]);

  if (loading) {
    return (
      <div>
        <Title size="h4" mb="md">Map</Title>
        <Skeleton height="calc(100vh - 160px)" radius="md" />
      </div>
    );
  }

  if (photos.length === 0) {
    return (
      <>
        <Title size="h4" mb="md">Map</Title>
        <EmptyState
          title="No geotagged photos"
          description="Photos with GPS coordinates in their EXIF data will appear here."
          icon={<IconMap size="2rem" />}
        />
      </>
    );
  }

  const bounds = new LatLngBounds(
    photos.map(p => [p.lat, p.lng] as [number, number])
  );
  const heatPoints = photos.map(p => [p.lat, p.lng] as [number, number]);

  return (
    <div style={{ height: 'calc(100vh - 160px)', width: '100%' }}>
      <Title size="h4" mb="md">Map ({photos.length} photos)</Title>
      <SegmentedControl
        value={mode}
        onChange={(v) => setMode(v as 'markers' | 'heatmap')}
        data={[
          { label: 'Markers', value: 'markers' },
          { label: 'Heatmap', value: 'heatmap' },
        ]}
        size="xs"
        mb="sm"
      />
      <MapContainer
        style={{ height: 'calc(100% - 40px)', width: '100%', borderRadius: 8 }}
        center={[photos[0].lat, photos[0].lng]}
        zoom={3}
        scrollWheelZoom={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <BoundsSetter bounds={bounds} />
        {mode === 'heatmap' ? (
          <HeatmapLayer points={heatPoints} />
        ) : (
          <MarkerClusterGroup chunkedLoading>
            {photos.map(photo => (
              <Marker key={photo.id} position={[photo.lat, photo.lng]}>
                <Popup>
                  <Text size="sm" fw={500}>
                    {photo.dateTaken ? new Date(photo.dateTaken).toLocaleDateString() : 'No date'}
                  </Text>
                </Popup>
              </Marker>
            ))}
          </MarkerClusterGroup>
        )}
      </MapContainer>
    </div>
  );
};
