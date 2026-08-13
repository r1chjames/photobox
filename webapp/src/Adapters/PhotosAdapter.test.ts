import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PhotosAdapter } from './PhotosAdapter';
import { IRestApiAdapter } from './RestApiAdapter';
import { Photo } from '../Models/Photo';
import { JobType } from '../Models/Job';

describe('PhotosAdapter', () => {
    let mockRestApiAdapter: IRestApiAdapter;
    let photosAdapter: PhotosAdapter;

    beforeEach(() => {
        mockRestApiAdapter = {
            getApiCall: vi.fn(),
            getPaginatedApiCall: vi.fn(),
            postApiCall: vi.fn(),
            getBinaryApiCall: vi.fn(),
            authHeader: vi.fn().mockReturnValue({ Authorization: 'Bearer test-token' }),
        } as unknown as IRestApiAdapter;

        photosAdapter = new PhotosAdapter(mockRestApiAdapter);
    });

    describe('getAllPhotosInfo', () => {
        it('should call getPaginatedApiCall with correct parameters', async () => {
            const mockPhotos: Photo[] = [];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockPhotos);

            const result = await photosAdapter.getAllPhotosInfo('photo-123', 30, true);

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'photos?fromId=photo-123&limit=30&thumbnail=true',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockPhotos);
        });

        it('should handle includeThumbnails=false', async () => {
            const mockPhotos: Photo[] = [];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockPhotos);

            await photosAdapter.getAllPhotosInfo('photo-123', 50, false);

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'photos?fromId=photo-123&limit=50&thumbnail=false',
                expect.any(Object),
                {}
            );
        });

        it('should append favorites=true when favorite flag is set', async () => {
            const mockPhotos: Photo[] = [];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockPhotos);

            const result = await photosAdapter.getAllPhotosInfo('photo-123', 30, true, undefined, undefined, true);

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'photos?fromId=photo-123&limit=30&thumbnail=true&favorites=true',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockPhotos);
        });

        it('should not append favorites param when favorite flag is false', async () => {
            const mockPhotos: Photo[] = [];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockPhotos);

            await photosAdapter.getAllPhotosInfo('photo-123', 30, true, undefined, undefined, false);

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'photos?fromId=photo-123&limit=30&thumbnail=true',
                expect.any(Object),
                {}
            );
        });
    });

    describe('getPhotoInfoById', () => {
        it('should call getApiCall with correct parameters', async () => {
            const mockPhoto = {} as Photo;
            vi.mocked(mockRestApiAdapter.getApiCall).mockResolvedValue(mockPhoto);

            const result = await photosAdapter.getPhotoInfoById('photo-456');

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                'photo/info/photo-456',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockPhoto);
        });
    });

    describe('getPhotosInfoInAlbum', () => {
        it('should call getPaginatedApiCall with correct parameters including albumId', async () => {
            const mockPhotos: Photo[] = [];
            vi.mocked(mockRestApiAdapter.getPaginatedApiCall).mockResolvedValue(mockPhotos);

            const result = await photosAdapter.getPhotosInfoInAlbum('album-789', 'photo-123', 20, true);

            expect(mockRestApiAdapter.getPaginatedApiCall).toHaveBeenCalledWith(
                'photos?fromId=photo-123&limit=20&thumbnail=true',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                { albumId: 'album-789' }
            );
            expect(result).toBe(mockPhotos);
        });
    });

    describe('getPhotoCountInAlbum', () => {
        it('should call getApiCall with correct parameters', async () => {
            const mockCount = 42;
            vi.mocked(mockRestApiAdapter.getApiCall).mockResolvedValue(mockCount);

            const result = await photosAdapter.getPhotoCountInAlbum('album-789');

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                'photos/count?albumId=album-789',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                { albumId: 'album-789' }
            );
            expect(result).toBe(mockCount);
        });
    });

    describe('getPhotoImage', () => {
        it('should call getBinaryApiCall with correct parameters', async () => {
            const mockBinary = new Blob();
            vi.mocked(mockRestApiAdapter.getBinaryApiCall).mockResolvedValue(mockBinary);

            const result = await photosAdapter.getPhotoImage('photo-456');

            expect(mockRestApiAdapter.getBinaryApiCall).toHaveBeenCalledWith(
                'photo/bin/photo-456',
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                },
                {}
            );
            expect(result).toBe(mockBinary);
        });
    });

    describe('uploadPhoto', () => {
        it('should call postApiCall with correct parameters', async () => {
            const mockBody = { name: 'test.jpg', data: 'base64data' };
            const mockResponse = { success: true };
            vi.mocked(mockRestApiAdapter.postApiCall).mockResolvedValue(mockResponse);

            const result = await photosAdapter.uploadPhoto(mockBody);

            expect(mockRestApiAdapter.postApiCall).toHaveBeenCalledWith(
                'photo',
                mockBody,
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                }
            );
            expect(result).toBe(mockResponse);
        });
    });

    describe('startJob', () => {
        it('should call postApiCall with job type', async () => {
            const mockResponse = { started: true };
            vi.mocked(mockRestApiAdapter.postApiCall).mockResolvedValue(mockResponse);

            const result = await photosAdapter.startJob(JobType.PhotoIndex);

            expect(mockRestApiAdapter.postApiCall).toHaveBeenCalledWith(
                'photos/jobs/Photo_index/start',
                {},
                {
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                }
            );
            expect(result).toBe(mockResponse);
        });
    });

    describe('buildHeaders', () => {
        it('should merge standard headers with auth headers', async () => {
            // Test indirectly through a public method
            await photosAdapter.getPhotoInfoById('test');

            expect(mockRestApiAdapter.getApiCall).toHaveBeenCalledWith(
                expect.any(String),
                expect.objectContaining({
                    'Content-Type': 'application/json',
                    Authorization: 'Bearer test-token',
                }),
                expect.any(Object)
            );
        });
    });
});