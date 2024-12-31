import React from 'react';
import {Photo} from '../../Models/Photo';
import {useDisclosure} from '@mantine/hooks';
import {Button, Dialog, Group, Image, Loader, Table} from '@mantine/core';
import {valueType} from "../../utils/TypeUtils.js";

interface IProps {
    photo: Photo;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
    const [opened, {toggle, close}] = useDisclosure(true);

    const tableRow = (key: string, value: string) => {
        return (
            <Table.Tr>
                <Table.Td>{key}</Table.Td>
                <Table.Td>{value}</Table.Td>
            </Table.Tr>
        );
    };

    const buildRows = (metadata: Record<string, any>) => {
        return Object.entries(metadata).map(([key, value]) => {
            switch (valueType(value)) {
                case ('json'):
                    return buildRows(JSON.parse(value));
                case ('object'):
                    return tableRow(Object.entries(value).find(e => typeof e !== 'undefined')[0], Object.entries(value).find(e => typeof e !== 'undefined')[1]);
                case ('string'):
                    return tableRow(key, value);
                default:
                    return;
            }
        })
    };

    const content = () => {
        if (props.photo !== undefined) {
            return (
                <>
                    <Image
                        radius={"md"}
                        src={props.photo.sourcePath}
                    />
                    <Group justify="center">
                        <Button onClick={toggle} mt={50}>Metadata</Button>
                    </Group>
                    <Dialog opened={opened} withCloseButton onClose={close} size="lg" radius="md"
                            position={{top: "30%", right: 50}}>
                        <Table>
                            <Table.Thead>
                                <Table.Tr>
                                    <Table.Th>Parameter</Table.Th>
                                    <Table.Th>Value</Table.Th>
                                </Table.Tr>
                            </Table.Thead>
                            <Table.Tbody>
                                {props.photo.metadata.map(metadataElement => buildRows(metadataElement))}
                            </Table.Tbody>
                        </Table>
                    </Dialog>
                </>
            );
        }

        return (
            <Loader size={"md"}/>
        );
    };

    return (
        <div>
            {content()}
        </div>
    );
}
