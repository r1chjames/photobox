import React, { useCallback, useEffect, useState } from 'react';
import { Setting } from '../../Models/Setting';
import { SettingModal } from '../SettingModal/SettingModal';
import { notifications } from '@mantine/notifications';
import { modals } from '@mantine/modals';
import {JobType} from "../../Models/Job";
import {ActionIcon, Button, Flex, Skeleton, Table, TextInput} from '@mantine/core';
import {
    IconDeviceFloppy,
    IconLayoutGridAdd,
    IconPencil,
    IconPencilCancel,
    IconSettings,
    IconPlayerStop
} from "@tabler/icons-react";
import {ISettingsAdapter} from "../../Adapters/ISettingsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {EmptyState} from "../EmptyState/EmptyState";

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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [indexingRunning, setIndexingRunning] = useState(false);
  const [regenerateRunning, setRegenerateRunning] = useState(false);
  const [analysisRunning, setAnalysisRunning] = useState(false);

  const loadSettings = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const retrievedSettings = await getAllSettings(props.settingsAdapter);
      setSettings(retrievedSettings);
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to load settings');
      setSettings([]);
    } finally {
      setLoading(false);
    }
  }, [props.settingsAdapter]);

  useEffect(() => {
    loadSettings();
  }, [loadSettings]);

  const handleValueChange = useCallback((setting: Setting, field: keyof Setting, value: string) => {
    setSettings(prev => prev.map(s =>
        s.key === setting.key ? { ...s, [field]: value } : s
    ));
  }, []);

  const resetForm = useCallback(async () => {
    setEditing(false);
    await loadSettings();
  }, [loadSettings]);

  const handleSave = useCallback(async () => {
    try {
      await handleSaveSettings(settings, props.settingsAdapter);
      setEditing(false);
      notifications.show({
        title: 'Settings saved',
        message: 'Your settings have been saved successfully',
        color: 'green',
      });
    } catch (e) {
      notifications.show({
        title: 'Save failed',
        message: e instanceof Error ? e.message : 'Failed to save settings',
        color: 'red',
      });
    }
  }, [settings, props.settingsAdapter]);

  const handleIndex = useCallback(() => {
    modals.openConfirmModal({
      title: 'Start photo indexing?',
      children: 'This will scan your photo directory and may take a while depending on the number of photos.',
      labels: { confirm: 'Start indexing', cancel: 'Cancel' },
      confirmProps: { color: 'blue' },
      onConfirm: async () => {
        try {
          await props.photosAdapter.startJob(JobType.PhotoIndex);
          setIndexingRunning(true);
          notifications.show({
            title: 'Indexing started',
            message: 'Photo indexing is running in the background',
            color: 'blue',
          });
        } catch (e) {
          notifications.show({
            title: 'Indexing failed',
            message: e instanceof Error ? e.message : 'Failed to start indexing',
            color: 'red',
          });
        }
      },
    });
  }, [props.photosAdapter]);

  const handleStopIndex = useCallback(async () => {
    try {
      await props.photosAdapter.stopJob(JobType.PhotoIndex);
      setIndexingRunning(false);
      notifications.show({
        title: 'Indexing stopped',
        message: 'Photo indexing has been stopped',
        color: 'orange',
      });
    } catch (e) {
      notifications.show({
        title: 'Stop failed',
        message: e instanceof Error ? e.message : 'Failed to stop indexing',
        color: 'red',
      });
    }
  }, [props.photosAdapter]);

  const handleRegenerateThumbnails = useCallback(() => {
    modals.openConfirmModal({
      title: 'Regenerate all thumbnails?',
      children: 'This will generate thumbnails for every photo that doesn\'t have one yet. This may take a while.',
      labels: { confirm: 'Regenerate', cancel: 'Cancel' },
      confirmProps: { color: 'blue' },
      onConfirm: async () => {
        try {
          await props.photosAdapter.startJob(JobType.ThumbnailRegenerate);
          setRegenerateRunning(true);
          notifications.show({
            title: 'Regeneration started',
            message: 'Thumbnail regeneration is running in the background',
            color: 'blue',
          });
        } catch (e) {
          notifications.show({
            title: 'Regeneration failed',
            message: e instanceof Error ? e.message : 'Failed to start thumbnail regeneration',
            color: 'red',
          });
        }
      },
    });
  }, [props.photosAdapter]);

  const handleStopRegenerate = useCallback(async () => {
    try {
      await props.photosAdapter.stopJob(JobType.ThumbnailRegenerate);
      setRegenerateRunning(false);
      notifications.show({
        title: 'Regeneration stopped',
        message: 'Thumbnail regeneration has been stopped',
        color: 'orange',
      });
    } catch (e) {
      notifications.show({
        title: 'Stop failed',
        message: e instanceof Error ? e.message : 'Failed to stop thumbnail regeneration',
        color: 'red',
      });
    }
  }, [props.photosAdapter]);

  const handleAnalyze = useCallback(() => {
    modals.openConfirmModal({
      title: 'Start AI photo analysis?',
      children: 'This will analyze photos without AI tags using the configured LLM. It may take a while.',
      labels: { confirm: 'Start analysis', cancel: 'Cancel' },
      confirmProps: { color: 'blue' },
      onConfirm: async () => {
        try {
          await props.photosAdapter.startJob(JobType.AIAnalysis);
          setAnalysisRunning(true);
          notifications.show({
            title: 'Analysis started',
            message: 'AI photo analysis is running in the background',
            color: 'blue',
          });
        } catch (e) {
          notifications.show({
            title: 'Analysis failed',
            message: e instanceof Error ? e.message : 'Failed to start AI analysis',
            color: 'red',
          });
        }
      },
    });
  }, [props.photosAdapter]);

  const handleStopAnalysis = useCallback(async () => {
    try {
      await props.photosAdapter.stopJob(JobType.AIAnalysis);
      setAnalysisRunning(false);
      notifications.show({
        title: 'Analysis stopped',
        message: 'AI photo analysis has been stopped',
        color: 'orange',
      });
    } catch (e) {
      notifications.show({
        title: 'Stop failed',
        message: e instanceof Error ? e.message : 'Failed to stop AI analysis',
        color: 'red',
      });
    }
  }, [props.photosAdapter]);

  const handleModalSave = useCallback((key: string, value: string, friendlyName: string, category: string, description: string) => {
    const updatedSettings = settings.concat({ key, value, friendlyName, category, description });
    setSettings(updatedSettings);
    setShowModal(false);
    notifications.show({
      title: 'Setting added',
      message: `Added ${friendlyName}`,
      color: 'green',
    });
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
              {loading
                ? Array.from({ length: 5 }).map((_, i) => (
                    <Table.Tr key={i}>
                      {Array.from({ length: 5 }).map((_, j) => (
                        <Table.Td key={j}>
                          <Skeleton height={20} radius="sm" />
                        </Table.Td>
                      ))}
                    </Table.Tr>
                  ))
                : error && settings.length === 0
                  ? (
                    <Table.Tr>
                      <Table.Td colSpan={5}>
                        <EmptyState
                          title="No settings available"
                          description="Settings could not be loaded. Check your connection and try again."
                          icon={<IconSettings size="2rem" />}
                          action={{ label: 'Retry', onClick: loadSettings }}
                        />
                      </Table.Td>
                    </Table.Tr>
                  )
                  : settings.map((setting: Setting) => (
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
            <ActionIcon variant="subtle" size="xl" m={"1rem"} onClick={editing ? resetForm : () => setEditing(true)}>
              {editing ? <IconPencilCancel size="1.5rem" /> : <IconPencil size="1.5rem" />}
            </ActionIcon>
            <ActionIcon variant="subtle" size="xl" m={"1rem"} onClick={editing ? handleSave : () => setShowModal(true)}>
              {editing ? <IconDeviceFloppy size="1.5rem" /> : <IconLayoutGridAdd size="1.5rem"/>}
            </ActionIcon>
        </Flex>
        <Flex direction="row" gap="md" style={{ width: "100%", justifyContent: "right" }}>
          {indexingRunning ? (
            <Button onClick={handleStopIndex} variant="outline" color="red" leftSection={<IconPlayerStop size={16} />}>
              Stop Indexing
            </Button>
          ) : (
            <Button onClick={handleIndex} variant="outline">
              Re-index Photos
            </Button>
          )}
          {regenerateRunning ? (
            <Button onClick={handleStopRegenerate} color="red" leftSection={<IconPlayerStop size={16} />}>
              Stop Regeneration
            </Button>
          ) : (
            <Button onClick={handleRegenerateThumbnails}>
              Regenerate Thumbnails
            </Button>
          )}
          {analysisRunning ? (
            <Button onClick={handleStopAnalysis} color="red" leftSection={<IconPlayerStop size={16} />}>
              Stop Analysis
            </Button>
          ) : (
            <Button onClick={handleAnalyze} variant="outline" color="teal">
              Analyze Photos
            </Button>
          )}
        </Flex>
    </div>
  );
};
