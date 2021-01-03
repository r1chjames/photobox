export class Photo {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _filesystemPath: string;
  private readonly _albumId: string;
  private readonly _tags: string;
  private readonly _metadata: any;

  constructor(id: string, name: string, filesystemPath: string, albumId: string, tags: string, metadata: any) {
    this._id = id;
    this._name = name;
    this._filesystemPath = filesystemPath;
    this._albumId = albumId;
    this._tags = tags;
    this._metadata = metadata;
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

  get albumId(): string {
    return this._albumId;
  }

  get tags(): string {
    return this._tags;
  }

  get metadata(): any {
    return this._metadata;
  }
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
