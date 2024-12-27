import React from 'react';
import './PhotoDetail.css';
import {Photo} from '../../Models/Photo';
import {Loader, Table} from '@mantine/core';

interface IProps {
  photo: Photo;
}

export const PhotoDetail: React.FunctionComponent<IProps> = (props) => {

  const content = () => {
    if (props.photo !== undefined) {
      return (
          <div className="photoDetail__contentWrapper">
              <img
                src={props.photo.sourcePath}
                alt={props.photo.name}
                className="photoDetail__mainImage"
              />
            <div className="photoDetail__imageMetadata">
              <Table aria-label="metadata table">
                <tr>
                  <th>Parameter</th>
                  <th>Value</th>
                </tr>
                {Object.entries(props.photo.metadata).map(([key, value]) => {
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
    <div>
      {content()}
    </div>
  );
};
