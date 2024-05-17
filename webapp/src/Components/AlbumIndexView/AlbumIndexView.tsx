import React, { useState } from 'react';
import './AlbumIndexView.css';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import { InputModal } from '../InputModal/InputModal';
import {Button, MantineProvider, TextInput} from '@mantine/core';
import {MdAddCircle} from 'react-icons/md';
import {useNavigate} from "react-router-dom";
import {AlbumsAdapter} from "../../Adapters/AlbumsAdapter";
import {PhotosAdapter} from "../../Adapters/PhotosAdapter";
import useAlbumIndexView from "./useAlbumIndexView";

interface IProps {
  albumsAdapter: AlbumsAdapter;
  photosAdapter: PhotosAdapter;
}

export const AlbumIndexView: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()
  const [showNewAlbumModal, setShowNewAlbumModal] = useState(Boolean);
  const [{albums, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}] = useAlbumIndexView(props.albumsAdapter);

  const newAlbumModalSaveClick = () => {
    navigate(`/album/new/${newAlbumName}`);
  };


  const handleCreateNewAlbum = () => {
    setShowNewAlbumModal(true);
  };

  return (
    <MantineProvider>
      <p>Albums</p>
        <InputModal
          isOpen={showNewAlbumModal}
          title="Create Album"
          handleSave={newAlbumModalSaveClick}
          handleClose={() => setShowNewAlbumModal(false)}
        >
          <TextInput
            required={true}
            id="album-name"
            label="Album Name"
            error={createAlbumModalAlbumNameErrorText.length !== 0}
            onChange={e => handleNewAlbumNameValueChange(e.target.value)}
          />
        </InputModal>
        <section className="albumIndexView__cardContainer">
          {albums.map((album: Album) => {
            return(
              <article key={album.id}>
                <AlbumItem
                  photosAdapter={props.photosAdapter}
                  source={album}
                  albumViewCallback={() => navigate(`album/${album.id}`)}
                />
              </article>
            );
          })}
        </section>
        <Button
          color="primary"
          aria-label="add"
          className="albumIndexView__addButton"
          onClick={() => handleCreateNewAlbum()}
        >
          <MdAddCircle/>
        </Button>
    </MantineProvider>
  );
};
