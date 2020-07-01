import { RestApiAdapter } from './RestApiAdapter';

export class AlbumsAdapter extends RestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private buildHeaders(additionalHeaders: {} = {}) {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public async getAllAlbumsInfo() {
    const getAllPhotosPath = `${this.baseApiPath}/albums`;
    return this.getApiCall(getAllPhotosPath, '', this.buildHeaders(), {});
  }

  public async getAlbumInfoById(albumId: string) {
    const getAllPhotosPath = `${this.baseApiPath}/albums`;
    const params = {
      albumId,
    };
    return this.getApiCall(getAllPhotosPath, '', this.buildHeaders(), params);
  }
}
