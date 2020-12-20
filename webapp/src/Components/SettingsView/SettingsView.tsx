import React, { useEffect, useState } from 'react';
import './SettingsView.css';
import { TitleBar } from '../TitleBar/TitleBar';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import EditIcon from '@material-ui/icons/Edit';
import SaveIcon from '@material-ui/icons/Save';
import { MainContent } from '../MainContent/MainContent';
import {
  Fab,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField
} from '@material-ui/core';
import { Setting } from '../../Models/Setting';
import { SettingsAdapter } from '../../Adapters/SettingsAdapter';

interface IProps {
  baseApiUrl: string;
}

const getAllSettings = async(baseApiUrl: string) => {
  const settingsAdapter = new SettingsAdapter(baseApiUrl);
  const allSettings: Setting[] = await settingsAdapter.getALlSettings();
  return allSettings;
};

const handleSaveSettings = async(settings: Setting[], baseApiUrl: string) => {
  const settingsAdapter = new SettingsAdapter(baseApiUrl);
  await settingsAdapter.updateSetting(settings);
};

export const SettingsView: React.FunctionComponent<IProps> = (props) => {

  const [settings, setSettings] = useState<Setting[]>([]);
  const [editing, setEditing] = useState(false);

  useEffect(() => {
    (async function retrieveAllSettings() {
      const retrievedSettings = await getAllSettings(props.baseApiUrl);
      setSettings(retrievedSettings);
    })();
  },        [setSettings, props.baseApiUrl]);

  const handleValueChange = (event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>, setting: Setting) => {
    const i = settings.findIndex(k => k.key === setting.key);
    settings[i].value = event.target.value;
  };

  const settingValue = (setting: Setting) => {
    if (editing) {
      return (
        <TextField
          defaultValue={setting.value}
          onChange={event => handleValueChange(event, setting)}
        />
      );
    }
    return (
      <div>
        {setting.value}
      </div>
    );
  };

  const editingButton = () => {
    if (editing) {
      return (
        <SaveIcon
          onClick={() => {
            handleSaveSettings(settings, props.baseApiUrl);
            setEditing(false);
          }}
        />
      );
    }
    return (
      <EditIcon
        onClick={() => setEditing(true)}
      />
    );
  };

  return (
    <MuiThemeProvider>
      <TitleBar />
      <MainContent title="Settings">
        <TableContainer component={Paper}>
          <Table aria-label="settings table">
            <TableHead>
              <TableRow>
                <TableCell>Setting</TableCell>
                <TableCell>Value</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {settings.map((setting: Setting) => (
                <TableRow key={setting.key}>
                  <TableCell component="th" scope="row">{setting.key}</TableCell>
                  <TableCell>
                    {settingValue(setting)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <Fab color="primary" aria-label="edit">
          {editingButton()}
        </Fab>
      </MainContent>
    </MuiThemeProvider>
  );
};
