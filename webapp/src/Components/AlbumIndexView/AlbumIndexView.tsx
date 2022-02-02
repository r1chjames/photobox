import React, { useEffect, useState } from 'react';
import './AlbumIndexView.css';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import history from '../../Routing/History';
import { MainContent } from '../MainContent/MainContent';
import { InputModal } from '../InputModal/InputModal';
import {Button, MantineProvider, TextInput} from '@mantine/core';
import {MdAddCircle} from 'react-icons/md';

interface IProps {
  baseApiUrl: string;
}

const getAllAlbums = async(baseApiUrl: string) => {
  const albumsAdapter = new AlbumsAdapter(baseApiUrl);
  const albumSources: Album[] = await albumsAdapter.getAllAlbumsInfo();
  return albumSources;
};

export const AlbumIndexView: React.FunctionComponent<IProps> = (props) => {

  const [albums, setAlbums] = useState<Album[]>([]);
  const [showNewAlbumModal, setShowNewAlbumModal] = useState(Boolean);
  const [createAlbumModalAlbumNameErrorText, setCreateAlbumModalAlbumNameErrorText] = useState('Required');
  const [newAlbumName, setNewAlbumName] = useState('');

  useEffect(() => {
    (async function retrieveAllAlbums() {
      const retrievedAlbums = await getAllAlbums(props.baseApiUrl);
      setAlbums(retrievedAlbums);
    })();
  },        [setAlbums, props.baseApiUrl]);

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
    history.push(`/album/new/${newAlbumName}`);
  };

  return (
    <MantineProvider>
      <MainContent title="Albums">
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
                  baseApiUrl={props.baseApiUrl}
                  source={album}
                  albumViewCallback={() => history.push(`album/${album.id}`)}
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
      </MainContent>
    </MantineProvider>
  );
};
