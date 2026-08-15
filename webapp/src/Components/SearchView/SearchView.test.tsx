import React from 'react';
import {render, screen, fireEvent, waitFor} from '../../test/test-utils';
import {vi} from 'vitest';
import {SearchView} from './SearchView';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {IAlbumsAdapter} from '../../Adapters/IAlbumsAdapter';

// Mock react-router-dom so useSearchParams returns a controllable value
// without needing a Router with custom initial entries. The set() mutates
// the params so clear-filters behaviour matches the real router.
const {searchParamsMock} = vi.hoisted(() => ({
    searchParamsMock: {params: new URLSearchParams(), set: vi.fn()},
}));

vi.mock('react-router-dom', async (importOriginal) => ({
    ...(await importOriginal<typeof import('react-router-dom')>()),
    useSearchParams: () => {
        const set = (next: Record<string, string> | URLSearchParams) => {
            if (next instanceof URLSearchParams) {
                searchParamsMock.params = next;
            } else {
                searchParamsMock.params = new URLSearchParams(next as Record<string, string>);
            }
        };
        return [searchParamsMock.params, set];
    },
}));

const renderView = () => render(
    <SearchView photosAdapter={photosAdapter} albumsAdapter={albumsAdapter} />
);

const photosAdapter = {
    getAllTags: vi.fn().mockResolvedValue(['beach', 'family']),
    getAllPhotosInfo: vi.fn().mockResolvedValue([]),
} as unknown as IPhotosAdapter;

const albumsAdapter = {} as unknown as IAlbumsAdapter;

describe('SearchView', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        searchParamsMock.params = new URLSearchParams();
    });

    it('shows empty state when no filters active', () => {
        renderView();
        expect(screen.getByText(/search your photos/i)).toBeInTheDocument();
    });

    it('renders advanced filter controls when query is active', async () => {
        searchParamsMock.params = new URLSearchParams('q=beach');
        renderView();
        await waitFor(() => {
            expect(screen.getByPlaceholderText(/camera/i)).toBeInTheDocument();
        });
        expect(screen.getByText(/has gps/i)).toBeInTheDocument();
        expect(screen.getByText(/tags:/i)).toBeInTheDocument();
    });

    it('clear filters resets state', async () => {
        searchParamsMock.params = new URLSearchParams('q=beach');
        renderView();
        await waitFor(() => {
            expect(screen.getByPlaceholderText(/camera/i)).toBeInTheDocument();
        });

        const cameraInput = screen.getByPlaceholderText(/camera/i) as HTMLInputElement;
        fireEvent.change(cameraInput, {target: {value: 'Canon'}});
        expect(cameraInput.value).toBe('Canon');

        fireEvent.click(screen.getByText(/clear all filters/i));

        // Clear resets the URL params (mocked) and local state, so the
        // empty-state hint returns.
        expect(searchParamsMock.params.toString()).toBe('');
        await waitFor(() => {
            expect(screen.getByText(/search your photos/i)).toBeInTheDocument();
        });
    });
});
