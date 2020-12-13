// tslint:disable-next-line:no-var-requires
const axios = require('axios');

export class RestApiAdapter {

  protected async postApiCall(path: string, body: {}, headers: {}) {
    return RestApiAdapter.apiCall('post', headers, body, path, {});
  }

  protected async putApiCall(path: string, body: {}, headers: {}) {
    return RestApiAdapter.apiCall('put', headers, body, path, {});
  }

  protected async getApiCall(path: string, body: {}, headers: {}, params: {}) {
    return RestApiAdapter.apiCall('get', headers, body, path, params);
  }

  private static async apiCall(method: string, headers: {}, data: {}, url: string, params: {}) {
    const options = {
      method,
      headers,
      data: JSON.stringify(data),
      url,
      params,
    };
    const resp = await axios(options);
    return await resp.data;
  }
}
