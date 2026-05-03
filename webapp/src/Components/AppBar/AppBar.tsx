import React, {useEffect, useState} from 'react';
import {
    ActionIcon,
    AppShell,
    Avatar,
    Breadcrumbs,
    Burger,
    Container,
    Drawer,
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
import {Link, useLocation, useNavigate, useParams} from "react-router-dom";
import {useAuth} from "../../Routing/AuthContext";
import {useAdapters} from "../../Routing/AdapterContext";
import {KeyboardShortcutsHelp} from "../KeyboardShortcutsHelp/KeyboardShortcutsHelp";
import {useShortcutsHint} from "./useShortcutsHint";
import cx from 'clsx';
import { AnimatePresence, motion } from 'framer-motion';
import {
    IconAlbum,
    IconChevronDown,
    IconCopy,
    IconHeart,
    IconHelp,
    IconHome,
    IconLibraryPhoto,
    IconLink,
    IconLogout,
    IconMap,
    IconMenu2,
    IconMessage,
    IconMoon,
    IconPhoto,
    IconSearch,
    IconSettings,
    IconSun,
    IconTag,
    IconTrash,
    IconUsers,
    IconVideo
} from "@tabler/icons-react";

interface IProps {
    children: React.ReactNode;
}

export enum Labels {
    Dashboard = "Dashboard",
    Photos = "Photos",
    Albums = "Albums",
    Favorites = "Favorites",
    Trash = "Trash",
    Shares = "Shares",
    Users = "Users",
    Map = "Map",
    Tags = "Tags",
    Duplicates = "Duplicates",
    Videos = "Videos",
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
        icon: IconLink,
        label: Labels.Shares,
        href: '/shares',
        description: 'Manage share links'
    },
    {
        icon: IconUsers,
        label: Labels.Users,
        href: '/users',
        description: 'Manage users'
    },
    {
        icon: IconMap,
        label: Labels.Map,
        href: '/map',
        description: 'Geotagged photos'
    },
    {
        icon: IconTag,
        label: Labels.Tags,
        href: '/tags',
        description: 'Browse by tag'
    },
    {
        icon: IconCopy,
        label: Labels.Duplicates,
        href: '/duplicates',
        description: 'Find duplicate photos'
    },
    {
        icon: IconVideo,
        label: Labels.Videos,
        href: '/videos',
        description: 'All videos'
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
    { icon: IconMap, label: 'Map', href: '/map' },
    { icon: IconMenu2, label: 'More', href: null },
];

const overflowNavLinks = navLinkData.filter(
    (item) => ![Labels.Dashboard, Labels.Photos, Labels.Albums, Labels.Map, Labels.Tags].includes(item.label)
);

const activeLinkIndex = (index: Labels) => navLinkData.map(item => item.label).indexOf(index);

const getActiveLinkFromPath = (pathname: string): Labels => {
    if (pathname === '/') return Labels.Dashboard;
    if (pathname.startsWith('/photo/')) return Labels.Photos;
    if (pathname.startsWith('/photos') || pathname.startsWith('/search') || pathname.startsWith('/favorites') || pathname.startsWith('/trash') || pathname.startsWith('/videos')) return Labels.Photos;
    if (pathname.startsWith('/album/')) return Labels.Albums;
    if (pathname === '/albums') return Labels.Albums;
    if (pathname.startsWith('/shares')) return Labels.Shares;
    if (pathname.startsWith('/users')) return Labels.Users;
    if (pathname.startsWith('/map')) return Labels.Map;
    if (pathname.startsWith('/tags')) return Labels.Tags;
    if (pathname.startsWith('/duplicates')) return Labels.Duplicates;
    if (pathname.startsWith('/settings')) return Labels.Settings;
    return Labels.Dashboard;
};

const buildBreadcrumbs = (location: string, activeLink: Labels, id?: string, albumName?: string) => {
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
            items.push({ label: albumName || 'Album' });
        }
    } else if (activeLink === Labels.Favorites) {
        items.push({ label: 'Favorites', icon: <IconHeart size="0.9rem" /> });
    } else if (activeLink === Labels.Trash) {
        items.push({ label: 'Trash', icon: <IconTrash size="0.9rem" /> });
    } else if (activeLink === Labels.Shares) {
        items.push({ label: 'Shares', icon: <IconLink size="0.9rem" /> });
    } else if (activeLink === Labels.Users) {
        items.push({ label: 'Users', icon: <IconUsers size="0.9rem" /> });
    } else if (activeLink === Labels.Map) {
        items.push({ label: 'Map', icon: <IconMap size="0.9rem" /> });
    } else if (activeLink === Labels.Tags) {
        items.push({ label: 'Tags', icon: <IconTag size="0.9rem" /> });
    } else if (activeLink === Labels.Duplicates) {
        items.push({ label: 'Duplicates', icon: <IconCopy size="0.9rem" /> });
    } else if (activeLink === Labels.Videos) {
        items.push({ label: 'Videos', icon: <IconVideo size="0.9rem" /> });
    } else if (activeLink === Labels.Settings) {
        items.push({ label: 'Settings', icon: <IconSettings size="0.9rem" /> });
    }

    return items;
};

