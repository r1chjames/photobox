import {Photo} from "../Models/Photo";
import {JobType} from "../Models/Job";
import {MemoryGroup} from "../Models/MemoryGroup";
import {IPhotosAdapter, PhotoGeoData, TimelineEntry} from "./IPhotosAdapter";

export class MockPhotosAdapter implements IPhotosAdapter {

  private _photos: Photo[] = [];

  public withPhotos(photos: Photo[]): this {
    this._photos = photos;
    return this;
  }

  public getAllPhotosInfo = async (fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string, favorite?: boolean, lowQuality?: boolean, camera?: string, hasGps?: boolean, orientation?: string): Promise<Photo[]> => {
    let photos = this._photos;
    if (favorite) {
      photos = photos.filter(p => p.favorite);
    }
    if (lowQuality) {
      photos = photos.filter(p => p.isLowQuality);
    }
    if (camera) {
      photos = photos.filter(p => (p.metadata as any)?.exif?.Make?.toLowerCase().includes(camera.toLowerCase()) || (p.metadata as any)?.exif?.Model?.toLowerCase().includes(camera.toLowerCase()));
    }
    if (hasGps) {
      photos = photos.filter(p => p.latitude && p.longitude);
    }
    if (orientation === 'landscape') {
      photos = photos.filter(p => (p.width ?? 0) > (p.height ?? 0));
    } else if (orientation === 'portrait') {
      photos = photos.filter(p => (p.height ?? 0) > (p.width ?? 0));
    } else if (orientation === 'square') {
      photos = photos.filter(p => (p.width ?? 0) === (p.height ?? 0));
    }
    const startIndex = fromId ? photos.findIndex(p => p.id === fromId) + 1 : 0;
    if (startIndex === -1) return []; // fromId not found
    const endIndex = startIndex + limit;
    return photos.slice(startIndex, endIndex);
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
  public getGeodata = async (north: number, south: number, east: number, west: number): Promise<PhotoGeoData[]> => {
    return [];
  }

  public getTimeline = async (): Promise<TimelineEntry[]> => {
    return [];
  }

  public getMemories = async (_date?: string): Promise<MemoryGroup[]> => {
    const now = new Date();
    const month = now.getMonth() + 1;
    const day = now.getDate();
    const groups = new Map<number, Photo[]>();
    this._photos.forEach(p => {
      if (!p.createdAt) return;
      const d = new Date(p.createdAt);
      if (d.getMonth() + 1 === month && d.getDate() === day && d.getFullYear() < now.getFullYear()) {
        const year = d.getFullYear();
        if (!groups.has(year)) groups.set(year, []);
        if (groups.get(year)!.length < 10) groups.get(year)!.push(p);
      }
    });
    return [...groups.entries()]
      .sort((a, b) => b[0] - a[0])
      .map(([year, photos]) => ({ year, yearsAgo: now.getFullYear() - year, photos }));
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public uploadPhoto = async (body: Record<string, unknown>) => {
    return null;
  }

  public startJob = async (_jobType: JobType) => {
    return null;
  }

  public stopJob = async (_jobType: JobType) => {
    return null;
  }

  public getAllJobStatuses = async () => {
    return { jobs: [] };
  }

  public stopAllJobs = async () => {
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

  public getPhotosByTag = async (tag: string, fromId: string, limit: number, _includeThumbnails: boolean): Promise<Photo[]> => {
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

  public batchSetFavorite = async (photoIds: string[], favorite: boolean): Promise<void> => {
    this._photos.forEach(p => {
      if (photoIds.includes(p.id)) p.favorite = favorite;
    });
  }

  public batchDeletePhotos = async (photoIds: string[]): Promise<void> => {
    this._photos = this._photos.filter(p => !photoIds.includes(p.id));
  }

  public batchAddToAlbum = async (photoIds: string[], albumId: string): Promise<void> => {
    this._photos.forEach(p => {
      if (photoIds.includes(p.id)) p.albumId = albumId;
    });
  }

  public rotatePhoto = async (photoId: string, direction: 'cw' | 'ccw'): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    photo.editParams = { ...(photo.editParams || {}), rotate: ((photo.editParams?.rotate || 0) + (direction === 'cw' ? 90 : 270)) % 360 };
    return photo;
  }

  public editPhoto = async (photoId: string, params: { rotate?: number; crop?: { x: number; y: number; width: number; height: number }; brightness?: number; contrast?: number; saturation?: number; autoEnhance?: boolean }): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    photo.editParams = { ...(photo.editParams || {}), ...params };
    return photo;
  }

  public clearEdits = async (photoId: string): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    photo.editParams = undefined;
    return photo;
  }

  public updatePhotoMetadata = async (photoId: string, updates: { description?: string; latitude?: number; longitude?: number; dateTaken?: string }): Promise<Photo> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo) {
      return Promise.reject(`Photo with id ${photoId} not found in mock adapter`);
    }
    if (updates.description !== undefined) photo.description = updates.description;
    if (updates.latitude !== undefined) photo.latitude = updates.latitude;
    if (updates.longitude !== undefined) photo.longitude = updates.longitude;
    if (updates.dateTaken !== undefined) photo.dateTaken = updates.dateTaken;
    return photo;
  }

  public getPhotoLocation = async (photoId: string): Promise<string> => {
    const photo = this._photos.find(p => p.id === photoId);
    if (!photo || (!photo.latitude && !photo.longitude)) return '';
    return 'Mock Location, United Kingdom';
  }
}
