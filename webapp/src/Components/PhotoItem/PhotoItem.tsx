import React, { useState } from 'react';
import { Photo } from '../../Models/Photo';
import {Modal, Group, Text, Image, Card, Overlay, Button, ActionIcon, Flex} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {useMediaQuery} from "@mantine/hooks";
import {IconArrowLeftDashed, IconArrowRightDashed} from "@tabler/icons-react";

interface IProps {
  src: string;
  thumbnail: string;
  groupKey?: number;
  key: number;
  source: Photo;
  photoSequence: number;
  previousPhoto: () => void;
  nextPhoto: () => void;
  lastInAlbum: boolean;
}

export const PhotoItem: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()
  const [isImageModalOpen, setImageModalOpen] = useState(false);
  const isMobile = useMediaQuery('(max-width: 50em)');

  const gridImage = () => {
    return(
      <Image
          radius="md"
          fit="contain"
          p={5}
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
          <>
            <Modal
                opened={isImageModalOpen}
                // withCloseButton={false}
                aria-labelledby="customized-dialog-title"
                size="auto"
                padding={"xs"}
                m={"xs"}
                overlayProps={{
                  backgroundOpacity: 0.55,
                }}
                fullScreen={isMobile}
                transitionProps={{ transition: 'fade', duration: 200 }}
                onClose={() => setImageModalOpen(false)}>
              <Card shadow="sm" radius="md" padding={"xs"}
                    style={{ width: '800px' }}
                  // className="albumItem__cardWrapper"
                    onClick={() => navigate(`/photo/${props.source.id}`)}>
                <Card.Section>
                  <Image
                      src={props.src}
                      alt={props.source.name}
                  />
                  <Overlay color="#000" backgroundOpacity={0} opacity={0.5}>
                    <Flex justify="left" align="center" direction="row">
                      <ActionIcon color="dark" size="xl">
                        <IconArrowLeftDashed size="2.125rem" onClick={() => props.previousPhoto()}/>
                      </ActionIcon>
                    </Flex>
                    <Flex justify="right" align="center" direction="row">
                      <ActionIcon color="dark" size="xl">
                        <IconArrowRightDashed size="2.125rem" onClick={() => props.nextPhoto()}/>
                      </ActionIcon>
                    </Flex>
                  </Overlay>
                </Card.Section>
                <Group justify="space-between" mt="md" mb="xs">
                  <Text fw={500}>{props.source.name}</Text>
                </Group>
            </Card>
              {/*<Group>*/}
              {/*  {previousButton()}*/}
              {/*  {nextButton()}*/}
              {/*</Group>*/}
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
