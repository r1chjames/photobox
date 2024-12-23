import React, {useState} from 'react';
import {AppShell, Burger, Flex, NavLink, Text} from "@mantine/core";
import {useDisclosure} from "@mantine/hooks";
import {IconAlbum, IconLibraryPhoto, IconPhoto} from "@tabler/icons-react";

interface IProps {
    activeLink: Labels;
}

export enum Labels {
    Dashboard = "Dashboard",
    Photos = "Photos",
    Albums = "Albums",
    Settings = "Settings"
}

const navLinkData = [
    {
        icon: IconLibraryPhoto,
        label: Labels.Dashboard,
        href: '/',
        description: 'All photos & albums'
    },
    {
        icon: IconPhoto,
        label: Labels.Photos,
        href: "/photos",
        description: 'All photos'
    },
    {
        icon: IconAlbum,
        label: Labels.Albums,
        href: '/albums',
        description: 'All albums'
    },
    {
        icon: IconAlbum,
        label: Labels.Settings,
        href: '/settings',
        description: 'Manage configuration'
    },
];

const activeLinkIndex = (index: Labels) => navLinkData.map(item => item.label).indexOf(index);

export const AppBar: React.FunctionComponent<IProps> = (props) => {
  const [opened, {toggle}] = useDisclosure(false);
  const [active, setActive] = useState(activeLinkIndex(props.activeLink));

    const navBarItems = navLinkData.map((item, index) => (
        <NavLink
            href={item.href}
            key={item.label}
            active={index === active}
            label={item.label.toString()}
            description={item.description}
            // rightSection={item.rightSection}
            leftSection={<item.icon size="1rem" stroke={1.5} />}
            onClick={() => setActive(index)}
        />
    ));

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
          <Flex align={ "center" } m={1}>
              <IconLibraryPhoto size="4rem" stroke={1.5} color={"#5474b4"}/>
              <Text size="xl" fw={900} c={"#5474b4"}>
                Photobox
              </Text>
          </Flex>
        </div>
        </AppShell.Header>

        <AppShell.Navbar p="md">
            {navBarItems}
        </AppShell.Navbar>

        <AppShell.Main>
            {props.children}
        </AppShell.Main>
      </AppShell>
  );
};
