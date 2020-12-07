import React from 'react';
import './LoadingScreen.css';

export const LoadingScreen: React.FunctionComponent = () => {

  const loadingScreen = (
            <div className="loadingScreen__outerWrapper">
                <div className="loadingScreen__innerWrapper">
                    <div className="loadingScreen__animation">
                        <div />
                        <div />
                        <div>
                            <div />
                        </div>
                        <div>
                            <div />
                        </div>
                    </div>
                </div>
            </div>
        );
  return (loadingScreen);
};
