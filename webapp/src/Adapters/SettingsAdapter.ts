import { Setting } from '../Models/Setting';
import {IRestApiAdapter} from "./RestApiAdapter";

export class SettingsAdapter {

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

  public getAllSettings = async () => {
    const getAllSettingsPath = "settings";
    return this.restApiAdapter.getApiCall(getAllSettingsPath, this.buildHeaders(), {});
  }

  public updateSettings = async (settings: Setting[]) => {
    const postAllSettingsPath = "settings";
    const body = {
      settings
    };
    return this.restApiAdapter.postApiCall(postAllSettingsPath, body, this.buildHeaders());
  }

  public updateSetting = async (setting: Setting) => {
    const postSettingPath = "setting";
    const body = {
      setting
    };
    return this.restApiAdapter.postApiCall(postSettingPath, body, this.buildHeaders());
  }
}
