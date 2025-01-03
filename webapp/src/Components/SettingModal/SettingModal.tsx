import React, { useState } from 'react';
import './SettingModal.css';
import { InputModal } from '../InputModal/InputModal';
import {TextInput} from '@mantine/core';

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
        <TextInput
          required={true}
          id="key"
          label="Setting"
          error={keyErrorText.length !== 0}
          onChange={e => handleKeyChange(e.target.value)}
        />
        <p/>
        <TextInput
          required={true}
          id="value"
          label="Value"
          error={valueErrorText.length !== 0}
          onChange={e => handleValueChange(e.target.value)}
        />
        <p/>
        <TextInput
          id="friendlyName"
          label="Friendly Name"
          onChange={e => setFriendlyName(e.target.value)}
        />
        <p/>
        <TextInput
          id="category"
          label="Category"
          error={categoryErrorText.length !== 0}
          onChange={e => handleCategoryChange(e.target.value)}
        />
        <p/>
        <TextInput
          id="description"
          label="Description"
          onChange={e => setDescription(e.target.value)}
        />
    </InputModal>
  );
};
