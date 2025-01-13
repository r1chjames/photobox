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
  public getAllPhotosInfo = async (offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const toReturn = paginate(this._photos, limit, offset);
    return toReturn;
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
    return this._photos.filter(p => `p${p.id}` === photoId)[0];
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getPhotosInfoInAlbum = async (albumId: string, offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const toReturn = paginate(this._photos.filter(p => `${p.albumId}` === albumId), limit, offset);
    return toReturn;
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `${p.albumId}` === albumId).length;
  }

  public getPhotoImage = async (photoId: string) => {
    return this._photos.filter(p => `p${p.id}` === photoId);
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }
}
