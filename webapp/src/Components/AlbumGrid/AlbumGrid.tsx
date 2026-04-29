import React, { useCallback, useState } from 'react';
import { Album } from '../../Models/Album';
import { AlbumCard } from '../AlbumCard/AlbumCard';
import { EmptyState } from '../EmptyState/EmptyState';
import { InputModal } from '../InputModal/InputModal';
import {Flex, TextInput} from '@mantine/core';
import {useNavigate} from "react-router-dom";
import useAlbumGrid from "./useAlbumGrid";
import {IAlbumsAdapter} from "../../Adapters/IAlbumsAdapter";
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import {IconAlbumOff} from '@tabler/icons-react';

interface IProps {
  albumsAdapter: IAlbumsAdapter;
  photosAdapter: IPhotosAdapter;
  maxDisplayed? : number;
}

const defaultProps = {
    maxDisplayed: 20000000,
}

export const AlbumGrid: React.FunctionComponent<IProps> = (propsIn) => {
    const props = {...defaultProps, ...propsIn};
    const navigate = useNavigate()
    const [showNewAlbumModal, setShowNewAlbumModal] = useState(false);
    const [{albums, createAlbumModalAlbumNameErrorText, newAlbumName, handleNewAlbumNameValueChange}] = useAlbumGrid(props.albumsAdapter);

  const newAlbumModalSaveClick = useCallback(() => {
    navigate(`/album/new/${newAlbumName}`);
  }, [navigate, newAlbumName]);

  return (
    <div>
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
          {albums && albums.length > 0 ? albums.slice(0, props.maxDisplayed).map((album: Album) => {
            return(
              <article key={album.id}>
                <AlbumCard
                  photosAdapter={props.photosAdapter}
                  source={album}
                  albumViewCallback={() => navigate(`../album/${album.id}`)}
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
        </section>
    </div>
  );
};
