import React, {useCallback, useState} from 'react';
import {ActionIcon, Button, Divider, Group, Slider, Stack, Switch, Text, Tooltip} from '@mantine/core';
import {IconPhotoOff, IconRotate} from '@tabler/icons-react';
import {Photo} from '../../Models/Photo';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {notifications} from '@mantine/notifications';

interface EditPanelProps {
    photo: Photo;
    photosAdapter: IPhotosAdapter;
    onSaved: () => void;
}

/**
 * Non-destructive edit controls (rotate, crop preset, brightness/contrast/
 * saturation, auto-enhance). Every change re-renders the edited copy server
 * side and the detail view refreshes via onSaved. Original is never touched.
 */
export const EditPanel: React.FC<EditPanelProps> = ({photo, photosAdapter, onSaved}) => {
    const [saving, setSaving] = useState(false);
    const [brightness, setBrightness] = useState(photo.editParams?.brightness ?? 0);
    const [contrast, setContrast] = useState(photo.editParams?.contrast ?? 0);
    const [saturation, setSaturation] = useState(photo.editParams?.saturation ?? 0);
    const [autoEnhance, setAutoEnhance] = useState(photo.editParams?.autoEnhance ?? false);

    const apply = useCallback(async (params: { rotate?: number; crop?: { x: number; y: number; width: number; height: number }; brightness?: number; contrast?: number; saturation?: number; autoEnhance?: boolean }) => {
        if (!photosAdapter) return;
        setSaving(true);
        try {
            await photosAdapter.editPhoto(photo.id, params);
            onSaved();
        } catch (e) {
            notifications.show({
                title: 'Edit failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setSaving(false);
        }
    }, [photosAdapter, photo.id, onSaved]);

    const handleRotate = useCallback((degrees: number) => {
        apply({rotate: ((photo.editParams?.rotate ?? 0) + degrees) % 360});
    }, [apply, photo.editParams?.rotate]);

    const handleSliders = useCallback((delta: { brightness?: number; contrast?: number; saturation?: number }) => {
        const params = {
            brightness: delta.brightness !== undefined ? delta.brightness : brightness,
            contrast: delta.contrast !== undefined ? delta.contrast : contrast,
            saturation: delta.saturation !== undefined ? delta.saturation : saturation,
            autoEnhance,
        };
        apply(params);
    }, [apply, brightness, contrast, saturation, autoEnhance]);

    const handleRevert = useCallback(async () => {
        if (!photosAdapter) return;
        setSaving(true);
        try {
            await photosAdapter.clearEdits(photo.id);
            setBrightness(0);
            setContrast(0);
            setSaturation(0);
            setAutoEnhance(false);
            onSaved();
        } catch (e) {
            notifications.show({
                title: 'Revert failed',
                message: e instanceof Error ? e.message : 'An error occurred',
                color: 'red',
            });
        } finally {
            setSaving(false);
        }
    }, [photosAdapter, photo.id, onSaved]);

    return (
        <Stack gap="sm">
            <Group justify="space-between">
                <Text size="sm" fw={600}>Edit</Text>
                {photo.editParams && (
                    <Tooltip label="Revert to original">
                        <ActionIcon size="sm" variant="subtle" color="red" onClick={handleRevert} disabled={saving} aria-label="Revert edits">
                            <IconPhotoOff size={14}/>
                        </ActionIcon>
                    </Tooltip>
                )}
            </Group>

            <Group gap="xs">
                <Button size="compact-xs" variant="light" leftSection={<IconRotate size={12}/>} onClick={() => handleRotate(90)} disabled={saving}>
                    Rotate 90°
                </Button>
                <Button size="compact-xs" variant="light" leftSection={<IconRotate size={12} style={{transform: 'scaleX(-1)'}}/>} onClick={() => handleRotate(270)} disabled={saving}>
                    Rotate 270°
                </Button>
            </Group>

            <Divider/>

            <Stack gap={4}>
                <Group justify="space-between">
                    <Text size="xs" c="dimmed">Brightness</Text>
                    <Text size="xs">{brightness > 0 ? `+${brightness}` : brightness}</Text>
                </Group>
                <Slider min={-50} max={50} step={5} value={brightness} onChange={setBrightness} onChangeEnd={(v) => handleSliders({brightness: v})} size="sm"/>
            </Stack>
            <Stack gap={4}>
                <Group justify="space-between">
                    <Text size="xs" c="dimmed">Contrast</Text>
                    <Text size="xs">{contrast > 0 ? `+${contrast}` : contrast}</Text>
                </Group>
                <Slider min={-50} max={50} step={5} value={contrast} onChange={setContrast} onChangeEnd={(v) => handleSliders({contrast: v})} size="sm"/>
            </Stack>
            <Stack gap={4}>
                <Group justify="space-between">
                    <Text size="xs" c="dimmed">Saturation</Text>
                    <Text size="xs">{saturation > 0 ? `+${saturation}` : saturation}</Text>
                </Group>
                <Slider min={-50} max={50} step={5} value={saturation} onChange={setSaturation} onChangeEnd={(v) => handleSliders({saturation: v})} size="sm"/>
            </Stack>

            <Switch
                label="Auto-enhance"
                size="sm"
                checked={autoEnhance}
                onChange={(e) => {
                    const v = e.currentTarget.checked;
                    setAutoEnhance(v);
                    apply({autoEnhance: v, brightness, contrast, saturation});
                }}
            />

            {saving && <Text size="xs" c="dimmed">Applying…</Text>}
        </Stack>
    );
};
