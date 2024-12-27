import React from 'react';
import './Dashboard.css';
import {AlbumGrid} from '../AlbumGrid/AlbumGrid';
import {PhotoGrid} from '../PhotoGrid/PhotoGrid';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Title} from "@mantine/core";

interface IProps {
    albumsAdapter: IAlbumsAdapter;
    photosAdapter: IPhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

    return (
        <>
            <div>
                <Title size="h4">Albums</Title>
                <AlbumGrid
                    albumsAdapter={props.albumsAdapter}
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={20}
                />
            </div>
            <div>
                <Title size="h4">Photos</Title>
                <PhotoGrid
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={50}
                />
            </div>
        </>
    );
};
