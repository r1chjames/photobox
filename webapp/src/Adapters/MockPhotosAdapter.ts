import {Photo} from "../Models/Photo";
import {paginate} from "../utils/ArraysUtils";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getAllPhotosInfo = async (fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const offset = parseInt(fromId, 10) || 0;
    const toReturn = paginate(this._photos, limit, offset);
    return toReturn;
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
    const photo = this._photos.find(p => `p${p.id}` === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    return photo;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getPhotosInfoInAlbum = async (albumId: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const offset = parseInt(fromId, 10) || 0;
    const toReturn = paginate(this._photos.filter(p => `${p.albumId}` === albumId), limit, offset);
    return toReturn;
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `${p.albumId}` === albumId).length;
  }

  public getPhotoImage = async (photoId: string) => {
    const photo = this._photos.find(p => `p${p.id}` === photoId);
    return photo?.thumbnail;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }

  public index = async () => {
    return null;
  }
}
