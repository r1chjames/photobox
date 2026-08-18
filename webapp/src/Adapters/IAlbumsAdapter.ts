import {Album} from "../Models/Album";
import {SmartAlbumRules} from "../Models/SmartAlbumRules";

export interface IAlbumsAdapter {

  getAllAlbumsInfo(): Promise<any>;
  getCountOfPhotosInAlbum(): Promise<any>;
  getAlbumInfoById(albumId: string): Promise<Album>;
  createAlbum(name: string, description?: string): Promise<Album>;
  createSmartAlbum(name: string, rules: SmartAlbumRules): Promise<Album>;
  updateSmartAlbum(albumId: string, rules: SmartAlbumRules): Promise<Album>;
  updateAlbum(albumId: string, updates: { name?: string; description?: string; coverPhotoId?: string }): Promise<Album>;
  deleteAlbum(albumId: string, deletePhotos?: boolean): Promise<void>;
  addPhotosToAlbum(albumId: string, photoIds: string[]): Promise<void>;
}
