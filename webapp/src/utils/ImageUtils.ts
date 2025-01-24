import {IPhotosAdapter} from "../Adapters/IPhotosAdapter";

export const fetchPhotoBinWithAuth = async (photosAdapter: IPhotosAdapter, id: string) => {
    const data = await photosAdapter.getPhotoImage(id)
    const blob = new Blob([data], {
        type: 'image/jpeg',
    });
    return URL.createObjectURL(blob)
}