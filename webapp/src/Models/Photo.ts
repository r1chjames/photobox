export class Photo {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _filesystemPath: string;
  private readonly _sourcePath: string;
  private readonly _albumId: string;
  private readonly _tags: string;
  private readonly _metadata: Record<string, any>[];
  private readonly _createdAt: string;
  private readonly _thumbnail: string;


  constructor(id: string, name: string, filesystemPath: string, sourcePath: string, thumbnailPath: string, albumId: string, tags: string,
              metadata: Record<string, any>[], createdAt: string, thumbnail: string) {
    this._id = id;
    this._name = name;
    this._filesystemPath = filesystemPath;
    this._sourcePath = sourcePath;
    this._albumId = albumId;
    this._tags = tags;
    this._metadata = metadata;
    this._createdAt = createdAt;
    this._thumbnail = thumbnail;
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


  get thumbnail(): string {
    return this._thumbnail;
  }

  getPhotoDate(): string {
    const exifVal = this._metadata?.find((k) => k && typeof k === 'object' && 'DateTime' in k);
    return exifVal?.DateTime ?? this._createdAt;
  };
}

export const previousPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex > 0 ? photos[currentPhotoIndex - 1] : photos[0];
}

export const nextPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex < photos.length - 1 ? photos[currentPhotoIndex + 1] : photos[currentPhotoIndex];
}