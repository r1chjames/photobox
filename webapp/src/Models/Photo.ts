export class Photo {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _filesystemPath: string;
  private readonly _sourcePath: string;
  private readonly _thumbnailPath: string;
  private readonly _albumId: string;
  private readonly _tags: string;
  private readonly _metadata: Record<string, any>[];
  private readonly _createdAt: string;

  constructor(id: string, name: string, filesystemPath: string, sourcePath: string, thumbnailPath: string, albumId: string, tags: string,
              metadata: Record<string, any>[], createdAt: string, apiBasePath: string) {
    this._id = id;
    this._name = name;
    this._filesystemPath = filesystemPath;
    this._sourcePath = `${apiBasePath}/${sourcePath}`;
    this._albumId = albumId;
    this._tags = tags;
    this._metadata = metadata;
    this._createdAt = createdAt;
    this._thumbnailPath = `${apiBasePath}/${thumbnailPath}`
  }

  get id(): string {
    return this._id;
  }

  get name(): string {
    return this._name;
  }

  get filesystemPath(): string {
    return this._filesystemPath;
  }

  get sourcePath(): string {
    return this._sourcePath;
  }

  get thumbnailPath(): string {
    return this._thumbnailPath;
  }

  get albumId(): string {
    return this._albumId;
  }

  get tags(): string {
    return this._tags;
  }

  get metadata(): Record<string, any>[] {
    return this._metadata;
  }

  get createdAt(): string {
    return this._createdAt;
  }

  getPhotoDate(): string {
    const exifVal = this._metadata.find((k) => k === 'exif');
    return exifVal !== null ? exifVal.DateTime : this._createdAt;
  };
}

export const previousPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex > 0 ? photos![currentPhotoIndex - 1] : photos![0];
}

export const nextPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex < photos!.length ? photos![currentPhotoIndex + 1] : photos![currentPhotoIndex];
}