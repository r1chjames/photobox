import axios from 'axios';

export interface IRestApiAdapter {

  postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>
  putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>): Promise<any>
  getApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>): Promise<any>

}

export class RestApiAdapter implements IRestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    this.baseApiPath = baseApiPath;
  }

  async postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return this.apiCall('post', headers, body, path, {});
  }

  async putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return this.apiCall('put', headers, body, path, {});
  }

  async getApiCall(path: string, headers: Record<string, string>, params: Record<string, any>) {
    return this.apiCall('get', headers, {}, path, params);
  }

  private async apiCall(method: string, headers: Record<string, string>, body: Record<string, any>, url: string, params: Record<string, any>) {
    const parsedUrl = `${this.baseApiPath}/${url}`;
    const options = {
      method,
      headers,
      data: JSON.stringify(body),
      url: parsedUrl,
      params,
    };
    const resp = await axios(options);
    return resp.data;
  }
}

export class MockRestApiAdapter implements IRestApiAdapter {

  public resp;

  constructor(resp: any) {
    this.resp = resp;
  }


  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return this.resp;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return this.resp;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  async getApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>) {
    return this.resp;
  }
}
