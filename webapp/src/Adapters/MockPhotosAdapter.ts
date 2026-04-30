import {Photo} from "../Models/Photo";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getAllPhotosInfo = async (fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const startIndex = fromId ? this._photos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return []; // fromId not found
    const endIndex = startIndex + limit;
    return this._photos.slice(startIndex, endIndex);
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    return photo;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getPhotosInfoInAlbum = async (albumId: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const albumPhotos = this._photos.filter(p => p.albumId === albumId);
    const startIndex = fromId ? albumPhotos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return []; // fromId not found
    const endIndex = startIndex + limit;
    return albumPhotos.slice(startIndex, endIndex);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `${p.albumId}` === albumId).length;
  }

  public getPhotoImage = async (photoId: string) => {
    const photo = this._photos.find(p => p.id === photoId);
    return photo?.thumbnail;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public downloadPhoto = async (photoId: string, filename: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public deletePhoto = async (photoId: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public favoritePhoto = async (photoId: string, favorite: boolean) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public downloadPhotosAsZip = async (photoIds: string[]) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public searchPhotos = async (query: string, fromId: string, limit: number): Promise<Photo[]> => {
    const lowerQuery = query.toLowerCase();
    const filtered = this._photos.filter(p =>
      p.name.toLowerCase().includes(lowerQuery) ||
      (p.albumId && p.albumId.toLowerCase().includes(lowerQuery))
    );
    const startIndex = fromId ? filtered.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return [];
    return filtered.slice(startIndex, startIndex + limit);
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getTrashedPhotos = async (fromId: string, limit: number): Promise<Photo[]> => {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public restorePhoto = async (photoId: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public rotatePhoto = async (photoId: string, direction: 'cw' | 'ccw'): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    return photo;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }

  public index = async () => {
    return null;
  }
}
