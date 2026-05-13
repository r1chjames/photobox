import {IPhotosAdapter} from "../Adapters/IPhotosAdapter";

export const fetchPhotoBinWithAuth = async (photosAdapter: IPhotosAdapter, id: string) => {
    const data = await photosAdapter.getPhotoImage(id);

    // If the data is a string (like a URL from the mock adapter), return it directly.
    if (typeof data === 'string') {
        return data;
    }

    // Otherwise, assume it's binary data and create a blob URL.
    const blob = new Blob([data], {
        type: 'image/jpeg',
    });
    return URL.createObjectURL(blob);
}

export const revokeBlobUrl = (url: string | undefined) => {
    if (url && url.startsWith('blob:')) {
        URL.revokeObjectURL(url);
    }
}