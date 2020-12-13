import React, { Component } from 'react';
import { Route, Router, Switch } from 'react-router-dom';
import history from './History';
import { PhotoIndexView } from '../Components/PhotoIndexView/PhotoIndexView';
import { Dashboard } from '../Components/Dashboard/Dashboard';
import { AlbumIndexView } from '../Components/AlbumIndexView/AlbumIndexView';
import { SettingsView } from '../Components/SettingsView/SettingsView';

interface IProps {
  baseApiUrl: string;
}

export default class Routes extends Component<IProps> {

  public render = () => {
    return(
      <Router history={history}>
        <Switch>
          <Route
              exact={true}
              path="/"
              render={props => <Dashboard {...props} baseApiUrl={this.props.baseApiUrl} />}
          />
          <Route
              exact={true}
              path="/photos"
              render={props => <PhotoIndexView {...props} baseApiUrl={this.props.baseApiUrl} />}
          />
          <Route
            exact={true}
            path="/albums"
            render={props => <AlbumIndexView {...props} baseApiUrl={this.props.baseApiUrl} />}
          />
          <Route
            exact={true}
            path="/albums/:id"
            render={props => <PhotoIndexView {...props} baseApiUrl={this.props.baseApiUrl} />}
          />
          <Route
            exact={true}
            path="/settings"
            render={props => <SettingsView {...props} baseApiUrl={this.props.baseApiUrl} />}
          />
        </Switch>
      </Router>
    );
  }
}
