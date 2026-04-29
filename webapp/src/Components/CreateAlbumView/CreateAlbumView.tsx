import React, {useState} from 'react';
import Dropzone from 'react-dropzone';
import './CreateAlbumView.css';
import {useParams} from 'react-router-dom';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {notifications} from '@mantine/notifications';

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
  const {name} = useParams<QueryParams>();

  const handleFileUpload = async (acceptedFiles: File[]) => {
    setError(null);
    setIsUploading(true);
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
                <p>{isUploading ? 'Uploading...' : 'Drag photos here to upload'}</p>
              </div>
            </section>
          )}
        </Dropzone>
      </div>
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
};
