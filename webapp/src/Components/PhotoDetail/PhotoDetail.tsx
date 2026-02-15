import React, {useEffect} from 'react';
import {useDisclosure} from '@mantine/hooks';
import {Button, Dialog, Group, Image, Loader, ScrollArea, Table} from '@mantine/core';
import {valueType} from "../../utils/TypeUtils.js";
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {useParams} from "react-router-dom";
import {fetchPhotoBinWithAuth, revokeBlobUrl} from "../../utils/ImageUtils";
import {useQuery} from "@tanstack/react-query";

interface IProps {
    photosAdapter: IPhotosAdapter;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const {id} = useParams();
    const [opened, {toggle, close}] = useDisclosure(true);

    const fetchPhoto = async () => {
        return props.photosAdapter.getPhotoInfoById(id!);
    }

    const fetchPhotoBin = async () => {
        return fetchPhotoBinWithAuth(props.photosAdapter, id!);
    }

    const {data: photo} = useQuery({
        queryKey: ['fetchPhoto', id],
        queryFn: fetchPhoto,
        enabled: !!id,
    });

    const {data: photoUrl} = useQuery({
        queryKey: ['fetchPhotoBin', id],
        queryFn: fetchPhotoBin,
        enabled: !!id,
    });

    useEffect(() => {
        return () => {
            revokeBlobUrl(photoUrl);
        };
    }, [photoUrl]);

    const tableRow = (key: string, value: string) =>
        <Table.Tr key={key}>
            <Table.Td>{key}</Table.Td>
            <Table.Td>{value}</Table.Td>
        </Table.Tr>

    const buildRows = (metadata: Record<string, any>): React.ReactNode => {
        if (!metadata) {
            return null;
        }
        return Object.entries(metadata)
            .filter(([key, value]) => (value !== undefined) && (![...Array(100).keys()].map(v => v.toString()).includes(key)) && (value !== ""))
            .map(([key, value]): React.ReactNode => {
                switch (valueType(value)) {
                    case ('json'):
                        return buildRows(JSON.parse(value));
                    case ('object'):
                        return buildRows(value);
                    case ('string'):
                        return tableRow(key, value);
                    default:
                        return null;
                }
            })
    };

    return (
        photo && photoUrl ?
            <>
                <Image
                    radius={"md"}
                    mah="600px"
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
                                {photo.metadata && buildRows(photo.metadata)}
                            </Table.Tbody>
                        </Table>
                    </ScrollArea>
                </Dialog>
            </>
            : <Loader size={"md"}/>
    );
}
