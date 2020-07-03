// @ts-ignore
import * as React from 'react';
import './ImageItem.css';
import { Dialog, DialogActions, DialogContent, DialogTitle } from '@material-ui/core';
import Button from '@material-ui/core/Button';
import { Photo } from '../../Models/Photo';

interface IProps {
  baseApiUrl: string;
  num: number;
  groupKey?: number;
  source: Photo;
}

interface IState {
  isOpen: boolean;
  url: string;
}

export class ImageItem extends React.Component<IProps, IState> {

  constructor(props: IProps) {
    super(props);
    this.state = {
      isOpen: false,
      url: `${props.baseApiUrl}/photo/bin?photoId=${props.source.id}`};
  }

  private handleShowImageModal = () => {
    this.setState({ isOpen: !this.state.isOpen });
    return this.state.isOpen;
  }

  private gridImage() {
    return(
        <div className="imageItem__item">
          <div className="imageItem__thumbnail">
            <img
              src={this.state.url}
              alt={this.props.source.name}
              onClick={this.handleShowImageModal}
            />
          </div>
          <div className="imageItem__info">{`egjs ${this.props.num}`}</div>
        </div>
    );
  }

  public render() {

    if (this.state && this.state.isOpen) {
      return (
        <div>
          <div>
            <Dialog
              className="imageItem__dialog"
              open={true}
              aria-labelledby="customized-dialog-title"
              onClose={this.handleShowImageModal}
            >
              <DialogTitle id="customized-dialog-title">
                {this.props.source.name}
              </DialogTitle>
              <DialogContent>
                <a href={this.state.url} >
                  <img
                    className="imageItem__dialog__image"
                    src={this.state.url}
                    alt={this.props.source.name}
                  />
                </a>
              </DialogContent>
              <DialogActions>
                <Button onClick={this.handleShowImageModal} color="primary">
                  Close
                </Button>
              </DialogActions>
            </Dialog>
          </div>
          {this.gridImage()}
        </div>
      );
    }

    return (
      <div>
        {this.gridImage()}
      </div>
    );
  }
}
