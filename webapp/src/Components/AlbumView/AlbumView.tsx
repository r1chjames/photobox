import React, { Component } from 'react';
import './AlbumView.css';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RaisedButton from 'material-ui/RaisedButton';
import history from '../../Routing/History';

interface IProps {
  baseApiUrl: string;
  albumId: string | null;
}

export class AlbumView extends Component<IProps> {

  constructor(props: IProps) {
    super(props);
  }

  private renderAlbumPhotoView() {
    return (
      <div>
        <RaisedButton onClick={() => history.push('/albums')} className="albumIndexView__button">
          Back
        </RaisedButton>
        <PhotoIndexView
          baseApiUrl={this.props.baseApiUrl}
          // albumId={this.props.albumId}
        />
      </div>
    );
  }

  public render() {
    return (
      <MuiThemeProvider>
        {this.renderAlbumPhotoView()}
      </MuiThemeProvider>
    );
  }
}
