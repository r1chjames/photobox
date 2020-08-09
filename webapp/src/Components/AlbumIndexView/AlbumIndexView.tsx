import React, { Component } from 'react';
import './AlbumIndexView.css';
import { TitleBar } from '../TitleBar/TitleBar';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { AlbumItem } from '../AlbumItem/AlbumItem';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RaisedButton from 'material-ui/RaisedButton';

interface IProps {
  baseApiUrl: string;
}

interface IState {
  albums: Album[];
  renderedAlbumId: string | null;
}

export class AlbumIndexView extends Component<IProps, IState> {

  constructor(props: IProps) {
    super(props);
    this.state = { albums: [], renderedAlbumId: null };
  }

  private async getAllAlbums() {
    const albumsAdapter = new AlbumsAdapter(this.props.baseApiUrl);
    const albumSources: Album[] = await albumsAdapter.getAllAlbumsInfo();
    this.setState({ albums: albumSources });
  }

  public async componentDidMount() {
    await this.getAllAlbums();
  }

  private handleAlbumViewClick(albumId: string) {
    this.setState({ renderedAlbumId: albumId });
  }

  private clearRenderedAlbumId() {
    this.setState({ renderedAlbumId: null });
  }

  private renderAlbumIndex() {
    return (
      <div>
        <TitleBar title="Albums"/>
          <div>
            <section className="albumIndexView__cardContainer">
              {/* tslint:disable-next-line:jsx-no-multiline-js */}
              {this.state.albums.map((album: Album) => {
                return(
                  <article key={album.id} className="albumIndexView__card">
                    <AlbumItem
                      baseApiUrl={this.props.baseApiUrl}
                      source={album}
                      albumViewCallback={this.handleAlbumViewClick.bind(this)}
                    />
                  </article>
                );
              })}
            </section>
          </div>
        </div>
    );
  }

  private renderAlbumPhotoView() {
    return (
      <div>
        <RaisedButton onClick={() => this.clearRenderedAlbumId()} className="albumIndexView__button">
          Back
        </RaisedButton>
        <PhotoIndexView
          baseApiUrl={this.props.baseApiUrl}
          // albumId={this.state.renderedAlbumId}
        />
      </div>
    );
  }

  public render() {
    let content;
    if (this.state && this.state.renderedAlbumId != null) {
      content = this.renderAlbumPhotoView();
    } else {
      content = this.renderAlbumIndex();
    }

    return (
      <MuiThemeProvider>
        {content}
      </MuiThemeProvider>
    );
  }
}
