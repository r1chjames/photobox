import React from 'react';
import {Skeleton, Stack} from '@mantine/core';

interface GridSkeletonProps {
    columns?: number;
    rows?: number;
    gap?: number;
    aspectRatio?: string;
}

/**
 * A skeleton photo-grid used while the first batch loads. Renders a grid of
 * rounded card skeletons with no spinner — matches the grid layout so the
 * real photos replace it without jarring layout shift (issues #143/#144).
 */
export const GridSkeleton: React.FC<GridSkeletonProps> = ({
    columns = 5,
    rows = 3,
    gap = 10,
    aspectRatio = '4 / 3',
}) => {
    return (
        <Stack gap="sm" style={{ flex: 1, overflow: 'hidden', padding: 8 }}>
            {Array.from({length: rows}).map((_, r) => (
                <div
                    key={r}
                    style={{
                        display: 'grid',
                        gridTemplateColumns: `repeat(${columns}, 1fr)`,
                        gap,
                    }}
                >
                    {Array.from({length: columns}).map((_, c) => (
                        <div key={c} style={{ aspectRatio, borderRadius: 8, overflow: 'hidden' }}>
                            <Skeleton height="100%" width="100%" />
                        </div>
                    ))}
                </div>
            ))}
        </Stack>
    );
};
