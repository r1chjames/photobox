import React from 'react';
import './AlbumView.css';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import RaisedButton from 'material-ui/RaisedButton';
import history from '../../Routing/History';
import { MainContent } from '../MainContent/MainContent';

interface IProps {
  baseApiUrl: string;
  albumId: string | null;
}

export const AlbumView: React.FunctionComponent<IProps> = (props) => {

  const renderAlbumPhotoView = () => {
    return (
      <MainContent>
        <RaisedButton onClick={() => history.push('/albums')} className="albumIndexView__button">
          Back
        </RaisedButton>
        <PhotoIndexView
          baseApiUrl={props.baseApiUrl}
        />
      </MainContent>
    );
  };

  return (
      <MuiThemeProvider>
        {renderAlbumPhotoView()}
      </MuiThemeProvider>
  );
};
