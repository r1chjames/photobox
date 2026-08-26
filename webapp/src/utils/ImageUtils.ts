import {IPhotosAdapter} from "../Adapters/IPhotosAdapter";

export const fetchPhotoBinWithAuth = async (photosAdapter: IPhotosAdapter, id: string, mediaType?: string) => {
    const data = await photosAdapter.getPhotoImage(id);

    // If the data is a string (like a URL from the mock adapter), return it directly.
    if (typeof data === 'string') {
        return data;
    }

    // Otherwise, assume it's binary data and create a blob URL.
    const blob = new Blob([data], {
        type: mediaType || 'image/jpeg',
    });
    return URL.createObjectURL(blob);
}

export const fetchLiveVideoWithAuth = async (photosAdapter: IPhotosAdapter, id: string) => {
    const data = await photosAdapter.getPhotoLiveVideo(id);

    if (!data) return undefined;
    if (typeof data === 'string') return data;

    const blob = new Blob([data], { type: 'video/quicktime' });
    return URL.createObjectURL(blob);
}

export const revokeBlobUrl = (url: string | undefined) => {
    if (url && url.startsWith('blob:')) {
        URL.revokeObjectURL(url);
    }
}