import React from 'react';
import './AlbumItem.css';
import { Album } from '../../Models/Album';
import {Badge, Text, Card, Group, Image, Loader} from '@mantine/core';
import useAlbumItem from "./useAlbumItem";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
  photosAdapter: IPhotosAdapter;
  source: Album;
  albumViewCallback: (albumId: string) => void;
}

export const AlbumItem: React.FunctionComponent<IProps> = (props) => {

  const [{thumbnailUrl, photoCount, isLoading}] = useAlbumItem(props.photosAdapter, props.source);

  const content = () => {
    if (isLoading) {
      return (
        <Loader size={"md"} />
      );
    }
    return (
      <Card shadow="sm" padding="lg" radius="md" withBorder
            style={{ width: '200px' , margin: '10px' }}
        // className="albumItem__cardWrapper"
        onClick={() => props.albumViewCallback(props.source.id)}
      >
        <Card.Section>
          <Image src={thumbnailUrl} />
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

    {/*  className="albumItem__cardImage"*/}
    {/*  component="img"*/}
    {/*  src={thumbnailUrl}*/}
    {/*  alt="No thumbnail"*/}
    {/*  onError={(e: any) => e.target.src = '/no_image.png'}*/}
    {/*  title={props.source.name}*/}
    {/*/>*/}


  return content();
};
