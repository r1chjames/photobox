import {Photo} from "../Models/Photo";
import {IPhotosAdapter, PhotoGeoData, TimelineEntry} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getAllPhotosInfo = async (fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]> => {
    const startIndex = fromId ? this._photos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return []; // fromId not found
    const endIndex = startIndex + limit;
    return this._photos.slice(startIndex, endIndex);
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    return photo;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getPhotosInfoInAlbum = async (albumId: string, fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]> => {
    const albumPhotos = this._photos.filter(p => p.albumId === albumId);
    const startIndex = fromId ? albumPhotos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return []; // fromId not found
    const endIndex = startIndex + limit;
    return albumPhotos.slice(startIndex, endIndex);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    return this._photos.filter(p => `${p.albumId}` === albumId).length;
  }

  public getPhotoImage = async (photoId: string) => {
    const photo = this._photos.find(p => p.id === photoId);
    return photo?.thumbnailUrl;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public downloadPhoto = async (photoId: string, filename: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public deletePhoto = async (photoId: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public favoritePhoto = async (photoId: string, favorite: boolean) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public downloadPhotosAsZip = async (photoIds: string[]) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public searchPhotos = async (query: string, fromId: string, limit: number): Promise<Photo[]> => {
    const lowerQuery = query.toLowerCase();
    const filtered = this._photos.filter(p =>
      p.name.toLowerCase().includes(lowerQuery) ||
      (p.albumId && p.albumId.toLowerCase().includes(lowerQuery))
    );
    const startIndex = fromId ? filtered.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return [];
    return filtered.slice(startIndex, startIndex + limit);
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getTrashedPhotos = async (fromId: string, limit: number): Promise<Photo[]> => {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public restorePhoto = async (photoId: string) => {
    return;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public rotatePhoto = async (photoId: string, direction: 'cw' | 'ccw'): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    return photo;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getGeodata = async (north: number, south: number, east: number, west: number): Promise<PhotoGeoData[]> => {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public getTimeline = async (): Promise<TimelineEntry[]> => {
    return [];
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }

  public index = async () => {
    return null;
  }

  public getAllTags = async (): Promise<string[]> => {
    const tagSet = new Set<string>();
    this._photos.forEach(p => {
      if (p.tags) {
        p.tags.split(',').forEach(t => {
          const trimmed = t.trim();
          if (trimmed) tagSet.add(trimmed);
        });
      }
    });
    return Array.from(tagSet);
  }

  public updatePhotoTags = async (photoId: string, tags: string[]): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    photo.tags = tags.join(',');
    return photo;
  }

  public batchUpdatePhotoTags = async (photoIds: string[], tags: string[], operation: 'add' | 'remove' | 'set'): Promise<void> => {
    photoIds.forEach(id => {
      const photo = this._photos.find(p => p.id === id);
      if (!photo) return;
      const existing = new Set(photo.tags ? photo.tags.split(',').map(t => t.trim()).filter(Boolean) : []);
      tags.forEach(tag => {
        const trimmed = tag.trim();
        if (!trimmed) return;
        if (operation === 'add') existing.add(trimmed);
        if (operation === 'remove') existing.delete(trimmed);
      });
      if (operation === 'set') {
        photo.tags = tags.map(t => t.trim()).filter(Boolean).join(',');
      } else {
        photo.tags = Array.from(existing).join(',');
      }
    });
  }

  public getPhotosByTag = async (tag: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const filtered = this._photos.filter(p => {
      if (!p.tags) return false;
      return p.tags.split(',').map(t => t.trim()).includes(tag);
    });
    const startIndex = fromId ? filtered.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return [];
    return filtered.slice(startIndex, startIndex + limit);
  }

  public getDuplicatePhotos = async (): Promise<Photo[]> => {
    const hashMap = new Map<string, Photo[]>();
    this._photos.forEach(p => {
      if (p.fileHash) {
        const existing = hashMap.get(p.fileHash) || [];
        existing.push(p);
        hashMap.set(p.fileHash, existing);
      }
    });
    const duplicates: Photo[] = [];
    hashMap.forEach(group => {
      if (group.length > 1) duplicates.push(...group);
    });
    return duplicates;
  }

  public getPhotoThumbnails = async (photoIds: string[]): Promise<Map<string, string>> => {
    const thumbnailMap = new Map<string, string>();
    for (const photoId of photoIds) {
      const photo = this._photos.find(p => p.id === photoId);
      if (photo?.thumbnailUrl) {
        thumbnailMap.set(photoId, photo.thumbnailUrl);
      }
    }
    return thumbnailMap;
  }

  public getThumbnailUrl = (photoId: string): string => {
    const photo = this._photos.find(p => p.id === photoId);
    return photo?.thumbnailUrl || 'no_image.png';
  }

  public getPhotoThumbnailBlob = async (photoId: string): Promise<Blob | string> => {
    const photo = this._photos.find(p => p.id === photoId);
    return photo?.thumbnailUrl || 'no_image.png';
  }

  public getVideos = async (fromId: string, limit: number): Promise<Photo[]> => {
    const videos = this._photos.filter(p => p.mediaType === 'video');
    const startIndex = fromId ? videos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return [];
    return videos.slice(startIndex, startIndex + limit);
  }
}
