import {IRestApiAdapter} from "./RestApiAdapter";

export class PhotosAdapter {

  private restApiAdapter: IRestApiAdapter;

  constructor(restApiAdapter: IRestApiAdapter) {
    this.restApiAdapter = restApiAdapter;
  }

  private static buildHeaders = (additionalHeaders: Record<string, string> = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllPhotosInfo = async () => {
    const getAllPhotosPath = "photos";
    return this.restApiAdapter.getApiCall(getAllPhotosPath, PhotosAdapter.buildHeaders(), {});
  }

  public getPhotoInfoById = async (photoId: string) => {
    const getAllPhotosPath = `photo/${photoId}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, PhotosAdapter.buildHeaders(), {});
  }

  public getPhotosInfoInAlbum = async (albumId: string, page: number, limit: number) => {
    const getPhotosInAlbumPath = `photos?page=${page}&limit=${limit}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, PhotosAdapter.buildHeaders(), params);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    const getPhotosInAlbumPath = `photos/count?albumId=${albumId}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, PhotosAdapter.buildHeaders(), params);
  }

  public getPhotoImage = async (photoId: string) => {
    const getPhotosImagePath = `photo/bin`;
    const params = {
      photoId,
    };
    return this.restApiAdapter.getApiCall(getPhotosImagePath, PhotosAdapter.buildHeaders(), params);
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    const postSettingPath = `photo`;
    return this.restApiAdapter.postApiCall(postSettingPath, body, PhotosAdapter.buildHeaders());
  }
}
