import {IRestApiAdapter} from "./RestApiAdapter";
import {IAlbumsAdapter} from "./IAlbumsAdapter";

export class AlbumsAdapter implements IAlbumsAdapter {

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
    return this.restApiAdapter.getPaginatedApiCall(getAllAlbumsPath, this.buildHeaders(this.restApiAdapter.authHeader()), { limit: 1000 });
  }

  public getCountOfPhotosInAlbum = async () => {
    const getAllPhotosPath = "album";
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getAlbumInfoById = async (albumId: string) => {
    const getAlbumInfoPath = `album/${albumId}`;
    return this.restApiAdapter.getApiCall(getAlbumInfoPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public createAlbum = async (name: string, description?: string) => {
    const createAlbumPath = 'albums';
    return this.restApiAdapter.postApiCall(createAlbumPath, { name, description }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public updateAlbum = async (albumId: string, updates: { name?: string; description?: string; coverPhotoId?: string }) => {
    const updateAlbumPath = `albums/${albumId}`;
    return this.restApiAdapter.putApiCall(updateAlbumPath, updates, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public deleteAlbum = async (albumId: string, deletePhotos?: boolean) => {
    const deleteAlbumPath = `albums/${albumId}?deletePhotos=${deletePhotos ?? false}`;
    return this.restApiAdapter.postApiCall(deleteAlbumPath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public addPhotosToAlbum = async (albumId: string, photoIds: string[]) => {
    const addPhotosPath = `albums/${albumId}/photos`;
    return this.restApiAdapter.postApiCall(addPhotosPath, { photoIds }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}
