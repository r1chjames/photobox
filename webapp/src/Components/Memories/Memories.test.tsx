import React from 'react';
import {render, screen, waitFor} from '../../test/test-utils';
import {vi} from 'vitest';
import {Memories} from './Memories';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

vi.mock('../../utils/ThumbnailUtils', () => ({
    fetchThumbnailWithAuth: vi.fn().mockResolvedValue('blob:mock'),
    getCachedThumbnail: vi.fn().mockReturnValue(null),
    revokeThumbnail: vi.fn(),
}));

const photo = (id: string, name: string) => ({
    id,
    name,
    filesystemPath: '',
    sourcePath: '',
    albumId: '',
    tags: '',
    metadata: {},
    createdAt: '2023-01-01T00:00:00Z',
    thumbnailUrl: '',
});

describe('Memories', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders memory groups with year badges', async () => {
        const adapter = {
            getMemories: vi.fn().mockResolvedValue([
                {year: 2023, yearsAgo: 3, photos: [photo('p1', 'old1.jpg')]},
                {year: 2022, yearsAgo: 4, photos: [photo('p2', 'old2.jpg')]},
            ]),
        } as unknown as IPhotosAdapter;

        render(<Memories photosAdapter={adapter} />);

        await waitFor(() => {
            expect(screen.getByText('On This Day')).toBeInTheDocument();
        });
        expect(screen.getByText('2023')).toBeInTheDocument();
        expect(screen.getByText('3 years ago')).toBeInTheDocument();
        expect(screen.getByText('2022')).toBeInTheDocument();
        expect(screen.getByText('4 years ago')).toBeInTheDocument();
        expect(screen.getByText('old1.jpg')).toBeInTheDocument();
        expect(screen.getByText('old2.jpg')).toBeInTheDocument();
    });

    it('shows empty state when no memories', async () => {
        const adapter = {
            getMemories: vi.fn().mockResolvedValue([]),
        } as unknown as IPhotosAdapter;

        render(<Memories photosAdapter={adapter} />);

        await waitFor(() => {
            expect(screen.getByText(/no memories for today/i)).toBeInTheDocument();
        });
    });
});
