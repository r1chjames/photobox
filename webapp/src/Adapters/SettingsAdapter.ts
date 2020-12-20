import { RestApiAdapter } from './RestApiAdapter';
import { Setting } from '../Models/Setting';

export class SettingsAdapter extends RestApiAdapter {

  private readonly baseApiPath: string;

  constructor(baseApiPath: string) {
    super();
    this.baseApiPath = `${baseApiPath}`;
  }

  private buildHeaders(additionalHeaders: {} = {}) {
    const standardHeaders = {
      'Content-Type': 'application/json'
    };
    return { ...standardHeaders, ...additionalHeaders };
  }

  public async getALlSettings() {
    const getAllSettingsPath = `${this.baseApiPath}/settings`;
    return this.getApiCall(getAllSettingsPath, '', this.buildHeaders(), {});
  }

  public async updateSettings(settings: Setting[]) {
    const postAllSettingsPath = `${this.baseApiPath}/settings`;
    const body = {
      settings
    };
    return this.postApiCall(postAllSettingsPath, body, this.buildHeaders());
  }

  public async updateSetting(setting: Setting) {
    const postSettingPath = `${this.baseApiPath}/setting`;
    const body = {
      setting
    };
    return this.postApiCall(postSettingPath, body, this.buildHeaders());
  }
}
