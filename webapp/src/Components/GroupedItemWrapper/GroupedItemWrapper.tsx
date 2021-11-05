import React from 'react';

interface IProps {
  groupId: string;
}

export const GroupedItemWrapper: React.FunctionComponent<IProps> = (props) => {

  return (
    <div>
      {props.groupId}
      {props.children}
    </div>
  );
};
