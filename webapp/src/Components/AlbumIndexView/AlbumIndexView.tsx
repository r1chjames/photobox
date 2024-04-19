import React, { useEffect, useState } from 'react';
import './AlbumIndexView.css';
import { ImplAlbumsAdapter } from '../../Adapters/ImplAlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import { InputModal } from '../InputModal/InputModal';
import {Button, MantineProvider, TextInput} from '@mantine/core';
import {MdAddCircle} from 'react-icons/md';
import {useNavigate} from "react-router-dom";
import {MockAlbumsAdapter} from "../../Adapters/MockAlbumsAdapter";

interface IProps {
  albumsAdapter: ImplAlbumsAdapter | MockAlbumsAdapter;
}

const getAllAlbums = async(propsAlbumsAdapter: ImplAlbumsAdapter | MockAlbumsAdapter) => {
  const albumsAdapter = propsAlbumsAdapter;
  const albumSources: Album[] = await albumsAdapter.getAllAlbumsInfo();
  return albumSources;
};

export const AlbumIndexView: React.FunctionComponent<IProps> = (props) => {
  const navigate = useNavigate()

  const [albums, setAlbums] = useState<Album[]>([]);
  const [showNewAlbumModal, setShowNewAlbumModal] = useState(Boolean);
  const [createAlbumModalAlbumNameErrorText, setCreateAlbumModalAlbumNameErrorText] = useState('Required');
  const [newAlbumName, setNewAlbumName] = useState('');

  useEffect(() => {
    (async function retrieveAllAlbums() {
      const retrievedAlbums = await getAllAlbums(props.albumsAdapter);
      setAlbums(retrievedAlbums);
    })();
  },[setAlbums, props.albumsAdapter]);

  const handleNewAlbumNameValueChange = (fieldValue: string) => {
    if (fieldValue.length > 1) {
      setCreateAlbumModalAlbumNameErrorText('');
      setNewAlbumName(fieldValue);
    } else {
      setCreateAlbumModalAlbumNameErrorText('Invalid length');
    }
  };
  const handleCreateNewAlbum = () => {
    setShowNewAlbumModal(true);
  };

  const newAlbumModalSaveClick = () => {
    navigate(`/album/new/${newAlbumName}`);
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
                  baseApiUrl={"TODO update this"}
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
