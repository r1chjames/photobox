import React, { useEffect, useState, useCallback } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import { LatLngBounds } from 'leaflet';
import { IPhotosAdapter, PhotoGeoData } from '../../Adapters/IPhotosAdapter';
import { EmptyState } from '../EmptyState/EmptyState';
import { Loader, Center, Title } from '@mantine/core';
import { IconMap } from '@tabler/icons-react';
import 'leaflet/dist/leaflet.css';
import { fetchThumbnailWithAuth, getCachedThumbnail, revokeThumbnail } from '../../utils/ThumbnailUtils';

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

export const MapView: React.FC<MapViewProps> = ({ photosAdapter }) => {
  const [photos, setPhotos] = useState<PhotoGeoData[]>([]);
  const [thumbnailUrls, setThumbnailUrls] = useState<Map<string, string>>(new Map());
  const [loading, setLoading] = useState(true);

  const loadGeodata = useCallback(async () => {
    try {
      const data = await photosAdapter.getGeodata(90, -90, 180, -180);
      setPhotos(data || []);

      const urls = new Map<string, string>();
      for (const photo of data || []) {
        const cached = getCachedThumbnail(photo.id);
        if (cached) {
          urls.set(photo.id, cached);
        } else {
          const url = await fetchThumbnailWithAuth(photosAdapter, photo.id);
          urls.set(photo.id, url);
        }
      }
      setThumbnailUrls(urls);
    } catch (e) {
      setPhotos([]);
    } finally {
      setLoading(false);
    }
  }, [photosAdapter]);

  useEffect(() => {
    loadGeodata();
    return () => {
      thumbnailUrls.forEach((_, id) => revokeThumbnail(id));
    };
  }, [loadGeodata]);

  if (loading) {
    return (
      <Center h="50vh">
        <Loader size="lg" />
      </Center>
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

  return (
    <div style={{ height: 'calc(100vh - 160px)', width: '100%' }}>
      <Title size="h4" mb="md">Map ({photos.length} photos)</Title>
      <MapContainer
        style={{ height: '100%', width: '100%', borderRadius: 8 }}
        center={[photos[0].lat, photos[0].lng]}
        zoom={3}
        scrollWheelZoom={true}
      >
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        />
        <BoundsSetter bounds={bounds} />
        {photos.map(photo => (
          <Marker key={photo.id} position={[photo.lat, photo.lng]}>
            <Popup>
              <div style={{ textAlign: 'center' }}>
                {thumbnailUrls.get(photo.id) && (
                  <img
                    src={thumbnailUrls.get(photo.id)}
                    alt=""
                    style={{ width: 120, height: 120, objectFit: 'cover', borderRadius: 4 }}
                  />
                )}
                <div style={{ fontSize: 12, marginTop: 4 }}>
                  {photo.dateTaken ? new Date(photo.dateTaken).toLocaleDateString() : ''}
                </div>
              </div>
            </Popup>
          </Marker>
        ))}
      </MapContainer>
    </div>
  );
};
