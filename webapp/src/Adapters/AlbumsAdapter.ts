import { RestApiAdapter } from './RestApiAdapter';

export class AlbumsAdapter extends RestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private buildHeaders = (additionalHeaders: {} = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllAlbumsInfo = async () => {
    const getAllAlbumsPath = `${this.baseApiPath}/albums`;
    return this.getApiCall(getAllAlbumsPath, '', this.buildHeaders(), {});
  }

  public getCountOfPhotosInAlbum = async () => {
    const getAllPhotosPath = `${this.baseApiPath}/album`;
    return this.getApiCall(getAllPhotosPath, '', this.buildHeaders(), {});
  }

  public getAlbumInfoById = async (albumId: string) => {
    const getAlbumInfoPath = `${this.baseApiPath}/album/${albumId}`;
    return this.getApiCall(getAlbumInfoPath, '', this.buildHeaders(), {});
  }
}
