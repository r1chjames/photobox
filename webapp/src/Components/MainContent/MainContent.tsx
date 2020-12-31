import React from 'react';
import './MainContent.css';
import { Typography } from '@material-ui/core';
import { TitleBar } from '../TitleBar/TitleBar';

interface IProps {
  title: string;
}

export const MainContent: React.FunctionComponent<IProps> = (props) => {

  return (
    <div>
      <TitleBar/>
      <div className="mainContent__mainPage">
        <div className="mainContent__title">
          <Typography variant="h4" component="h1">
            {props.title}
          </Typography>
        </div>
        <div className="mainContent__contentWrapper">
          {props.children}
        </div>
      </div>
    </div>
  );
};
