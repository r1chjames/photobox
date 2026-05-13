import { useEffect, useRef } from 'react';
import { notifications } from '@mantine/notifications';

const STORAGE_KEY = 'photobox-seen-shortcuts-hint';

export function useShortcutsHint() {
    const hasShownRef = useRef(false);

    useEffect(() => {
        const seen = localStorage.getItem(STORAGE_KEY);
        if (seen || hasShownRef.current) return;

        const timer = setTimeout(() => {
            notifications.show({
                title: 'Tip',
                message: 'Press ? anytime to see keyboard shortcuts',
                color: 'blue',
                withCloseButton: true,
                autoClose: 8000,
            });
            localStorage.setItem(STORAGE_KEY, 'true');
            hasShownRef.current = true;
        }, 3000);

        return () => clearTimeout(timer);
    }, []);
}
