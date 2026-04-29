import React, { useEffect, useState } from 'react';
import { Setting } from '../../Models/Setting';
import { SettingModal } from '../SettingModal/SettingModal';
import {InfoSnackbar} from '../Snackbar/InfoSnackbar';
import {ActionIcon, Button, Flex, Table, TextInput} from '@mantine/core';
import {
    IconDeviceFloppy,
    IconLayoutGridAdd,
    IconPencil,
    IconPencilCancel
} from "@tabler/icons-react";
import {ISettingsAdapter} from "../../Adapters/ISettingsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";

interface IProps {
    settingsAdapter: ISettingsAdapter;
    photosAdapter: IPhotosAdapter;
}

const getAllSettings = async (settingsAdapter: ISettingsAdapter) => {
  return settingsAdapter.getAllSettings();
};

const handleSaveSettings = async (settings: Setting[], settingsAdapter: ISettingsAdapter) => {
  await settingsAdapter.updateSettings(settings);
};

export const SettingsView: React.FunctionComponent<IProps> = (props) => {

  const [settings, setSettings] = useState<Setting[]>([]);
  const [editing, setEditing] = useState(false);
  const [showModal, setShowModal] = useState(false);
  const [showSnackbar, setShowSnackbar] = useState(false);

  useEffect(() => {
    (async function retrieveAllSettings() {
      const retrievedSettings = await getAllSettings(props.settingsAdapter);
      setSettings(retrievedSettings);
    })();
  },[setSettings]);

  const handleValueChange = (event: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>, setting: Setting) => {
    const i = settings.findIndex(k => k.key === setting.key);
    settings[i].value = event.target.value;
  };

  const tableRow = (setting: Setting) => {
    if (editing) {
      return (
        <Table.Tr>
          <Table.Td>
            {setting.key}
          </Table.Td>
          <Table.Td>
            <TextInput
              defaultValue={setting.value}
              onChange={event => handleValueChange(event, setting)}
            />
          </Table.Td>
          <Table.Td>
            <TextInput
              defaultValue={setting.friendlyName}
              onChange={event => handleValueChange(event, setting)}
            />
          </Table.Td>
          <Table.Td>
            <TextInput
              defaultValue={setting.category}
              onChange={event => handleValueChange(event, setting)}
            />
          </Table.Td>
          <Table.Td>
            <TextInput
              defaultValue={setting.description}
              onChange={event => handleValueChange(event, setting)}
            />
          </Table.Td>
        </Table.Tr>
      );
    }
    return (
      <Table.Tr>
        <Table.Td>
          {setting.key}
        </Table.Td>
        <Table.Td>
          {setting.value}
        </Table.Td>
        <Table.Td>
          {setting.friendlyName}
        </Table.Td>
        <Table.Td>
          {setting.category}
        </Table.Td>
        <Table.Td>
          {setting.description}
        </Table.Td>
      </Table.Tr>
    );
  };

  const resetForm = async () => {
    setEditing(false);
    const refreshedSettings = await getAllSettings(props.settingsAdapter);
    setSettings(refreshedSettings);
  };

  const addCancelButton = () => {
    if (editing) {
      return (
        <IconPencilCancel onClick={() => resetForm()} />
      );
    }
    return (
      <IconPencil onClick={() => setEditing(true)} />
    );
  };

  const editingButton = () => {
    if (editing) {
      return (
          <IconDeviceFloppy size="2.125rem"
              onClick={() => {
                handleSaveSettings(settings, props.settingsAdapter);
                setEditing(false);
                setShowSnackbar(true);
              }}
          />
      );
    }
    return (
        <IconLayoutGridAdd size="2.125rem" onClick={() => setShowModal(true)}/>
    );
  };

  const handleIndex = async () => {
      return await props.photosAdapter.index();
  }

  const snackbar = () => {
      if (showSnackbar) {
          return(<InfoSnackbar text={'Settings saved'} show={showSnackbar} handleStopShowing={() => setShowSnackbar(false)}/>)
      }
  };

  const handleModalSave = (key: string, value: string, friendlyName: string, category: string, description: string) => {
    const updatedSettings = settings.concat({ key, value, friendlyName, category, description });
    setSettings(updatedSettings);
    setShowModal(false);
    setShowSnackbar(true);
    setEditing(true);
  };

  return (
    <div>
        <SettingModal
          isOpen={showModal}
          handleSave={handleModalSave}
          handleClose={() => setShowModal(false)}
        />
        <Table>
            <Table.Thead>
                <Table.Tr>
                    <Table.Th>Setting</Table.Th>
                    <Table.Th>Value</Table.Th>
                    <Table.Th>Friendly Name</Table.Th>
                    <Table.Th>Category</Table.Th>
                    <Table.Th>Description</Table.Th>
                </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {settings.map((setting: Setting) => (
                tableRow(setting)
              ))}
            </ Table.Tbody>
          </Table>
        <Flex direction="row" style={{width: "100%", justifyContent: "right"}}>
            <ActionIcon color="dark" size="xl" m={"1rem"}>
              {addCancelButton()}
            </ActionIcon>
            <ActionIcon color="dark" size="xl" m={"1rem"}>
              {editingButton()}
            </ActionIcon>
        </Flex>
        <Button onClick={() => handleIndex()}>
            Index
        </Button>
        {snackbar()}
    </div>
  );
};
