import React, { useState } from 'react';
import { Photo } from '../../Models/Photo';
import {Button, Dialog, Group, Text, Image} from '@mantine/core';
import {useNavigate} from "react-router-dom";

interface IProps {
  src: string;
  thumbnail: string;
  groupKey?: number;
  source: Photo;
  photoSequence: number;
  previousPhoto: () => void;
  nextPhoto: () => void;
  lastInAlbum: boolean;
}


export const PhotoItem: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()
  const [isImageModalOpen, setImageModalOpen] = useState(false);

  const gridImage = () => {
    return(
      <Image
          radius="md"
          h={100}
          w="auto"
          fit="contain"
          src={props.thumbnail}
          alt={props.source.name}
          onClick={() => setImageModalOpen(!isImageModalOpen)}
          onError={(e: any) => e.target.src = '/no_image.png'}
      />
    );
  };

  const previousButton = () => {
    if (props.photoSequence !== 0) {
      return (
        <Button onClick={() => props.previousPhoto()} color="primary">
          Previous
        </Button>
      );
    }
  };

  const nextButton = () => {
    if (!props.lastInAlbum) {
      return (
        <Button onClick={() => props.nextPhoto()} color="primary">
          Next
        </Button>
      );
    }
  };

  const content = () => {
    if (isImageModalOpen) {
      return (
        <div>
          <div>
            <Dialog
              opened={isImageModalOpen}
              aria-labelledby="customized-dialog-title"
              onClose={() => setImageModalOpen(false)}>
              <Text size="sm" style={{ marginBottom: 10 }}>
                {props.source.name}
              </Text>
              <Group>
                <Image
                  onClick={() => navigate(`/photo/${props.source.id}`)}
                  src={props.src}
                  alt={props.source.name}
                />
              </Group>
              <Group>
                {previousButton()}
                <Button onClick={() => setImageModalOpen(!isImageModalOpen)} color="primary">
                  Close
                </Button>
                {nextButton()}
              </Group>
            </Dialog>
          </div>
          { gridImage() }
        </div>
      );
    }

    return (
      <div>
        {gridImage()}
      </div>
    );
  };

  return content();
};
