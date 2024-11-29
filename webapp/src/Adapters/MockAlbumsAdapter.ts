import {IAlbumsAdapter} from "./IAlbumsAdapter";
import {Album} from "../Models/Album";

export class MockAlbumsAdapter implements IAlbumsAdapter {

  private _getAllAlbumsInfo: Album[] = [];
  private _getCountOfPhotosInAlbum: number = 0;
  private _getAlbumInfoById: string = "";

  public withGetAllAlbumsResponse(getAllAlbumsInfo: Album[]): this {
    this._getAllAlbumsInfo = getAllAlbumsInfo
    return this;
  }

  public withGetCountOfPhotosInAlbum(getCountOfPhotosInAlbum: number): this {
    this._getCountOfPhotosInAlbum = getCountOfPhotosInAlbum
    return this;
  }

  public withGetAlbumInfoById(getAlbumInfoById: string): this {
    this._getAlbumInfoById = getAlbumInfoById
    return this;
  }

  public getAllAlbumsInfo = async () => {
    return this._getAllAlbumsInfo;
  }

  public getCountOfPhotosInAlbum = async () => {
    return {
      "photoCount" : this._getCountOfPhotosInAlbum
    };
  }

  public getAlbumInfoById = async (albumId: string) => {
    return this._getAlbumInfoById;
  }
}
