import {IRestApiAdapter} from "./RestApiAdapter";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class PhotosAdapter implements IPhotosAdapter {

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

  public getAllPhotosInfo = async (includeThumbnails: boolean) => {
    const getAllPhotosPath = `photos?thumbnail=${includeThumbnails}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(), {});
  }

  public getPhotoInfoById = async (photoId: string) => {
    const getAllPhotosPath = `photo/${photoId}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(), {});
  }

  public getPhotosInfoInAlbum = async (albumId: string, page: number, limit: number, includeThumbnails: boolean) => {
    const getPhotosInAlbumPath = `photos?page=${page}&limit=${limit}&thumbnail=${includeThumbnails}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, this.buildHeaders(), params);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    const getPhotosInAlbumPath = `photos/count?albumId=${albumId}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, this.buildHeaders(), params);
  }

  public getPhotoImage = async (photoId: string) => {
    const getPhotosImagePath = `photo/bin`;
    const params = {
      photoId,
    };
    return this.restApiAdapter.getApiCall(getPhotosImagePath, this.buildHeaders(), params);
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    const postSettingPath = `photo`;
    return this.restApiAdapter.postApiCall(postSettingPath, body, this.buildHeaders());
  }
}
