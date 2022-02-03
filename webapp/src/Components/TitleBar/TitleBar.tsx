import React from 'react';
import {AppShell, Navbar, Header, Title} from '@mantine/core';
import { useNavigate } from 'react-router-dom';
import './TitleBar.css';

export const TitleBar: React.FunctionComponent = (props) => {
  const navigate = useNavigate()
  return (
    <div className="titleBar__mainBar">
      <AppShell
        padding="md"
        navbar={
          <Navbar width={{base: 200}} height={500} padding="xs">
            <Navbar.Section>
              <Title onClick={() => navigate('/')}>Dashboard</Title>
              <Title onClick={() => navigate('/photos')}>Photos</Title>
              <Title onClick={() => navigate('/albums')}>Albums</Title>
              <Title onClick={() => navigate('/settings')}>Settings</Title>
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
