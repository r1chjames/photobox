import React, {useState} from 'react';
import './SettingModal.css';
import Modal from 'react-bootstrap/Modal';
import {Button} from 'react-bootstrap';
import {TextField} from '@material-ui/core';

interface IProps {
  isOpen: boolean;
  handleSave: (key: string, value: string, friendlyName: string, category: string, description: string) => void;
  handleClose: () => void;
}

export const SettingModal: React.FunctionComponent<IProps> = (props) => {

  const [key, setKey] = useState<string>();
  const [value, setValue] = useState<string>();
  const [friendlyName, setFriendlyName] = useState<string>();
  const [category, setCategory] = useState<string>();
  const [description, setDescription] = useState<string>();

  const handleKeyChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
    setKey(e.target.value);
  };

  const handleValueChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
    setValue(e.target.value);
  };
  const handleFriendlyNameChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
    setFriendlyName(e.target.value);
  };
  const handleCategoryChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
    setCategory(e.target.value);
  };
  const handleDescriptionChange = (e: React.ChangeEvent<HTMLTextAreaElement | HTMLInputElement>) => {
    setDescription(e.target.value);
  };

  return (
    <Modal show={props.isOpen} dialogClassName="settingModal__content">
      <Modal.Header className="settingModal__header">
        <Modal.Title className="settingModal__title">Add Setting</Modal.Title>
      </Modal.Header>
      <Modal.Body className="settingModal__body">
        <TextField
          required={true}
          id="outlined-required"
          label="Setting"
          variant="outlined"
          onChange={e => handleKeyChange(e)}
        />
        <p/>
        <TextField
          required={true}
          id="outlined-required"
          label="Value"
          variant="outlined"
          onChange={e => handleValueChange(e)}
        />
        <p/>
        <TextField
          id="outlined"
          label="Friendly Name"
          variant="outlined"
          onChange={e => handleFriendlyNameChange(e)}
        />
        <p/>
        <TextField
          id="outlined"
          label="Category"
          variant="outlined"
          onChange={e => handleCategoryChange(e)}
        />
        <p/>
        <TextField
          id="outlined"
          label="Description"
          variant="outlined"
          onChange={e => handleDescriptionChange(e)}
        />
      </Modal.Body>
      <Modal.Footer className="settingModal__footer">
        <Button
          className="settingModal__saveButton"
          variant="primary"
          onClick={() => props.handleSave(key!, value!, friendlyName!, category!, description!)}
        >
          Save
        </Button>
        <Button
          className="settingModal__closeButton"
          variant="secondary"
          onClick={() => props.handleClose()}
        >
          Close
        </Button>
      </Modal.Footer>
    </Modal>
  );
};
