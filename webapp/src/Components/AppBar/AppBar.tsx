import React, {useState} from 'react';
import {
    ActionIcon,
    AppShell,
    Avatar,
    Breadcrumbs,
    Burger,
    Container,
    Group,
    Menu,
    NavLink,
    Text,
    TextInput,
    UnstyledButton,
    useComputedColorScheme,
    useMantineColorScheme
} from "@mantine/core";
import classes from './AppBar.module.css';
import {useDisclosure} from "@mantine/hooks";
import {useLocation, useNavigate, useParams} from "react-router-dom";
import {useAuth} from "../../Routing/AuthContext";
import {KeyboardShortcutsHelp} from "../KeyboardShortcutsHelp/KeyboardShortcutsHelp";
import cx from 'clsx';
import {
    IconAlbum,
    IconChevronDown,
    IconHeart,
    IconHelp,
    IconHome,
    IconLibraryPhoto,
    IconLogout,
    IconMessage,
    IconMoon,
    IconPhoto,
    IconSearch,
    IconSettings,
    IconSun,
    IconTrash
} from "@tabler/icons-react";

interface IProps {
    activeLink: Labels;
    children: React.ReactNode;
}

export enum Labels {
    Dashboard = "Dashboard",
    Photos = "Photos",
    Albums = "Albums",
    Favorites = "Favorites",
    Trash = "Trash",
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
        icon: IconHeart,
        label: Labels.Favorites,
        href: '/favorites',
        description: 'Your favorite photos'
    },
    {
        icon: IconTrash,
        label: Labels.Trash,
        href: '/trash',
        description: 'Deleted photos'
    },
    {
        icon: IconSettings,
        label: Labels.Settings,
        href: '/settings',
        description: 'Manage configuration'
    },
];

const bottomNavData = [
    { icon: IconLibraryPhoto, label: 'Home', href: '/' },
    { icon: IconPhoto, label: 'Photos', href: '/photos' },
    { icon: IconAlbum, label: 'Albums', href: '/albums' },
    { icon: IconSettings, label: 'Settings', href: '/settings' },
];

const activeLinkIndex = (index: Labels) => navLinkData.map(item => item.label).indexOf(index);

const buildBreadcrumbs = (location: string, activeLink: Labels, id?: string) => {
    const items: { label: string; href?: string; icon?: React.ReactNode }[] = [];

    if (activeLink === Labels.Dashboard) {
        items.push({ label: 'Dashboard', icon: <IconHome size="0.9rem" /> });
    } else if (activeLink === Labels.Photos) {
        items.push({ label: 'Photos', href: '/photos', icon: <IconPhoto size="0.9rem" /> });
        if (location.startsWith('/photo/') && id) {
            items.push({ label: 'Detail' });
        } else if (location.startsWith('/search')) {
            items.push({ label: 'Search' });
        }
    } else if (activeLink === Labels.Albums) {
        items.push({ label: 'Albums', href: '/albums', icon: <IconAlbum size="0.9rem" /> });
        if (location.startsWith('/album/') && id) {
            items.push({ label: 'Album' });
        }
    } else if (activeLink === Labels.Favorites) {
        items.push({ label: 'Favorites', icon: <IconHeart size="0.9rem" /> });
    } else if (activeLink === Labels.Trash) {
        items.push({ label: 'Trash', icon: <IconTrash size="0.9rem" /> });
    } else if (activeLink === Labels.Settings) {
        items.push({ label: 'Settings', icon: <IconSettings size="0.9rem" /> });
    }

    return items;
};

