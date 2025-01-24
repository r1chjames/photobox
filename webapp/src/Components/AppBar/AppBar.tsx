import React, {useState} from 'react';
import {
    ActionIcon,
    AppShell,
    Avatar,
    Burger,
    Container,
    Group,
    Menu,
    NavLink,
    Text,
    UnstyledButton,
    useComputedColorScheme,
    useMantineColorScheme
} from "@mantine/core";
import classes from './AppBar.module.css';
import {useDisclosure} from "@mantine/hooks";
import cx from 'clsx';
import {
    IconAlbum,
    IconChevronDown,
    IconLibraryPhoto,
    IconLogout,
    IconMessage,
    IconMoon,
    IconPhoto,
    IconSettings,
    IconSun
} from "@tabler/icons-react";

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
    const {colorScheme, setColorScheme} = useMantineColorScheme();
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

    const colourSchemeButton = (colourScheme: string) => {
        return colourScheme === 'light' ?
            <IconMoon className={cx(classes.icon, classes.light)} stroke={1.5}/>
            : <IconSun className={cx(classes.icon, classes.dark)} stroke={1.5}/>
    }

    return (
        <AppShell
            header={{
                height: 80
            }}
            navbar={{
                width: 200,
                breakpoint: 'sm',
                collapsed: {mobile: !opened},
            }}
            padding="md"
        >
            <div className={classes.header}>
                <Container fluid>
                    <Group justify="space-between" gap="xl">
                        <Burger
                            opened={opened}
                            onClick={toggle}
                            hiddenFrom="sm"
                            size="sm"
                        />
                        <Group>
                            <IconLibraryPhoto size="4rem" stroke={1.5} color={"#5474b4"}/>
                            <Text size="xl" fw={900} c={"#5474b4"}>
                                Photobox
                            </Text>
                        </Group>

                        <Group justify="space-between" gap="xl">
                            <Menu
                                width={260}
                                position="bottom-end"
                                transitionProps={{transition: 'pop-top-right'}}
                                onClose={() => console.log("close")}
                                onOpen={() => console.log("open")}
                                withinPortal
                            >
                                <Menu.Target>
                                    <UnstyledButton
                                        className={cx(classes.user, {[classes.userActive]: true})}
                                    >
                                        <Group gap={7}>
                                            <Avatar src={"user.image"} alt={"user.name"} radius="xl" size={20}/>
                                            <Text fw={500} size="sm" lh={1} mr={3}>
                                                {"Rich"}
                                            </Text>
                                            <IconChevronDown size={12} stroke={1.5}/>
                                        </Group>
                                    </UnstyledButton>
                                </Menu.Target>
                                <Menu.Dropdown>
                                    <Menu.Item
                                        leftSection={<IconMessage size={16} color={"blue"} stroke={1.5}/>}
                                    >
                                        Your comments
                                    </Menu.Item>
                                    <Menu.Label>Settings</Menu.Label>
                                    <Menu.Item leftSection={<IconSettings size={16} stroke={1.5}/>}>
                                        Account settings
                                    </Menu.Item>
                                    <Menu.Divider/>
                                    <Menu.Item leftSection={<IconLogout size={16} stroke={1.5}/>}>Logout</Menu.Item>
                                </Menu.Dropdown>
                            </Menu>
                            <ActionIcon
                                onClick={() => setColorScheme(computedColorScheme === 'light' ? 'dark' : 'light')}
                                variant="default"
                                size="xl"
                                aria-label="Toggle color scheme"
                            >
                                {colourSchemeButton(colorScheme)}
                            </ActionIcon>
                        </Group>
                    </Group>
                </Container>

            </div>

            <AppShell.Navbar p="md">
                {navBarItems}
            </AppShell.Navbar>

            <AppShell.Main>
                {props.children}
            </AppShell.Main>
        </AppShell>
    );
};
