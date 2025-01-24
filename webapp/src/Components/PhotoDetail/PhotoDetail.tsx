import React, {useEffect, useState} from 'react';
import {Photo} from '../../Models/Photo';
import {useDisclosure} from '@mantine/hooks';
import {Button, Dialog, Group, Image, Loader, ScrollArea, Table} from '@mantine/core';
import {valueType} from "../../utils/TypeUtils.js";
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {useParams} from "react-router-dom";
import {fetchPhotoBinWithAuth} from "../../utils/ImageUtils";

interface IProps {
    photosAdapter: IPhotosAdapter;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const {id} = useParams();
    const [photo, setPhoto] = useState<Photo>();
    const [photoUrl, setPhotoUrl] = useState<string>();
    const [opened, {toggle, close}] = useDisclosure(true);
    const img: React.Ref<HTMLImageElement> = React.createRef();

    const fetchPhoto = async (id: string) => {
        const photo = await props.photosAdapter.getPhotoInfoById(id);
        setPhoto(photo);
    }

    const fetchPhotoBin = async (id: string) => {
        const fetchedPhotoUrl = await fetchPhotoBinWithAuth(props.photosAdapter, id);
        setPhotoUrl(fetchedPhotoUrl);
    }

    useEffect(() => {
        if (id !== undefined) {
            fetchPhoto(id);
            fetchPhotoBin(id);
        }
    }, []);

    const tableRow = (key: string, value: string) =>
        <Table.Tr>
            <Table.Td>{key}</Table.Td>
            <Table.Td>{value}</Table.Td>
        </Table.Tr>

    const buildRows = (metadata: Record<string, any>) => {
        return Object.entries(metadata)
            .filter(([key, value]) => (typeof value !== 'undefined') && (![...Array(100).keys()].map(v => v.toString()).includes(key)) && (value !== ""))
            .map(([key, value]) => {
                switch (valueType(value)) {
                    case ('json'):
                        return buildRows(JSON.parse(value));
                    case ('object'):
                        return buildRows(value);
                    case ('string'):
                        return tableRow(key, value);
                    default:
                        return;
                }
            })
    };

    return (
        photo && photoUrl ?
            <>
                <Image
                    radius={"md"}
                    ref={img}
                    // mah="100"
                    mah="600px"
                    // maw={"100%"}
                    fit="scale-down"
                    src={photoUrl}
                />
                <Group justify="center">
                    <Button onClick={toggle} mt={50}>Metadata</Button>
                </Group>
                <Dialog opened={opened} withCloseButton onClose={close} size="lg" radius="md" mah="50%"
                        position={{top: "30%", right: 50, bottom: 50}}>
                    <ScrollArea h={400}>
                        <Table>
                            <Table.Thead>
                                <Table.Tr>
                                    <Table.Th>Parameter</Table.Th>
                                    <Table.Th>Value</Table.Th>
                                </Table.Tr>
                            </Table.Thead>
                            <Table.Tbody>
                                {buildRows(photo.metadata)}
                            </Table.Tbody>
                        </Table>
                    </ScrollArea>
                </Dialog>
            </>
            : <Loader size={"md"}/>
    );
}
