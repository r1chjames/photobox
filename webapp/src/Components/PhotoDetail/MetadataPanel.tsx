import React from 'react';
import {Group, Text, Badge, Divider, Stack} from '@mantine/core';
import {Photo} from '../../Models/Photo';

interface MetadataPanelProps {
    photo: Photo;
}

interface MetadataRow {
    label: string;
    value: string;
}

interface MetadataSection {
    title: string;
    rows: MetadataRow[];
}

/**
 * Formats a rational EXIF value like "52/1" or "334/100" into a decimal
 * string, returning the original when it cannot be parsed.
 */
function formatRational(value: unknown): string {
    if (typeof value !== 'string') return String(value ?? '');
    const match = value.split('/');
    if (match.length === 2) {
        const num = Number(match[0]);
        const den = Number(match[1]);
        if (!Number.isNaN(num) && !Number.isNaN(den) && den !== 0) {
            return (num / den).toFixed(num / den < 1 ? 2 : 1).replace(/\.0+$/, '').replace(/\.$/, '');
        }
    }
    return value;
}

/**
 * Formats GPS rational components ("52/1","2/1","24/1") into decimal degrees.
 */
function formatGps(coords: unknown, ref: unknown): string {
    if (!Array.isArray(coords) || coords.length < 3) return '';
    const deg = Number(coords[0].toString().split('/')[0]) / Number(coords[0].toString().split('/')[1] || 1);
    const min = Number(coords[1].toString().split('/')[0]) / Number(coords[1].toString().split('/')[1] || 1);
    const sec = Number(coords[2].toString().split('/')[0]) / Number(coords[2].toString().split('/')[1] || 1);
    const decimal = deg + min / 60 + sec / 3600;
    const hemisphere = typeof ref === 'string' && ref.trim() ? ` ${ref.trim()}` : '';
    return `${decimal.toFixed(6)}${hemisphere}`;
}

/**
 * Normalises EXIF date strings ("2023:01:01 21:43:14") to a friendly display.
 */
