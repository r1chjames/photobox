import React from 'react';
import './Snackbar.css';
import {IconButton, Snackbar} from '@material-ui/core';

interface IProps {
  text: string;
  show: boolean;
  handleStopShowing: () => void;
}

export const InfoSnackbar: React.FunctionComponent<IProps> = (props) => {

  return (
    <Snackbar
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'center',
      }}
      open={props.show}
      autoHideDuration={5000}
      onClose={() => props.handleStopShowing()}
      message={props.text}
      action={
        <React.Fragment>
          <IconButton size="small" aria-label="close" color="inherit" onClick={() => props.handleStopShowing()} />
        </React.Fragment>
      }
    />
  );
};
