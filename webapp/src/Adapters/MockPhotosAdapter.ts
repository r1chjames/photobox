import {Photo} from "../Models/Photo";
import {paginate} from "../utils/ArraysUtils";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  public getAllPhotosInfo = async (offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const toReturn = paginate(this._photos, limit, offset);
    console.log("Returning: " + JSON.stringify(toReturn));
    return toReturn;
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo[]> => {
    return this._photos.filter(p => `p${p.id}` === photoId);
  }

  public getPhotosInfoInAlbum = async (albumId: string, offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const toReturn = paginate(this._photos.filter(p => `${p.albumId}` === albumId), limit, offset);
    console.log("Returning: " + JSON.stringify(toReturn));
    return toReturn;
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `a${p.albumId}` === albumId).length;
  }

  public getPhotoImage = async (photoId: string) => {
    return this._photos.filter(p => `p${p.id}` === photoId);
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }
}
