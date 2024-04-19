
interface IAlbumsAdapter {

  getAllAlbumsInfo(): Promise<any>;
  getCountOfPhotosInAlbum(): Promise<any>;
  getAlbumInfoById(): Promise<any>;
}

