import React from 'react';
import './PhotoDetail.css';
import {Photo} from '../../Models/Photo';
import {Loader, Table} from '@mantine/core';
import {object} from "prop-types";

interface IProps {
    photo: Photo;
}

const isJson = (str: string) => {
    try {
        return JSON.parse(str) && !!str;
    } catch (e) {
        return false;
    }
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {

    const tableRow = (key: string, value: string) => {
        return (
            <Table.Tr>
                <Table.Td>{key}</Table.Td>
                <Table.Td>{value}</Table.Td>
            </Table.Tr>
        );
    }

    const valueType = (value: any) => {
        if (isJson(value)) return 'json'
        if (typeof value === 'object' || Array.isArray(value)) return 'object'
        if (typeof value === 'string') return 'string'
    }

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

            // if (isJson(value)) {
            //
            // } else if (typeof value === 'object') {
            //
            // } else {
            //     return tableRow(key, value);
            // }
        })
    }

    const content = () => {
        if (props.photo !== undefined) {
            return (
                <div className="photoDetail__contentWrapper">
                    <img
                        src={props.photo.sourcePath}
                        alt={props.photo.name}
                        className="photoDetail__mainImage"
                    />
                    <div className="photoDetail__imageMetadata">
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
                    </div>
                </div>
            )
                ;
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
