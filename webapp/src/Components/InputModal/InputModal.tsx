import React from 'react';
import {Button, Group, Modal} from '@mantine/core';

interface IProps {
  isOpen: boolean;
  title: string;
  handleSave: () => void;
  handleClose: () => void;
  children?: React.ReactNode;
}

export const InputModal: React.FunctionComponent<IProps> = (props) => {

  return (
    <Modal
      opened={props.isOpen}
      onClose={props.handleClose}
      title={props.title}
      centered
    >
      {props.children}
      <Group justify="flex-end" mt="md">
        <Button variant="default" onClick={props.handleClose}>
          Close
        </Button>
        <Button onClick={props.handleSave}>
          Save
        </Button>
      </Group>
    </Modal>
  );
};
