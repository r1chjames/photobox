import React, {createContext, useCallback, useContext, useEffect, useMemo, useRef, useState} from 'react';

const TOKEN_STORAGE_KEY = 'token';
const TOKEN_EXPIRY_STORAGE_KEY = 'tokenExpiry';
const MAX_TIMEOUT_MS = 2 ** 31 - 1; // largest setTimeout delay browsers honor (~24.8 days)

interface AuthContextType {
    token: string | null;
    login: (token: string, expiresAt?: string) => void;
    logout: () => void;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

/**
 * Parses an ISO-8601 expires_at value (or the epoch-ms numeric string we
 * persist ourselves) into epoch milliseconds.
 * Returns null when the value is missing or unparseable.
 */
function parseExpiry(expiresAt?: string | null): number | null {
    if (!expiresAt) return null;
    if (/^\d+$/.test(expiresAt)) {
        const ms = Number(expiresAt);
        return Number.isNaN(ms) ? null : ms;
    }
    const ms = Date.parse(expiresAt);
    return Number.isNaN(ms) ? null : ms;
}

/**
 * Reads a stored token from localStorage, discarding it when its recorded
 * expiry is already in the past (a token whose lifetime is over must never
 * leave the user "authenticated" inside an empty app).
 */
function readStoredToken(): string | null {
    const stored = localStorage.getItem(TOKEN_STORAGE_KEY);
    if (!stored) return null;
    const expiryMs = parseExpiry(localStorage.getItem(TOKEN_EXPIRY_STORAGE_KEY));
    if (expiryMs !== null && expiryMs <= Date.now()) {
        localStorage.removeItem(TOKEN_STORAGE_KEY);
        localStorage.removeItem(TOKEN_EXPIRY_STORAGE_KEY);
        return null;
    }
    return stored;
}

export const AuthProvider: React.FunctionComponent<{ children: React.ReactNode }> = ({children}) => {
    const [token, setToken] = useState<string | null>(readStoredToken);
    const expiryMsRef = useRef<number | null>(parseExpiry(localStorage.getItem(TOKEN_EXPIRY_STORAGE_KEY)));
    const logoutTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    const login = useCallback((newToken: string, expiresAt?: string) => {
        const expiryMs = parseExpiry(expiresAt);
        localStorage.setItem(TOKEN_STORAGE_KEY, newToken);
        if (expiryMs !== null) {
            localStorage.setItem(TOKEN_EXPIRY_STORAGE_KEY, String(expiryMs));
        } else {
            localStorage.removeItem(TOKEN_EXPIRY_STORAGE_KEY);
        }
        expiryMsRef.current = expiryMs;
        setToken(newToken);
    }, []);

    const logout = useCallback(() => {
        localStorage.removeItem(TOKEN_STORAGE_KEY);
        localStorage.removeItem(TOKEN_EXPIRY_STORAGE_KEY);
        expiryMsRef.current = null;
        setToken(null);
    }, []);

    // Proactively log out the instant the token expires, so the route guard
    // redirects to /login even when the user is idle or navigating cached
    // views with no live API traffic (a 401 would only arrive reactively).
    useEffect(() => {
        if (logoutTimerRef.current) {
            clearTimeout(logoutTimerRef.current);
            logoutTimerRef.current = null;
        }
        if (!token) return;
        const expiryMs = expiryMsRef.current;
        if (expiryMs === null) return;

        const armLogoutTimer = () => {
            const remaining = expiryMs - Date.now();
            if (remaining <= 0) {
                logout();
                return;
            }
            // Browsers fire setTimeout immediately when the delay overflows
            // the 32-bit signed range (~24.8 days). A long-lived token (e.g.
            // 30 days) must therefore be armed in sub-max chunks that re-arm
            // until the token genuinely expires.
            logoutTimerRef.current = setTimeout(armLogoutTimer, Math.min(remaining, MAX_TIMEOUT_MS));
        };
        armLogoutTimer();

        return () => {
            if (logoutTimerRef.current) {
                clearTimeout(logoutTimerRef.current);
                logoutTimerRef.current = null;
            }
        };
    }, [token, logout]);

    const value = useMemo(() => ({
        token,
        login,
        logout,
        isAuthenticated: token !== null,
    }), [token, login, logout]);

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
