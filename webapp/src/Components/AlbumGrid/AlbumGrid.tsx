import React, { useState } from 'react';
import { Album } from '../../Models/Album';
import { AlbumCard } from '../AlbumCard/AlbumCard';
import { InputModal } from '../InputModal/InputModal';
import {Button, Flex, TextInput} from '@mantine/core';
import {MdAddCircle} from 'react-icons/md';
import {useNavigate} from "react-router-dom";
import useAlbumGrid from "./useAlbumGrid";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
  albumsAdapter: IAlbumsAdapter;
  photosAdapter: IPhotosAdapter;
  maxDisplayed? : number;
}

const defaultProps = {
    maxDisplayed: 20000000,
}

export const AlbumGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const navigate = useNavigate()
    const [showNewAlbumModal, setShowNewAlbumModal] = useState(Boolean);
    const [{albums, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}] = useAlbumGrid(props.albumsAdapter);

  const newAlbumModalSaveClick = () => {
    navigate(`/album/new/${newAlbumName}`);
  };

  const handleCreateNewAlbum = () => {
    setShowNewAlbumModal(true);
  };

  return (
    <div>
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
          <Flex
              direction={{ base: 'column', sm: 'row' }}
              gap={{ base: 'lg', sm: 'lg' }}
              justify={{ sm: 'flex-start' }}
              mih={50}
              bg="rgba(0, 0, 0, 0)"
              align="flex-start"
              wrap="wrap"
          >
          {albums.slice(0, props.maxDisplayed).map((album: Album) => {
            return(
              <article key={album.id}>
                <AlbumCard
                  photosAdapter={props.photosAdapter}
                  source={album}
                  albumViewCallback={() => navigate(`album/${album.id}`)}
                />
              </article>
            );
          })}
          </Flex>
        </section>
        <Button
          color="primary"
          aria-label="add"
          className="albumIndexView__addButton"
          onClick={() => handleCreateNewAlbum()}
        >
          <MdAddCircle/>
        </Button>
    </div>
  );
};
