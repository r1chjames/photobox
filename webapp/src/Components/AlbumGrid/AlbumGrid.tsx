import React, {useCallback, useState} from 'react';
import { Album } from '../../Models/Album';
import { AlbumCard } from '../AlbumCard/AlbumCard';
import { AlbumCardSkeleton } from './AlbumCardSkeleton';
import { EmptyState } from '../EmptyState/EmptyState';
import { InputModal } from '../InputModal/InputModal';
import {Button, TextInput} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import useAlbumGrid from "./useAlbumGrid";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IconAlbumOff} from '@tabler/icons-react';
import { notifications } from '@mantine/notifications';

interface IProps {
  albumsAdapter: IAlbumsAdapter;
  photosAdapter: IPhotosAdapter;
  initialDisplayCount?: number;
  loadMoreIncrement?: number;
  showTitleBar?: boolean;
}

const defaultProps = {
    initialDisplayCount: undefined,
    loadMoreIncrement: 30,
}

export const AlbumGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const navigate = useNavigate();
    const [showNewAlbumModal, setShowNewAlbumModal] = useState(false);
    const [albumKey, setAlbumKey] = useState(0);
    const [{albums, isLoading, isError, refetch, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}] = useAlbumGrid(props.albumsAdapter);
    // Paging applies only when a caller supplies an initial page size (e.g.
    // the Dashboard's 30); the standalone /albums view always shows every
    // album. displayCount starts at the page size and grows via "Load More".
    const [displayCount, setDisplayCount] = useState(props.initialDisplayCount ?? 0);

    const isPaged = props.initialDisplayCount !== undefined;
    // The slice shown at any moment; recomputed each render so "Load More"
    // visibility and the rendered cards can never disagree, and an empty
    // grid never flashes between albums arriving and state catching up.
    const visibleAlbums = albums && !isLoading
        ? albums.slice(0, isPaged ? displayCount : albums.length)
        : [];

  const newAlbumModalSaveClick = useCallback(async () => {
    try {
      await props.albumsAdapter.createAlbum(newAlbumName);
      notifications.show({
        title: 'Album created',
        message: `Created album "${newAlbumName}"`,
        color: 'green',
      });
      setShowNewAlbumModal(false);
      setAlbumKey(prev => prev + 1); // force refresh
    } catch (e) {
      notifications.show({
        title: 'Failed to create album',
        message: e instanceof Error ? e.message : 'An error occurred',
        color: 'red',
      });
    }
  }, [props.albumsAdapter, newAlbumName]);

  const handleAlbumUpdated = useCallback(() => {
    setAlbumKey(prev => prev + 1);
  }, []);

  const renderGridContent = () => {
    if (isLoading && !albums) {
      // First load (no cached albums yet): skeleton cards in the same 3-col
      // grid so real cards replace them without layout shift (issue #173).
      const skeletonCount = Math.min(props.initialDisplayCount ?? 6, 6);
      return Array.from({length: skeletonCount}).map((_, i) => (
        <article key={`skeleton-${i}`}>
          <AlbumCardSkeleton />
        </article>
      ));
    }

    if (isError && !albums) {
      return (
        <EmptyState
            title="Failed to load albums"
            description="Albums could not be loaded. Check your connection and try again."
            icon={<IconAlbumOff size="2rem" />}
            action={{
                label: "Try again",
                onClick: () => { refetch(); },
            }}
        />
      );
    }

    if (albums && albums.length > 0) {
      return visibleAlbums.map((album: Album) => (
        <article key={album.id}>
          <AlbumCard
            photosAdapter={props.photosAdapter}
            albumsAdapter={props.albumsAdapter}
            source={album}
            albumViewCallback={() => navigate(`../album/${album.id}`)}
            onAlbumUpdated={handleAlbumUpdated}
          />
        </article>
      ));
    }

    // Loaded successfully but no albums yet.
    return (
        <EmptyState
            title="No albums yet"
            description="Albums appear automatically from your photo folders, or create one manually."
            icon={<IconAlbumOff size="2rem" />}
            action={{
                label: "Create album",
                onClick: () => setShowNewAlbumModal(true),
            }}
        />
    );
  };

  return (
    <div key={albumKey}>
        <InputModal
          isOpen={showNewAlbumModal}
          title="Create Album"
          handleSave={newAlbumModalSaveClick}
          handleClose={() => setShowNewAlbumModal(false)}
        >
          <TextInput
            required={true}
            id="album-name"
            label="Album Name"
            error={createAlbumModalAlbumNameErrorText.length !== 0}
            onChange={e => handleNewAlbumNameValueChange(e.target.value)}
          />
        </InputModal>
        {props.showTitleBar && (
            <div className="view-title-bar">
                <h1 className="view-title">Albums</h1>
                <div className="view-actions">
                    <button type="button" className="action-btn primary" onClick={() => setShowNewAlbumModal(true)}>
                        + New Album
                    </button>
                </div>
            </div>
        )}
        <section className="albumIndexView__cardContainer">
          <div className="albums-grid">
            {renderGridContent()}
          </div>
          {!isLoading && isPaged && albums && albums.length > displayCount && (
              <Button
                  variant="light"
                  fullWidth
                  mt="md"
                  onClick={() => setDisplayCount((prev: number) => Math.min(prev + props.loadMoreIncrement, albums.length))}
              >
                  Load More
              </Button>
          )}
        </section>
    </div>
  );
};
