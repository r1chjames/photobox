import React from 'react';
import './MainContent.css';
import { TitleBar } from '../TitleBar/TitleBar';

interface IProps {
  title: string;
}

export const MainContent: React.FunctionComponent<IProps> = (props) => {

  return (
    <div>
      <TitleBar />
      <div className="mainContent__mainPage">
        <div className="mainContent__title">
          {props.title}
        </div>
        <div className="mainContent__contentWrapper">
          {props.children}
        </div>
      </div>
    </div>
  );
};
