import React from 'react';
import {render, screen, waitFor} from '../../test/test-utils';
import {vi} from 'vitest';
import {TrashView} from './TrashView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {ISettingsAdapter} from '../../Adapters/ISettingsAdapter';

vi.mock('../../utils/ThumbnailUtils', () => ({
    fetchThumbnailWithAuth: vi.fn().mockResolvedValue('blob:mock'),
    getCachedThumbnail: vi.fn().mockReturnValue(null),
    revokeThumbnail: vi.fn(),
}));

const daysAgo = (days: number) => new Date(Date.now() - days * 24 * 60 * 60 * 1000).toISOString();

const photosAdapter = {
    getTrashedPhotos: vi.fn(),
    restorePhoto: vi.fn().mockResolvedValue(undefined),
    deletePhoto: vi.fn().mockResolvedValue(undefined),
} as unknown as IPhotosAdapter;

const settingsAdapter = {
    getAllSettings: vi.fn().mockResolvedValue([{ key: 'trash_retention_days', value: '30', friendlyName: '', category: '', description: '' }]),
} as unknown as ISettingsAdapter;

describe('TrashView', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('shows empty state when trash is empty', async () => {
        photosAdapter.getTrashedPhotos = vi.fn().mockResolvedValue([]);
        render(<TrashView photosAdapter={photosAdapter} settingsAdapter={settingsAdapter} />);
        await waitFor(() => {
            expect(screen.getByText('Trash is empty')).toBeInTheDocument();
        });
        expect(screen.getByText(/automatically deleted after 30 days/i)).toBeInTheDocument();
    });

    it('shows trashed photos with retention countdown', async () => {
        photosAdapter.getTrashedPhotos = vi.fn().mockResolvedValue([
            { id: 'p1', name: 'photo1.jpg', filesystemPath: '', sourcePath: '', albumId: '', tags: '', metadata: {}, createdAt: '', thumbnailUrl: '', deletedAt: daysAgo(10) },
        ]);
        render(<TrashView photosAdapter={photosAdapter} settingsAdapter={settingsAdapter} />);
        await waitFor(() => {
            expect(screen.getByText('photo1.jpg')).toBeInTheDocument();
        });
        // 10 days ago + 30-day retention => ~20 days left
        expect(screen.getByText('20d')).toBeInTheDocument();
        expect(screen.getByText(/automatically deleted after 30 days/i)).toBeInTheDocument();
    });

    it('hides countdown when retention is disabled', async () => {
        photosAdapter.getTrashedPhotos = vi.fn().mockResolvedValue([
            { id: 'p1', name: 'photo1.jpg', filesystemPath: '', sourcePath: '', albumId: '', tags: '', metadata: {}, createdAt: '', thumbnailUrl: '', deletedAt: daysAgo(10) },
        ]);
        const disabledSettings = {
            getAllSettings: vi.fn().mockResolvedValue([{ key: 'trash_retention_days', value: '0', friendlyName: '', category: '', description: '' }]),
        } as unknown as ISettingsAdapter;
        render(<TrashView photosAdapter={photosAdapter} settingsAdapter={disabledSettings} />);
        await waitFor(() => {
            expect(screen.getByText('photo1.jpg')).toBeInTheDocument();
        });
        expect(screen.queryByText(/^\d+d$/)).not.toBeInTheDocument();
        expect(screen.getByText(/automatic cleanup is disabled/i)).toBeInTheDocument();
    });
});