export const AppBar: React.FunctionComponent<IProps> = (props) => {
    const [opened, {toggle}] = useDisclosure(false);
    const [active, setActive] = useState(activeLinkIndex(props.activeLink));
    const [showHelp, setShowHelp] = useState(false);
    const [searchValue, setSearchValue] = useState('');
    const {colorScheme, setColorScheme} = useMantineColorScheme();
    const computedColorScheme = useComputedColorScheme('light', {getInitialValueInEffect: true});
    const navigate = useNavigate();
    const location = useLocation();
    const { id } = useParams<{ id: string }>();
    const {logout} = useAuth();
    const breadcrumbItems = buildBreadcrumbs(location.pathname, props.activeLink, id);

    const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && searchValue.trim()) {
            navigate(`/search?q=${encodeURIComponent(searchValue.trim())}`);
        }
    };

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

    const handleLogout = () => {
        logout();
        navigate('/login');
    }

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
            <AppShell.Header>
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
                                <IconLibraryPhoto size="4rem" stroke={1.5} color={"var(--mantine-primary-color-filled)"}/>
                                <Text size="xl" fw={900} c={"var(--mantine-primary-color-filled)"}>
                                    Photobox
                                </Text>
                                <Breadcrumbs separator="›" ml="md" visibleFrom="sm">
                                    {breadcrumbItems.map((item, index) => (
                                        <Text
                                            key={index}
                                            size="sm"
                                            c={index === breadcrumbItems.length - 1 ? 'var(--mantine-primary-color-filled)' : 'dimmed'}
                                            fw={index === breadcrumbItems.length - 1 ? 600 : 400}
                                            style={{ cursor: item.href ? 'pointer' : 'default' }}
                                            onClick={() => item.href && navigate(item.href)}
                                        >
                                            <Group gap={4}>
                                                {item.icon}
                                                {item.label}
                                            </Group>
                                        </Text>
                                    ))}
                                </Breadcrumbs>
                            </Group>

                            <Group justify="space-between" gap="xl">
                                <TextInput
                                    placeholder="Search photos..."
                                    leftSection={<IconSearch size="1rem" />}
                                    value={searchValue}
                                    onChange={(e) => setSearchValue(e.target.value)}
                                    onKeyDown={handleSearch}
                                    size="sm"
                                    style={{ width: 220 }}
                                    visibleFrom="sm"
                                />
                                <Menu
                                    width={260}
                                    position="bottom-end"
                                    transitionProps={{transition: 'pop-top-right'}}
                                    onClose={() => {}}
                                    onOpen={() => {}}
                                    withinPortal
                                >
                                    <Menu.Target>
                                        <UnstyledButton
                                            className={cx(classes.user, {[classes.userActive]: true})}
                                        >
                                            <Group gap={7}>
                                                <Avatar src={undefined} alt="User" radius="xl" size={20}/>
                                                <Text fw={500} size="sm" lh={1} mr={3}>
                                                    User
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
                                        <Menu.Item onClick={handleLogout} leftSection={<IconLogout size={16} stroke={1.5}/>}>Logout</Menu.Item>
                                    </Menu.Dropdown>
                                </Menu>
                                <ActionIcon
                                    onClick={() => setShowHelp(true)}
                                    variant="default"
                                    size="xl"
                                    aria-label="Keyboard shortcuts"
                                >
                                    <IconHelp size="1.25rem" />
                                </ActionIcon>
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
            </AppShell.Header>

            <AppShell.Navbar p="md">
                {navBarItems}
            </AppShell.Navbar>

            <AppShell.Main>
                {props.children}
            </AppShell.Main>

            {/* Mobile bottom navigation */}
            <div
                className="mobile-bottom-nav"
                style={{
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    height: 64,
                    background: 'var(--mantine-color-body)',
                    borderTop: '1px solid var(--mantine-color-default-border)',
                    display: 'flex',
                    justifyContent: 'space-around',
                    alignItems: 'center',
                    zIndex: 100,
                    paddingBottom: 'env(safe-area-inset-bottom)',
                }}
            >
                {bottomNavData.map((item) => {
                    const isActive = location.pathname === item.href || (item.href !== '/' && location.pathname.startsWith(item.href));
                    return (
                        <div
                            key={item.href}
                            onClick={() => navigate(item.href)}
                            style={{
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'center',
                                gap: 4,
                                cursor: 'pointer',
                                color: isActive ? 'var(--mantine-primary-color-filled)' : 'var(--mantine-color-dimmed)',
                                fontSize: 12,
                                padding: '8px 16px',
                            }}
                        >
                            <item.icon size="1.5rem" stroke={isActive ? 2 : 1.5} />
                            <span>{item.label}</span>
                        </div>
                    );
                })}
            </div>

            <KeyboardShortcutsHelp opened={showHelp} onClose={() => setShowHelp(false)} />
        </AppShell>
    );
};
