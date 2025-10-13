import {Album} from "../Models/Album";

export interface IAlbumsAdapter {

  getAllAlbumsInfo(): Promise<any>;
  getCountOfPhotosInAlbum(): Promise<any>;
  getAlbumInfoById(albumId: string): Promise<Album>;
}
