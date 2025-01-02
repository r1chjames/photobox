import React, {useState} from 'react';
import {
    ActionIcon,
    AppShell,
    Burger,
    Flex,
    Menu,
    NavLink,
    Text,
    useComputedColorScheme,
    useMantineColorScheme
} from "@mantine/core";
import {useDisclosure} from "@mantine/hooks";
import cx from 'clsx';
import {IconAlbum, IconLibraryPhoto, IconMoon, IconPhoto, IconSun} from "@tabler/icons-react";
import classes = Menu.classes;

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
    const {setColorScheme} = useMantineColorScheme();
    const computedColorScheme = useComputedColorScheme('light', {getInitialValueInEffect: true});

    const navBarItems = navLinkData.map((item, index) => (
        <NavLink
            href={item.href}
            key={item.label}
            active={index === active}
            label={item.label.toString()}
            description={item.description}
            // rightSection={item.rightSection}
            leftSection={<item.icon size="1rem" stroke={1.5}/>}
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
                    <Flex align={"center"} m={1}>
                        <IconLibraryPhoto size="4rem" stroke={1.5} color={"#5474b4"}/>
                        <Text size="xl" fw={900} c={"#5474b4"}>
                            Photobox
                        </Text>
                    </Flex>
                </div>

                <Flex align={"center"} m={1}>
                    <ActionIcon
                        onClick={() => setColorScheme(computedColorScheme === 'light' ? 'dark' : 'light')}
                        variant="default"
                        size="xl"
                        aria-label="Toggle color scheme"
                    >
                        <IconSun className={cx(classes.icon, classes.light)} stroke={1.5} />
                        <IconMoon className={cx(classes.icon, classes.dark)} stroke={1.5} />
                    </ActionIcon>
                </Flex>
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
