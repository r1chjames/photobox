import React, {useEffect, useRef, useState} from 'react';
import {Alert, Transition} from '@mantine/core';
import {IconWifiOff, IconServerOff} from '@tabler/icons-react';
import {useWebSocket} from '../../hooks/useWebSocket';

const DISMISSED_STORAGE_KEY = 'network-banner-dismissed';
const DISMISS_COOLDOWN_MS = 5 * 60 * 1000; // 5 minutes
const DISCONNECT_THRESHOLD_MS = 10_000; // 10s threshold before showing banner

type NetworkStatus = 'online' | 'offline' | 'backend-unreachable';

interface IProps {
    baseApiUrl: string;
}

export const NetworkStatusBanner: React.FunctionComponent<IProps> = ({baseApiUrl}) => {
    const [status, setStatus] = useState<NetworkStatus>('online');
    const [dismissed, setDismissed] = useState(false);
    const disconnectedSinceRef = useRef<number | null>(null);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const token = localStorage.getItem('token');
    const ws = useWebSocket(baseApiUrl, token);

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

    // Track WebSocket connection state
    useEffect(() => {
        if (ws.status === 'connected') {
            disconnectedSinceRef.current = null;
            setStatus('online');
            setDismissed(false);
        } else if (ws.status === 'disconnected') {
            if (disconnectedSinceRef.current === null) {
                disconnectedSinceRef.current = Date.now();
            }
        }
    }, [ws.status]);

    // Periodically check if we've been disconnected long enough to show the banner
    useEffect(() => {
        timerRef.current = setInterval(() => {
            if (disconnectedSinceRef.current !== null) {
                const elapsed = Date.now() - disconnectedSinceRef.current;
                if (elapsed >= DISCONNECT_THRESHOLD_MS && navigator.onLine) {
                    setStatus('backend-unreachable');
                    if (!isDismissedRecently()) {
                        setDismissed(false);
                    }
                }
            }
        }, 2000);

        return () => {
            if (timerRef.current) clearInterval(timerRef.current);
        };
    }, []);

    // Listen to browser online/offline events
    useEffect(() => {
        const handleOnline = () => {
            disconnectedSinceRef.current = null;
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

        return () => {
            window.removeEventListener('online', handleOnline);
            window.removeEventListener('offline', handleOffline);
        };
    }, []);

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
