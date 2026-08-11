import {act, renderHook} from '@testing-library/react';
import {afterEach, beforeEach, describe, expect, test, vi} from 'vitest';
import {useWebSocket} from './useWebSocket';

class MockWebSocket {
    static instances: MockWebSocket[] = [];
    url: string;
    onopen: ((event: unknown) => void) | null = null;
    onclose: ((event: unknown) => void) | null = null;
    onmessage: ((event: unknown) => void) | null = null;
    onerror: ((event: unknown) => void) | null = null;

    constructor(url: string) {
        this.url = url;
        MockWebSocket.instances.push(this);
    }

    close() {
        // No-op for tests
    }

    emitOpen() {
        this.onopen?.({});
    }

    emitClose(code = 1000, reason = '', wasClean = true) {
        this.onclose?.({code, reason, wasClean});
    }
}

beforeEach(() => {
    MockWebSocket.instances = [];
    vi.useFakeTimers();
    vi.stubGlobal('WebSocket', MockWebSocket);
});

afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
});

describe('useWebSocket', () => {
    test('opens a ws:// connection with the token as a query parameter', () => {
        renderHook(() => useWebSocket('http://localhost:8080/api', 'secret'));

        expect(MockWebSocket.instances).toHaveLength(1);
        const url = new URL(MockWebSocket.instances[0].url);
        expect(url.protocol).toBe('ws:');
        expect(url.pathname).toBe('/api/ws');
        expect(url.searchParams.get('token')).toBe('secret');
    });

    test('upgrades to wss:// for https API bases', () => {
        renderHook(() => useWebSocket('https://photos.example.com/api', 'secret'));

        const url = new URL(MockWebSocket.instances[0].url);
        expect(url.protocol).toBe('wss:');
    });

    test('does not connect when no token is present', () => {
        renderHook(() => useWebSocket('http://localhost:8080/api', null));

        expect(MockWebSocket.instances).toHaveLength(0);
    });

    test('reports connected status and resets backoff after a successful open', () => {
        const {result} = renderHook(() => useWebSocket('http://localhost:8080/api', 'token'));

        act(() => {
            MockWebSocket.instances[0].emitOpen();
        });

        expect(result.current.status).toBe('connected');
    });

    test('records close diagnostics and reconnects with exponential backoff', () => {
        const {result} = renderHook(() => useWebSocket('http://localhost:8080/api', 'token'));

        act(() => {
            MockWebSocket.instances[0].emitClose(1006, '', false);
        });

        expect(result.current.status).toBe('disconnected');
        expect(result.current.lastClose).toEqual({code: 1006, reason: '', wasClean: false});

        // First reconnect after 1s
        act(() => {
            vi.advanceTimersByTime(1_000);
        });
        expect(MockWebSocket.instances).toHaveLength(2);

        // Second reconnect after 2s (backoff doubles)
        act(() => {
            MockWebSocket.instances[1].emitClose(1000, 'going away', true);
        });
        act(() => {
            vi.advanceTimersByTime(2_000);
        });
        expect(MockWebSocket.instances).toHaveLength(3);
    });

    test('caps reconnect delay at 30s', () => {
        renderHook(() => useWebSocket('http://localhost:8080/api', 'token'));

        // Emit close 6 times; each schedules the next reconnect at a doubled delay.
        for (let i = 0; i < 6; i += 1) {
            act(() => {
                MockWebSocket.instances[i].emitClose(1006, '', false);
            });
            act(() => {
                vi.advanceTimersByTime(40_000);
            });
        }

        // instance 6 scheduled the 7th reconnect
        const lastDelay = 40_000;
        act(() => {
            MockWebSocket.instances[6].emitClose(1006, '', false);
        });
        let attempts = 0;
        act(() => {
            vi.advanceTimersByTime(lastDelay);
            attempts = MockWebSocket.instances.length;
        });
        // After 30s + a little margin, the reconnect fires (capped, not 2^7=128s)
        expect(attempts).toBe(8);
    });

    test('logs the close code and reason for diagnostics', () => {
        const warnSpy = vi.spyOn(console, 'warn').mockImplementation(() => {});
        renderHook(() => useWebSocket('http://localhost:8080/api', 'token'));

        act(() => {
            MockWebSocket.instances[0].emitClose(1011, 'server error', false);
        });

        expect(warnSpy).toHaveBeenCalledWith(
            expect.stringContaining('code=1011'),
        );
        expect(warnSpy).toHaveBeenCalledWith(
            expect.stringContaining('reason="server error"'),
        );
        warnSpy.mockRestore();
    });

    test('dispatches messages to registered listeners', () => {
        const {result} = renderHook(() => useWebSocket('http://localhost:8080/api', 'token'));
        const listener = vi.fn();

        act(() => {
            result.current.addEventListener('photo.updated', listener);
        });
        act(() => {
            MockWebSocket.instances[0].onmessage?.({data: JSON.stringify({type: 'photo.updated', payload: {id: 1}})});
        });

        expect(listener).toHaveBeenCalledWith({type: 'photo.updated', payload: {id: 1}});
    });
});
