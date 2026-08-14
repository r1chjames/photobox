import React from 'react';
import {render, screen, waitFor, fireEvent} from '../../test/test-utils';
import {vi} from 'vitest';
import {MetadataPanel} from './MetadataPanel';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';

const photo = {
    id: 'p1',
    name: 'photo.jpg',
    filesystemPath: '/photos/photo.jpg',
    sourcePath: '/photos/photo.jpg',
    albumId: '',
    tags: '',
    metadata: {
        exif: {Make: 'Canon', Model: 'EOS R5', GPSLatitude: ['52/1', '2/1', '24/1'], GPSLongitude: ['0/1', '5/1', '40/1'], GPSLatitudeRef: 'N', GPSLongitudeRef: 'E'},
        size: 1000,
        width: 100,
        height: 100,
        mime: 'image/jpeg',
    },
    createdAt: '2023-01-01T00:00:00Z',
    thumbnailUrl: '',
    latitude: 52.04,
    longitude: 0.094,
};

describe('MetadataPanel', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders grouped sections', () => {
        render(<MetadataPanel photo={photo as never} />);
        expect(screen.getByText('Camera')).toBeInTheDocument();
        expect(screen.getByText('Location')).toBeInTheDocument();
        expect(screen.getByText('File')).toBeInTheDocument();
        expect(screen.getByText('Canon')).toBeInTheDocument();
    });

    it('fetches reverse-geocoded location when GPS present', async () => {
        const adapter = {
            getPhotoLocation: vi.fn().mockResolvedValue('Cambridge, United Kingdom'),
        } as unknown as IPhotosAdapter;

        render(<MetadataPanel photo={photo as never} photosAdapter={adapter} />);

        await waitFor(() => {
            expect(screen.getByText('Cambridge, United Kingdom')).toBeInTheDocument();
        });
        expect(adapter.getPhotoLocation).toHaveBeenCalledWith('p1');
    });

    it('saves description edits via the adapter', async () => {
        const adapter = {
            getPhotoLocation: vi.fn().mockResolvedValue(''),
            updatePhotoMetadata: vi.fn().mockResolvedValue(photo),
        } as unknown as IPhotosAdapter;

        render(<MetadataPanel photo={{...photo, description: ''} as never} photosAdapter={adapter} />);

        // Enter edit mode
        const editButton = screen.getByRole('button', {name: /edit metadata/i});
        fireEvent.click(editButton);

        const input = screen.getByLabelText('Description') as HTMLInputElement;
        fireEvent.change(input, {target: {value: 'Beach day'}});

        const saveButton = screen.getByRole('button', {name: /save/i});
        fireEvent.click(saveButton);

        await waitFor(() => {
            expect(adapter.updatePhotoMetadata).toHaveBeenCalledWith('p1', {description: 'Beach day'});
        });
    });

    it('does not call adapter when no changes made', async () => {
        const adapter = {
            getPhotoLocation: vi.fn().mockResolvedValue(''),
            updatePhotoMetadata: vi.fn().mockResolvedValue(photo),
        } as unknown as IPhotosAdapter;

        render(<MetadataPanel photo={photo as never} photosAdapter={adapter} />);

        const editButton = screen.getByRole('button', {name: /edit metadata/i});
        fireEvent.click(editButton);
        const saveButton = screen.getByRole('button', {name: /save/i});
        fireEvent.click(saveButton);

        await waitFor(() => {
            expect(adapter.updatePhotoMetadata).not.toHaveBeenCalled();
        });
    });
});
