export class Photo {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _filesystemPath: string;
  private readonly _albumId: string;
  private readonly _tags: string;
  private readonly _metadata: string;

  constructor(id: string, name: string, filesystemPath: string, albumId: string, tags: string, metadata: string) {
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

  get metadata(): string {
    return this._metadata;
  }
}

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
