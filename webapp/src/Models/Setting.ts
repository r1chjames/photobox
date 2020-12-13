export class Setting {

  // tslint:disable:variable-name
  private readonly _key: string;
  private _value: string;

  constructor(key: string, value: string) {
    this._key = key;
    this._value = value;
  }

  get key(): string {
    return this._key;
  }

  get value(): string {
    return this._value;
  }

  set value(value: string) {
    this._value = value;
  }
}
