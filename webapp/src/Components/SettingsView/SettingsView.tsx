import React, { useEffect, useState } from 'react';
import './SettingsView.css';
import { TitleBar } from '../TitleBar/TitleBar';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import { MainContent } from '../MainContent/MainContent';
import { Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core';
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

export const SettingsView: React.FunctionComponent<IProps> = (props) => {

  const [settings, setSettings] = useState<Setting[]>([]);

  useEffect(() => {
    (async function retrieveAllSettings() {
      const retrievedSettings = await getAllSettings(props.baseApiUrl);
      setSettings(retrievedSettings);
    })();
  },        [setSettings, props.baseApiUrl]);

  const renderSettings = () => {
    return (
      <div>
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
                      <TableCell>{setting.value}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </TableContainer>
          </MainContent>
        </div>
    );
  };

  return (
    <MuiThemeProvider>
      {renderSettings()}
    </MuiThemeProvider>
  );
};
