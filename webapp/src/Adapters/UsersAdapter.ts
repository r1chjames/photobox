import {IRestApiAdapter} from "./RestApiAdapter";
import {IUsersAdapter} from "./IUsersAdapter";
import {User} from "../Models/User";

export class UsersAdapter implements IUsersAdapter {

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

  public login = async (user: User): Promise<Token> => {
    const loginPath = "login";
    const body = {
      username: user.username,
      password: user.password,
    };
    return this.restApiAdapter.postApiCall(loginPath, body, this.buildHeaders());
  }

  public register = async (user: User): Promise<Token> => {
    const registerPath = "user/register";
    const body = {
      username: user.username,
      email: user.email,
      password: user.password,
    };
    return this.restApiAdapter.postApiCall(registerPath, body, this.buildHeaders());
  }
}

export interface Token {
  token: string;
}
