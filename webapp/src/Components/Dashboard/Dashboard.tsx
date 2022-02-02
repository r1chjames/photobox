import React from 'react';
import { MantineProvider } from '@mantine/core';
import './Dashboard.css';
import { AlbumIndexView } from '../AlbumIndexView/AlbumIndexView';
import { PhotoIndexView } from '../PhotoIndexView/PhotoIndexView';

interface IProps {
  baseApiUrl: string;
}


export const Dashboard: React.FunctionComponent<IProps> = (props) => {

  return (
      <MantineProvider>
        <div>
          <AlbumIndexView baseApiUrl={props.baseApiUrl} />
          <PhotoIndexView baseApiUrl={props.baseApiUrl} />
         </div>
       </MantineProvider>
  );
};
