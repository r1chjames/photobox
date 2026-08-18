import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AlbumsAdapter } from './AlbumsAdapter';
import { IRestApiAdapter } from './RestApiAdapter';

describe('AlbumsAdapter', () => {
    let mockRestApiAdapter: IRestApiAdapter;
    let albumsAdapter: AlbumsAdapter;

    beforeEach(() => {
        mockRestApiAdapter = {
            getApiCall: vi.fn(),
            getPaginatedApiCall: vi.fn(),
            postApiCall: vi.fn(),
            getBinaryApiCall: vi.fn(),
            authHeader: vi.fn().mockReturnValue({ Authorization: 'Bearer test-token' }),
        } as unknown as IRestApiAdapter;

        albumsAdapter = new AlbumsAdapter(mockRestApiAdapter);
    });

    describe('getAllAlbumsInfo', () => {
        it('should call getPaginatedApiCall with correct parameters', async () => {
            const mockAlbums = [
                { id: 'album-1', name: 'Album 1' },
                { id: 'album-2', name: 'Album 2' },
            ];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockAlbums);

            const result = await albumsAdapter.getAllAlbumsInfo();

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'albums',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                { limit: 1000 }
            );
            expect(result).toBe(mockAlbums);
        });

        it('should return empty array when no albums exist', async () => {
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue([]);

            const result = await albumsAdapter.getAllAlbumsInfo();

            expect(result).toEqual([]);
        });
    });

    describe('getCountOfPhotosInAlbum', () => {
        it('should call getApiCall with correct parameters', async () => {
            const mockCount = { count: 42 };
            vi.mocked(mockRestApiAdapter.getApiCall).mockResolvedValue(mockCount);

            const result = await albumsAdapter.getCountOfPhotosInAlbum();

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                'album',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockCount);
        });
    });

    describe('getAlbumInfoById', () => {
        it('should call getApiCall with correct album ID', async () => {
            const mockAlbum = { id: 'album-123', name: 'Test Album', photoCount: 10 };
            vi.mocked(mockRestApiAdapter.getApiCall).mockResolvedValue(mockAlbum);

            const result = await albumsAdapter.getAlbumInfoById('album-123');

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                'album/album-123',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockAlbum);
        });

        it('should handle special characters in album ID', async () => {
            const mockAlbum = { id: 'album-with-special-chars-@#$', name: 'Special Album' };
            vi.mocked(mockRestApiAdapter.getApiCall).mockResolvedValue(mockAlbum);

            const result = await albumsAdapter.getAlbumInfoById('album-with-special-chars-@#$');

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                'album/album-with-special-chars-@#$',
                expect.any(Object),
                {}
            );
            expect(result).toBe(mockAlbum);
        });
    });

    describe('buildHeaders', () => {
        it('should include standard Content-Type header', async () => {
            await albumsAdapter.getAllAlbumsInfo();

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                expect.any(String),
                expect.objectContaining({
                    'Content-Type': 'application/json',
                }),
                expect.any(Object)
            );
        });

        it('should include authorization header from authHeader()', async () => {
            await albumsAdapter.getAllAlbumsInfo();

            expect(mockRestApiAdapter.authHeader).toHaveBeenCalled();
            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                expect.any(String),
                expect.objectContaining({
                    Authorization: 'Bearer test-token',
                }),
                expect.any(Object)
            );
        });
    });

    describe('error handling', () => {
        it('should propagate errors from API calls', async () => {
            const mockError = new Error('API Error');
            vi.mocked(mockRestApiAdapter.getApiCall).mockRejectedValue(mockError);

            await expect(albumsAdapter.getAlbumInfoById('album-123')).rejects.toThrow('API Error');
        });

        it('should propagate network errors', async () => {
            const networkError = new Error('Network Error');
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockRejectedValue(networkError);

            await expect(albumsAdapter.getAllAlbumsInfo()).rejects.toThrow('Network Error');
        });
    });
});
describe('AlbumsAdapter smart albums', () => {
    let mockRestApiAdapter: IRestApiAdapter;
    let albumsAdapter: AlbumsAdapter;

    beforeEach(() => {
        mockRestApiAdapter = {
            getApiCall: vi.fn(),
            getPaginatedApiCall: vi.fn(),
            postApiCall: vi.fn(),
            patchApiCall: vi.fn(),
            getBinaryApiCall: vi.fn(),
            authHeader: vi.fn().mockReturnValue({ Authorization: 'Bearer test-token' }),
        } as unknown as IRestApiAdapter;

        albumsAdapter = new AlbumsAdapter(mockRestApiAdapter);
    });

    it('createSmartAlbum posts rules to albums/smart', async () => {
        await albumsAdapter.createSmartAlbum('Best 2024', { camera: 'Canon', favorite: true });

        expect(mockRestApiAdapter.postApiCall).toHaveBeenCalledWith(
            'albums/smart',
            { name: 'Best 2024', rules: { camera: 'Canon', favorite: true } },
            expect.objectContaining({ Authorization: 'Bearer test-token' })
        );
    });

    it('updateSmartAlbum patches rules to albums/:id/smart', async () => {
        await albumsAdapter.updateSmartAlbum('smart-1', { lowQuality: true });

        expect(mockRestApiAdapter.patchApiCall).toHaveBeenCalledWith(
            'albums/smart-1/smart',
            { rules: { lowQuality: true } },
            expect.objectContaining({ Authorization: 'Bearer test-token' })
        );
    });
});
