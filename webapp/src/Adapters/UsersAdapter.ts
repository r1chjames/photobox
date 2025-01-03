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

  public login = async (user: User): Promise<string> => {
    const loginPath = "user/login";
    const body = {
      username: user.getUsername(),
      password: user.getPassword(),
    };
    return this.restApiAdapter.postApiCall(loginPath, body, this.buildHeaders());
  }

  public register = async (user: User) => {
    const registerPath = "user/register";
    const body = {
      username: user.getUsername(),
      email: user.getEmail(),
      password: user.getPassword(),
    };
    return this.restApiAdapter.postApiCall(registerPath, body, this.buildHeaders());
  }
}
