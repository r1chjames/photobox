import React from 'react';
import { useParams } from 'react-router-dom';
import {MasonryInfiniteGrid} from '@egjs/react-infinitegrid';
import './PhotoIndexView.css';
import {Loader, MantineProvider} from '@mantine/core';
import {IPhotosAdapter} from "../../Adapters/IPhotosAdapter";
import usePhotoIndexView from "./usePhotoIndexView";

interface IProps {
  photosAdapter: IPhotosAdapter
}


export const PhotoIndexView: React.FunctionComponent<IProps> = (props) => {

  const { id } = useParams();
  const [{photoItems, photosLoaded}] = usePhotoIndexView(props.photosAdapter, id!);

  const apiCallsCompleted = () => (photosLoaded);

  const content = () => {
    if (apiCallsCompleted() && !photoItems) {
      return (
        <div>
          This album is empty
        </div>
      );
    }

    if (apiCallsCompleted() && photoItems) {
      return (
        <MasonryInfiniteGrid
          options={{ isConstantSize: false, transitionDuration: 0.2, useFit: true }}
          layoutOptions={{ margin: 5, column: [0, 5] }}
          // onAppend={onAppend}
          // onLayoutComplete={onLayoutComplete}
        >
          {photoItems}
        </MasonryInfiniteGrid>
      );
    }

    return (
      <Loader size={"md"} />
    );
  };

  return (
    <MantineProvider>
        <div className="photoIndexView__photoIndex">
          {content()}
        </div>
    </MantineProvider>
  );
};
