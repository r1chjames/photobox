import {useCallback, useState} from 'react';
import {notifications} from '@mantine/notifications';
import {IPhotosAdapter} from "../Adapters/IPhotosAdapter";

const readUploadedFileAsText = (inputFile: File): Promise<string> => {
    const temporaryFileReader = new FileReader();

    return new Promise<string>((resolve, reject) => {
        temporaryFileReader.onerror = () => {
            temporaryFileReader.abort();
            reject(new DOMException('Problem parsing input file.'));
        };

        temporaryFileReader.onload = () => {
            resolve(temporaryFileReader.result as string);
        };
        temporaryFileReader.readAsDataURL(inputFile);
    });
};

export interface PhotoUploadState {
    isUploading: boolean;
    progress: { current: number; total: number };
    handleFileUpload: (files: File[], albumName?: string) => Promise<void>;
}

/**
 * Shared photo-upload pipeline used by the drag-onto-page overlay and the
 * "Add Photos" button. Each file is base64-encoded and POSTed to the upload
 * endpoint sequentially, with progress surfaced through Mantine notifications.
 */
export const usePhotoUpload = (photosAdapter: IPhotosAdapter): PhotoUploadState => {
    const [isUploading, setIsUploading] = useState(false);
    const [progress, setProgress] = useState({ current: 0, total: 0 });

    const handleFileUpload = useCallback(async (acceptedFiles: File[], albumName = 'General') => {
        if (acceptedFiles.length === 0) return;
        setIsUploading(true);
        setProgress({ current: 0, total: acceptedFiles.length });
        let uploadedCount = 0;
        try {
            for (const file of acceptedFiles) {
                const fileContent = await readUploadedFileAsText(file);
                const photoContent = {
                    name: file.name,
                    albumName,
                    binaryContent: fileContent,
                };
                await photosAdapter.uploadPhoto(photoContent);
                uploadedCount++;
                setProgress({ current: uploadedCount, total: acceptedFiles.length });
            }
            notifications.show({
                title: 'Upload complete',
                message: `${uploadedCount} photo${uploadedCount !== 1 ? 's' : ''} uploaded successfully`,
                color: 'green',
            });
        } catch (e) {
            const message = e instanceof Error ? e.message : 'Upload failed';
            notifications.show({
                title: 'Upload failed',
                message,
                color: 'red',
            });
        } finally {
            setIsUploading(false);
            setProgress({ current: 0, total: 0 });
        }
    }, [photosAdapter]);

    return { isUploading, progress, handleFileUpload };
};
