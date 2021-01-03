import { RestApiAdapter } from './RestApiAdapter';

export class PhotosAdapter extends RestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private static buildHeaders = (additionalHeaders: {} = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllPhotosInfo = async () => {
    const getAllPhotosPath = `${this.baseApiPath}/photos`;
    return this.getApiCall(getAllPhotosPath, PhotosAdapter.buildHeaders(), {});
  }

  public getPhotoInfoById = async (photoId: string) => {
    const getAllPhotosPath = `${this.baseApiPath}/photo/${photoId}`;
    return this.getApiCall(getAllPhotosPath, PhotosAdapter.buildHeaders(), {});
  }

  public getPhotosInfoInAlbum = async (albumId: string, page: number, limit: number) => {
    const getPhotosInAlbumPath = `${this.baseApiPath}/photos?page=${page}&limit=${limit}`;
    const params = {
      albumId,
    };
    return this.getApiCall(getPhotosInAlbumPath, PhotosAdapter.buildHeaders(), params);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    const getPhotosInAlbumPath = `${this.baseApiPath}/photos/count?albumId=${albumId}`;
    const params = {
      albumId,
    };
    return this.getApiCall(getPhotosInAlbumPath, PhotosAdapter.buildHeaders(), params);
  }

  public getPhotoImage = async (photoId: string) => {
    const getPhotosImagePath = `${this.baseApiPath}/photo/bin`;
    const params = {
      photoId,
    };
    return this.getApiCall(getPhotosImagePath, PhotosAdapter.buildHeaders(), params);
  }

  public uploadPhoto = async (body: {}) => {
    const postSettingPath = `${this.baseApiPath}/photo`;
    return this.postApiCall(postSettingPath, body, PhotosAdapter.buildHeaders());
  }
}
