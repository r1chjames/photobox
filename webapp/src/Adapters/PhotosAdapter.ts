import {IRestApiAdapter} from "./RestApiAdapter";
import {IPhotosAdapter} from "./IPhotosAdapter";
import {Photo} from "../Models/Photo";

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

  public getAllPhotosInfo = async (offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const getAllPhotosPath = `photos?page=${offset}&limit=${limit}&thumbnail=${includeThumbnails}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo[]> => {
    const getAllPhotosPath = `photo/${photoId}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getPhotosInfoInAlbum = async (albumId: string, offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const getPhotosInAlbumPath = `photos?page=${offset}&limit=${limit}&thumbnail=${includeThumbnails}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, this.buildHeaders(this.restApiAdapter.authHeader()), params);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    const getPhotosInAlbumPath = `photos/count?albumId=${albumId}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, this.buildHeaders(this.restApiAdapter.authHeader()), params);
  }

  public getPhotoImage = async (photoId: string) => {
    const getPhotosImagePath = `photo/bin`;
    const params = {
      photoId,
    };
    return this.restApiAdapter.getApiCall(getPhotosImagePath, this.buildHeaders(this.restApiAdapter.authHeader()), params);
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    const postSettingPath = `photo`;
    return this.restApiAdapter.postApiCall(postSettingPath, body, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}


