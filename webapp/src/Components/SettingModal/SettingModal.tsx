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

  const [key, setKey] = useState('');
  const [value, setValue] = useState('');
  const [friendlyName, setFriendlyName] = useState('');
  const [category, setCategory] = useState('');
  const [description, setDescription] = useState('');
  const [keyErrorText, setKeyErrorText] = useState('Required');
  const [valueErrorText, setValueErrorText] = useState('Required');
  const [categoryErrorText, setCategoryErrorText] = useState('Required');

  const handleKeyChange = (fieldValue: string) => {
    if (fieldValue.length > 1) {
      setKeyErrorText('');
      setKey(fieldValue);
    } else {
      setKeyErrorText('Invalid length');
    }
  };

  const handleValueChange = (fieldValue: string) => {
    if (fieldValue.length > 1) {
      setValueErrorText('');
      setValue(fieldValue);
    } else {
      setValueErrorText('Invalid length');
    }
  };

  const handleCategoryChange = (fieldValue: string) => {
    if (fieldValue.length > 1) {
      setCategoryErrorText('');
      setCategory(fieldValue);
    } else {
      setCategoryErrorText('Invalid length');
    }
  };

  const handleSave = () => {
    if (key!.length > 1 &&
      value!.length > 1 &&
      friendlyName!.length > 1 &&
      category!.length > 1 &&
      description!.length > 1
    ) {
      props.handleSave(key!, value!, friendlyName!, category!, description!);
    }
  };

  return (
    <Modal show={props.isOpen} dialogClassName="settingModal__content" onHide={() => props.handleClose()}>
      <Modal.Header className="settingModal__header">
        <Modal.Title className="settingModal__title">Add Setting</Modal.Title>
      </Modal.Header>
      <Modal.Body className="settingModal__body">
        <TextField
          required={true}
          id="key"
          label="Setting"
          variant="outlined"
          error={keyErrorText.length !== 0}
          helperText={keyErrorText}
          onChange={e => handleKeyChange(e.target.value)}
        />
        <p/>
        <TextField
          required={true}
          id="value"
          label="Value"
          variant="outlined"
          error={valueErrorText.length !== 0}
          helperText={valueErrorText}
          onChange={e => handleValueChange(e.target.value)}
        />
        <p/>
        <TextField
          id="friendlyName"
          label="Friendly Name"
          variant="outlined"
          onChange={e => setFriendlyName(e.target.value)}
        />
        <p/>
        <TextField
          id="category"
          label="Category"
          variant="outlined"
          error={categoryErrorText.length !== 0}
          helperText={categoryErrorText}
          onChange={e => handleCategoryChange(e.target.value)}
        />
        <p/>
        <TextField
          id="description"
          label="Description"
          variant="outlined"
          onChange={e => setDescription(e.target.value)}
        />
      </Modal.Body>
      <Modal.Footer className="settingModal__footer">
        <Button
          className="settingModal__saveButton"
          variant="primary"
          onClick={() => handleSave()}
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
