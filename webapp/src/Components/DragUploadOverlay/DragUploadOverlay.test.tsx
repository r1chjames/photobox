import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, waitFor } from '../../test/test-utils';
import { DragUploadOverlay } from './DragUploadOverlay';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

// jsdom's DataTransfer exposes an empty `types` list, so the overlay's
// isFileDrag() check (which looks for the 'Files' type) must be stubbed.
const stubTypes = (dt: DataTransfer, types: string[] = ['Files']) => {
    Object.defineProperty(dt, 'types', { value: types, configurable: true });
    return dt;
};

const makeFileDragEvent = (type: string, files?: File[]): DragEvent => {
    const dt = new DataTransfer();
    if (files) {
        files.forEach(f => dt.items.add(f));
    } else {
        dt.items.add(new File(['x'], 'x.txt', { type: 'text/plain' }));
    }
    stubTypes(dt);
    // happy-dom's DragEvent constructor drops the dataTransfer option, so
    // attach it after construction (real browsers accept it via init).
    const ev = new DragEvent(type, { bubbles: true, cancelable: true });
    Object.defineProperty(ev, 'dataTransfer', { value: dt, configurable: true });
    return ev;
};

const overlayActive = (container: HTMLElement): boolean =>
    !!container.querySelector('.drag-overlay')?.classList.contains('is-active');

describe('DragUploadOverlay', () => {
    let mockPhotosAdapter: IPhotosAdapter;

    beforeEach(() => {
        mockPhotosAdapter = {
            uploadPhoto: vi.fn().mockResolvedValue({}),
        } as unknown as IPhotosAdapter;
    });

    it('should activate the overlay when a file drag enters the window', async () => {
        const { container } = render(<DragUploadOverlay photosAdapter={mockPhotosAdapter} />);
        expect(overlayActive(container)).toBe(false);

        window.dispatchEvent(makeFileDragEvent('dragenter'));
        await waitFor(() => expect(overlayActive(container)).toBe(true));
    });

    it('should deactivate the overlay when the file drag leaves the window', async () => {
        const { container } = render(<DragUploadOverlay photosAdapter={mockPhotosAdapter} />);

        window.dispatchEvent(makeFileDragEvent('dragenter'));
        await waitFor(() => expect(overlayActive(container)).toBe(true));

        window.dispatchEvent(makeFileDragEvent('dragleave'));
        await waitFor(() => expect(overlayActive(container)).toBe(false));
    });

    it('should ignore non-file drags', () => {
        const { container } = render(<DragUploadOverlay photosAdapter={mockPhotosAdapter} />);

        const dt = new DataTransfer();
        stubTypes(dt, ['text/plain']);
        const ev = new DragEvent('dragenter', { bubbles: true, cancelable: true });
        Object.defineProperty(ev, 'dataTransfer', { value: dt, configurable: true });
        window.dispatchEvent(ev);
        expect(overlayActive(container)).toBe(false);
    });

    it('should upload dropped files to the album and deactivate the overlay', async () => {
        const { container } = render(<DragUploadOverlay photosAdapter={mockPhotosAdapter} albumName="General" />);

        window.dispatchEvent(makeFileDragEvent('dragenter'));
        const file = new File(['data:image/png;base64,AAA='], 'photo.png', { type: 'image/png' });
        window.dispatchEvent(makeFileDragEvent('drop', [file]));

        await waitFor(() => {
            expect(mockPhotosAdapter.uploadPhoto).toHaveBeenCalledWith(
                expect.objectContaining({ name: 'photo.png', albumName: 'General' })
            );
        });
        expect(overlayActive(container)).toBe(false);
    });
});
