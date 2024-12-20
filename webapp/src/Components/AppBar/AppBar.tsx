import React from 'react';
import {AppShell, Burger, Flex, NavLink, Text} from "@mantine/core";
import {useDisclosure} from "@mantine/hooks";
import {IconAlbum, IconLibraryPhoto, IconPhoto} from "@tabler/icons-react";

export const AppBar: React.FunctionComponent = (props) => {
  const [opened, {toggle}] = useDisclosure(false);
  return (
      <AppShell
          header={{height: 60}}
          navbar={{
            width: 300,
            breakpoint: 'sm',
            collapsed: {mobile: !opened},
          }}
          padding="md"
      >
        <AppShell.Header>
          <Burger
              opened={opened}
              onClick={toggle}
              hiddenFrom="sm"
              size="sm"
          />
        <div>
          <Flex align={ "center" }>
              <IconLibraryPhoto size="4rem" stroke={1.5} color={"#228be6"}/>
              <Text size="xl" fw={900} c={"#228be6"}>
                Photobox
              </Text>
          </Flex>
        </div>
        </AppShell.Header>

        <AppShell.Navbar p="md">
            <NavLink
                href="/photos"
                label="Photos"
                leftSection={<IconPhoto size="2rem" stroke={1.5} />}
            />
            <NavLink
                href="/albums"
                label="Albums"
                leftSection={<IconAlbum size="2rem" stroke={1.5} />}
            />
        </AppShell.Navbar>

        <AppShell.Main>
            {props.children}
        </AppShell.Main>
      </AppShell>
  );
};
