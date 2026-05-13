import React, {useState} from 'react';
import Dropzone from 'react-dropzone';
import './CreateAlbumView.css';
import {useParams} from 'react-router-dom';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {notifications} from '@mantine/notifications';
import {Progress, Text, Card, Stack} from '@mantine/core';

interface IProps {
  photosAdapter: IPhotosAdapter;
}

type QueryParams = {
  name: string;
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

export const CreateAlbumView: React.FunctionComponent<IProps> = (props) => {

  const [error, setError] = useState<string | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState({ current: 0, total: 0 });
  const {name} = useParams<QueryParams>();

  const handleFileUpload = async (acceptedFiles: File[]) => {
    setError(null);
    setIsUploading(true);
    setUploadProgress({ current: 0, total: acceptedFiles.length });
    let uploadedCount = 0;
    try {
      for (const file of acceptedFiles) {
        const fileContent = await readUploadedFileAsText(file);
        const photoContent = {
          name: file.name,
          albumName: name,
          binaryContent: fileContent,
        };
        await props.photosAdapter.uploadPhoto(photoContent);
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
      setError(message);
      notifications.show({
        title: 'Upload failed',
        message,
        color: 'red',
      });
    } finally {
      setIsUploading(false);
      setUploadProgress({ current: 0, total: 0 });
    }
  };

  return (
    <div>
      <div className="createAlbumView__dropzone">
        <Dropzone onDrop={acceptedFiles => handleFileUpload(acceptedFiles)} disabled={isUploading}>
          {({getRootProps, getInputProps}) => (
            <section>
              <div {...getRootProps()}>
                <input {...getInputProps()} />
                <p>{isUploading ? `Uploading ${uploadProgress.current} of ${uploadProgress.total}...` : 'Drag photos here to upload'}</p>
              </div>
            </section>
          )}
        </Dropzone>
      </div>
      {isUploading && uploadProgress.total > 0 && (
        <Card mt="md" p="sm" withBorder>
          <Stack gap="xs">
            <Text size="sm">Uploading {uploadProgress.current} of {uploadProgress.total} photos...</Text>
            <Progress value={(uploadProgress.current / uploadProgress.total) * 100} size="lg" />
          </Stack>
        </Card>
      )}
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
};
