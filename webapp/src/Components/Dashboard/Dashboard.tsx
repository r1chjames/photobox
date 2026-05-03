import React, {useCallback, useState} from 'react';
import {useDropzone} from 'react-dropzone';
import {AlbumGrid} from '../AlbumGrid/AlbumGrid';
import {PhotoGrid} from '../PhotoGrid/PhotoGrid';
import {Memories} from '../Memories/Memories';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Card, Paper, Progress, Space, Stack, Text, Title} from "@mantine/core";
import {notifications} from '@mantine/notifications';

interface IProps {
    albumsAdapter: IAlbumsAdapter;
    photosAdapter: IPhotosAdapter;
}

const readUploadedFileAsText = (inputFile: File) => {
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

const DashboardUploadDropzone: React.FunctionComponent<{ photosAdapter: IPhotosAdapter }> = ({ photosAdapter }) => {
    const [isUploading, setIsUploading] = useState(false);
    const [uploadProgress, setUploadProgress] = useState({ current: 0, total: 0 });

    const handleFileUpload = useCallback(async (acceptedFiles: File[]) => {
        setIsUploading(true);
        setUploadProgress({ current: 0, total: acceptedFiles.length });
        let uploadedCount = 0;
        try {
            for (const file of acceptedFiles) {
                const fileContent = await readUploadedFileAsText(file);
                const photoContent = {
                    name: file.name,
                    albumName: 'General',
                    binaryContent: fileContent,
                };
                await photosAdapter.uploadPhoto(photoContent);
                uploadedCount++;
                setUploadProgress({ current: uploadedCount, total: acceptedFiles.length });
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
            setUploadProgress({ current: 0, total: 0 });
        }
    }, [photosAdapter]);

    const onDrop = useCallback((acceptedFiles: File[]) => {
        handleFileUpload(acceptedFiles);
    }, [handleFileUpload]);

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { 'image/*': [] },
        disabled: isUploading,
    });

    return (
        <Card mb="md" p="md" withBorder>
            <Paper
                {...getRootProps()}
                p="xl"
                withBorder
                style={{
                    border: isDragActive ? '2px dashed var(--mantine-color-blue-6)' : '2px dashed var(--mantine-color-gray-4)',
                    borderRadius: 'var(--mantine-radius-md)',
                    textAlign: 'center',
                    cursor: isUploading ? 'default' : 'pointer',
                    background: isDragActive ? 'var(--mantine-color-blue-light)' : 'transparent',
                    transition: 'all 0.2s ease',
                }}
            >
                <input {...getInputProps()} />
                <Text size="lg" c={isDragActive ? 'blue' : 'dimmed'}>
                    {isDragActive ? 'Drop photos here...' : 'Drag photos here or click to upload'}
                </Text>
            </Paper>
            {isUploading && uploadProgress.total > 0 && (
                <Stack gap="xs" mt="md">
                    <Text size="sm">Uploading {uploadProgress.current} of {uploadProgress.total} photos...</Text>
                    <Progress value={(uploadProgress.current / uploadProgress.total) * 100} size="lg" />
                </Stack>
            )}
        </Card>
    );
};

export const Dashboard: React.FunctionComponent<IProps> = (props) => {

    return (
        <>
            <DashboardUploadDropzone photosAdapter={props.photosAdapter} />
            <Memories photosAdapter={props.photosAdapter} />
            <div>
                <Title size="h4">Albums</Title>
                <Space h="md" />
                <AlbumGrid
                    albumsAdapter={props.albumsAdapter}
                    photosAdapter={props.photosAdapter}
                    maxDisplayed={20}
                />
            </div>
            <div>
                <Space h="md" />
                <Title size="h4">Photos</Title>
                <Space h="md" />
                <PhotoGrid
                    photosAdapter={props.photosAdapter}
                    albumsAdapter={props.albumsAdapter}
                    maxDisplayed={50}
                />
            </div>
        </>
    );
};
