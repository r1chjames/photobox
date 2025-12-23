import React from 'react';
import {AlbumGrid} from '../AlbumGrid/AlbumGrid';
import {PhotoGrid} from '../PhotoGrid/PhotoGrid';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Space, Title} from "@mantine/core";

interface IProps {
    albumsAdapter: IAlbumsAdapter;
    photosAdapter: IPhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

    return (
        <>
            <div>
                <Title size="h4">Albums</Title>
                <Space h="md" />
                <AlbumGrid
                    albumsAdapter={props.albumsAdapter}
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={20}
                />
            </div>
            <div>
                <Space h="md" />
                <Title size="h4">Photos</Title>
                <Space h="md" />
                <PhotoGrid
                    photosAdapter={props.photosAdapter}
                    albumsAdapter={props.albumsAdapter}
                    maxDisplayed={50}
                />
            </div>
        </>
    );
};
