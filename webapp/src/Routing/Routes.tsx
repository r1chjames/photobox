import React, { Component } from 'react';
import { Route, Router } from 'react-router-dom';
import history from './History';
import { PhotoIndexView } from '../Components/PhotoIndexView/PhotoIndexView';
import { Dashboard } from '../Components/Dashboard/Dashboard';
import { AlbumIndexView } from '../Components/AlbumIndexView/AlbumIndexView';

interface IProps {
  baseApiUrl: string;
}

export default class Routes extends Component<IProps> {

  public render = () => {
    return(
      <Router history={history}>
        <Route
            exact={true}
            path="/"
            render={() => <Dashboard />}
        />
        <Route
            exact={true}
            path="/photos"
            render={props => <PhotoIndexView {...props} baseApiUrl={this.props.baseApiUrl} albumId={null}/>}
        />
        <Route
          exact={true}
          path="/albums"
          render={props => <AlbumIndexView {...props} baseApiUrl={this.props.baseApiUrl} />}
        />
      </Router>
    );
  }
}
