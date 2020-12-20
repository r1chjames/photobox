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
}

export const PhotoItem: React.FunctionComponent<IProps> = (props) => {

  const [isImageModalOpen, setImageModalOpen] = useState(false);
  const thumbnailUrl = `${props.baseApiUrl}/photo/${props.source.id}/thumbnail`;
  const photoUrl = `${props.baseApiUrl}/photo/${props.source.id}/bin`;

  const gridImage = () => {
    return(
        <div className="imageItem__item">
          <div className="imageItem__thumbnail">
            <img
              src={thumbnailUrl}
              alt={props.source.name}
              onClick={() => setImageModalOpen(!isImageModalOpen)}
            />
          </div>
          <div className="imageItem__info">{`egjs ${props.num}`}</div>
        </div>
    );
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
                <a href={photoUrl} >
                  <img
                    className="imageItem__dialog__image"
                    src={photoUrl}
                    alt={props.source.name}
                  />
                </a>
              </DialogContent>
              <DialogActions>
                <Button onClick={() => setImageModalOpen(!isImageModalOpen)} color="primary">
                  Close
                </Button>
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