export const AppBar: React.FunctionComponent<IProps> = (props) => {
    useShortcutsHint();
    const [opened, {toggle}] = useDisclosure(false);
    const location = useLocation();
    const activeLink = getActiveLinkFromPath(location.pathname);
    const [active, setActive] = useState(activeLinkIndex(activeLink));
    const [bottomSheetOpen, setBottomSheetOpen] = useState(false);
    const [showHelp, setShowHelp] = useState(false);
    const [searchValue, setSearchValue] = useState('');
    const [albumName, setAlbumName] = useState<string | undefined>();
    const {colorScheme, setColorScheme} = useMantineColorScheme();
    const computedColorScheme = useComputedColorScheme('light', {getInitialValueInEffect: true});
    const navigate = useNavigate();
    const { id } = useParams<{ id: string }>();
    const {logout} = useAuth();
    const {albumsAdapter} = useAdapters();
    const breadcrumbItems = buildBreadcrumbs(location.pathname, activeLink, id, albumName);

    useEffect(() => {
        setActive(activeLinkIndex(activeLink));
    }, [activeLink]);

    useEffect(() => {
        if (activeLink === Labels.Albums && id) {
            albumsAdapter.getAlbumInfoById(id)
                .then(album => setAlbumName(album?.name))
                .catch(() => setAlbumName(undefined));
        } else {
            setAlbumName(undefined);
        }
    }, [activeLink, id, albumsAdapter]);

    const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && searchValue.trim()) {
            navigate(`/search?q=${encodeURIComponent(searchValue.trim())}`);
        }
    };

    const navBarItems = navLinkData.map((item, index) => (
        <NavLink
            component={Link}
            to={item.href}
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
                                        <span
                                            key={index}
                                            style={{
                                                cursor: item.href ? 'pointer' : 'default',
                                                display: 'inline-flex',
                                                alignItems: 'center',
                                                gap: 4,
                                                color: index === breadcrumbItems.length - 1 ? 'var(--mantine-primary-color-filled)' : 'var(--mantine-color-gray-7)',
                                                fontWeight: index === breadcrumbItems.length - 1 ? 600 : 400,
                                                fontSize: 'var(--mantine-font-size-sm)',
                                            }}
                                            onClick={() => item.href && navigate(item.href)}
                                        >
                                            {item.icon}
                                            {item.label}
                                        </span>
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
                <AnimatePresence mode="wait">
                    <motion.div
                        key={location.pathname}
                        initial={{ opacity: 0, y: 8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.2 }}
                        style={{ height: '100%' }}
                    >
                        {props.children}
                    </motion.div>
                </AnimatePresence>
            </AppShell.Main>

            {/* Mobile bottom navigation */}
            <div
                className="mobile-bottom-nav"
                hidden-from="sm"
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
                    const isMore = item.href === null;
                    const isActive = isMore
                        ? overflowNavLinks.some((link) => location.pathname === link.href || location.pathname.startsWith(link.href))
                        : location.pathname === item.href || (item.href !== '/' && location.pathname.startsWith(item.href));
                    return (
                        <div
                            key={item.label}
                            onClick={() => isMore ? setBottomSheetOpen(true) : navigate(item.href)}
                            style={{
                                display: 'flex',
                                flexDirection: 'column',
                                alignItems: 'center',
                                gap: 4,
                                cursor: 'pointer',
                                color: isActive ? 'var(--mantine-primary-color-filled)' : 'var(--mantine-color-gray-7)',
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

            {/* Mobile "More" bottom sheet */}
            <Drawer
                position="bottom"
                opened={bottomSheetOpen}
                onClose={() => setBottomSheetOpen(false)}
                size="100%"
                radius="lg"
                withCloseButton={false}
                styles={{
                    body: { padding: '8px 16px 24px' },
                }}
            >
                {overflowNavLinks.map((item) => {
                    const isActive = location.pathname === item.href || location.pathname.startsWith(item.href);
                    return (
                        <NavLink
                            component={Link}
                            to={item.href}
                            key={item.label}
                            active={isActive}
                            label={item.label.toString()}
                            description={item.description}
                            leftSection={<item.icon size="1rem" stroke={isActive ? 2 : 1.5} />}
                            onClick={() => setBottomSheetOpen(false)}
                        />
                    );
                })}
            </Drawer>

            <KeyboardShortcutsHelp opened={showHelp} onClose={() => setShowHelp(false)} />
        </AppShell>
    );
};
