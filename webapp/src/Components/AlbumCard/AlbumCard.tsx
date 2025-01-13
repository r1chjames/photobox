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
      <Card shadow="sm" radius="md" withBorder mb={'10px'} mr={'10px'} w={"200px"}
        onClick={() => props.albumViewCallback(props.source.id)}
      >
        <Card.Section>
          <Image src={`data:image/png;base64,${thumbnailUrl}`} wah={"200px"} h={"150px"} />
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
