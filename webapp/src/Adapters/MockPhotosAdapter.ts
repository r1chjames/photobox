import {Photo} from "../Models/Photo";
import {IPhotosAdapter} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _getAllPhotosInfo: Photo[] | any = null;
  private _getPhotoInfoById: Photo | any = null;
  private _getPhotosInfoInAlbum: Photo[] | any = null;
  private _getPhotoCountInAlbum: number = 0;
  private _getPhotoImage: string | any = null;

  public withGetAllPhotosInfo(getAllPhotosInfo: Photo[]): this {
    this._getAllPhotosInfo = getAllPhotosInfo;
    return this;
  }

  public withGetPhotoInfoById(getPhotoInfoById: Photo): this {
    this._getPhotoInfoById = getPhotoInfoById;
    return this;
  }

  public withGetPhotosInfoInAlbum(getPhotosInfoInAlbum: Photo[]): this {
    this._getPhotosInfoInAlbum = getPhotosInfoInAlbum;
    return this;
  }

  public withGetPhotoCountInAlbum(getPhotoCountInAlbum: number): this {
    this._getPhotoCountInAlbum = getPhotoCountInAlbum;
    return this;
  }

  public withGetPhotoImage(getPhotoImage: string): this {
    this._getPhotoImage = getPhotoImage;
    return this;
  }


  public getAllPhotosInfo = async (includeThumbnails: boolean) => {
    return this._getAllPhotosInfo;
  }

  public getPhotoInfoById = async (photoId: string) => {
    return this._getPhotoInfoById;
  }

  public getPhotosInfoInAlbum = async (albumId: string, page: number, limit: number, includeThumbnails: boolean) => {
    return this._getPhotosInfoInAlbum;
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._getPhotoCountInAlbum;
  }

  public getPhotoImage = async (photoId: string) => {
    return this._getPhotoImage;
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }
}
