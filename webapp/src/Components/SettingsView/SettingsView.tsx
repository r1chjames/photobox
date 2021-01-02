import React, { useEffect, useState } from 'react';
import './SettingsView.css';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import EditIcon from '@material-ui/icons/Edit';
import SaveIcon from '@material-ui/icons/Save';
import AddIcon from '@material-ui/icons/Add';
import CancelIcon from '@material-ui/icons/Cancel';
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
import { SettingModal } from '../SettingModal/SettingModal';

interface IProps {
  baseApiUrl: string;
}

const getAllSettings = async (baseApiUrl: string) => {
  const settingsAdapter = new SettingsAdapter(baseApiUrl);
  const allSettings: Setting[] = await settingsAdapter.getAllSettings();
  return allSettings;
};

const handleSaveSettings = async (settings: Setting[], baseApiUrl: string) => {
  const settingsAdapter = new SettingsAdapter(baseApiUrl);
  console.table(settings);
  await settingsAdapter.updateSettings(settings);
};

export const SettingsView: React.FunctionComponent<IProps> = (props) => {

  const [settings, setSettings] = useState<Setting[]>([]);
  const [editing, setEditing] = useState(false);
  const [showModal, setShowModal] = useState(false);

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

  const tableRow = (setting: Setting) => {
    if (editing) {
      return (
        <TableRow key={setting.key}>
          <TableCell component="th" scope="row">{setting.key}</TableCell>
          <TableCell>
            <TextField
              defaultValue={setting.value}
              onChange={event => handleValueChange(event, setting)}
            />
          </TableCell>
          <TableCell>
            <TextField
              defaultValue={setting.friendlyName}
              onChange={event => handleValueChange(event, setting)}
            />
          </TableCell>
          <TableCell>
            <TextField
              defaultValue={setting.category}
              onChange={event => handleValueChange(event, setting)}
            />
          </TableCell>
          <TableCell>
            <TextField
              defaultValue={setting.description}
              onChange={event => handleValueChange(event, setting)}
            />
          </TableCell>
        </TableRow>
      );
    }
    return (
      <TableRow key={setting.key}>
        <TableCell component="th" scope="row">{setting.key}</TableCell>
        <TableCell>
          {setting.value}
        </TableCell>
        <TableCell>
          {setting.friendlyName}
        </TableCell>
        <TableCell>
          {setting.category}
        </TableCell>
        <TableCell>
          {setting.description}
        </TableCell>
      </TableRow>
    );
  };

  const resetForm = () => {
    setEditing(false);
    window.location.reload(false); // not very elegant
  };

  const addCancelButton = () => {
    if (editing) {
      return (
        <CancelIcon onClick={() => resetForm()} />
      );
    }
    return (
      <EditIcon onClick={() => setEditing(true)} />
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
      <AddIcon onClick={() => setShowModal(true)} />
    );
  };

  const handleModalSave = (key: string, value: string, friendlyName: string, category: string, description: string) => {
    const updatedSettings = settings.concat(new Setting(key, value, friendlyName, category, description));
    setSettings(updatedSettings);
    setShowModal(false);
    setEditing(true);
  };

  return (
    <MuiThemeProvider>
      <MainContent title="Settings">
        <SettingModal
          isOpen={showModal}
          handleSave={handleModalSave}
          handleClose={() => setShowModal(false)}
        />
        <TableContainer component={Paper}>
          <Table aria-label="settings table">
            <TableHead>
              <TableRow>
                <TableCell>Setting</TableCell>
                <TableCell>Value</TableCell>
                <TableCell>Friendly Name</TableCell>
                <TableCell>Category</TableCell>
                <TableCell>Description</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {settings.map((setting: Setting) => (
                tableRow(setting)
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <Fab color="primary" aria-label="add" className="settingsView__addButton">
          {addCancelButton()}
        </Fab>
        <Fab color="primary" aria-label="edit" className="settingsView__saveEditButton">
          {editingButton()}
        </Fab>
      </MainContent>
    </MuiThemeProvider>
  );
};
