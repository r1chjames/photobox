import { useCallback, useEffect, useRef, useState } from 'react';

type ConnectionStatus = 'connecting' | 'connected' | 'disconnected';

interface WSEvent {
    type: string;
    payload: unknown;
}

type EventListener = (event: WSEvent) => void;

const RECONNECT_BASE_MS = 1000;
const RECONNECT_MAX_MS = 30000;

/**
 * Builds a WebSocket URL from the API base URL.
 * Handles both absolute URLs (local dev: http://localhost:8080/api)
 * and relative paths (production: /api behind a reverse proxy).
 * Appends the auth token as a query parameter since the browser
 * WebSocket API does not support custom headers.
 */
function wsUrl(apiBaseUrl: string, token: string): string {
    const base = apiBaseUrl.startsWith('http')
        ? apiBaseUrl
        : `${window.location.origin}${apiBaseUrl}`;
    const url = new URL(base);
    url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
    url.pathname = url.pathname.replace(/\/?$/, '/ws');
    url.searchParams.set('token', token);
    return url.toString();
}

/**
 * Manages a single WebSocket connection with auto-reconnect.
 * Returns connection status and methods to subscribe to events.
 */
export function useWebSocket(apiBaseUrl: string, token: string | null) {
    const [status, setStatus] = useState<ConnectionStatus>('disconnected');
    const listenersRef = useRef<Map<string, Set<EventListener>>>(new Map());
    const wsRef = useRef<WebSocket | null>(null);
    const reconnectAttemptRef = useRef(0);
    const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const mountedRef = useRef(true);

    const addEventListener = useCallback((eventType: string, listener: EventListener) => {
        const map = listenersRef.current;
        if (!map.has(eventType)) {
            map.set(eventType, new Set());
        }
        map.get(eventType)!.add(listener);

        return () => {
            map.get(eventType)?.delete(listener);
        };
    }, []);

    const connect = useCallback(() => {
        if (!token || !mountedRef.current) return;

        // Close existing connection
        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }

        setStatus('connecting');

        const socket = new WebSocket(wsUrl(apiBaseUrl, token));

        socket.onopen = () => {
            if (!mountedRef.current) {
                socket.close();
                return;
            }
            setStatus('connected');
            reconnectAttemptRef.current = 0;
        };

        socket.onmessage = (event) => {
            try {
                const data: WSEvent = JSON.parse(event.data);
                const listeners = listenersRef.current.get(data.type);
                if (listeners) {
                    listeners.forEach((fn) => fn(data));
                }
            } catch {
                // Ignore malformed messages
            }
        };

        socket.onclose = () => {
            if (!mountedRef.current) return;
            setStatus('disconnected');
            wsRef.current = null;

            // Auto-reconnect with exponential backoff
            const delay = Math.min(
                RECONNECT_BASE_MS * 2 ** reconnectAttemptRef.current,
                RECONNECT_MAX_MS
            );
            reconnectAttemptRef.current += 1;
            reconnectTimerRef.current = setTimeout(connect, delay);
        };

        socket.onerror = () => {
            // onclose will fire after onerror, so reconnect is handled there
        };

        wsRef.current = socket;
    }, [apiBaseUrl, token]);

    // Connect on mount / token change
    useEffect(() => {
        mountedRef.current = true;
        if (token) {
            connect();
        }
        return () => {
            mountedRef.current = false;
            if (reconnectTimerRef.current) {
                clearTimeout(reconnectTimerRef.current);
            }
            if (wsRef.current) {
                wsRef.current.close();
                wsRef.current = null;
            }
        };
    }, [connect, token]);

    return { status, addEventListener };
}
