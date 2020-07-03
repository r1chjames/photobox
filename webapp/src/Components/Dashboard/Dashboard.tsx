import React, { Component } from 'react';
import { TitleBar } from '../TitleBar/TitleBar';
import MuiThemeProvider from 'material-ui/styles/MuiThemeProvider';
import './Dashboard.css';

export class Dashboard extends Component {

  public render = () => {
    return (
        <div>
            <MuiThemeProvider>
                <TitleBar title="Dashboard"/>
            </MuiThemeProvider>
        </div>
    );

  }
}
