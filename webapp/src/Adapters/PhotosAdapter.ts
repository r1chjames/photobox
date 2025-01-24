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

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
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
    const getPhotosImagePath = `photo/${photoId}/bin`;
    const resp = await this.restApiAdapter.getBinaryApiCall(getPhotosImagePath, this.buildHeaders(this.restApiAdapter.authHeader()), {responseType: 'blob'});
    return resp;
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    const uploadPhotoPath = `photo`;
    return this.restApiAdapter.postApiCall(uploadPhotoPath, body, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public index = async () => {
    const indexPath = 'photos/index';
    return this.restApiAdapter.postApiCall(indexPath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}


