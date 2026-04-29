import React, { useCallback, useEffect, useState } from 'react';
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

interface SettingsTableRowProps {
    setting: Setting;
    editing: boolean;
    onChange: (setting: Setting, field: keyof Setting, value: string) => void;
}

const SettingsTableRow = React.memo(({ setting, editing, onChange }: SettingsTableRowProps) => {
    if (editing) {
        return (
            <Table.Tr key={setting.key}>
                <Table.Td>{setting.key}</Table.Td>
                <Table.Td>
                    <TextInput
                        value={setting.value}
                        onChange={event => onChange(setting, 'value', event.target.value)}
                    />
                </Table.Td>
                <Table.Td>
                    <TextInput
                        value={setting.friendlyName}
                        onChange={event => onChange(setting, 'friendlyName', event.target.value)}
                    />
                </Table.Td>
                <Table.Td>
                    <TextInput
                        value={setting.category}
                        onChange={event => onChange(setting, 'category', event.target.value)}
                    />
                </Table.Td>
                <Table.Td>
                    <TextInput
                        value={setting.description}
                        onChange={event => onChange(setting, 'description', event.target.value)}
                    />
                </Table.Td>
            </Table.Tr>
        );
    }
    return (
        <Table.Tr key={setting.key}>
            <Table.Td>{setting.key}</Table.Td>
            <Table.Td>{setting.value}</Table.Td>
            <Table.Td>{setting.friendlyName}</Table.Td>
            <Table.Td>{setting.category}</Table.Td>
            <Table.Td>{setting.description}</Table.Td>
        </Table.Tr>
    );
});

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
  },[props.settingsAdapter]);

  const handleValueChange = useCallback((setting: Setting, field: keyof Setting, value: string) => {
    setSettings(prev => prev.map(s =>
        s.key === setting.key ? { ...s, [field]: value } : s
    ));
  }, []);

  const resetForm = useCallback(async () => {
    setEditing(false);
    const refreshedSettings = await getAllSettings(props.settingsAdapter);
    setSettings(refreshedSettings);
  }, [props.settingsAdapter]);

  const handleSave = useCallback(() => {
    handleSaveSettings(settings, props.settingsAdapter);
    setEditing(false);
    setShowSnackbar(true);
  }, [settings, props.settingsAdapter]);

  const handleIndex = useCallback(async () => {
      return await props.photosAdapter.index();
  }, [props.photosAdapter]);

  const snackbar = () => {
      if (showSnackbar) {
          return(<InfoSnackbar text={'Settings saved'} show={showSnackbar} handleStopShowing={() => setShowSnackbar(false)}/>)
      }
  };

  const handleModalSave = useCallback((key: string, value: string, friendlyName: string, category: string, description: string) => {
    const updatedSettings = settings.concat({ key, value, friendlyName, category, description });
    setSettings(updatedSettings);
    setShowModal(false);
    setShowSnackbar(true);
    setEditing(true);
  }, [settings]);

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
                <SettingsTableRow
                    key={setting.key}
                    setting={setting}
                    editing={editing}
                    onChange={handleValueChange}
                />
              ))}
            </Table.Tbody>
          </Table>
        <Flex direction="row" style={{width: "100%", justifyContent: "right"}}>
            <ActionIcon color="dark" size="xl" m={"1rem"}>
              {editing ? <IconPencilCancel onClick={resetForm} /> : <IconPencil onClick={() => setEditing(true)} />}
            </ActionIcon>
            <ActionIcon color="dark" size="xl" m={"1rem"}>
              {editing ? <IconDeviceFloppy size="2.125rem" onClick={handleSave} /> : <IconLayoutGridAdd size="2.125rem" onClick={() => setShowModal(true)}/>}
            </ActionIcon>
        </Flex>
        <Button onClick={handleIndex}>
            Index
        </Button>
        {snackbar()}
    </div>
  );
};
