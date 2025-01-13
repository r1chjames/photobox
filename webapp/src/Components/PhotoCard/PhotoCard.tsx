import React, {useEffect, useState} from 'react';
import {Photo} from '../../Models/Photo';
import {ActionIcon, Button, Card, Flex, Group, Image, Overlay, Text} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {IconArrowLeftDashed, IconArrowRightDashed, IconX} from "@tabler/icons-react";
import {useHotkeys} from "@mantine/hooks";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
    photosAdapter: IPhotosAdapter
    source: Photo;
    previousPhoto: () => void;
    nextPhoto: () => void;
    firstInAlbum: boolean;
    lastInAlbum: boolean;
    closeModal: () => void;
}

export const PhotoCard: React.FunctionComponent<IProps> = (props) => {
    const navigate = useNavigate()
    const [fetchedImage, setFetchedImage] = useState<string | undefined>();
    const img: React.Ref<HTMLImageElement> = React.createRef();

    const fetchImage = async () => {
        const data = await props.photosAdapter.getPhotoImage(props.source.id)
        const blob = new Blob([data], {
            type: 'image/jpeg',
        });
        const objectURL = URL.createObjectURL(blob)
        setFetchedImage(objectURL);
    }

    useEffect(() => {
        fetchImage();
    }, [props.source])

    useHotkeys([
        ['ArrowLeft', () => props.previousPhoto()],
        ['ArrowRight', () => props.nextPhoto()],
    ]);

    const previousButton = () => {
        if (!props.firstInAlbum) {
            return (
                <ActionIcon color="dark" size="xl">
                    <IconArrowLeftDashed size="2.125rem" onClick={() => props.previousPhoto()}/>
                </ActionIcon>
            );
        }
    };

    const nextButton = () => {
        if (!props.lastInAlbum) {
            return (
                <ActionIcon color="dark" size="xl">
                    <IconArrowRightDashed size="2.125rem" onClick={() => props.nextPhoto()}/>
                </ActionIcon>
            );
        }
    };

    return (
        <Card shadow="sm" radius="md" padding={"xs"}
              onClick={() => navigate(`/photo/${props.source.id}`)}>
            <Card.Section>
                {fetchedImage ?
                    <Image
                        h={"500px"}
                        fit={"cover"}
                        w={"auto"}
                        ref={img}
                        src={fetchedImage}
                        alt={props.source.name}
                    /> :
                    'Loading...'}
                <Overlay color="#000" backgroundOpacity={0} opacity={0.5}>
                    <Flex direction="row" style={{width: "100%", justifyContent: "right"}}>
                        <ActionIcon color="dark" size="l" opacity={1}>
                            <IconX size="1.75rem" onClick={() => props.closeModal()}/>
                        </ActionIcon>
                    </Flex>
                    <Flex direction="row" style={{
                        width: "100%",
                        height: "100%",
                        justifyContent: "space-between",
                        alignItems: "center"
                    }}>
                        {previousButton()}
                        {nextButton()}
                    </Flex>
                </Overlay>
            </Card.Section>
            <Group justify="space-between" mt="md" mb="xs">
                <Text fw={500}>{props.source.name}</Text>
                <Button onClick={() => navigate(`/photo/${props.source.id}`)} mt={50}>View</Button>
            </Group>
        </Card>
    );
};
