export class Album {

  // tslint:disable:variable-name
  private readonly _id: string;
  private readonly _name: string;
  private readonly _description: string;
  private readonly _tags: string;
  private readonly _metadata: string;

  constructor(id: string, name: string, description: string, tags: string, metadata: string) {
    this._id = id;
    this._name = name;
    this._description = description;
    this._tags = tags;
    this._metadata = metadata;
  }

  get id(): string {
    return this._id;
  }

  get name(): string {
    return this._name;
  }

  get description(): string {
    return this._description;
  }

  get tags(): string {
    return this._tags;
  }

  get metadata(): string {
    return this._metadata;
  }
}
