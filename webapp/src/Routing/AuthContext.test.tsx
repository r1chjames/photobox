import React from 'react';
import {act, render, screen} from '../test/test-utils';
import {vi} from 'vitest';
import {AuthProvider, useAuth} from './AuthContext';

const TOKEN_STORAGE_KEY = 'token';
const TOKEN_EXPIRY_STORAGE_KEY = 'tokenExpiry';

const AuthProbe = () => {
    const {token, isAuthenticated, login, logout} = useAuth();
    return (
        <div>
            <span data-testid="authenticated">{String(isAuthenticated)}</span>
            <span data-testid="token">{token ?? 'none'}</span>
            <button data-testid="logout" onClick={logout}>logout</button>
            <button data-testid="login" onClick={() => login('test-token', new Date(Date.now() + 60_000).toISOString())}>login</button>
            <button data-testid="login-long" onClick={() => login('long-token', new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString())}>login long</button>
        </div>
    );
};

const renderProvider = () => render(
    <AuthProvider>
        <AuthProbe/>
    </AuthProvider>
);

describe('AuthContext', () => {
    beforeEach(() => {
        localStorage.clear();
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
        vi.restoreAllMocks();
    });

    it('is not authenticated when no token is stored', () => {
        renderProvider();
        expect(screen.getByTestId('authenticated').textContent).toBe('false');
    });

    it('is authenticated when a valid token is stored', () => {
        const expiry = new Date(Date.now() + 60_000).toISOString();
        localStorage.setItem(TOKEN_STORAGE_KEY, 'stored-token');
        localStorage.setItem(TOKEN_EXPIRY_STORAGE_KEY, String(Date.parse(expiry)));
        renderProvider();
        expect(screen.getByTestId('authenticated').textContent).toBe('true');
        expect(screen.getByTestId('token').textContent).toBe('stored-token');
    });

    it('discards an expired stored token on init', () => {
        localStorage.setItem(TOKEN_STORAGE_KEY, 'expired-token');
        localStorage.setItem(TOKEN_EXPIRY_STORAGE_KEY, String(Date.now() - 1_000));
        renderProvider();
        expect(screen.getByTestId('authenticated').textContent).toBe('false');
        expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBeNull();
        expect(localStorage.getItem(TOKEN_EXPIRY_STORAGE_KEY)).toBeNull();
    });

    it('logs out proactively when the token expires while the app is open', () => {
        renderProvider();
        act(() => {
            screen.getByTestId('login').click();
        });
        expect(screen.getByTestId('authenticated').textContent).toBe('true');
        act(() => {
            vi.advanceTimersByTime(61_000);
        });
        expect(screen.getByTestId('authenticated').textContent).toBe('false');
        expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBeNull();
    });

    it('stays authenticated for long-lived tokens whose delay exceeds the max setTimeout range', () => {
        renderProvider();
        act(() => {
            screen.getByTestId('login-long').click();
        });
        expect(screen.getByTestId('authenticated').textContent).toBe('true');
        // A 30-day token used to schedule setTimeout with a delay beyond the
        // 32-bit signed range, which browsers clamp to ~1ms — logging the
        // user out the instant they logged in. No timers must fire early.
        act(() => {
            vi.advanceTimersByTime(1000);
        });
        expect(screen.getByTestId('authenticated').textContent).toBe('true');
        expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBe('long-token');
    });

    it('logout clears token and expiry', () => {
        localStorage.setItem(TOKEN_STORAGE_KEY, 'stored-token');
        localStorage.setItem(TOKEN_EXPIRY_STORAGE_KEY, String(Date.now() + 60_000));
        renderProvider();
        act(() => {
            screen.getByTestId('logout').click();
        });
        expect(screen.getByTestId('authenticated').textContent).toBe('false');
        expect(localStorage.getItem(TOKEN_STORAGE_KEY)).toBeNull();
        expect(localStorage.getItem(TOKEN_EXPIRY_STORAGE_KEY)).toBeNull();
    });

    it('treats a token without expiry as authenticated (backwards compatible)', () => {
        localStorage.setItem(TOKEN_STORAGE_KEY, 'legacy-token');
        renderProvider();
        expect(screen.getByTestId('authenticated').textContent).toBe('true');
    });
});
