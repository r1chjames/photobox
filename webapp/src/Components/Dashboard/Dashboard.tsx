import React from 'react';
import { TitleBar } from '../TitleBar/TitleBar';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import './Dashboard.css';

export const Dashboard: React.FunctionComponent = (props) => {

  return (
        <div>
            <MuiThemeProvider>
                <TitleBar />
            </MuiThemeProvider>
        </div>
  );
};
