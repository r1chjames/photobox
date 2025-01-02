import {IRestApiAdapter} from "./RestApiAdapter";
import {IUsersAdapter} from "./IUsersAdapter";

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

  public login = async (): Promise<string> => {
    const loginPath = "/api/login";
    return this.restApiAdapter.postApiCall(loginPath, this.buildHeaders(), {});
  }

  public register = async () => {
    const registerPath = "/api/register";
    return this.restApiAdapter.postApiCall(registerPath, this.buildHeaders(), {});
  }
}
