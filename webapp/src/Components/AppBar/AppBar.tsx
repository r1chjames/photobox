import React, {useEffect, useRef, useState} from 'react';
import {Modal, Text, Group} from "@mantine/core";
import {
    IconAlbum,
    IconChevronDown,
    IconCopy,
    IconDots,
    IconHeart,
    IconHelp,
    IconHome2,
    IconInfoCircle,
    IconLibraryPhoto,
    IconLink,
    IconLogout,
    IconMap,
    IconMoon,
    IconPhoto,
    IconPhotoOff,
    IconClock,
    IconSearch,
    IconSettings,
    IconSun,
    IconTag,
    IconTrash,
    IconUsers,
    IconVideo
} from "@tabler/icons-react";
import {useLocation, useNavigate} from "react-router-dom";
import {useMantineColorScheme} from "@mantine/core";
import {useAuth} from "../../Routing/AuthContext";
import {useAdapters} from "../../Routing/AdapterContext";
import {DragUploadOverlay} from "../DragUploadOverlay/DragUploadOverlay";
import {KeyboardShortcutsHelp} from "../KeyboardShortcutsHelp/KeyboardShortcutsHelp";

interface IProps {
    children: React.ReactNode;
}

export enum Labels {
    Home = "Home",
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
    Quality = "Quality",
    Memories = "Memories",
    Settings = "Settings",
    More = "More"
}

const primaryNavLinkData = [
    { icon: IconHome2, label: Labels.Home, href: '/' },
    { icon: IconPhoto, label: Labels.Photos, href: "/photos" },
    { icon: IconAlbum, label: Labels.Albums, href: '/albums' },
    { icon: IconClock, label: Labels.Memories, href: '/memories' },
    { icon: IconSettings, label: Labels.Settings, href: '/settings' },
];

const moreNavLinkData = [
    { icon: IconTrash, label: Labels.Trash, href: '/trash' },
    { icon: IconPhotoOff, label: Labels.Quality, href: '/quality' },
    { icon: IconLink, label: Labels.Shares, href: '/shares' },
    { icon: IconMap, label: Labels.Map, href: '/map' },
    { icon: IconTag, label: Labels.Tags, href: '/tags' },
    { icon: IconCopy, label: Labels.Duplicates, href: '/duplicates' },
    { icon: IconVideo, label: Labels.Videos, href: '/videos' },
];

const getActiveLinkFromPath = (pathname: string): Labels => {
    if (pathname === '/') return Labels.Home;
    if (pathname.startsWith('/photo/')) return Labels.Photos;
    if (pathname.startsWith('/photos') || pathname.startsWith('/search') || pathname.startsWith('/favorites') || pathname.startsWith('/trash') || pathname.startsWith('/videos')) return Labels.Photos;
    if (pathname.startsWith('/album/')) return Labels.Albums;
    if (pathname === '/albums') return Labels.Albums;
    if (pathname.startsWith('/shares')) return Labels.Shares;
    if (pathname.startsWith('/users')) return Labels.Users;
    if (pathname.startsWith('/map')) return Labels.Map;
    if (pathname.startsWith('/tags')) return Labels.Tags;
    if (pathname.startsWith('/duplicates')) return Labels.Duplicates;
    if (pathname.startsWith('/quality')) return Labels.Quality;
    if (pathname.startsWith('/memories')) return Labels.Memories;
    if (pathname.startsWith('/settings') || pathname.startsWith('/account-settings')) return Labels.Settings;
    return Labels.Home;
};

const USERNAME_STORAGE_KEY = 'pb-username';

const initialsOf = (name: string): string => {
    const clean = name.trim();
    if (!clean) return 'U';
    return clean.slice(0, 2).toUpperCase();
};

