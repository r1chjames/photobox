import {IRestApiAdapter} from "./RestApiAdapter";

export class AlbumsAdapter {

  private restApiAdapter: IRestApiAdapter;

  constructor(restApiAdapter: IRestApiAdapter) {
    this.restApiAdapter = restApiAdapter;
  }

  private buildHeaders = (additionalHeaders: Record<string, string> = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllAlbumsInfo = async () => {
    const getAllAlbumsPath = "albums";
    return this.restApiAdapter.getApiCall(getAllAlbumsPath, this.buildHeaders(), {});
  }

  public getCountOfPhotosInAlbum = async () => {
    const getAllPhotosPath = "album";
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(), {});
  }

  public getAlbumInfoById = async (albumId: string) => {
    const getAlbumInfoPath = `/album/${albumId}`;
    return this.restApiAdapter.getApiCall(getAlbumInfoPath, this.buildHeaders(), {});
  }
}
