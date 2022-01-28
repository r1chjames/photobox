import React from 'react';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';

interface IProps {
  baseApiUrl: string;
}

const theme = createTheme();


export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <ThemeProvider theme={theme}>
        <div>
          <AlbumIndexView baseApiUrl={props.baseApiUrl} />
          <PhotoIndexView baseApiUrl={props.baseApiUrl} />
         </div>
       </ThemeProvider>
  );
};
