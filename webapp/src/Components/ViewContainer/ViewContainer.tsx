import React from 'react';
import './ViewContainer.css';
import { TitleBar } from '../TitleBar/TitleBar';


export const ViewContainer: React.FunctionComponent = (props) => {
  return (
    <div>
      <TitleBar />
      {props.children}
    </div>
  );
};
