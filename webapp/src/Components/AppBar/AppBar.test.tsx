import { describe, it, expect, beforeEach, vi } from 'vitest';
import { render, screen } from '../../test/test-utils';
import userEvent from '@testing-library/user-event';
import { AppBar } from './AppBar';

describe('AppBar', () => {
    beforeEach(() => {
        localStorage.clear();
        vi.clearAllMocks();
    });

    it('should render app name "Photobox"', () => {
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.getByText('Photobox')).toBeInTheDocument();
    });

    it('should render primary navigation links', () => {
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.getByRole('button', { name: 'Home' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Photos' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Albums' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Memories' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Settings' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'More' })).toBeInTheDocument();
    });

    it('should render secondary navigation items inside the More submenu', async () => {
        const user = userEvent.setup();
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.queryByRole('button', { name: 'Trash' })).not.toBeInTheDocument();
        await user.click(screen.getByRole('button', { name: 'More' }));
        expect(screen.getByRole('button', { name: 'Trash' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Quality' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Shares' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Map' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Tags' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Duplicates' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Videos' })).toBeInTheDocument();
    });

    it('should highlight active navigation link for the root path', () => {
        render(<AppBar><div>content</div></AppBar>);
        const home = screen.getByRole('button', { name: 'Home' });
        expect(home.className).toContain('active');
        expect(screen.getByRole('button', { name: 'Albums' }).className).not.toContain('active');
    });

    it('should render user pill with stored username', () => {
        localStorage.setItem('pb-username', 'rich');
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.getByText('rich')).toBeInTheDocument();
    });

    it('should render children content', () => {
        render(<AppBar><div>page body</div></AppBar>);
        expect(screen.getByText('page body')).toBeInTheDocument();
    });

    it('should render color scheme toggle button', () => {
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.getByRole('button', { name: 'Toggle color scheme' })).toBeInTheDocument();
    });

    it('should render keyboard shortcuts button', () => {
        render(<AppBar><div>content</div></AppBar>);
        expect(screen.getByRole('button', { name: 'Keyboard shortcuts' })).toBeInTheDocument();
    });

    it('should open the user menu with logout and Users entry', async () => {
        const user = userEvent.setup();
        render(<AppBar><div>content</div></AppBar>);
        await user.click(screen.getByRole('button', { name: /^US User$/ }));
        expect(screen.getByText('Logout')).toBeInTheDocument();
        expect(screen.getByText('Account settings')).toBeInTheDocument();
        expect(screen.getByText('Users')).toBeInTheDocument();
    });
});
