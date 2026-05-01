import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { Center, Loader, Title, Text, Image, Badge, Group } from '@mantine/core';
import { IconCalendar, IconPhotoOff } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';
import { Album } from '../../Models/Album';

interface SharedResourceData {
  share: {
    token: string;
    resourceType: 'photo' | 'album';
    resourceId: string;
    viewCount: number;
    createdAt: string;
  };
  resource: Photo | Album;
}

export const SharedView: React.FC = () => {
  const { token } = useParams<{ token: string }>();
  const [data, setData] = useState<SharedResourceData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchShared = async () => {
      try {
        const apiBase = import.meta.env.VITE_API_BASE_URL || '/api';
        const response = await fetch(`${apiBase}/shared/${token}/resource`);
        if (response.status === 410) {
          setError('This shared link has expired.');
          return;
        }
        if (response.status === 401) {
          setError('This shared link requires a password.');
          return;
        }
        if (!response.ok) {
          setError('Failed to load shared content.');
          return;
        }
        const result = await response.json();
        if (result.success) {
          setData(result.data);
        } else {
          setError(result.message || 'Failed to load shared content.');
        }
      } catch (e) {
        setError('Failed to load shared content.');
      } finally {
        setLoading(false);
      }
    };
    fetchShared();
  }, [token]);

  if (loading) {
    return (
      <Center h="100vh">
        <Loader size="lg" />
      </Center>
    );
  }

  if (error) {
    return (
      <Center h="100vh">
        <div style={{ textAlign: 'center' }}>
          <IconPhotoOff size="3rem" style={{ marginBottom: 16 }} />
          <Title size="h4">{error}</Title>
        </div>
      </Center>
    );
  }

  if (!data) return null;

  const isPhoto = data.share.resourceType === 'photo';
  const resource = data.resource as Photo;
  const album = data.resource as Album;

  return (
    <div style={{ maxWidth: 1200, margin: '0 auto', padding: 24 }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 24 }}>
        <Title size="h3">Shared {isPhoto ? 'Photo' : 'Album'}</Title>
        <Text size="sm" c="dimmed">{data.share.viewCount} views</Text>
      </div>

      {isPhoto ? (
        <div>
          <Image
            src={`${import.meta.env.VITE_API_BASE_URL || '/api'}/photo/bin/${resource.id}`}
            alt={resource.name}
            radius="md"
            style={{ maxHeight: '70vh', objectFit: 'contain' }}
          />
          <Group gap="xs" mt="md">
            {resource.createdAt && (
              <Badge leftSection={<IconCalendar size={14} />} variant="light" color="blue">
                {new Date(resource.createdAt).toLocaleString()}
              </Badge>
            )}
          </Group>
          <Text mt="md" fw={500}>{resource.name}</Text>
        </div>
      ) : (
        <div>
          <Text fw={500} size="lg" mb="md">{album.name}</Text>
          <Text c="dimmed" mb="md">{album.description || ''}</Text>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: 16 }}>
            {/* Album photos would need to be fetched separately; showing placeholder for now */}
            <Center h={200} style={{ border: '1px dashed var(--mantine-color-gray-4)', borderRadius: 8 }}>
              <Text size="sm" c="dimmed">Album photos view coming soon</Text>
            </Center>
          </div>
        </div>
      )}
    </div>
  );
};
