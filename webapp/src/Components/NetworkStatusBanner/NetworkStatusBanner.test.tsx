import React from 'react';
import {act, render, screen} from '../../test/test-utils';
import {vi} from 'vitest';
import {NetworkStatusBanner} from './NetworkStatusBanner';

const {wsState} = vi.hoisted(() => ({
    wsState: {
        status: 'disconnected',
    },
}));

vi.mock('../../hooks/useWebSocket', () => ({
    useWebSocket: () => ({
        status: wsState.status,
        lastClose: null,
        addEventListener: () => () => undefined,
    }),
}));

const BASE_API_URL = 'http://localhost:8080/api';

let mockFetch: ReturnType<typeof vi.fn>;

describe('NetworkStatusBanner', () => {
    beforeEach(() => {
        wsState.status = 'disconnected';
        localStorage.clear();
        mockFetch = vi.fn().mockResolvedValue({ok: true});
        vi.stubGlobal('fetch', mockFetch);
        vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(true);
        vi.useFakeTimers();
    });

    afterEach(() => {
        vi.useRealTimers();
        vi.unstubAllGlobals();
        vi.restoreAllMocks();
    });

    const renderBanner = () => render(<NetworkStatusBanner baseApiUrl={BASE_API_URL}/>);

    const advance = async (ms: number) => {
        await act(async () => {
            vi.advanceTimersByTime(ms);
            // Flush microtasks from fetch continuations fired during the advance
            await Promise.resolve();
            await Promise.resolve();
        });
    };

    const flushHealthCheck = async () => {
        await act(async () => {
            await Promise.resolve();
        });
    };

    test('renders nothing when the WebSocket is connected', () => {
        wsState.status = 'connected';
        renderBanner();

        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();
        expect(screen.queryByText('No Internet Connection')).not.toBeInTheDocument();
    });

    test('shows "Backend Unavailable" only when the WS is down AND the HTTP health check is failing past the threshold', async () => {
        mockFetch.mockResolvedValue({ok: false});
        renderBanner();
        await flushHealthCheck();

        // Before the threshold, no banner even with a failing health check
        await advance(20_000);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();

        await advance(11_000);
        expect(screen.getByText('Backend Unavailable')).toBeInTheDocument();
    });

    test('does NOT show the banner when the WS is down but the REST API health check succeeds', async () => {
        mockFetch.mockResolvedValue({ok: true});
        renderBanner();
        await flushHealthCheck();

        await advance(60_000);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();
    });

    test('does NOT show the banner while the health check result is still unknown', async () => {
        // fetch never resolves -> healthOk stays null
        mockFetch.mockReturnValue(new Promise(() => {}));
        renderBanner();
        await flushHealthCheck();

        await advance(60_000);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();
    });

    test('shows the offline banner when the browser reports no network', async () => {
        vi.spyOn(navigator, 'onLine', 'get').mockReturnValue(false);
        renderBanner();
        await flushHealthCheck();

        act(() => {
            window.dispatchEvent(new Event('offline'));
        });

        expect(screen.getByText('No Internet Connection')).toBeInTheDocument();
    });

    test('dismiss hides the banner and the cooldown suppresses it after a flapping reconnect', async () => {
        mockFetch.mockResolvedValue({ok: false});
        const {rerender} = renderBanner();
        await flushHealthCheck();
        await advance(30_000);
        expect(screen.getByText('Backend Unavailable')).toBeInTheDocument();

        act(() => {
            screen.getByRole('button', {name: 'Dismiss'}).click();
        });
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();

        // Flap: reconnect, then drop again quickly. Banner must stay dismissed (cooldown).
        wsState.status = 'connected';
        rerender(<NetworkStatusBanner baseApiUrl={BASE_API_URL}/>);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();

        wsState.status = 'disconnected';
        rerender(<NetworkStatusBanner baseApiUrl={BASE_API_URL}/>);
        await flushHealthCheck();
        await advance(60_000);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();
    });

    test('banner disappears once the WebSocket reconnects', async () => {
        mockFetch.mockResolvedValue({ok: false});
        const {rerender} = renderBanner();
        await flushHealthCheck();
        await advance(30_000);
        expect(screen.getByText('Backend Unavailable')).toBeInTheDocument();

        wsState.status = 'connected';
        rerender(<NetworkStatusBanner baseApiUrl={BASE_API_URL}/>);
        expect(screen.queryByText('Backend Unavailable')).not.toBeInTheDocument();
    });
});
