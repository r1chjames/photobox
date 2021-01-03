import React, {useEffect, useState} from 'react';
import './PhotoDetail.css';
import {Photo} from '../../Models/Photo';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { MainContent } from '../MainContent/MainContent';
import { useParams } from 'react-router-dom';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { LoadingScreen } from '../LoadingScreen/LoadingScreen';
import {Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow} from '@material-ui/core';

interface IProps {
  baseApiUrl: string;
}

const getPhoto = async(baseApiUrl: string, photoId: string) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  return photosAdapter.getPhotoInfoById(photoId);
};

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
  const [photo, setPhoto] = useState<Photo>();
  const { id } = useParams();

  useEffect(() => {
    if (id !== undefined) {
      (async function retrievePhoto() {
        const retrievedPhoto = await getPhoto(props.baseApiUrl, id);
        setPhoto(retrievedPhoto);
      })();
    }
  },        [setPhoto, id, props.baseApiUrl]);

  const content = () => {
    if (photo !== undefined) {
      return (
        <MainContent title={photo.name}>
          <div className="photoDetail__contentWrapper">
              <img
                src={`${props.baseApiUrl}/photo/${id}/bin`}
                alt={photo.name}
                className="photoDetail__mainImage"
              />
            <div className="photoDetail__imageMetadata">
              <TableContainer component={Paper}>
                <Table aria-label="metadata table">
                  <TableHead>
                    <TableRow>
                      <TableCell>Parameter</TableCell>
                      <TableCell>Value</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {/* eslint-disable-next-line array-callback-return */}
                    {Object.entries(photo.metadata).map(([key, value]) => {
                      if (value !== null) {
                        return (
                          <TableRow key={key}>
                            <TableCell component="th" scope="row">{key}</TableCell>
                            <TableCell>{JSON.stringify(value)}</TableCell>
                          </TableRow>
                        );
                      }
                    })}
                  </TableBody>
                </Table>
              </TableContainer>
            </div>
          </div>
        </MainContent>
      );
    }
    return (
      <LoadingScreen />
    );
  };

  return (
    <MuiThemeProvider>
      {content()}
    </MuiThemeProvider>
  );
};
