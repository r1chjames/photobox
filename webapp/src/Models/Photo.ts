export class Photo {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _filesystemPath: string;
  private readonly _sourcePath: string;
  private readonly _thumbnailPath: string;
  private readonly _albumId: string;
  private readonly _tags: string;
  private readonly _metadata: Record<string, Record<string, string>>;
  private readonly _createdAt: string;

  constructor(id: string, name: string, filesystemPath: string, sourcePath: string, thumbnailPath: string, albumId: string, tags: string,
              metadata: Record<string, Record<string, string>>, createdAt: string, apiBasePath: string) {
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

  get metadata(): Record<string, unknown> {
    return this._metadata;
  }

  get createdAt(): string {
    return this._createdAt;
  }

  getPhotoDate(): string {
    const exifVal = this._metadata.exif;
    // @ts-ignore
    return exifVal !== null ? exifVal.DateTime : this._createdAt;
  };
}

// tslint:disable-next-line:max-classes-per-file
// export class Metadata {
//   // tslint:disable:variable-name
//   private readonly _directory: string;
//   private readonly _exif: any;
//   private readonly _extension: string;
//   private readonly _id: string;
//   private readonly _md5: string;
//   private readonly _mime: string;
//   private readonly _name: string;
//   private readonly _path: string;
//   private readonly _size: number;
//   private readonly _thumbnail: string | null;
//
//   constructor(directory: string, exif: string, extension: string, id: string, md5: string, mime: string, name: string,
//               path: string, size: number, thumbnail: string | null) {
//     this._directory = directory;
//     this._exif = exif;
//     this._extension = extension;
//     this._id = id;
//     this._md5 = md5;
//     this._mime = mime;
//     this._name = name;
//     this._path = path;
//     this._size = size;
//     this._thumbnail = thumbnail;
//   }
// }

// tslint:disable-next-line:max-classes-per-file
export class PhotoCount {
  // tslint:disable:variable-name
  private readonly _photoCount: number;

  constructor(photoCount: number) {
    this._photoCount = photoCount;
  }

  get photoCount(): number {
    return this._photoCount;
  }
}

export const previousPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex > 0 ? photos![currentPhotoIndex - 1] : photos![0];
}

export const nextPhotoInAlbum = (photos: Photo[], currentPhotoIndex: number) => {
  return currentPhotoIndex < photos!.length ? photos![currentPhotoIndex + 1] : photos![currentPhotoIndex];
}