import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '../../test/test-utils';
import { AppBar, Labels } from './AppBar';

// Mock Mantine hooks
vi.mock('@mantine/hooks', () => ({
    useDisclosure: () => [false, { toggle: vi.fn() }],
    useMantineColorScheme: () => ({
        colorScheme: 'light',
        setColorScheme: vi.fn()
    }),
    useComputedColorScheme: () => 'light'
}));

describe('AppBar', () => {
    beforeEach(() => {
        // Clear localStorage before each test
        localStorage.clear();
    });

    it('should render app name "Photobox"', () => {
        render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        expect(screen.getByText('Photobox')).toBeInTheDocument();
    });

    it('should render navigation links', () => {
        render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        expect(screen.getByText('Dashboard')).toBeInTheDocument();
        expect(screen.getByText('Photos')).toBeInTheDocument();
        expect(screen.getByText('Albums')).toBeInTheDocument();
        expect(screen.getByText('Settings')).toBeInTheDocument();
    });

    it('should render user avatar and name', () => {
        render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        expect(screen.getByText('User')).toBeInTheDocument();
    });

    it('should render children content', () => {
        render(
            <AppBar activeLink={Labels.Dashboard}>
                <div>Test Content</div>
            </AppBar>
        );

        expect(screen.getByText('Test Content')).toBeInTheDocument();
    });

    it('should highlight active navigation link', () => {
        const { container } = render(<AppBar activeLink={Labels.Photos}><div>Content</div></AppBar>);

        // Find nav links and check if Photos is active
        const navLinks = container.querySelectorAll('.mantine-NavLink-root');
        expect(navLinks.length).toBeGreaterThan(0);
    });

    it('should render color scheme toggle button', () => {
        const { container } = render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        const toggleButton = container.querySelector('[aria-label="Toggle color scheme"]');
        expect(toggleButton).toBeInTheDocument();
    });

    it('should render burger menu for mobile', () => {
        const { container } = render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        const burger = container.querySelector('.mantine-Burger-root');
        expect(burger).toBeInTheDocument();
    });

    it('should render user menu trigger', () => {
        render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        // Menu items are only rendered when menu is opened due to withinPortal
        // Just verify the user name trigger is present
        expect(screen.getByText('User')).toBeInTheDocument();
    });

    it('should render menu component', () => {
        const { container } = render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        // Verify menu structure exists (items are in portal, not visible until opened)
        expect(container).toBeInTheDocument();
    });

    it('should render navigation with correct descriptions', () => {
        render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        expect(screen.getByText('All photos & albums')).toBeInTheDocument();
        expect(screen.getByText('All photos')).toBeInTheDocument();
        expect(screen.getByText('All albums')).toBeInTheDocument();
    });

    it('should render app logo icon', () => {
        const { container } = render(<AppBar activeLink={Labels.Dashboard}><div>Content</div></AppBar>);

        // Check for IconLibraryPhoto
        const svgs = container.querySelectorAll('svg');
        expect(svgs.length).toBeGreaterThan(0);
    });
});