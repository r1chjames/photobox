import {Photo} from "../Models/Photo";

export interface PhotoGeoData {
  id: string;
  lat: number;
  lng: number;
  dateTaken: string;
}

export interface TimelineEntry {
  year: number;
  month: number;
  count: number;
}

export interface IPhotosAdapter {

  getAllPhotosInfo(fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]>;
  getPhotoInfoById(photoId: string): Promise<Photo>;
  getPhotosInfoInAlbum(albumId: string, fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]>;
  getPhotoCountInAlbum(albumId: string): Promise<any>;
  getPhotoImage(photoId: string): Promise<any>;
  downloadPhoto(photoId: string, filename: string): Promise<void>;
  deletePhoto(photoId: string): Promise<void>;
  favoritePhoto(photoId: string, favorite: boolean): Promise<void>;
  downloadPhotosAsZip(photoIds: string[]): Promise<void>;
  searchPhotos(query: string, fromId: string, limit: number): Promise<Photo[]>;
  getTrashedPhotos(fromId: string, limit: number): Promise<Photo[]>;
  restorePhoto(photoId: string): Promise<void>;
  rotatePhoto(photoId: string, direction: 'cw' | 'ccw'): Promise<Photo>;
  getGeodata(north: number, south: number, east: number, west: number): Promise<PhotoGeoData[]>;
  getTimeline(): Promise<TimelineEntry[]>;
  getAllTags(): Promise<string[]>;
  updatePhotoTags(photoId: string, tags: string[]): Promise<Photo>;
  batchUpdatePhotoTags(photoIds: string[], tags: string[], operation: 'add' | 'remove' | 'set'): Promise<void>;
  getPhotosByTag(tag: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]>;
  getDuplicatePhotos(): Promise<Photo[]>;
  getPhotoThumbnailBlob(photoId: string): Promise<Blob | string>;
  uploadPhoto(body: Record<string, unknown>): Promise<any>;
  index(): Promise<any>;
  getVideos(fromId: string, limit: number): Promise<Photo[]>;
}
