import {Photo} from "../Models/Photo";

export interface IPhotosAdapter {

  getAllPhotosInfo(fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoInfoById(photoId: string): Promise<Photo>;
  getPhotosInfoInAlbum(albumId: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoCountInAlbum(albumId: string): Promise<any>;
  getPhotoImage(photoId: string): Promise<any>;
  downloadPhoto(photoId: string, filename: string): Promise<void>;
  deletePhoto(photoId: string): Promise<void>;
  favoritePhoto(photoId: string, favorite: boolean): Promise<void>;
  downloadPhotosAsZip(photoIds: string[]): Promise<void>;
  searchPhotos(query: string, fromId: string, limit: number): Promise<Photo[]>;
  getTrashedPhotos(fromId: string, limit: number): Promise<Photo[]>;
  restorePhoto(photoId: string): Promise<void>;
  rotatePhoto(photoId: string, direction: 'cw' | 'ccw'): Promise<Photo>;
  uploadPhoto(body: Record<string, unknown>): Promise<any>;
  index(): Promise<any>;
}
