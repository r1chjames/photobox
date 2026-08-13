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
    return this._photos.filter(p => `${p.albumId}` === albumId).length;
  }

  public getCountOfPhotosInAlbum = async () => {
    return {
      "photoCount" : this._photos.filter(p => `${p.albumId}` === "albumId").length
    };
  }

  public getAlbumInfoById = async (albumIdentifier: string): Promise<Album> => {
    // In Storybook, the identifier might be a name like "Album 1" from the route.
    // In real scenarios, it would be an ID. This mock handles both.
    const album = this._albums.find(a => a.id === albumIdentifier || a.name === albumIdentifier);
    if (!album) {
      return Promise.reject(`Album with identifier ${albumIdentifier} not found in mock adapter`);
    }
    return album;
  }

  public createAlbum = async (name: string, description?: string): Promise<Album> => {
    return { id: 'new-album', name, description: description ?? '', tags: '', metadata: {} as Record<string, unknown>, createdAt: new Date().toISOString(), updatedAt: new Date().toISOString() } as unknown as Album;
  }

  public updateAlbum = async (albumId: string, updates: { name?: string; description?: string; coverPhotoId?: string }): Promise<Album> => {
    const album = this._albums.find(a => a.id === albumId);
    if (!album) {
      return Promise.reject(`Album not found`);
    }
    return { ...album, ...updates };
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public deleteAlbum = async (albumId: string, deletePhotos?: boolean): Promise<void> => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public addPhotosToAlbum = async (albumId: string, photoIds: string[]): Promise<void> => {
    return;
  }
}
