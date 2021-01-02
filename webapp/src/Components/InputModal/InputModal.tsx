import React from 'react';
import './InputModal.css';
import Modal from 'react-bootstrap/Modal';
import { Button } from 'react-bootstrap';

interface IProps {
  isOpen: boolean;
  title: string;
  handleSave: () => void;
  handleClose: () => void;
}

export const InputModal: React.FunctionComponent<IProps> = (props) => {

  return (
    <Modal show={props.isOpen} dialogClassName="inputModal__content" onHide={() => props.handleClose()}>
      <Modal.Header className="inputModal__header">
        <Modal.Title className="inputModal__title">{props.title}</Modal.Title>
      </Modal.Header>
      <Modal.Body className="inputModal__body">
        {props.children}
      </Modal.Body>
      <Modal.Footer className="inputModal__footer">
        <Button
          className="inputModal__saveButton"
          variant="primary"
          onClick={() => props.handleSave()}
        >
          Save
        </Button>
        <Button
          className="inputModal__closeButton"
          variant="secondary"
          onClick={() => props.handleClose()}
        >
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
