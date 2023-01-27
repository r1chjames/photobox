import React, { useEffect, useState } from 'react';
import './SettingsView.css';
import { Setting } from '../../Models/Setting';
import { SettingsAdapter } from '../../Adapters/SettingsAdapter';
import { SettingModal } from '../SettingModal/SettingModal';
import {InfoSnackbar} from '../Snackbar/InfoSnackbar';
import {ActionIcon, MantineProvider, Table, TextInput} from '@mantine/core';
import {MdAdd, MdCancel, MdSave} from 'react-icons/md';
import {ImPencil} from 'react-icons/im';

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
  await settingsAdapter.updateSettings(settings);
};

export const SettingsView: React.FunctionComponent<IProps> = (props) => {

  const [settings, setSettings] = useState<Setting[]>([]);
  const [editing, setEditing] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [showSnackbar, setShowSnackbar] = useState(false);

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
        <tr key={setting.key}>
          <td>
            {setting.key}
          </td>
          <td>
            <TextInput
              defaultValue={setting.value}
              onChange={event => handleValueChange(event, setting)}
            />
          </td>
          <td>
            <TextInput
              defaultValue={setting.friendlyName}
              onChange={event => handleValueChange(event, setting)}
            />
          </td>
          <td>
            <TextInput
              defaultValue={setting.category}
              onChange={event => handleValueChange(event, setting)}
            />
          </td>
          <td>
            <TextInput
              defaultValue={setting.description}
              onChange={event => handleValueChange(event, setting)}
            />
          </td>
        </tr>
      );
    }
    return (
      <tr key={setting.key}>
        <td>
          {setting.key}
        </td>
        <td>
          {setting.value}
        </td>
        <td>
          {setting.friendlyName}
        </td>
        <td>
          {setting.category}
        </td>
        <td>
          {setting.description}
        </td>
      </tr>
    );
  };

  const resetForm = () => {
    setEditing(false);
    window.location.reload(); // TODO: not very elegant
  };

  const addCancelButton = () => {
    if (editing) {
      return (
        <MdCancel onClick={() => resetForm()} />
      );
    }
    return (
      <ImPencil onClick={() => setEditing(true)} />
    );
  };

  const editingButton = () => {
    if (editing) {
      return (
        <MdSave
          onClick={() => {
            handleSaveSettings(settings, props.baseApiUrl);
            setEditing(false);
            setShowSnackbar(true);
          }}
        />
      );
    }
    return (
      <MdAdd onClick={() => setShowModal(true)} />
    );
  };

  const handleModalSave = (key: string, value: string, friendlyName: string, category: string, description: string) => {
    const updatedSettings = settings.concat(new Setting(key, value, friendlyName, category, description));
    setSettings(updatedSettings);
    setShowModal(false);
    setShowSnackbar(true);
    setEditing(true);
  };

  return (
    <MantineProvider>
        <SettingModal
          isOpen={showModal}
          handleSave={handleModalSave}
          handleClose={() => setShowModal(false)}
        />
          <Table aria-label="settings table">
            <thead>
              <tr>
                <th>Setting</th>
                <th>Value</th>
                <th>Friendly Name</th>
                <th>Category</th>
                <th>Description</th>
              </tr>
            </thead>
            <tbody>
              {settings.map((setting: Setting) => (
                tableRow(setting)
              ))}
            </tbody>
          </Table>
        <InfoSnackbar text={'Settings saved'} show={showSnackbar} handleStopShowing={() => setShowSnackbar(false)}/>
        <ActionIcon color="primary" aria-label="add" className="settingsView__addButton">
          {addCancelButton()}
        </ActionIcon>
        <ActionIcon color="primary" aria-label="edit" className="settingsView__saveEditButton">
          {editingButton()}
        </ActionIcon>
    </MantineProvider>
  );
};
