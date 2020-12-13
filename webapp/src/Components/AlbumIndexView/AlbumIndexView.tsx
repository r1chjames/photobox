import React, { useEffect, useState } from 'react';
import './AlbumIndexView.css';
import { TitleBar } from '../TitleBar/TitleBar';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RaisedButton from 'material-ui/RaisedButton';
import history from '../../Routing/History';
import { MainContent } from '../MainContent/MainContent';

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
  const [renderedAlbumId, setRenderedAlbumId] = useState('');

  useEffect(() => {
    (async function retrieveAllAlbums() {
      const retrievedAlbums = await getAllAlbums(props.baseApiUrl);
      setAlbums(retrievedAlbums);
    })();
  },        [setAlbums, setRenderedAlbumId, props.baseApiUrl]);

  const renderAlbumIndex = () => {
    return (
      <div>
        <TitleBar />
          <MainContent title="Albums">
            <section className="albumIndexView__cardContainer">
              {albums.map((album: Album) => {
                return(
                  <article key={album.id} className="albumIndexView__card">
                    <AlbumItem
                      baseApiUrl={props.baseApiUrl}
                      source={album}
                      albumViewCallback={() => history.push(`albums/${album.id}`)}
                    />
                  </article>
                );
              })}
            </section>
          </MainContent>
        </div>
    );
  };

  const renderAlbumPhotoView = () => {
    return (
      <div>
        <RaisedButton onClick={() => setRenderedAlbumId('')} className="albumIndexView__button">
          Back
        </RaisedButton>
        <PhotoIndexView
          baseApiUrl={props.baseApiUrl}
        />
      </div>
    );
  };

  const content = () => {
    if (renderedAlbumId !== '') {
      return renderAlbumPhotoView();
    }
    return renderAlbumIndex();
  };

  return (
    <MuiThemeProvider>
      {content()}
    </MuiThemeProvider>
  );
};
