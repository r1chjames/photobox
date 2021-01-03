import React, { useState } from 'react';
import './PhotoItem.css';
import {Dialog, DialogActions, DialogContent, DialogTitle} from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { Photo } from '../../Models/Photo';
import history from '../../Routing/History';

interface IProps {
  baseApiUrl: string;
  num: number;
  groupKey?: number;
  source: Photo;
  allPhotos: Photo[];
}

const findImageIndexInArray = (photos: Photo[], photoId: string, defaultVal: number) => {
  for (let i = 0; i < photos.length; i += 1) {
    if (photos[i].id === photoId) {
      return i;
    }
  }
  return defaultVal;
};

export const PhotoItem: React.FunctionComponent<IProps> = (props) => {

  const [currentId, setCurrentId] = useState(props.source.id);
  const [isImageModalOpen, setImageModalOpen] = useState(false);

  const currentPhotoIdIndex = (): number => {
    return findImageIndexInArray(props.allPhotos, currentId, 0);
  };

  const previousPhotoId = (currentPhotoIndex: number) => setCurrentId(props.allPhotos![currentPhotoIndex - 1].id);
  const nextPhotoId = (currentPhotoIndex: number) => setCurrentId(props.allPhotos![currentPhotoIndex + 1].id);

  const gridImage = () => {
    return(
        <div className="photoItem__item">
          <div className="photoItem__thumbnail">
            <img
              src={`${props.baseApiUrl}/photo/${props.source.id}/thumbnail`}
              alt={props.source.name}
              onClick={() => setImageModalOpen(!isImageModalOpen)}
              onError={(e: any) => e.target.src = '/no_image.png'}
            />
          </div>
          <div className="photoItem__info">{`egjs ${props.num}`}</div>
        </div>
    );
  };

  const previousButton = () => {
    const currentIdIndex = currentPhotoIdIndex() ;
    if (currentIdIndex !== 0) {
      return (
        <Button onClick={() => previousPhotoId(currentIdIndex)} color="primary">
          Previous
        </Button>
      );
    }
  };

  const nextButton = () => {
    const currentIdIndex = currentPhotoIdIndex() ;
    if (currentIdIndex !== props.allPhotos.length - 1) {
      return (
        <Button onClick={() => nextPhotoId(currentIdIndex)} color="primary">
          Next
        </Button>
      );
    }
  };

  const content = () => {
    if (isImageModalOpen) {
      return (
        <div>
          <div>
            <Dialog
              className="photoItem__dialog"
              open={true}
              aria-labelledby="customized-dialog-title"
              onClose={() => setImageModalOpen(false)}
            >
              <DialogTitle id="customized-dialog-title">
                {props.source.name}
              </DialogTitle>
              <DialogContent>
                <img
                  className="photoItem__dialog__image"
                  onClick={() => history.push(`photo/${currentId}`)}
                  src={`${props.baseApiUrl}/photo/${currentId}/bin`}
                  alt={props.source.name}
                />
              </DialogContent>
              <DialogActions>
                {previousButton()}
                <Button onClick={() => setImageModalOpen(!isImageModalOpen)} color="primary">
                  Close
                </Button>
                {nextButton()}
              </DialogActions>
            </Dialog>
          </div>
          {gridImage()}
        </div>
      );
    }

    return (
      <div>
        {gridImage()}
      </div>
    );
  };

  return content();
};
