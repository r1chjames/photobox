import React, {useEffect, useRef, useState} from 'react';
import {Alert, Transition} from '@mantine/core';
import {IconWifiOff, IconServerOff} from '@tabler/icons-react';

const DISMISSED_STORAGE_KEY = 'network-banner-dismissed';
const DISMISS_COOLDOWN_MS = 5 * 60 * 1000; // 5 minutes
const HEALTH_CHECK_INTERVAL_MS = 30_000; // 30 seconds
const CONSECUTIVE_FAILURES_THRESHOLD = 2;

type NetworkStatus = 'online' | 'offline' | 'backend-unreachable';

interface IProps {
    baseApiUrl: string;
}

export const NetworkStatusBanner: React.FunctionComponent<IProps> = ({baseApiUrl}) => {
    const [status, setStatus] = useState<NetworkStatus>('online');
    const [dismissed, setDismissed] = useState(false);
    const consecutiveFailuresRef = useRef(0);
    const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
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
            const timeoutId = setTimeout(() => controller.abort(), 5000);

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

        // Reset failure counter on each effect run (handles StrictMode double-mount)
        consecutiveFailuresRef.current = 0;

        // Delay the first health check to let service worker and network settle after hard refresh
        const initialTimeoutId = setTimeout(() => {
            if (navigator.onLine) {
                checkBackendHealth(true);
            }
            hasCheckedRef.current = true;
        }, 3000);

        // Set up periodic health checks
        intervalRef.current = setInterval(() => {
            if (navigator.onLine) {
                checkBackendHealth();
            }
        }, HEALTH_CHECK_INTERVAL_MS);

        return () => {
            window.removeEventListener('online', handleOnline);
            window.removeEventListener('offline', handleOffline);
            clearTimeout(initialTimeoutId);
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
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
