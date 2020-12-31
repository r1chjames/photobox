export class Setting {

  // tslint:disable:variable-name
  public readonly key: string;
  public value: string;
  public friendlyName: string;
  public category: string;
  public description: string;

  constructor(key: string, value: string, friendlyName: string, category: string, description: string) {
    this.key = key;
    this.value = value;
    this.friendlyName = friendlyName;
    this.category = category;
    this.description = description;
  }

  public getKey(): string {
    return this.key;
  }

  public getValue(): string {
    return this.value;
  }

  public setValue(value: string) {
    this.value = value;
  }
}
