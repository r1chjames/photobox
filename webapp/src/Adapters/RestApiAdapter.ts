// eslint-disable-next-line @typescript-eslint/no-var-requires
const axios = require('axios').default;

export class RestApiAdapter {

  protected async postApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return RestApiAdapter.apiCall('post', headers, body, path, {});
  }

  protected async putApiCall(path: string, body: Record<string, unknown>, headers: Record<string, string>) {
    return RestApiAdapter.apiCall('put', headers, body, path, {});
  }

  protected async getApiCall(path: string, headers: Record<string, string>, params: Record<string, unknown>) {
    return RestApiAdapter.apiCall('get', headers, {}, path, params);
  }

  private static async apiCall(method: string, headers: Record<string, string>, body: Record<string, unknown>, url: string, params: Record<string, unknown>) {
    const options = {
      method,
      headers,
      data: JSON.stringify(body),
      url,
      params,
    };
    const resp = await axios(options);
    return await resp.data;
  }
}
