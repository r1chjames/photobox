import {Photo} from "../Models/Photo";

export interface IPhotosAdapter {

  getAllPhotosInfo(fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoInfoById(photoId: string): Promise<Photo>;
  getPhotosInfoInAlbum(albumId: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoCountInAlbum(albumId: string): Promise<any>;
  getPhotoImage(photoId: string): Promise<any>;
  downloadPhoto(photoId: string, filename: string): Promise<void>;
  uploadPhoto(body: Record<string, unknown>): Promise<any>;
  index(): Promise<any>;
}
