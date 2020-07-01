import * as React from 'react';
import './AlbumItem.css';
import { Album } from '../../Models/Album';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import { LoadingScreen } from '../LoadingScreen/LoadingScreen';
import { Button, Card, CardActionArea, CardActions, CardContent, CardMedia, Typography } from '@material-ui/core';

interface IProps {
  baseApiUrl: string;
  source: Album;
  albumViewCallback: (albumId: string) => void;
}

interface IState {
  thumbnailUrl: string;
}

export class AlbumItem extends React.Component<IProps, IState> {

  constructor(props: IProps) {
    super(props);
    this.albumViewCallback = this.albumViewCallback.bind(this);
  }

  private async getUrlOfFirstImageInAlbum() {
    const photosAdapter = new PhotosAdapter(this.props.baseApiUrl);
    const photos: Photo[] = await photosAdapter.getPhotosInfoInAlbum(this.props.source.id);
    return `${this.props.baseApiUrl}/photo/bin?photoId=${photos[0].id}`;
  }

  public async componentDidMount() {
    const url = await this.getUrlOfFirstImageInAlbum();
    this.setState({ thumbnailUrl: url });
  }

  private albumViewCallback() {
    this.props.albumViewCallback(this.props.source.id);
  }

  public render() {
    if (this.state && this.state.thumbnailUrl) {
      return (
        <Card className="albumItem__cardWapper">
          <CardActionArea>
            <CardMedia
              className="albumItem__cardImage"
              image={this.state.thumbnailUrl}
              title={this.props.source.name}
            />
            <CardContent>
              <Typography variant="h5" component="h2">
                {this.props.source.name}
              </Typography>
            </CardContent>
          </CardActionArea>
          <CardActions>
            <Button size="small" color="primary">
              Share
            </Button>
            <Button
              size="small"
              color="primary"
              onClick={() => this.albumViewCallback()}
            >
              View
            </Button>
          </CardActions>
        </Card>
      );
    }
    return (
      <LoadingScreen />
    );
  }
}
