import React from 'react';
import {AppShell, Navbar, Header, Title} from '@mantine/core';
import history from '../../Routing/History';
import './TitleBar.css';

export const TitleBar: React.FunctionComponent = (props) => {
  return (
    <div className="titleBar__mainBar">
      <AppShell
        padding="md"
        navbar={
          <Navbar width={{base: 300}} height={500} padding="xs">
            <Navbar.Section>
              <Title onClick={(() => history.push('/'))}>Dashboard</Title>
              <Title onClick={(() => history.push('/photos'))}>Photos</Title>
              <Title onClick={(() => history.push('/albums'))}>Albums</Title>
              <Title onClick={(() => history.push('/settings'))}>Settings</Title>
            </Navbar.Section>
          </Navbar>
        }
        header={
          <Header height={60} padding="xs">
            Photobox
          </Header>
        }
        styles={(theme) => ({
          main: {backgroundColor: theme.colorScheme === 'dark' ? theme.colors.dark[8] : theme.colors.gray[0]},
        })}
      >
        {props.children}
      </AppShell>
    </div>
  );
};
