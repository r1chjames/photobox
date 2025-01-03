import {Photo} from "../Models/Photo";

export interface IPhotosAdapter {

  getAllPhotosInfo(offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoInfoById(photoId: string): Promise<Photo>;
  getPhotosInfoInAlbum(albumId: string, offset: number, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getPhotoCountInAlbum(albumId: string): Promise<any>;
  getPhotoImage(photoId: string): Promise<any>;
  uploadPhoto(body: Record<string, unknown>): Promise<any>;
}
