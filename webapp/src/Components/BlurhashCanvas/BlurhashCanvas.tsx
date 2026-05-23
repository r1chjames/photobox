import React, { useEffect, useRef } from 'react';
import { decode } from 'blurhash';

interface BlurhashCanvasProps {
    hash: string;
    width?: number;
    height?: number;
    punch?: number;
    style?: React.CSSProperties;
    className?: string;
}

export const BlurhashCanvas: React.FC<BlurhashCanvasProps> = ({
    hash,
    width = 32,
    height = 24,
    punch = 1,
    style,
    className,
}) => {
    const canvasRef = useRef<HTMLCanvasElement>(null);

    useEffect(() => {
        if (!canvasRef.current || !hash) return;
        try {
            const pixels = decode(hash, width, height, punch);
            const canvas = canvasRef.current;
            const ctx = canvas.getContext('2d');
            if (!ctx) return;
            const imageData = ctx.createImageData(width, height);
            imageData.data.set(pixels);
            ctx.putImageData(imageData, 0, 0);
        } catch {
            // Invalid blurhash — silently fail
        }
    }, [hash, width, height, punch]);

    return (
        <canvas
            ref={canvasRef}
            width={width}
            height={height}
            className={className}
            style={{
                width: '100%',
                height: '100%',
                objectFit: 'cover',
                ...style,
            }}
        />
    );
};
