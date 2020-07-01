import { RestApiAdapter } from './RestApiAdapter';

export class PhotosAdapter extends RestApiAdapter {

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

  public async getAllPhotosInfo() {
    const getAllPhotosPath = `${this.baseApiPath}/photos`;
    return this.getApiCall(getAllPhotosPath, '', this.buildHeaders(), {});
  }

  public async getPhotoInfoById(photoId: string) {
    const getAllPhotosPath = `${this.baseApiPath}/photos`;
    const params = {
      photoId,
    };
    return this.getApiCall(getAllPhotosPath, '', this.buildHeaders(), params);
  }

  public async getPhotosInfoInAlbum(albumId: string) {
    const getPhotosInAlbumPath = `${this.baseApiPath}/photos`;
    const params = {
      albumId,
    };
    return this.getApiCall(getPhotosInAlbumPath, '', this.buildHeaders(), params);
  }

  public async getPhotoImage(photoId: string) {
    const getPhotosImagePath = `${this.baseApiPath}/photo/bin`;
    const params = {
      photoId,
    };
    return this.getApiCall(getPhotosImagePath, '', this.buildHeaders(), params);
  }
}
