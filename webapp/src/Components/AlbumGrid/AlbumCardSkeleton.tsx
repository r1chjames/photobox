import React from 'react';
import {Skeleton} from '@mantine/core';

/**
 * Loading placeholder shown while the album list is being fetched for the
 * first time. Mirrors the real AlbumCard DOM/classes (album-card >
 * album-media-grid with one main thumb + two sub-thumbs, plus the overlay
 * title bar) so fetched albums replace it without layout shift. The Mantine
 * skeleton shimmer stands in for photo content — no spinner (issue #173).
 */
export const AlbumCardSkeleton: React.FC = () => {
    return (
        <div className="album-card" aria-busy="true" data-testid="album-card-skeleton">
            <div className="album-media-grid">
                <div className="album-main-thumb">
                    <Skeleton height={184} width="100%" radius="md" />
                </div>
                <div className="album-sub-thumbs">
                    <div className="album-sub-thumb">
                        <Skeleton height={90} width="100%" radius="md" />
                    </div>
                    <div className="album-sub-thumb">
                        <Skeleton height={90} width="100%" radius="md" />
                    </div>
                </div>
            </div>
            <div className="album-info-overlay">
                <Skeleton height={16} width="55%" radius="sm" />
                <Skeleton height={22} width={72} radius="xl" />
            </div>
        </div>
    );
};
