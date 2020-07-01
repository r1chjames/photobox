import React, { Component } from 'react';
import './LoadingScreen.css';

export class LoadingScreen extends Component {

  public render = () => {
    const loadingScreen = (
            <div className="loadingScreen__outerWrapper">
                <div className="loadingScreen__innerWrapper">
                    <div className="loadingScreen__animation">
                        <div></div>
                        <div></div>
                        <div>
                            <div></div>
                        </div>
                        <div>
                            <div></div>
                        </div>
                    </div>
                </div>
            </div>
        );
    return (loadingScreen);
  }
}
