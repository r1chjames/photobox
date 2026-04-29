import React from 'react';
import { Notification } from '@mantine/core'

interface IProps {
  text: string;
  show: boolean;
  handleStopShowing: () => void;
}

export const InfoSnackbar: React.FunctionComponent<IProps> = (props) => {
  if (!props.show) {
    return null;
  }

  return (
    <Notification title="Notification" onClick={() => props.handleStopShowing()} onClose={() => props.handleStopShowing()}>
      {props.text}
    </Notification>
  );
};
