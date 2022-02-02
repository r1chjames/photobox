import React from 'react';
import { Notification } from '@mantine/core'
import './Snackbar.css';

interface IProps {
  text: string;
  show: boolean;
  handleStopShowing: () => void;
}

export const InfoSnackbar: React.FunctionComponent<IProps> = (props) => {

  return (
    <Notification title="Notification" onClick={() => props.handleStopShowing()} onClose={() => props.handleStopShowing()}>
      {props.text}
    </ Notification>
      // anchorOrigin={{
      //   vertical: 'bottom',
      //   horizontal: 'center',
      // }}
      // open={props.show}
      // autoHideDuration={5000}
      // onClose={() => props.handleStopShowing()}
      // message={props.text}
      // action={
      //   <React.Fragment>
      //     <IconButton size="small" aria-label="close" color="inherit" onClick={() => props.handleStopShowing()} />
      //   </React.Fragment>
      // }
    // />
  );
};
