
export interface IPhotosAdapter {

  getAllPhotosInfo(): Promise<any>;
  getPhotoInfoById(photoId: string): Promise<any>;
  getPhotosInfoInAlbum(albumId: string, page: number, limit: number): Promise<any>;
  getPhotoCountInAlbum(albumId: string): Promise<any>;
  getPhotoImage(photoId: string): Promise<any>;
  uploadPhoto(body: Record<string, unknown>): Promise<any>;
}
