import React, { Component, ReactElement } from 'react';
import { JustifiedLayout, OnAppend, OnLayoutComplete } from '@egjs/react-infinitegrid';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import { Photo } from '../../Models/Photo';
import './PhotoIndexView.css';
import { ImageItem } from '../ImageItem/ImageItem';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { TitleBar } from '../TitleBar/TitleBar';
import { AlbumsAdapter } from '../../Adapters/AlbumsAdapter';
import { Album } from '../../Models/Album';
import { RouteComponentProps } from 'react-router-dom';

interface IRouteInfo {
  id: string;
}

interface IComponentProps extends RouteComponentProps<IRouteInfo> {
}

interface IProps {
  baseApiUrl: string;
}

interface IState {
  images: ReactElement[];
  albumTitle: string;
  albumId: string;
}

export class PhotoIndexView extends Component<IProps, IState, IComponentProps> {

  constructor(props: IProps, compProps: IComponentProps) {
    super(props);
    const albumId = compProps.match.params.id;
    this.state = { images: [], albumTitle: 'Photos', albumId };
  }

  private async getAllPhotos() {
    const photosAdapter = new PhotosAdapter(this.props.baseApiUrl);
    let photos: Photo[];
    if (this.state.albumId != null) {
      photos = await photosAdapter.getPhotosInfoInAlbum(this.state.albumId);
    } else {
      photos = await photosAdapter.getAllPhotosInfo();
    }
    return photos;
  }

  private loadItems(groupKey: number, num: number, imageSources: Photo[]) {
    const items = [];
    const start = this.state.images.length;

    for (let i = 0; i < num; i += 1) {
      const imageSource = imageSources[start + i];
      if (typeof imageSource !== 'undefined') {
        items.push(
          <ImageItem
            baseApiUrl={this.props.baseApiUrl}
            groupKey={groupKey}
            num={1 + start + i}
            key={start + i}
            source={imageSource}
          />
        );
      }
    }
    return items;
  }

  public async componentDidMount() {
    const albumsAdapter = new AlbumsAdapter(this.props.baseApiUrl);
    if (this.state.albumId != null) {
      const album: Album = await albumsAdapter.getAlbumInfoById(this.state.albumId);
      this.setState({ albumTitle: album.name });
    }
  }

  private onAppend = async (params: OnAppend) => {
    // @ts-ignore
    params.startLoading();
    const imageSources = await this.getAllPhotos();
    if (imageSources.length === this.state.images.length) {
      return;
    }
    const items = this.loadItems(
      (parseFloat(params.groupKey as string) || 0) + 1,
      20,
      imageSources
    );

    this.setState({ images: this.state.images.concat(items) });
  }

  public onLayoutComplete = (params: OnLayoutComplete) => {
    // @ts-ignore
    return !params.isLayout && params.endLoading();
  }

  public render() {
    return (
      <MuiThemeProvider>
        <div>
          <TitleBar title={this.state.albumTitle} />
          <div className="photoIndexView__photoIndex">
            <JustifiedLayout
              options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
              layoutOptions={{ margin: 5, column: [0, 5] }}
              onAppend={this.onAppend}
              onLayoutComplete={this.onLayoutComplete}
            >
              {this.state.images}
            </JustifiedLayout>
          </div>
        </div>
      </MuiThemeProvider>
    );
  }
}