function formatExifDate(value: unknown): string {
    if (typeof value !== 'string') return String(value ?? '');
    const match = value.match(/^(\d{4}):(\d{2}):(\d{2})[ T](\d{2}):(\d{2})(?::(\d{2}))?/);
    if (!match) return value;
    const [, y, mo, d, h, mi, s] = match;
    const date = new Date(Date.UTC(Number(y), Number(mo) - 1, Number(d), Number(h), Number(mi), Number(s || 0)));
    if (Number.isNaN(date.getTime())) return value;
    return date.toLocaleString(undefined, { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
}

/**
 * Formats a byte count into a human-readable size.
 */
function formatBytes(value: unknown): string {
    const bytes = Number(value);
    if (!Number.isFinite(bytes) || bytes <= 0) return String(value ?? '');
    const units = ['B', 'KB', 'MB', 'GB'];
    let i = 0;
    let size = bytes;
    while (size >= 1024 && i < units.length - 1) {
        size /= 1024;
        i++;
    }
    return `${size.toFixed(size >= 10 || i === 0 ? 0 : 1)} ${units[i]}`;
}

/**
 * Builds the structured metadata sections for a photo from its stored
 * metadata JSONB payload and top-level photo fields.
 */
function buildSections(photo: Photo): MetadataSection[] {
    const meta = (photo.metadata ?? {}) as Record<string, any>;
    const exif = (meta.exif ?? {}) as Record<string, any>;
    const sections: MetadataSection[] = [];

    // Camera
    const cameraRows: MetadataRow[] = [];
    if (exif.Make) cameraRows.push({ label: 'Make', value: String(exif.Make) });
    if (exif.Model) cameraRows.push({ label: 'Model', value: String(exif.Model) });
    if (exif.LensModel) cameraRows.push({ label: 'Lens', value: String(exif.LensModel) });
    if (cameraRows.length) sections.push({ title: 'Camera', rows: cameraRows });

    // Settings
    const settingsRows: MetadataRow[] = [];
    if (exif.ISOSpeedRatings !== undefined) settingsRows.push({ label: 'ISO', value: Array.isArray(exif.ISOSpeedRatings) ? String(exif.ISOSpeedRatings[0]) : String(exif.ISOSpeedRatings) });
    if (exif.FNumber !== undefined) settingsRows.push({ label: 'Aperture', value: `f/${formatRational(exif.FNumber)}` });
    if (exif.ExposureTime !== undefined) settingsRows.push({ label: 'Shutter speed', value: `${formatRational(exif.ExposureTime)}s` });
    if (exif.FocalLength !== undefined) settingsRows.push({ label: 'Focal length', value: `${formatRational(exif.FocalLength)}mm` });
    if (exif.ExposureBiasValue !== undefined) settingsRows.push({ label: 'Exposure bias', value: `${formatRational(exif.ExposureBiasValue)} EV` });
    if (exif.Flash !== undefined) settingsRows.push({ label: 'Flash', value: formatRational(exif.Flash) === '1' ? 'Fired' : 'Did not fire' });
    if (settingsRows.length) sections.push({ title: 'Settings', rows: settingsRows });

    // Date
    const dateRows: MetadataRow[] = [];
    if (exif.DateTimeOriginal) dateRows.push({ label: 'Date taken', value: formatExifDate(exif.DateTimeOriginal) });
    if (photo.dateTaken) dateRows.push({ label: 'Date taken (record)', value: new Date(photo.dateTaken).toLocaleString() });
    if (exif.DateTime) dateRows.push({ label: 'Date modified', value: formatExifDate(exif.DateTime) });
    if (dateRows.length) sections.push({ title: 'Date', rows: dateRows });

    // Location
    const locationRows: MetadataRow[] = [];
    if (exif.GPSLatitude && exif.GPSLongitude) {
        locationRows.push({ label: 'Latitude', value: formatGps(exif.GPSLatitude, exif.GPSLatitudeRef) });
        locationRows.push({ label: 'Longitude', value: formatGps(exif.GPSLongitude, exif.GPSLongitudeRef) });
    }
    if (meta.latitude !== undefined && meta.longitude !== undefined && !locationRows.length) {
        locationRows.push({ label: 'Latitude', value: Number(meta.latitude).toFixed(6) });
        locationRows.push({ label: 'Longitude', value: Number(meta.longitude).toFixed(6) });
    }
    if (locationRows.length) sections.push({ title: 'Location', rows: locationRows });

    // File
    const fileRows: MetadataRow[] = [];
    if (photo.name) fileRows.push({ label: 'File name', value: photo.name });
    if (meta.size !== undefined) fileRows.push({ label: 'File size', value: formatBytes(meta.size) });
    if (meta.mime) fileRows.push({ label: 'Type', value: String(meta.mime) });
    if (meta.width && meta.height) fileRows.push({ label: 'Dimensions', value: `${meta.width} × ${meta.height}` });
    if (meta.extension) fileRows.push({ label: 'Extension', value: String(meta.extension) });
    if (meta.md5) fileRows.push({ label: 'File hash', value: String(meta.md5).slice(0, 16) + '…' });
    if (fileRows.length) sections.push({ title: 'File', rows: fileRows });

    // Orientation
    const orientationRows: MetadataRow[] = [];
    if (exif.Orientation !== undefined) orientationRows.push({ label: 'Orientation', value: formatRational(exif.Orientation) });
    if (orientationRows.length) sections.push({ title: 'Orientation', rows: orientationRows });

    return sections;
}

export const MetadataPanel: React.FC<MetadataPanelProps> = ({photo}) => {
    const sections = buildSections(photo);

    if (sections.length === 0) {
        return (
            <Text size="sm" c="dimmed" ta="center" py="md">
                No metadata available for this photo.
            </Text>
        );
    }

    return (
        <Stack gap="md">
            {sections.map(section => (
                <div key={section.title}>
                    <Group justify="space-between" mb={4}>
                        <Text size="sm" fw={600}>{section.title}</Text>
                        <Badge size="xs" variant="light" color="gray">{section.rows.length}</Badge>
                    </Group>
                    <Divider mb={6}/>
                    {section.rows.map(row => (
                        <Group key={row.label} justify="space-between" wrap="nowrap" gap="md" mb={2}>
                            <Text size="sm" c="dimmed" style={{whiteSpace: 'nowrap'}}>{row.label}</Text>
                            <Text size="sm" ta="right" style={{wordBreak: 'break-word'}}>{row.value}</Text>
                        </Group>
                    ))}
                </div>
            ))}
        </Stack>
    );
};
