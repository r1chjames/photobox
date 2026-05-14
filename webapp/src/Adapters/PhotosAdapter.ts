import {IRestApiAdapter} from "./RestApiAdapter";
import {IPhotosAdapter, PhotoGeoData, TimelineEntry} from "./IPhotosAdapter";
import {Photo} from "../Models/Photo";

export class PhotosAdapter implements IPhotosAdapter {

  private restApiAdapter: IRestApiAdapter;

  constructor(restApiAdapter: IRestApiAdapter) {
    this.restApiAdapter = restApiAdapter;
  }

  private buildHeaders = (additionalHeaders: Record<string, string> = {}) => {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public getAllPhotosInfo = async (fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]> => {
    let getAllPhotosPath = `photos?fromId=${fromId}&limit=${limit}&thumbnail=${includeThumbnails}`;
    if (startDate) {
      getAllPhotosPath += `&startDate=${encodeURIComponent(startDate)}`;
    }
    if (endDate) {
      getAllPhotosPath += `&endDate=${encodeURIComponent(endDate)}`;
    }
    return this.restApiAdapter.getPaginatedApiCall(getAllPhotosPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getPhotoInfoById = async (photoId: string): Promise<Photo> => {
    const getAllPhotosPath = `photo/info/${photoId}`;
    return this.restApiAdapter.getApiCall(getAllPhotosPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getPhotosInfoInAlbum = async (albumId: string, fromId: string, limit: number, includeThumbnails: boolean, startDate?: string, endDate?: string): Promise<Photo[]> => {
    let getPhotosInAlbumPath = `photos?fromId=${fromId}&limit=${limit}&thumbnail=${includeThumbnails}`;
    if (startDate) {
      getPhotosInAlbumPath += `&startDate=${encodeURIComponent(startDate)}`;
    }
    if (endDate) {
      getPhotosInAlbumPath += `&endDate=${encodeURIComponent(endDate)}`;
    }
    const params = {
      albumId,
    };
    return this.restApiAdapter.getPaginatedApiCall(getPhotosInAlbumPath, this.buildHeaders(this.restApiAdapter.authHeader()), params);
  }

  public getPhotoCountInAlbum = async (albumId: string) => {
    const getPhotosInAlbumPath = `photos/count?albumId=${albumId}`;
    const params = {
      albumId,
    };
    return this.restApiAdapter.getApiCall(getPhotosInAlbumPath, this.buildHeaders(this.restApiAdapter.authHeader()), params);
  }

  public getPhotoImage = async (photoId: string) => {
    const getPhotosImagePath = `photo/bin/${photoId}`;
    return await this.restApiAdapter.getBinaryApiCall(getPhotosImagePath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public downloadPhoto = async (photoId: string, filename: string) => {
    const data = await this.getPhotoImage(photoId);
    const url = typeof data === 'string' ? data : URL.createObjectURL(new Blob([data], { type: 'image/jpeg' }));
    const link = document.createElement('a');
    link.href = url;
    link.download = filename;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    if (typeof data !== 'string') {
      URL.revokeObjectURL(url);
    }
  }

  public deletePhoto = async (photoId: string) => {
    const deletePhotoPath = `photos/${photoId}`;
    return this.restApiAdapter.postApiCall(deletePhotoPath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public favoritePhoto = async (photoId: string, favorite: boolean) => {
    const favoritePath = `photos/${photoId}/favorite`;
    return this.restApiAdapter.patchApiCall(favoritePath, { favorite }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public downloadPhotosAsZip = async (photoIds: string[]) => {
    const downloadPath = 'photos/download';
    const data = await this.restApiAdapter.postApiCall(downloadPath, { photoIds }, { ...this.buildHeaders(this.restApiAdapter.authHeader()), 'Content-Type': 'application/json' });
    const url = typeof data === 'string' ? data : URL.createObjectURL(new Blob([data], { type: 'application/zip' }));
    const link = document.createElement('a');
    link.href = url;
    link.download = 'photobox-download.zip';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    if (typeof data !== 'string') {
      URL.revokeObjectURL(url);
    }
  }

  public searchPhotos = async (query: string, fromId: string, limit: number): Promise<Photo[]> => {
    const searchPath = `search?q=${encodeURIComponent(query)}&fromId=${fromId}&limit=${limit}`;
    return this.restApiAdapter.getPaginatedApiCall(searchPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getTrashedPhotos = async (fromId: string, limit: number): Promise<Photo[]> => {
    const trashPath = `photos/trash?fromId=${fromId}&limit=${limit}`;
    return this.restApiAdapter.getPaginatedApiCall(trashPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public restorePhoto = async (photoId: string) => {
    const restorePath = `photos/trash/restore/${photoId}`;
    return this.restApiAdapter.postApiCall(restorePath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public rotatePhoto = async (photoId: string, direction: 'cw' | 'ccw'): Promise<Photo> => {
    const rotatePath = `photos/${photoId}/rotate?direction=${direction}`;
    return this.restApiAdapter.postApiCall(rotatePath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public getGeodata = async (north: number, south: number, east: number, west: number): Promise<PhotoGeoData[]> => {
    const geodataPath = `photos/geodata?north=${north}&south=${south}&east=${east}&west=${west}`;
    return this.restApiAdapter.getApiCall(geodataPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getTimeline = async (): Promise<TimelineEntry[]> => {
    const timelinePath = 'photos/timeline';
    return this.restApiAdapter.getApiCall(timelinePath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getAllTags = async (): Promise<string[]> => {
    const tagsPath = 'photos/tags';
    const response = await this.restApiAdapter.getApiCall(tagsPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
    return response.data || [];
  }

  public updatePhotoTags = async (photoId: string, tags: string[]): Promise<Photo> => {
    const tagsPath = `photos/${photoId}/tags`;
    return this.restApiAdapter.postApiCall(tagsPath, { tags }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public batchUpdatePhotoTags = async (photoIds: string[], tags: string[], operation: 'add' | 'remove' | 'set'): Promise<void> => {
    const batchPath = 'photos/tags/batch';
    return this.restApiAdapter.postApiCall(batchPath, { photoIds, tags, operation }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public getPhotosByTag = async (tag: string, fromId: string, limit: number, includeThumbnails: boolean): Promise<Photo[]> => {
    const path = `photos?fromId=${fromId}&limit=${limit}&thumbnail=${includeThumbnails}&tags=${encodeURIComponent(tag)}`;
    return this.restApiAdapter.getPaginatedApiCall(path, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public getDuplicatePhotos = async (): Promise<Photo[]> => {
    const path = 'photos/duplicates';
    return this.restApiAdapter.getApiCall(path, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public uploadPhoto = async (body: Record<string, unknown>) => {
    const uploadPhotoPath = `photo`;
    return this.restApiAdapter.postApiCall(uploadPhotoPath, body, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public index = async () => {
    const indexPath = 'photos/index';
    return this.restApiAdapter.postApiCall(indexPath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public regenerateThumbnails = async () => {
    const regeneratePath = 'photos/regenerate-thumbnails';
    return this.restApiAdapter.postApiCall(regeneratePath, {}, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public getPhotoThumbnails = async (photoIds: string[]): Promise<Map<string, string>> => {
    const thumbnailsPath = 'photos/thumbnails';
    const blob = await this.restApiAdapter.postBinaryApiCall(
      thumbnailsPath,
      { photoIds },
      this.buildHeaders(this.restApiAdapter.authHeader())
    );

    const JSZip = (await import('jszip')).default;
    const zip = await JSZip.loadAsync(blob);
    const thumbnailMap = new Map<string, string>();

    for (const photoId of photoIds) {
      const file = zip.file(photoId);
      if (file) {
        const thumbnailBlob = await file.async('blob');
        const url = URL.createObjectURL(thumbnailBlob);
        thumbnailMap.set(photoId, url);
      }
    }

    return thumbnailMap;
  }

  public getThumbnailUrl = (photoId: string): string => {
    return `${this.restApiAdapter.getBaseApiPath()}/photo/thumbnail/${photoId}`;
  }

  public getPhotoThumbnailBlob = async (photoId: string): Promise<Blob | string> => {
    const thumbnailPath = `photo/thumbnail/${photoId}`;
    return await this.restApiAdapter.getBinaryApiCall(
      thumbnailPath,
      this.buildHeaders(this.restApiAdapter.authHeader()),
      {}
    );
  }

  public getVideos = async (fromId: string, limit: number): Promise<Photo[]> => {
    const path = `photos?fromId=${fromId}&limit=${limit}&thumbnail=false&mediaType=video`;
    return this.restApiAdapter.getPaginatedApiCall(path, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }
}


