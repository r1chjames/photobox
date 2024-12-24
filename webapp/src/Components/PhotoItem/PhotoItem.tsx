import React, {useEffect, useState} from 'react';
import { Photo } from '../../Models/Photo';
import {Modal, Group, Text, Image, Card, Overlay, ActionIcon, Flex} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {useMediaQuery} from "@mantine/hooks";
import {IconArrowLeftDashed, IconArrowRightDashed, IconX} from "@tabler/icons-react";

interface IProps {
  groupKey?: number;
  key: number;
  source: Photo;
  photoSequence: number;
  previousPhoto: (currentIndex: number) => Photo;
  nextPhoto: (currentIndex: number) => Photo;
  firstInAlbum: boolean;
  lastInAlbum: boolean;
}

export const PhotoItem: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()
  const [currentPhoto, setCurrentPhoto] = useState(props.source);
  const [currentIndex, setCurrentIndex] = useState(props.photoSequence);
  const [isImageModalOpen, setImageModalOpen] = useState(false);
  const isMobile = useMediaQuery('(max-width: 50em)');

  useEffect(() => {
    document.addEventListener('keydown', keyPress);
    return ()=> {
      document.removeEventListener('keydown', keyPress);
    }
  }, []);

  const keyPress = (e: any) => {
    switch (e.which) {
      case 37: {
        // left
        setPreviousPhoto();
        break
      }
      case 39: {
        // right
        setNextPhoto();
        break
      }
      default:
    }
  }

  const gridImage = () => {
    return(
        <Image
            radius="md"
            fit="contain"
            p={5}
            src={currentPhoto.thumbnailPath}
            alt={currentPhoto.name}
            onClick={() => setImageModalOpen(!isImageModalOpen)}
            onError={(e: any) => e.target.src = '/no_image.png'}
        />
    );
  };

  const setPreviousPhoto = () => {
    setCurrentPhoto(props.previousPhoto(currentIndex))
    setCurrentIndex(currentIndex - 1);
  }

  const setNextPhoto = () => {
    setCurrentPhoto(props.nextPhoto(currentIndex))
    setCurrentIndex(currentIndex + 1);
  }

  const previousButton = () => {
    if (!props.firstInAlbum) {
      return (
          <ActionIcon color="dark" size="xl">
            <IconArrowLeftDashed size="2.125rem" onClick={() => setPreviousPhoto()}/>
          </ActionIcon>
      );
    }
  };

  const nextButton = () => {
    if (!props.lastInAlbum) {
      return (
          <ActionIcon color="dark" size="xl">
            <IconArrowRightDashed size="2.125rem" onClick={() => setNextPhoto()}/>
          </ActionIcon>
      );
    }
  };

  const content = () => {
    if (isImageModalOpen) {
      return (
          <>
            <Modal
                opened={isImageModalOpen}
                withCloseButton={false}
                aria-labelledby="customized-dialog-title"
                size="auto"
                padding={"0"}
                m={"0"}
                overlayProps={{
                  backgroundOpacity: 0.55,
                }}
                fullScreen={isMobile}
                transitionProps={{ transition: 'fade', duration: 200 }}
                onClose={() => setImageModalOpen(false)}>
              <Card shadow="sm" radius="md" padding={"xs"}
                    // style={{ width: '800px' }}
                    onClick={() => navigate(`/photo/${currentPhoto.id}`)}>
                <Card.Section>
                  <Image
                      h={"500px"}
                      fit={"cover"}
                      w={"auto"}
                      src={currentPhoto.sourcePath}
                      alt={currentPhoto.name}
                  />
                  <Overlay color="#000" backgroundOpacity={0} opacity={0.5}>
                    <Flex direction="row" style={{ width: "100%", justifyContent: "right" }}>
                      <ActionIcon color="dark" size="l" opacity={1}>
                        <IconX size="1.75rem" onClick={() => setImageModalOpen(false)}/>
                      </ActionIcon>
                    </Flex>
                    <Flex direction="row" style={{ width: "100%", height: "100%", justifyContent: "space-between", alignItems: "center" }}>
                      {previousButton()}
                      {nextButton()}
                    </Flex>
                  </Overlay>
                </Card.Section>
                <Group justify="space-between" mt="md" mb="xs">
                  <Text fw={500}>{currentPhoto.name}</Text>
                </Group>
              </Card>
            </Modal>
          { gridImage() }
        </>
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
