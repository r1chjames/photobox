import React, {useEffect, useState} from 'react';
import './PhotoDetail.css';
import {Photo} from '../../Models/Photo';
import { useParams } from 'react-router-dom';
import { PhotosAdapter } from '../../Adapters/PhotosAdapter';
import {Loader, MantineProvider, Table} from '@mantine/core';

interface IProps {
  baseApiUrl: string;
}

const getPhoto = async(baseApiUrl: string, photoId: string) => {
  const photosAdapter = new PhotosAdapter(baseApiUrl);
  return photosAdapter.getPhotoInfoById(photoId);
};

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {
  const [photo, setPhoto] = useState<Photo>();
  const { id } = useParams();

  useEffect(() => {
    if (id !== undefined) {
      (async function retrievePhoto() {
        const retrievedPhoto = await getPhoto(props.baseApiUrl, id);
        setPhoto(retrievedPhoto);
      })();
    }
  },        [setPhoto, id, props.baseApiUrl]);

  const content = () => {
    if (photo !== undefined) {
      return (
          <div className="photoDetail__contentWrapper">
              <img
                src={`${props.baseApiUrl}/photo/${id}/bin`}
                alt={photo.name}
                className="photoDetail__mainImage"
              />
            <div className="photoDetail__imageMetadata">
              <Table aria-label="metadata table">
                <tr>
                  <th>Parameter</th>
                  <th>Value</th>
                </tr>
                {Object.entries(photo.metadata).map(([key, value]) => {
                  if (value !== null) {
                    return (
                      <tr key={key}>
                        <td>{key}</td>
                        <td>{JSON.stringify(value)}</td>
                      </tr>
                    );
                  }
                  return;
                })}
              </Table>
            </div>
          </div>
      );
    }
    return (
      <Loader size={"md"} />
    );
  };

  return (
    <MantineProvider>
      {content()}
    </MantineProvider>
  );
};
