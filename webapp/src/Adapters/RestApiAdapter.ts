// tslint:disable-next-line:no-var-requires
const axios = require('axios');

export class RestApiAdapter {

  protected async postApiCall(path: string, body: {}, headers: {}) {
    return RestApiAdapter.apiCall('post', headers, body, path, {});
  }

  protected async putApiCall(path: string, body: {}, headers: {}) {
    return RestApiAdapter.apiCall('put', headers, body, path, {});
  }

  protected async getApiCall(path: string, headers: {}, params: {}) {
    return RestApiAdapter.apiCall('get', headers, {}, path, params);
  }

  private static async apiCall(method: string, headers: {}, body: {}, url: string, params: {}) {
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
