import { useCallback, useEffect, useRef, useState } from 'react';

export interface UseBlobUrlResult {
    url: string | undefined;
    loading: boolean;
    error: Error | null;
    retry: () => void;
}

interface CacheEntry {
    url: string;
    refCount: number;
}

// Module-level cache shared across all hook instances.
// Keyed by the caller-provided identifier (e.g. photo ID).
const blobCache = new Map<string, CacheEntry>();

/**
 * Resolve a fetcher result into a blob URL string.
 * Handles both raw string URLs and Blob objects.
 */
function resolveBlobUrl(value: string | Blob): string {
    if (typeof value === 'string') return value;
    return URL.createObjectURL(value);
}

/**
 * Decrement the ref count for a cached entry and revoke the blob URL
 * when the last consumer unmounts.
 */
function release(key: string): void {
    const entry = blobCache.get(key);
    if (!entry) return;

    entry.refCount -= 1;

    if (entry.refCount <= 0) {
        if (entry.url.startsWith('blob:')) {
            URL.revokeObjectURL(entry.url);
        }
        blobCache.delete(key);
    }
}

/**
 * Centralized blob URL lifecycle hook.
 *
 * Creates, caches, and cleans up blob URLs with ref-counting so that
 * multiple components sharing the same key reuse a single URL and it
 * is only revoked when the last consumer unmounts.
 *
 * @param key     - Unique identifier for the blob (e.g. photo ID).
 * @param fetcher - Async function that returns a Blob or a string URL.
 * @returns `{ url, loading, error, retry }`
 */
export function useBlobUrl(
    key: string,
    fetcher: () => Promise<string | Blob>
): UseBlobUrlResult {
    const [url, setUrl] = useState<string | undefined>(() => {
        const entry = blobCache.get(key);
        return entry?.url;
    });
    const [loading, setLoading] = useState<boolean>(() => {
        return !blobCache.has(key);
    });
    const [error, setError] = useState<Error | null>(null);

    // Keep the latest fetcher in a ref so we don't re-trigger effects
    // when the fetcher identity changes (common with inline lambdas).
    const fetcherRef = useRef(fetcher);
    fetcherRef.current = fetcher;

    const execute = useCallback(async () => {
        // If we already have a cached entry, just use it.
        const cached = blobCache.get(key);
        if (cached) {
            setUrl(cached.url);
            setLoading(false);
            setError(null);
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const result = await fetcherRef.current();
            const blobUrl = resolveBlobUrl(result);

            blobCache.set(key, { url: blobUrl, refCount: 0 });
            setUrl(blobUrl);
            setLoading(false);
        } catch (err) {
            setError(err instanceof Error ? err : new Error(String(err)));
            setLoading(false);
        }
    }, [key]);

    const retry = useCallback(() => {
        // If there is a cached entry, revoke it first so we fetch fresh.
        const cached = blobCache.get(key);
        if (cached) {
            if (cached.url.startsWith('blob:')) {
                URL.revokeObjectURL(cached.url);
            }
            blobCache.delete(key);
        }
        setUrl(undefined);
        execute();
    }, [key, execute]);

    useEffect(() => {
        // Increment ref count (or initialize).
        const entry = blobCache.get(key);
        if (entry) {
            entry.refCount += 1;
        } else {
            blobCache.set(key, { url: '', refCount: 1 });
        }

        // Fetch if not already cached with a real URL.
        if (!entry || !entry.url) {
            execute();
        }

        // Cleanup: decrement ref count on unmount.
        return () => {
            release(key);
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [key]);

    return { url, loading, error, retry };
}

/**
 * Imperative helpers for non-React code (e.g. ThumbnailUtils) to
 * interact with the same cache.
 */

export function getCachedBlobUrl(key: string): string | undefined {
    return blobCache.get(key)?.url;
}

export function cacheBlobUrl(key: string, value: string | Blob): string {
    const existing = blobCache.get(key);
    if (existing) {
        existing.refCount += 1;
        return existing.url;
    }

    const blobUrl = resolveBlobUrl(value);
    blobCache.set(key, { url: blobUrl, refCount: 1 });
    return blobUrl;
}

export function releaseBlobUrl(key: string): void {
    release(key);
}

export function clearBlobCache(): void {
    blobCache.forEach((entry) => {
        if (entry.url.startsWith('blob:')) {
            URL.revokeObjectURL(entry.url);
        }
    });
    blobCache.clear();
}
