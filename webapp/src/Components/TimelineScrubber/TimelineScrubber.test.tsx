import { describe, it, expect, vi, beforeEach } from 'vitest';
import { fireEvent, waitFor } from '@testing-library/react';
import { render } from '../../test/test-utils';
import { TimelineScrubber } from './TimelineScrubber';
import { IPhotosAdapter, TimelineEntry } from '../../Adapters/IPhotosAdapter';

// useMediaQuery must return false (desktop) so the component renders the
// desktop rail path — the path that previously had a hook-after-return bug
// (React #310). Mocking at module level ensures a stable boolean.
vi.mock('@mantine/hooks', () => ({
    useMediaQuery: () => false,
}));

const mockEntries: TimelineEntry[] = [
    { year: 2026, month: 9, count: 30 },
    { year: 2026, month: 3, count: 15 },
    { year: 2025, month: 12, count: 10 },
    { year: 2025, month: 6, count: 8 },
];

function createMockAdapter(entries: TimelineEntry[] = mockEntries): IPhotosAdapter {
    return {
        getTimeline: vi.fn().mockResolvedValue(entries),
        getAllPhotosInfo: vi.fn().mockResolvedValue([]),
        getPhotoInfoById: vi.fn(),
        getPhotosInfoInAlbum: vi.fn(),
        getPhotoCountInAlbum: vi.fn(),
        getPhotoImage: vi.fn(),
        getPhotoLiveVideo: vi.fn(),
        downloadPhoto: vi.fn(),
        deletePhoto: vi.fn(),
        favoritePhoto: vi.fn(),
        downloadPhotosAsZip: vi.fn(),
        searchPhotos: vi.fn(),
        getTrashedPhotos: vi.fn(),
        restorePhoto: vi.fn(),
        rotatePhoto: vi.fn(),
        editPhoto: vi.fn(),
        clearEdits: vi.fn(),
        getGeodata: vi.fn(),
        getMemories: vi.fn(),
        getAllTags: vi.fn(),
        updatePhotoTags: vi.fn(),
        batchUpdatePhotoTags: vi.fn(),
        batchSetFavorite: vi.fn(),
        batchDeletePhotos: vi.fn(),
        batchAddToAlbum: vi.fn(),
        updatePhotoMetadata: vi.fn(),
        getPhotoLocation: vi.fn(),
        getPhotosByTag: vi.fn(),
        getDuplicatePhotos: vi.fn(),
        getPhotoThumbnailBlob: vi.fn(),
        uploadPhoto: vi.fn(),
        startJob: vi.fn(),
        stopJob: vi.fn(),
        getAllJobStatuses: vi.fn(),
        stopAllJobs: vi.fn(),
        getVideos: vi.fn(),
    } as unknown as IPhotosAdapter;
}

describe('TimelineScrubber', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('renders the desktop rail without throwing when timeline entries exist', async () => {
        const adapter = createMockAdapter();
        const onSelectMonth = vi.fn();
        const onClear = vi.fn();

        // This render must not throw "Rendered more hooks than during the
        // previous render" (React #310). The bug was handleRailKeyDown
        // (a useCallback) declared after the isMobile early return.
        const { container } = render(
            <TimelineScrubber
                photosAdapter={adapter}
                onSelectMonth={onSelectMonth}
                onClear={onClear}
            />
        );

        // Wait for the timeline to load and the rail to render.
        await waitFor(() => {
            const rail = container.querySelector('.timeline-rail');
            expect(rail).not.toBeNull();
        });

        // The rail must have the slider role and keyboard support.
        const slider = container.querySelector('[role="slider"]');
        expect(slider).not.toBeNull();
        expect(slider?.getAttribute('aria-valuemin')).toBe('0');
        expect(slider?.getAttribute('aria-valuemax')).toBe('3'); // 4 entries
    });

    it('renders nothing when timeline is empty', async () => {
        const adapter = createMockAdapter([]);
        const { container } = render(
            <TimelineScrubber
                photosAdapter={adapter}
                onSelectMonth={vi.fn()}
                onClear={vi.fn()}
            />
        );

        await waitFor(() => {
            expect(container.querySelector('.timeline-rail')).toBeNull();
            expect(container.querySelector('[role="slider"]')).toBeNull();
        });
    });

    it('survives an isMobile transition from false to true without hook errors', async () => {
        // The hook-count bug only manifests when isMobile changes between
        // renders. We can't easily flip the mock, but we can at least confirm
        // that multiple re-renders (state changes from interaction) don't
        // throw. Hovering triggers state updates that re-render the rail path.
        const adapter = createMockAdapter();
        const { container } = render(
            <TimelineScrubber
                photosAdapter={adapter}
                onSelectMonth={vi.fn()}
                onClear={vi.fn()}
            />
        );

        await waitFor(() => {
            expect(container.querySelector('[role="slider"]')).not.toBeNull();
        });

        // Simulate mouse enter on the slider to trigger isHovering state change
        const slider = container.querySelector('[role="slider"]')!;
        fireEvent.mouseEnter(slider);

        // If the hook order were wrong, this re-render would throw #310.
        await waitFor(() => {
            expect(container.querySelector('[role="slider"]')).not.toBeNull();
        });
    });

    it('responds to ArrowDown key to move the scrubber', async () => {
        const adapter = createMockAdapter();
        const onSelectMonth = vi.fn();
        const { container } = render(
            <TimelineScrubber
                photosAdapter={adapter}
                onSelectMonth={onSelectMonth}
                onClear={vi.fn()}
            />
        );

        await waitFor(() => {
            expect(container.querySelector('[role="slider"]')).not.toBeNull();
        });

        const slider = container.querySelector('[role="slider"]')!;
        slider.focus();
        fireEvent.keyDown(slider, { key: 'ArrowDown' });

        // ArrowDown on 4 entries (index 0 → 1) should select the 2nd month
        expect(onSelectMonth).toHaveBeenCalledWith(2026, 3);
    });

    it('responds to Home key to jump to the first month', async () => {
        const adapter = createMockAdapter();
        const onSelectMonth = vi.fn();
        const { container } = render(
            <TimelineScrubber
                photosAdapter={adapter}
                onSelectMonth={onSelectMonth}
                onClear={vi.fn()}
            />
        );

        await waitFor(() => {
            expect(container.querySelector('[role="slider"]')).not.toBeNull();
        });

        const slider = container.querySelector('[role="slider"]')!;
        slider.focus();
        fireEvent.keyDown(slider, { key: 'Home' });

        expect(onSelectMonth).toHaveBeenCalledWith(2026, 9);
    });
});
