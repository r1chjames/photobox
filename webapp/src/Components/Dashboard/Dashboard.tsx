import React from 'react';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';

interface IProps {
  baseApiUrl: string;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <MuiThemeProvider>
        <div>
          <AlbumIndexView baseApiUrl={props.baseApiUrl} />
          <PhotoIndexView baseApiUrl={props.baseApiUrl} />
        </div>
      </MuiThemeProvider>
  );
};
