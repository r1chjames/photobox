import React, {useState} from 'react';
import Dropzone from 'react-dropzone';
import './CreateAlbumView.css';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { MainContent } from '../MainContent/MainContent';
import {Button, Fab} from '@material-ui/core';
import SaveIcon from '@material-ui/icons/Save';
import { useParams } from 'react-router-dom';
import AddIcon from '@material-ui/icons/Add';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import {InfoSnackbar} from '../Snackbar/InfoSnackbar';

interface IProps {
  baseApiUrl: string;
}

const uploadPhotoToApi = async (baseApiUrl: string, body: {}) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  await photosAdapter.uploadPhoto(body);
};

const readUploadedFileAsText = (inputFile: File) => {
  const temporaryFileReader = new FileReader();

  return new Promise((resolve, reject) => {
    temporaryFileReader.onerror = () => {
      temporaryFileReader.abort();
      reject(new DOMException('Problem parsing input file.'));
    };

    temporaryFileReader.onload = () => {
      resolve(temporaryFileReader.result);
    };
    temporaryFileReader.readAsDataURL(inputFile);
  });
};

export const CreateAlbumView: React.FunctionComponent<IProps> = (props) => {

  const [showSnackbar, setShowSnackbar] = useState(false);
  const { name } = useParams();

  const handleFileUpload = async (acceptedFiles: File[]) => {
    for (const file of acceptedFiles) {
      const fileContent = await readUploadedFileAsText(file);
      const photoContent = {
        name: file.name,
        albumName: name,
        binaryContent: fileContent,
      };
      uploadPhotoToApi(props.baseApiUrl, photoContent).then(() => setShowSnackbar(true));
    }
  };

  return (
    <MuiThemeProvider>
      <MainContent title={name}>
        <div className="createAlbumView__dropzone">
          <Dropzone onDrop={acceptedFiles => handleFileUpload(acceptedFiles)}>
            {({ getRootProps, getInputProps }) => (
              <section>
                <div {...getRootProps()}>
                  <input {...getInputProps()} />
                  <div className="createAlbumView__dropzone__uploadButtonContainer">
                    <AddIcon className="createAlbumView__dropzone__uploadButtonIcon"/>
                    <Button>Select photos</Button>
                  </div>
                  <p>Drag photos here to upload</p>
                </div>
              </section>
            )}
          </Dropzone>
        </div>
        <Fab color="primary" aria-label="add" className="createAlbumView__addPhotoButton">
          <AddIcon onClick={() => console.log('save')}/>
        </Fab>
        <Fab color="primary" aria-label="add" className="createAlbumView__addButton">
          <SaveIcon onClick={() => console.log('save')}/>
        </Fab>
        <InfoSnackbar text={'Photo Uploaded'} show={showSnackbar} handleStopShowing={() => setShowSnackbar(false)}/>
      </MainContent>
    </MuiThemeProvider>
  );
};
