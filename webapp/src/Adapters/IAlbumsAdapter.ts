import {Album} from "../Models/Album";

export interface IAlbumsAdapter {

  getAllAlbumsInfo(): Promise<any>;
  getCountOfPhotosInAlbum(): Promise<any>;
  getAlbumInfoById(albumId: string): Promise<Album>;
  createAlbum(name: string, description?: string): Promise<Album>;
  updateAlbum(albumId: string, updates: { name?: string; description?: string; coverPhotoId?: string }): Promise<Album>;
  deleteAlbum(albumId: string, deletePhotos?: boolean): Promise<void>;
  addPhotosToAlbum(albumId: string, photoIds: string[]): Promise<void>;
}
