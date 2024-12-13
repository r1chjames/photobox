import {IAlbumsAdapter} from "./IAlbumsAdapter";
import {Album} from "../Models/Album";
import {Photo} from "../Models/Photo";

export class MockAlbumsAdapter implements IAlbumsAdapter {

  private _photos: Photo[] = [];
  private _albums: Album[] = [];

  public withAlbums(albums: Album[]): this {
    this._albums = albums
    return this;
  }

  public getAllAlbumsInfo = async () => {
    return this._albums;
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `a${p.albumId}` === albumId).length;
  }

  public getCountOfPhotosInAlbum = async () => {
    return {
      "photoCount" : this._photos.filter(p => `a${p.albumId}` === "albumId").length
    };
  }

  public getAlbumInfoById = async (albumId: string) => {
    return this._albums.filter(a => `a${a.id}` === albumId);
  }
}
