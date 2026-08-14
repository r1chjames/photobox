import React from 'react';
import {render, screen, fireEvent, waitFor} from '../../test/test-utils';
import {vi} from 'vitest';
import {EditPanel} from './EditPanel';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

const basePhoto = {
    id: 'p1',
    name: 'photo.jpg',
    filesystemPath: '/photos/photo.jpg',
    sourcePath: '/photos/photo.jpg',
    albumId: '',
    tags: '',
    metadata: {},
    createdAt: '2023-01-01T00:00:00Z',
    thumbnailUrl: '',
};

describe('EditPanel', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('rotates 90° via the adapter', async () => {
        const adapter = {
            editPhoto: vi.fn().mockResolvedValue(basePhoto),
        } as unknown as IPhotosAdapter;
        const onSaved = vi.fn();

        render(<EditPanel photo={basePhoto as never} photosAdapter={adapter} onSaved={onSaved} />);

        fireEvent.click(screen.getByRole('button', {name: /rotate 90/i}));

        await waitFor(() => {
            expect(adapter.editPhoto).toHaveBeenCalledWith('p1', {rotate: 90});
        });
        await waitFor(() => {
            expect(onSaved).toHaveBeenCalled();
        });
    });

    it('reverts edits via clearEdits', async () => {
        const adapter = {
            clearEdits: vi.fn().mockResolvedValue(basePhoto),
        } as unknown as IPhotosAdapter;
        const onSaved = vi.fn();

        render(<EditPanel photo={{...basePhoto, editParams: {rotate: 90}} as never} photosAdapter={adapter} onSaved={onSaved} />);

        fireEvent.click(screen.getByRole('button', {name: /revert edits/i}));

        await waitFor(() => {
            expect(adapter.clearEdits).toHaveBeenCalledWith('p1');
        });
        await waitFor(() => {
            expect(onSaved).toHaveBeenCalled();
        });
    });

    it('renders adjustment sliders', async () => {
        const adapter = {
            editPhoto: vi.fn().mockResolvedValue(basePhoto),
        } as unknown as IPhotosAdapter;
        const onSaved = vi.fn();

        render(<EditPanel photo={basePhoto as never} photosAdapter={adapter} onSaved={onSaved} />);

        expect(screen.getByText('Brightness')).toBeInTheDocument();
        expect(screen.getByText('Contrast')).toBeInTheDocument();
        expect(screen.getByText('Saturation')).toBeInTheDocument();
        expect(screen.getAllByRole('slider').length).toBeGreaterThanOrEqual(3);
    });

    it('applies auto-enhance toggle', async () => {
        const adapter = {
            editPhoto: vi.fn().mockResolvedValue(basePhoto),
        } as unknown as IPhotosAdapter;
        const onSaved = vi.fn();

        render(<EditPanel photo={basePhoto as never} photosAdapter={adapter} onSaved={onSaved} />);

        fireEvent.click(screen.getByRole('switch', {name: /auto-enhance/i}));

        await waitFor(() => {
            const call = adapter.editPhoto.mock.calls[0];
            expect(call[0]).toBe('p1');
            expect(call[1].autoEnhance).toBe(true);
        });
    });
});
