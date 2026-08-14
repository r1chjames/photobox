import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {MetadataPanel} from './MetadataPanel';

const fullMetadataPhoto = {
    id: 'p1',
    name: '20230101_214314_0459.jpg',
    filesystemPath: '/photos/Album1/20230101_214314_0459.jpg',
    sourcePath: '/photos/Album1/20230101_214314_0459.jpg',
    albumId: 'album1',
    tags: 'holiday,beach',
    metadata: {
        exif: {
            Make: 'Google',
            Model: 'Google Pixel 9',
            DateTimeOriginal: '2023:01:01 21:43:14',
            FNumber: '17/10',
            ISOSpeedRatings: [52],
            ExposureTime: '1/200',
            FocalLength: '24/10',
            ExposureBiasValue: '0/1',
            GPSLatitude: ['52/1', '2/1', '24/1'],
            GPSLongitude: ['0/1', '5/1', '40/1'],
            GPSLatitudeRef: 'N',
            GPSLongitudeRef: 'E',
            Orientation: '1',
        },
        mime: 'image/jpeg',
        size: 72731,
        width: 2400,
        height: 1200,
        extension: '.jpg',
        md5: '1031db5affaee4fb204553478faec521',
    },
    createdAt: '2023-01-01T21:43:14Z',
    thumbnailUrl: '',
};

const minimalMetadataPhoto = {
    id: 'p2',
    name: 'scan_001.png',
    filesystemPath: '/photos/scan_001.png',
    sourcePath: '/photos/scan_001.png',
    albumId: 'album1',
    tags: '',
    metadata: {
        mime: 'image/png',
        size: 1024,
        width: 100,
        height: 100,
    },
    createdAt: '2024-01-01T00:00:00Z',
    thumbnailUrl: '',
};

const noMetadataPhoto = {
    id: 'p3',
    name: '',
    filesystemPath: '',
    sourcePath: '',
    albumId: '',
    tags: '',
    metadata: {},
    createdAt: '2024-01-01T00:00:00Z',
    thumbnailUrl: '',
};

const meta: Meta<typeof MetadataPanel> = {
    title: 'PhotoDetail/MetadataPanel',
    component: MetadataPanel,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem', maxWidth: 420}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof MetadataPanel>;

export const FullMetadata: Story = {
    args: {
        photo: fullMetadataPhoto as never,
    },
};

export const Minimal: Story = {
    args: {
        photo: minimalMetadataPhoto as never,
    },
};

export const Empty: Story = {
    args: {
        photo: noMetadataPhoto as never,
    },
};
