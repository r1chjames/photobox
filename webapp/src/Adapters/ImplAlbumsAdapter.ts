import { RestApiAdapter } from './RestApiAdapter';
import { IAlbumsAdapter } from './IAlbumsAdapter';

export class ImplAlbumsAdapter extends RestApiAdapter  implements IAlbumsAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private static buildHeaders = (additionalHeaders: Record<string, string> = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllAlbumsInfo = async () => {
    const getAllAlbumsPath = `${this.baseApiPath}/albums`;
    return this.getApiCall(getAllAlbumsPath, ImplAlbumsAdapter.buildHeaders(), {});
  }

  public getCountOfPhotosInAlbum = async () => {
    const getAllPhotosPath = `${this.baseApiPath}/album`;
    return this.getApiCall(getAllPhotosPath, ImplAlbumsAdapter.buildHeaders(), {});
  }

  public getAlbumInfoById = async (albumId: string) => {
    const getAlbumInfoPath = `${this.baseApiPath}/album/${albumId}`;
    return this.getApiCall(getAlbumInfoPath, ImplAlbumsAdapter.buildHeaders(), {});
  }
}
