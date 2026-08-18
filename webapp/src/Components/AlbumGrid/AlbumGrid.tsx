import React, { useCallback, useState } from 'react';
import { Album } from '../../Models/Album';
import { AlbumCard } from '../AlbumCard/AlbumCard';
import { EmptyState } from '../EmptyState/EmptyState';
import { InputModal } from '../InputModal/InputModal';
import {Button, Flex, TextInput} from '@mantine/core';
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
}

const defaultProps = {
    initialDisplayCount: undefined,
    loadMoreIncrement: 30,
}

export const AlbumGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const navigate = useNavigate()
    const [showNewAlbumModal, setShowNewAlbumModal] = useState(false);
    const [albumKey, setAlbumKey] = useState(0);
    const [{albums, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}] = useAlbumGrid(props.albumsAdapter);
    const [displayCount, setDisplayCount] = useState(props.initialDisplayCount ?? albums?.length ?? 0);

    // Sync display count once albums arrive. Previously initialized to 0 when
    // albums was still loading, leaving the /albums view empty ("Load More"
    // shown but no cards) until a manual increment — the Dashboard passed an
    // initialDisplayCount so it was masked there.
    React.useEffect(() => {
        if (props.initialDisplayCount === undefined && albums && albums.length > 0) {
            setDisplayCount(albums.length);
        }
    }, [albums, props.initialDisplayCount]);

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
        <section className="albumIndexView__cardContainer">
          <Flex
              direction={{ base: 'column', sm: 'row' }}
              gap={{ base: 'lg', sm: 'lg' }}
              justify={{ sm: 'flex-start' }}
              mih={50}
              bg="rgba(0, 0, 0, 0)"
              align="flex-start"
              wrap="wrap"
          >
          {albums && albums.length > 0 ? albums.slice(0, displayCount).map((album: Album) => {
            return(
              <article key={album.id}>
                <AlbumCard
                  photosAdapter={props.photosAdapter}
                  albumsAdapter={props.albumsAdapter}
                  source={album}
                  albumViewCallback={() => navigate(`../album/${album.id}`)}
                  onAlbumUpdated={handleAlbumUpdated}
                />
              </article>
            );
          }) : (
              <EmptyState
                  title="No albums yet"
                  description="Albums appear automatically from your photo folders, or create one manually."
                  icon={<IconAlbumOff size="2rem" />}
                  action={{
                      label: "Create album",
                      onClick: () => setShowNewAlbumModal(true),
                  }}
              />
          )}
          </Flex>
          {albums && albums.length > displayCount && (
              <Button
                  variant="light"
                  fullWidth
                  mt="md"
                  onClick={() => setDisplayCount(prev => Math.min(prev + props.loadMoreIncrement, albums.length))}
              >
                  Load More
              </Button>
          )}
        </section>
    </div>
  );
};
