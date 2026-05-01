import {IRestApiAdapter} from "./RestApiAdapter";
import {ISharesAdapter, ShareLink} from "./ISharesAdapter";

export class SharesAdapter implements ISharesAdapter {

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

  public createShare = async (resourceType: 'photo' | 'album', resourceId: string, expiry?: string, password?: string): Promise<ShareLink> => {
    const createSharePath = 'share';
    return this.restApiAdapter.postApiCall(createSharePath, { resourceType, resourceId, expiry, password }, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public getShares = async (): Promise<ShareLink[]> => {
    const getSharesPath = 'shares';
    return this.restApiAdapter.getApiCall(getSharesPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public revokeShare = async (token: string): Promise<void> => {
    const revokeSharePath = `shares/${token}`;
    return this.restApiAdapter.deleteApiCall(revokeSharePath, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}
