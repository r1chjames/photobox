export class User {

  // tslint:disable:variable-name
  private readonly _username: string;
  private readonly _email: string;
  private readonly _password: string;

  constructor(username: string, email: string, password: string) {
    this._username = username;
    this._email = email;
    this._password = password;
  }

  get username(): string {
    return this._username;
  }

  get email(): string {
    return this._email;
  }

  get password(): string {
    return this._password;
  }

}
