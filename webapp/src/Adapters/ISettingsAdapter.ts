import {Setting} from "../Models/Setting";

export interface ISettingsAdapter {

  getAllSettings(): Promise<Setting[]>;
  updateSettings(settings: Setting[]): Promise<any>;
  updateSetting(setting: Setting): Promise<any>;
}