export const AppBar: React.FunctionComponent<IProps> = (props) => {
    const location = useLocation();
    const navigate = useNavigate();
    const activeLink = getActiveLinkFromPath(location.pathname);
    const {logout} = useAuth();
    const {photosAdapter} = useAdapters();
    const {colorScheme, setColorScheme} = useMantineColorScheme();
    const [searchValue, setSearchValue] = useState('');
    const [showHelp, setShowHelp] = useState(false);
    const [showAbout, setShowAbout] = useState(false);
    const [appVersion, setAppVersion] = useState<string>('');
    const [menuOpen, setMenuOpen] = useState(false);
    const [moreOpen, setMoreOpen] = useState(false);
    const [username] = useState<string>(() => localStorage.getItem(USERNAME_STORAGE_KEY) || 'User');
    const clusterRef = useRef<HTMLDivElement>(null);
    const moreRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        // Fetch app version from health endpoint
        fetch('/api/health')
            .then(r => r.json())
            .then(data => {
                if (data.version) {
                    setAppVersion(data.version);
                }
            })
            .catch(() => {});
    }, []);

    // Global keyboard shortcut: ? opens help dialog
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            const tag = (e.target as HTMLElement)?.tagName;
            if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;
            if (e.key === '?') {
                e.preventDefault();
                setShowHelp(true);
            }
        };
        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, []);

    // Close the user menu on outside click
    useEffect(() => {
        if (!menuOpen) return;
        const onDown = (e: MouseEvent) => {
            if (clusterRef.current && !clusterRef.current.contains(e.target as Node)) {
                setMenuOpen(false);
            }
        };
        document.addEventListener('mousedown', onDown);
        return () => document.removeEventListener('mousedown', onDown);
    }, [menuOpen]);

    // Close the More submenu on outside click
    useEffect(() => {
        if (!moreOpen) return;
        const onDown = (e: MouseEvent) => {
            if (moreRef.current && !moreRef.current.contains(e.target as Node)) {
                setMoreOpen(false);
            }
        };
        document.addEventListener('mousedown', onDown);
        return () => document.removeEventListener('mousedown', onDown);
    }, [moreOpen]);

    const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && searchValue.trim()) {
            navigate(`/search?q=${encodeURIComponent(searchValue.trim())}`);
        }
    };

    const handleLogout = () => {
        setMenuOpen(false);
        logout();
        navigate('/login');
    };

    const isDark = colorScheme !== 'light';

    const moreActive = moreNavLinkData.some(item => activeLink === item.label);

    return (
        <div className="app-shell">
            <aside className="sidebar">
                <a
                    className="brand-logo"
                    role="button"
                    tabIndex={0}
                    aria-label="Photobox home"
                    onClick={() => navigate('/')}
                    onKeyDown={(e) => e.key === 'Enter' && navigate('/')}
                >
                    <div className="brand-icon">
                        <IconLibraryPhoto size={18} stroke={2.2} color="#fff" />
                    </div>
                    <span className="brand-name">Photobox</span>
                </a>

                <nav className="nav-group">
                    {primaryNavLinkData.map((item) => (
                        <button
                            key={item.label}
                            type="button"
                            className={activeLink === item.label ? 'nav-item active' : 'nav-item'}
                            title={item.label}
                            aria-label={item.label}
                            onClick={() => navigate(item.href)}
                        >
                            <item.icon size={20} stroke={1.5} />
                            <span>{item.label}</span>
                        </button>
                    ))}

                    <div className="nav-more" data-open={moreOpen} ref={moreRef}>
                        <button
                            type="button"
                            className={moreActive ? 'nav-item active' : 'nav-item'}
                            title="More"
                            aria-label="More"
                            aria-expanded={moreOpen}
                            onClick={() => setMoreOpen((open) => !open)}
                        >
                            <IconDots size={20} stroke={1.5} />
                            <span>More</span>
                            <IconChevronDown size={14} stroke={2} className="nav-more-chevron" />
                        </button>
                        {moreOpen && (
                            <div className="nav-more-menu">
                                {moreNavLinkData.map((item) => (
                                    <button
                                        key={item.label}
                                        type="button"
                                        className={activeLink === item.label ? 'nav-more-item active' : 'nav-more-item'}
                                        title={item.label}
                                        aria-label={item.label}
                                        onClick={() => { setMoreOpen(false); navigate(item.href); }}
                                    >
                                        <item.icon size={18} stroke={1.5} />
                                        <span>{item.label}</span>
                                    </button>
                                ))}
                            </div>
                        )}
                    </div>
                </nav>

                <div className="sidebar-footer">
                    <div className="footer-links">
                        <a href="#shortcuts" onClick={(e) => { e.preventDefault(); setShowHelp(true); }}>Shortcuts</a>
                        <a href="#about" onClick={(e) => { e.preventDefault(); setShowAbout(true); }}>About</a>
                    </div>
                </div>
            </aside>

            <div className="main-wrapper">
                <div className="float-cluster" ref={clusterRef}>
                    <div className="search-box">
                        <IconSearch size={16} stroke={2} />
                        <input
                            type="text"
                            placeholder="Search photos, locations, tags..."
                            value={searchValue}
                            onChange={(e) => setSearchValue(e.target.value)}
                            onKeyDown={handleSearch}
                            aria-label="Search"
                        />
                    </div>
                    <div className="float-cluster-actions">
                        <button
                            type="button"
                            className="icon-btn"
                            aria-label="Keyboard shortcuts"
                            title="Keyboard shortcuts (?)"
                            onClick={() => setShowHelp(true)}
                        >
                            <IconHelp size={18} stroke={1.8} />
                        </button>
                        <button
                            type="button"
                            className="icon-btn"
                            aria-label="Toggle color scheme"
                            title="Toggle theme"
                            onClick={() => setColorScheme(isDark ? 'light' : 'dark')}
                        >
                            {isDark ? <IconMoon size={18} stroke={1.8} /> : <IconSun size={18} stroke={1.8} />}
                        </button>
                        <div
                            className={menuOpen ? 'user-pill menu-open' : 'user-pill'}
                            role="button"
                            tabIndex={0}
                            aria-haspopup="true"
                            aria-expanded={menuOpen}
                            onClick={() => setMenuOpen((open) => !open)}
                            onKeyDown={(e) => e.key === 'Enter' && setMenuOpen((open) => !open)}
                        >
                            <div className="user-avatar">{initialsOf(username)}</div>
                            <span className="user-name">{username}</span>
                            <IconChevronDown size={14} stroke={2} className="user-menu-chevron" />
                        </div>
                    </div>

                    {menuOpen && (
                        <div className="user-menu">
                            <div className="user-menu-header">
                                <div className="user-menu-avatar">{initialsOf(username)}</div>
                                <div className="user-menu-id">
                                    <span className="user-menu-name">{username}</span>
                                    <span className="user-menu-email">Signed in to Photobox</span>
                                </div>
                            </div>
                            <div className="user-menu-divider" />
                            <button type="button" className="user-menu-item" onClick={() => { setMenuOpen(false); navigate('/users'); }}>
                                <IconUsers size={16} stroke={1.8} /> Users
                            </button>
                            <button type="button" className="user-menu-item" onClick={() => { setMenuOpen(false); navigate('/account-settings'); }}>
                                <IconSettings size={16} stroke={1.8} /> Account settings
                            </button>
                            <button type="button" className="user-menu-item" onClick={() => { setMenuOpen(false); navigate('/favorites'); }}>
                                <IconHeart size={16} stroke={1.8} /> Favorites
                            </button>
                            <div className="user-menu-divider" />
                            <button type="button" className="user-menu-item" onClick={() => { setMenuOpen(false); setShowAbout(true); }}>
                                <IconInfoCircle size={16} stroke={1.8} /> About
                            </button>
                            <button type="button" className="user-menu-item user-menu-item-danger" onClick={handleLogout}>
                                <IconLogout size={16} stroke={1.8} /> Logout
                            </button>
                        </div>
                    )}
                </div>

                <main className="content-area">
                    {props.children}
                </main>
            </div>

            <KeyboardShortcutsHelp opened={showHelp} onClose={() => setShowHelp(false)} />

            <Modal opened={showAbout} onClose={() => setShowAbout(false)} title="About Photobox" centered>
                <Group justify="center" gap="xs">
                    <Text size="lg" fw={600}>Photobox</Text>
                </Group>
                <Text c="dimmed" size="sm" ta="center" mt="xs">
                    Version {appVersion || 'unknown'}
                </Text>
            </Modal>

            <DragUploadOverlay photosAdapter={photosAdapter} />
        </div>
    );
};
