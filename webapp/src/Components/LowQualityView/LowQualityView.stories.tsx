import React from 'react';
import {Meta, StoryObj} from '@storybook/react';
import {LowQualityView} from './LowQualityView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {IAlbumsAdapter} from '../../Adapters/IAlbumsAdapter';

const photosAdapter = {
    getAllPhotosInfo: async () => [
        { id: 'p1', name: 'blurry_shot.jpg', filesystemPath: '', sourcePath: '', albumId: 'a1', tags: '', metadata: {}, createdAt: '2023-01-01T00:00:00Z', thumbnailUrl: '', isLowQuality: true, qualityScore: 12 },
        { id: 'p2', name: 'lens_cap.jpg', filesystemPath: '', sourcePath: '', albumId: 'a1', tags: '', metadata: {}, createdAt: '2023-01-02T00:00:00Z', thumbnailUrl: '', isLowQuality: true, qualityScore: 5 },
        { id: 'p3', name: 'overexposed.jpg', filesystemPath: '', sourcePath: '', albumId: 'a1', tags: '', metadata: {}, createdAt: '2023-01-03T00:00:00Z', thumbnailUrl: '', isLowQuality: true, qualityScore: 28 },
    ],
    getPhotosInfoInAlbum: async () => [],
    searchPhotos: async () => [],
    getVideos: async () => [],
    getPhotosByTag: async () => [],
    getPhotoInfoById: async (id: string) => ({ id, name: 'x.jpg', filesystemPath: '', sourcePath: '', albumId: '', tags: '', metadata: {}, createdAt: '', thumbnailUrl: '' }),
    getPhotoThumbnailBlob: async () => '',
    getTrashedPhotos: async () => [],
} as unknown as IPhotosAdapter;

const albumsAdapter = {
    getAlbumInfoById: async () => null,
} as unknown as IAlbumsAdapter;

const meta: Meta<typeof LowQualityView> = {
    title: 'Views/LowQualityView',
    component: LowQualityView,
    decorators: [
        (Story) => (
            <div style={{padding: '1rem'}}>
                <Story/>
            </div>
        ),
    ],
};

export default meta;
type Story = StoryObj<typeof LowQualityView>;

export const Populated: Story = {
    args: {
        photosAdapter,
        albumsAdapter,
    },
};

export const Empty: Story = {
    args: {
        photosAdapter: {
            ...photosAdapter,
            getAllPhotosInfo: async () => [],
        } as unknown as IPhotosAdapter,
        albumsAdapter,
    },
};
