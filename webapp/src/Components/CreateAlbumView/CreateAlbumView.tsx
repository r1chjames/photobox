import React, {useRef, useState} from 'react';
import './CreateAlbumView.css';
import {useParams} from 'react-router-dom';
import {IPhotosAdapter} from '../../Adapters/IPhotosAdapter';
import {Loader, Text} from '@mantine/core';
import {IconPhotoPlus} from '@tabler/icons-react';
import {usePhotoUpload} from '../../hooks/usePhotoUpload';

interface IProps {
  photosAdapter: IPhotosAdapter;
}

type QueryParams = {
  name: string;
}

export const CreateAlbumView: React.FunctionComponent<IProps> = (props) => {
  const [error, setError] = useState<string | null>(null);
  const {name} = useParams<QueryParams>();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const {handleFileUpload, isUploading} = usePhotoUpload(props.photosAdapter);

  const handleFilesSelected = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files ?? []);
    e.target.value = '';
    if (files.length === 0) return;
    setError(null);
    try {
      await handleFileUpload(files, name);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Upload failed';
      setError(message);
    }
  };

  return (
    <div>
      <input
        ref={fileInputRef}
        type="file"
        accept="image/*"
        multiple
        style={{ display: 'none' }}
        onChange={handleFilesSelected}
        aria-hidden="true"
        tabIndex={-1}
      />
      <div className="createAlbumView__dropzone">
        <button
          type="button"
          className="action-btn primary"
          onClick={() => fileInputRef.current?.click()}
          disabled={isUploading}
        >
          <IconPhotoPlus size="1rem" /> Add photos to {name ?? 'album'}
        </button>
        <p style={{ color: 'var(--muted)', marginTop: 12 }}>
          {isUploading ? 'Uploading…' : 'or drag photos anywhere on this window to upload'}
        </p>
        {isUploading && (
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, justifyContent: 'center', marginTop: 12 }}>
            <Loader size="sm" />
            <Text size="sm">Uploading photos…</Text>
          </div>
        )}
      </div>
      {error && <p style={{ color: 'red', textAlign: 'center' }}>{error}</p>}
    </div>
  );
};
