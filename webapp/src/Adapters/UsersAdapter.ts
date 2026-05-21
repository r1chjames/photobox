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

  public getAllUsers = async (): Promise<User[]> => {
    const response = await this.restApiAdapter.getApiCall("users", this.buildHeaders(this.restApiAdapter.authHeader()), {});
    const users = response?.users;
    return Array.isArray(users) ? users : [];
  }

  public updateUser = async (userId: string, updates: Partial<User>): Promise<User> => {
    return this.restApiAdapter.patchApiCall(`users/${userId}`, updates as Record<string, unknown>, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public updateSelf = async (updates: Partial<User>): Promise<User> => {
    return this.restApiAdapter.postApiCall('user/update', updates as Record<string, unknown>, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public deleteUser = async (userId: string): Promise<void> => {
    return this.restApiAdapter.deleteApiCall(`users/${userId}`, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}

export interface Token {
  token: string;
}
