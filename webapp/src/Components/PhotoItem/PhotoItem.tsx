import React, { useState } from 'react';
import './PhotoItem.css';
import { Dialog, DialogActions, DialogContent, DialogTitle } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { Photo } from '../../Models/Photo';

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
        <div className="imageItem__item">
          <div className="imageItem__thumbnail">
            <img
              src={`${props.baseApiUrl}/photo/${props.source.id}/thumbnail`}
              alt={props.source.name}
              onClick={() => setImageModalOpen(!isImageModalOpen)}
            />
          </div>
          <div className="imageItem__info">{`egjs ${props.num}`}</div>
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
              className="imageItem__dialog"
              open={true}
              aria-labelledby="customized-dialog-title"
              onClose={() => setImageModalOpen(false)}
            >
              <DialogTitle id="customized-dialog-title">
                {props.source.name}
              </DialogTitle>
              <DialogContent>
                <a href={`${props.baseApiUrl}/photo/${currentId}/bin`} >
                  <img
                    className="imageItem__dialog__image"
                    src={`${props.baseApiUrl}/photo/${currentId}/bin`}
                    alt={props.source.name}
                  />
                </a>
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
