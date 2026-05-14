import React, {useEffect, useRef, useState} from 'react';
import {Alert, Transition} from '@mantine/core';
import {IconWifiOff, IconServerOff} from '@tabler/icons-react';

const DISMISSED_STORAGE_KEY = 'network-banner-dismissed';
const DISMISS_COOLDOWN_MS = 5 * 60 * 1000; // 5 minutes
const HEALTH_CHECK_INTERVAL_MS = 30_000; // 30 seconds
const CONSECUTIVE_FAILURES_THRESHOLD = 3; // requires 3 consecutive failures (was 2)
const HEALTH_CHECK_TIMEOUT_MS = 10_000; // 10s timeout (was 5s) to handle throttled background tabs

type NetworkStatus = 'online' | 'offline' | 'backend-unreachable';

interface IProps {
    baseApiUrl: string;
}

export const NetworkStatusBanner: React.FunctionComponent<IProps> = ({baseApiUrl}) => {
    const [status, setStatus] = useState<NetworkStatus>('online');
    const [dismissed, setDismissed] = useState(false);
    const consecutiveFailuresRef = useRef(0);
    const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
    const mountedRef = useRef(true);
    const hasCheckedRef = useRef(false);

    const isDismissedRecently = (): boolean => {
        const dismissedAt = localStorage.getItem(DISMISSED_STORAGE_KEY);
        if (!dismissedAt) return false;
        const elapsed = Date.now() - parseInt(dismissedAt, 10);
        return elapsed < DISMISS_COOLDOWN_MS;
    };

    const dismiss = () => {
        localStorage.setItem(DISMISSED_STORAGE_KEY, Date.now().toString());
        setDismissed(true);
    };

    const checkBackendHealth = async (isInitialCheck = false) => {
        try {
            // Use a lightweight HEAD request to the base API path
            const controller = new AbortController();
            const timeoutId = setTimeout(() => controller.abort(), HEALTH_CHECK_TIMEOUT_MS);

            const response = await fetch(baseApiUrl, {
                method: 'HEAD',
                signal: controller.signal,
            });

            clearTimeout(timeoutId);

            if (response.ok || response.status === 401) {
                // 401 is fine — it means the backend is reachable, just unauthenticated
                consecutiveFailuresRef.current = 0;
                if (status !== 'online') {
                    setStatus('online');
                    setDismissed(false);
                }
            } else {
                if (!isInitialCheck) {
                    consecutiveFailuresRef.current += 1;
                }
                if (consecutiveFailuresRef.current >= CONSECUTIVE_FAILURES_THRESHOLD) {
                    setStatus('backend-unreachable');
                    if (!isDismissedRecently()) {
                        setDismissed(false);
                    }
                }
            }
        } catch {
            if (!isInitialCheck) {
                consecutiveFailuresRef.current += 1;
            }
            if (consecutiveFailuresRef.current >= CONSECUTIVE_FAILURES_THRESHOLD) {
                setStatus('backend-unreachable');
                if (!isDismissedRecently()) {
                    setDismissed(false);
                }
            }
        }
    };

    useEffect(() => {
        // Listen to browser online/offline events
        const handleOnline = () => {
            consecutiveFailuresRef.current = 0;
            setStatus('online');
            setDismissed(false);
        };

        const handleOffline = () => {
            setStatus('offline');
            if (!isDismissedRecently()) {
                setDismissed(false);
            }
        };

        window.addEventListener('online', handleOnline);
        window.addEventListener('offline', handleOffline);

        // When tab becomes visible (user returns), reset failures and check health immediately.
        // Browsers throttle setInterval in background tabs, so visible check ensures
        // any false failures from throttled checks are cleared when the user comes back.
        const handleVisibility = () => {
            if (document.visibilityState === 'visible' && navigator.onLine) {
                consecutiveFailuresRef.current = 0;
                checkBackendHealth();
            }
        };
        document.addEventListener('visibilitychange', handleVisibility);

        // Reset failure counter on each effect run (handles StrictMode double-mount)
        consecutiveFailuresRef.current = 0;
        mountedRef.current = true;

        // Schedule the next periodic check. Using setTimeout chaining instead of
        // setInterval avoids aggressive browser throttling in background/idle tabs.
        const scheduleNext = (delayMs: number) => {
            if (!mountedRef.current) return;
            timeoutRef.current = setTimeout(async () => {
                if (!mountedRef.current) return;
                if (navigator.onLine) {
                    await checkBackendHealth();
                }
                scheduleNext(HEALTH_CHECK_INTERVAL_MS);
            }, delayMs);
        };

        // Delay initial health check to let service worker and network settle
        timeoutRef.current = setTimeout(async () => {
            if (!mountedRef.current) return;
            if (navigator.onLine) {
                await checkBackendHealth(true);
            }
            hasCheckedRef.current = true;
            scheduleNext(HEALTH_CHECK_INTERVAL_MS);
        }, 3000);

        return () => {
            mountedRef.current = false;
            window.removeEventListener('online', handleOnline);
            window.removeEventListener('offline', handleOffline);
            document.removeEventListener('visibilitychange', handleVisibility);
            if (timeoutRef.current) {
                clearTimeout(timeoutRef.current);
            }
        };
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [baseApiUrl]);

    if (status === 'online' || dismissed) {
        return null;
    }

    const isOffline = status === 'offline';

    return (
        <Transition mounted={true} transition="slide-down" duration={300}>
            {(styles) => (
                <div style={{position: 'sticky', top: 0, zIndex: 1001, ...styles}}>
                    <Alert
                        color={isOffline ? 'red' : 'yellow'}
                        variant="filled"
                        title={isOffline ? 'No Internet Connection' : 'Backend Unavailable'}
                        icon={isOffline ? <IconWifiOff size={18}/> : <IconServerOff size={18}/>}
                        onClose={dismiss}
                        styles={{
                            root: {
                                borderRadius: 0,
                                borderLeft: 'none',
                                borderRight: 'none',
                                borderTop: 'none',
                            },
                            message: {
                                fontWeight: 500,
                            },
                        }}
                    >
                        {isOffline
                            ? 'You are currently offline. Please check your network connection.'
                            : 'Unable to reach the server. Some features may not work properly.'}
                    </Alert>
                </div>
            )}
        </Transition>
    );
};
