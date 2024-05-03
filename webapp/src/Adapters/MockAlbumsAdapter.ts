import {IAlbumsAdapter} from "./IAlbumsAdapter";

export class MockAlbumsAdapter implements IAlbumsAdapter {

  public getAllAlbumsInfo = async () => {
    return "";
  }

  public getCountOfPhotosInAlbum = async () => {
    return "";
  }

  public getAlbumInfoById = async (albumId: string) => {
    return "";
  }
}
