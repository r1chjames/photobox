import React, {useEffect} from 'react';
import {Photo} from '../../Models/Photo';
import {ActionIcon, Card, Flex, Group, Image, Overlay, Text} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import {IconArrowLeftDashed, IconArrowRightDashed, IconX} from "@tabler/icons-react";

interface IProps {
    source: Photo;
    previousPhoto: () => void;
    nextPhoto: () => void;
    firstInAlbum: boolean;
    lastInAlbum: boolean;
    closeModal: () => void;
}

export const PhotoCard: React.FunctionComponent<IProps> = (props) => {
    const navigate = useNavigate()

    useEffect(() => {
        document.addEventListener('keydown', keyPress);
        return () => {
            document.removeEventListener('keydown', keyPress);
        }
    }, []);

    const keyPress = (e: any) => {
        switch (e.which) {
            case 37: {
                // left
                props.previousPhoto();
                break
            }
            case 39: {
                // right
                props.nextPhoto();
                break
            }
            default:
        }
    }

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
                <Image
                    h={"500px"}
                    fit={"cover"}
                    w={"auto"}
                    src={props.source.sourcePath}
                    alt={props.source.name}
                />
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
            </Group>
        </Card>
    );
};
