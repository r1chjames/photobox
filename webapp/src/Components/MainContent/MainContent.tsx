import React from 'react';
import './MainContent.css';

export const MainContent: React.FunctionComponent = (props) => {

  return (
        <div className="mainContent__mainPage">
            {props.children}
        </div>
  );
};
