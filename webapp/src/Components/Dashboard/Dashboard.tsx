import React, {useRef} from 'react';
import {AlbumGrid} from '../AlbumGrid/AlbumGrid';
import {PhotoGrid} from '../PhotoGrid/PhotoGrid';
import {Memories} from '../Memories/Memories';
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {Loader, Space, Text, Title} from "@mantine/core";
import {IconPhotoPlus} from "@tabler/icons-react";
import {usePhotoUpload} from "../../hooks/usePhotoUpload";

interface IProps {
    albumsAdapter: IAlbumsAdapter;
    photosAdapter: IPhotosAdapter;
}

export const Dashboard: React.FunctionComponent<IProps> = (props) => {
    const fileInputRef = useRef<HTMLInputElement>(null);
    const { handleFileUpload, isUploading } = usePhotoUpload(props.photosAdapter);

    const handleFilesSelected = (e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files ?? []);
        e.target.value = '';
        if (files.length > 0) {
            handleFileUpload(files, 'General');
        }
    };

    return (
        <>
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
            {isUploading && (
                <div style={{ display: 'flex', alignItems: 'center', gap: 12, padding: '8px 12px', marginBottom: 12, background: 'var(--surface-card)', borderRadius: 'var(--radius-md)' }}>
                    <Loader size="sm" />
                    <Text size="sm">Uploading photos…</Text>
                </div>
            )}
            <Memories photosAdapter={props.photosAdapter} />
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 12 }}>
                <Title size="h4">Albums</Title>
                <button type="button" className="action-btn primary" onClick={() => fileInputRef.current?.click()} aria-label="Add photos">
                    <IconPhotoPlus size="1rem" /> Add Photos
                </button>
            </div>
            <AlbumGrid
                albumsAdapter={props.albumsAdapter}
                photosAdapter={props.photosAdapter}
                initialDisplayCount={30}
                loadMoreIncrement={30}
            />
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
