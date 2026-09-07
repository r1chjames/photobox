import React, { useEffect, useMemo, useRef } from 'react';
import { decode, isBlurhashValid } from 'blurhash';
import { Skeleton } from '@mantine/core';

// The blurhash base83 alphabet (mirrors the package's char set; kept local
// because isBlurhashValid only checks structure/length, not membership —
// a length-valid hash containing foreign characters must still fall through).
const BASE83 = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz#$%*+,-.:;=?@[]^_{|}~';

// Decode the blurhash to a 32x24 RGBA pixel buffer. Memoized on the hash so
// parent re-renders (opacity crossfade, selection state) never re-decode the
// same hash, and invalid hashes fall through to dominantColor/skeleton.
const decodeBlurhashPixels = (blurhash: string | null | undefined): Uint8ClampedArray | null => {
    if (!blurhash) return null;
    try {
        for (let i = 0; i < blurhash.length; i++) {
            if (BASE83.indexOf(blurhash[i]) === -1) return null;
        }
        if (!isBlurhashValid(blurhash).result) return null;
        return decode(blurhash, 32, 24, 1);
    } catch {
        // Invalid blurhash — never leave a blank tile; callers fall through.
        return null;
    }
};

export interface PhotoMediaPlaceholderProps {
    blurhash?: string | null;
    dominantColor?: string | null;
    className?: string;
    style?: React.CSSProperties;
}

/**
 * Single placeholder layer shown beneath a photo until its thumbnail loads.
 * Precedence: decoded blurhash → dominantColor → skeleton, so an invalid or
 * missing blurhash always falls back to something visible (never a blank
 * frame, issue #173). The decoded canvas renders absolutely-fillable at
 * 100% width/height — the sibling <img> crossfades over it via opacity
 * without any layout shift.
 */
export const PhotoMediaPlaceholder: React.FC<PhotoMediaPlaceholderProps> = ({
    blurhash,
    dominantColor,
    className,
    style,
}) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);

    const pixels = useMemo(
        () => decodeBlurhashPixels(blurhash),
        [blurhash]
    );

    useEffect(() => {
        // Drawing must never crash the tile — environments without a 2D
        // context (tests, odd embeds) keep the sized canvas as-is.
        try {
            const canvas = canvasRef.current;
            if (!canvas || !pixels) return;
            const ctx = canvas.getContext('2d');
            if (!ctx) return;
            const imageData = ctx.createImageData(32, 24);
            imageData.data.set(pixels);
            ctx.putImageData(imageData, 0, 0);
        } catch {
            // Ignore: the sized canvas remains a non-blank placeholder.
        }
    }, [pixels]);

    if (pixels) {
        return (
            <canvas
                ref={canvasRef}
                className={className}
                data-testid="photo-media-placeholder"
                style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: '100%',
                    objectFit: 'cover',
                    ...style,
                }}
                width={32}
                height={24}
            />
        );
    }

    if (dominantColor) {
        return (
            <div
                className={className}
                data-testid="photo-media-placeholder"
                style={{
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: '100%',
                    backgroundColor: dominantColor,
                    ...style,
                }}
            />
        );
    }

    return (
        <div
            className={className}
            data-testid="photo-media-placeholder"
            style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: '100%',
                ...style,
            }}
        >
            <Skeleton height="100%" width="100%" />
        </div>
    );
};
