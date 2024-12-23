import {Photo} from "../Models/Photo";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  public getAllPhotosInfo = async (page: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    return this._photos;
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo[]> => {
    return this._photos.filter(p => `p${p.id}` === photoId);
  }

  public getPhotosInfoInAlbum = async (albumId: string, page: number, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    return this._photos.filter(p => `${p.albumId}` === albumId);
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
