import React, {useState} from 'react';
import {AppShell, Burger, Flex, NavLink, Text} from "@mantine/core";
import {useDisclosure} from "@mantine/hooks";
import {IconAlbum, IconLibraryPhoto, IconPhoto} from "@tabler/icons-react";

interface IProps {
    activeLink: string;
}

const navLinkData = [
    {
        icon: IconLibraryPhoto,
        label: 'Dashboard',
        href: '/',
        description: 'All photos & albums'
    },
    {
        icon: IconPhoto,
        label: 'Photos',
        href: "/photos",
        description: 'All photos'
    },
    {
        icon: IconAlbum,
        label: 'Albums',
        href: '/albums',
        description: 'All albums'
    },
];

const activeLinkIndex = (index: string) => navLinkData.map(item => item.label).indexOf(index)

export const AppBar: React.FunctionComponent<IProps> = (props) => {
  const [opened, {toggle}] = useDisclosure(false);
  const [active, setActive] = useState(activeLinkIndex(props.activeLink));

    const navBarItems = navLinkData.map((item, index) => (
        <NavLink
            href="#required-for-focus"
            key={item.label}
            active={index === active}
            label={item.label}
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
          <Flex align={ "center" }>
              <IconLibraryPhoto size="4rem" stroke={1.5} color={"#228be6"}/>
              <Text size="xl" fw={900} c={"#228be6"}>
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
