import React, {useEffect, useRef, useState} from 'react';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {usePhotoUpload} from "../../hooks/usePhotoUpload";

interface IProps {
    photosAdapter: IPhotosAdapter;
    /** Album name uploaded files land in; defaults to 'General'. */
    albumName?: string;
}

const isFileDrag = (e: DragEvent): boolean => {
    if (!e.dataTransfer) return false;
    // Real browsers expose 'Files' in dataTransfer.types for file drags.
    const types = e.dataTransfer.types;
    if (types && types.length > 0) {
        for (let i = 0; i < types.length; i++) {
            if (types[i] === 'Files') return true;
        }
        return false;
    }
    // Test environments (happy-dom) leave types empty even for file drags —
    // fall back to the presence of files/items.
    if (e.dataTransfer.files && e.dataTransfer.files.length > 0) return true;
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) return true;
    return false;
};

/**
 * Window-wide drag-and-drop upload surface (replaces the inline dropzones).
 * While any file drag enters the window a full-screen "Drop photos to upload"
 * overlay appears; releasing the files uploads them to the current album.
 */
export const DragUploadOverlay: React.FunctionComponent<IProps> = ({ photosAdapter, albumName }) => {
    const [isActive, setIsActive] = useState(false);
    const dragDepth = useRef(0);
    const { handleFileUpload } = usePhotoUpload(photosAdapter);

    useEffect(() => {
        const onDragEnter = (e: DragEvent) => {
            if (!isFileDrag(e)) return;
            e.preventDefault();
            dragDepth.current++;
            setIsActive(true);
        };

        const onDragOver = (e: DragEvent) => {
            if (!isFileDrag(e)) return;
            e.preventDefault();
            if (e.dataTransfer) {
                e.dataTransfer.dropEffect = 'copy';
            }
        };

        const onDragLeave = (e: DragEvent) => {
            if (!isFileDrag(e)) return;
            dragDepth.current = Math.max(0, dragDepth.current - 1);
            if (dragDepth.current === 0) {
                setIsActive(false);
            }
        };

        const onDrop = (e: DragEvent) => {
            if (!isFileDrag(e)) return;
            e.preventDefault();
            dragDepth.current = 0;
            setIsActive(false);
            const files = e.dataTransfer?.files;
            if (files && files.length > 0) {
                handleFileUpload(Array.from(files), albumName);
            }
        };

        window.addEventListener('dragenter', onDragEnter);
        window.addEventListener('dragover', onDragOver);
        window.addEventListener('dragleave', onDragLeave);
        window.addEventListener('drop', onDrop);
        return () => {
            window.removeEventListener('dragenter', onDragEnter);
            window.removeEventListener('dragover', onDragOver);
            window.removeEventListener('dragleave', onDragLeave);
            window.removeEventListener('drop', onDrop);
        };
    }, [handleFileUpload, albumName]);

    return (
        <div className={isActive ? 'drag-overlay is-active' : 'drag-overlay'}>
            <div className="drag-overlay-inner">
                <svg className="drag-overlay-icon" width="56" height="56" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                    <polyline points="17 8 12 3 7 8" />
                    <line x1="12" y1="3" x2="12" y2="15" />
                </svg>
                <div className="drag-overlay-title">Drop photos to upload</div>
                <div className="drag-overlay-hint">Release to add to your library</div>
            </div>
        </div>
    );
};
