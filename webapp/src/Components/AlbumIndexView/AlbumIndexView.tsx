import React, { useEffect, useState } from 'react';
import './AlbumIndexView.css';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
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

  useEffect(() => {
    (async function retrieveAllAlbums() {
      const retrievedAlbums = await getAllAlbums(props.baseApiUrl);
      setAlbums(retrievedAlbums);
    })();
  },        [setAlbums, props.baseApiUrl]);

  return (
    <MuiThemeProvider>
      <MainContent title="Albums">
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
      </MainContent>
    </MuiThemeProvider>
  );
};
