import React, {useState} from 'react';
import './SettingModal.css';
import {TextField} from '@material-ui/core';
import {InputModal} from '../InputModal/InputModal';

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

  const checkForNullFields = () => {
    return (key!.length > 1 &&
      value!.length > 1 &&
      friendlyName!.length > 1 &&
      category!.length > 1 &&
      description!.length > 1);
  };

  const handleSave = () => {
    if (checkForNullFields()) {
      props.handleSave(key!, value!, friendlyName!, category!, description!);
    }
  };

  return (
    <InputModal
      isOpen={props.isOpen}
      title="Add Setting"
      handleSave={handleSave}
      handleClose={props.handleClose}
    >
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
    </InputModal>
  );
};
