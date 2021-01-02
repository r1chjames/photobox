import React, { useEffect, useState } from 'react';
import './AlbumIndexView.css';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import history from '../../Routing/History';
import { MainContent } from '../MainContent/MainContent';
import { Fab, TextField } from '@material-ui/core';
import AddIcon from '@material-ui/icons/Add';
import { InputModal } from '../InputModal/InputModal';

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
    <MuiThemeProvider>
      <MainContent title="Albums">
        <InputModal
          isOpen={showNewAlbumModal}
          title="Create Album"
          handleSave={newAlbumModalSaveClick}
          handleClose={() => setShowNewAlbumModal(false)}
        >
          <TextField
            required={true}
            id="album-name"
            label="Album Name"
            variant="outlined"
            error={createAlbumModalAlbumNameErrorText.length !== 0}
            helperText={createAlbumModalAlbumNameErrorText}
            onChange={e => handleNewAlbumNameValueChange(e.target.value)}
          />
        </InputModal>
        <section className="albumIndexView__cardContainer">
          {albums.map((album: Album) => {
            return(
              <article key={album.id} className="albumIndexView__card">
                <AlbumItem
                  baseApiUrl={props.baseApiUrl}
                  source={album}
                  albumViewCallback={() => history.push(`album/${album.id}`)}
                />
              </article>
            );
          })}
        </section>
        <Fab color="primary" aria-label="add" className="albumIndexView__addButton">
          <AddIcon onClick={() => handleCreateNewAlbum()} />
        </Fab>
      </MainContent>
    </MuiThemeProvider>
  );
};
