import {ISettingsAdapter} from "./ISettingsAdapter";
import {Setting} from "../Models/Setting";

export class MockSettingsAdapter implements ISettingsAdapter {

  private _settings: Setting[] = [];

  public withSettings(settings: Setting[]): this {
    this._settings = settings
    return this;
  }
  public getAllSettings = async (): Promise<Setting[]> => {
    return this._settings;
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public updateSettings = async (settings: Setting[]) => {
    return {};
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public updateSetting = async (setting: Setting) => {
    return {};
  }
}
