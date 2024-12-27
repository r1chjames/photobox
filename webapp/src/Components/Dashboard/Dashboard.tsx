import React from 'react';
import './Dashboard.css';
import {AlbumIndexView} from '../AlbumIndexView/AlbumIndexView';
import {PhotoIndexView} from '../PhotoIndexView/PhotoIndexView';
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
                <AlbumIndexView
                    albumsAdapter={props.albumsAdapter}
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={20}
                />
            </div>
            <div>
                <Title size="h4">Photos</Title>
                <PhotoIndexView
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={50}
                />
            </div>
        </>
    );
};
