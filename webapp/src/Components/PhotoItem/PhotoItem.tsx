import React, { useState } from 'react';
import { Photo } from '../../Models/Photo';
import {Button, Dialog, Group, Text, Image} from '@mantine/core';
import {useNavigate} from "react-router-dom";

interface IProps {
  src: string;
  thumbnail: string;
  num: number;
  groupKey?: number;
  source: Photo;
  allPhotos: Photo[];
}

const findImageIndexInArray = (photos: Photo[], photoId: string, defaultVal: number) => {
  for (let i = 0; i < photos.length; i += 1) {
    if (photos[i].id === photoId) {
      return i;
    }
  }
  return defaultVal;
};

export const PhotoItem: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()
  const [currentId, setCurrentId] = useState(props.source.id);
  const [isImageModalOpen, setImageModalOpen] = useState(false);

  const currentPhotoIdIndex = (): number => {
    return findImageIndexInArray(props.allPhotos, currentId, 0);
  };

  const previousPhotoId = (currentPhotoIndex: number) => setCurrentId(props.allPhotos![currentPhotoIndex - 1].id);
  const nextPhotoId = (currentPhotoIndex: number) => setCurrentId(props.allPhotos![currentPhotoIndex + 1].id);

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
    const currentIdIndex = currentPhotoIdIndex() ;
    if (currentIdIndex !== 0) {
      return (
        <Button onClick={() => previousPhotoId(currentIdIndex)} color="primary">
          Previous
        </Button>
      );
    }
  };

  const nextButton = () => {
    const currentIdIndex = currentPhotoIdIndex() ;
    if (currentIdIndex !== props.allPhotos.length - 1) {
      return (
        <Button onClick={() => nextPhotoId(currentIdIndex)} color="primary">
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
                  onClick={() => navigate(`/photo/${currentId}`)}
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
