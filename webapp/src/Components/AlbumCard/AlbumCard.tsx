import React from 'react';
import { Album } from '../../Models/Album';
import {Badge, Text, Card, Group, Image, Loader} from '@mantine/core';
import useAlbumCard from "./useAlbumCard";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
  photosAdapter: IPhotosAdapter;
  source: Album;
  albumViewCallback: (albumId: string) => void;
}

export const AlbumCard: React.FunctionComponent<IProps> = (props) => {

  const [{thumbnailUrl, photoCount, isLoading}] = useAlbumCard(props.photosAdapter, props.source);

  const content = () => {
    if (isLoading) {
      return (
        <Loader size={"md"} />
      );
    }
    return (
      <Card shadow="sm" padding="lg" radius="md" withBorder
            style={{ width: '200px' , margin: '10px'}}
        onClick={() => props.albumViewCallback(props.source.id)}
      >
        <Card.Section>
          <Image src={thumbnailUrl} wah={"200px"} mah={"140px"} />
        </Card.Section>
        <Group justify="space-between" mt="md" mb="xs">
          <Text fw={500}>{props.source.name}</Text>
          <Badge color="blue" variant="light" size="xs">
            {photoCount} photos
          </Badge>
        </Group>
      </Card>
    );
  };

  return content();
};
