import { RestApiAdapter } from './RestApiAdapter';

export class PhotosAdapter extends RestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private static buildHeaders(additionalHeaders: {} = {}) {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public async getAllPhotosInfo() {
    const getAllPhotosPath = `${this.baseApiPath}/photos`;
    return this.getApiCall(getAllPhotosPath, '', PhotosAdapter.buildHeaders(), {});
  }

  public async getPhotoInfoById(photoId: string) {
    const getAllPhotosPath = `${this.baseApiPath}/photos`;
    const params = {
      photoId,
    };
    return this.getApiCall(getAllPhotosPath, '', PhotosAdapter.buildHeaders(), params);
  }

  public async getPhotosInfoInAlbum(albumId: string, page: number, limit: number) {
    const getPhotosInAlbumPath = `${this.baseApiPath}/photos?page=${page}&limit=${limit}`;
    const params = {
      albumId,
    };
    return this.getApiCall(getPhotosInAlbumPath, '', PhotosAdapter.buildHeaders(), params);
  }

  public async getPhotoImage(photoId: string) {
    const getPhotosImagePath = `${this.baseApiPath}/photo/bin`;
    const params = {
      photoId,
    };
    return this.getApiCall(getPhotosImagePath, '', PhotosAdapter.buildHeaders(), params);
  }
}
