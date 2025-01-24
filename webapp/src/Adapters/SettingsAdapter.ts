import { Setting } from '../Models/Setting';
import {IRestApiAdapter} from "./RestApiAdapter";
import {ISettingsAdapter} from "./ISettingsAdapter";

export class SettingsAdapter implements ISettingsAdapter {

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

  public getAllSettings = async (): Promise<Setting[]> => {
    const getAllSettingsPath = "settings";
    return this.restApiAdapter.getApiCall(getAllSettingsPath, this.buildHeaders(this.restApiAdapter.authHeader()), {});
  }

  public updateSettings = async (settings: Setting[]) => {
    const postAllSettingsPath = "settings";
    const body = {
      settings
    };
    return this.restApiAdapter.postApiCall(postAllSettingsPath, body, this.buildHeaders(this.restApiAdapter.authHeader()));
  }

  public updateSetting = async (setting: Setting) => {
    const postSettingPath = "setting";
    const body = {
      setting
    };
    return this.restApiAdapter.postApiCall(postSettingPath, body, this.buildHeaders(this.restApiAdapter.authHeader()));
  }
}
