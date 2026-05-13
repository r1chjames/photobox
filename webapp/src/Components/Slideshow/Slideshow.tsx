import React, { useCallback, useEffect, useState } from 'react';
import { ActionIcon, Overlay, Transition } from '@mantine/core';
import { useInterval } from '@mantine/hooks';
import { IconPlayerPause, IconPlayerPlay, IconX } from '@tabler/icons-react';
import { Photo } from '../../Models/Photo';
import { fetchPhotoBinWithAuth, revokeBlobUrl } from '../../utils/ImageUtils';
import { IPhotosAdapter } from '../../Adapters/IPhotosAdapter';

interface SlideshowProps {
    photos: Photo[];
    startIndex: number;
    photosAdapter: IPhotosAdapter;
    onClose: () => void;
}

export const Slideshow: React.FC<SlideshowProps> = ({ photos, startIndex, photosAdapter, onClose }) => {
    const [currentIndex, setCurrentIndex] = useState(startIndex);
    const [imageUrl, setImageUrl] = useState<string | undefined>();
    const [isPlaying, setIsPlaying] = useState(true);
    const [showControls, setShowControls] = useState(true);
    const interval = useInterval(() => {
        setCurrentIndex(prev => (prev + 1) % photos.length);
    }, 5000);

    const loadImage = useCallback(async () => {
        if (!photos[currentIndex]) return;
        const url = await fetchPhotoBinWithAuth(photosAdapter, photos[currentIndex].id);
        setImageUrl(prev => {
            if (prev && prev.startsWith('blob:')) revokeBlobUrl(prev);
            return url;
        });
    }, [photosAdapter, photos, currentIndex]);

    useEffect(() => {
        void loadImage();
    }, [loadImage]);

    useEffect(() => {
        if (isPlaying) {
            interval.start();
        } else {
            interval.stop();
        }
        return () => interval.stop();
    }, [isPlaying, interval]);

    useEffect(() => {
        const timer = setTimeout(() => setShowControls(false), 3000);
        return () => clearTimeout(timer);
    }, [showControls]);

    const handleMouseMove = () => {
        setShowControls(true);
    };

    const photo = photos[currentIndex];

    return (
        <div
            style={{
                position: 'fixed',
                top: 0,
                left: 0,
                width: '100vw',
                height: '100vh',
                background: 'black',
                zIndex: 9999,
                cursor: showControls ? 'default' : 'none',
            }}
            onMouseMove={handleMouseMove}
            onClick={() => setShowControls(prev => !prev)}
        >
            {imageUrl && photo?.mediaType === 'video' ? (
                <video
                    key={imageUrl}
                    src={imageUrl}
                    controls={false}
                    autoPlay
                    muted
                    style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'contain',
                    }}
                />
            ) : imageUrl && (
                <img
                    key={imageUrl}
                    src={imageUrl}
                    alt={photo?.name}
                    style={{
                        width: '100%',
                        height: '100%',
                        objectFit: 'contain',
                        transition: 'transform 10s ease',
                        transform: isPlaying ? 'scale(1.05)' : 'scale(1)',
                    }}
                />
            )}
            <Transition mounted={showControls} transition="fade" duration={200}>
                {(styles) => (
                    <Overlay style={styles} backgroundOpacity={0}>
                        <div style={{ position: 'absolute', top: 16, right: 16, display: 'flex', gap: 8 }}>
                            <ActionIcon color="white" variant="transparent" onClick={(e) => { e.stopPropagation(); setIsPlaying(prev => !prev); }}>
                                {isPlaying ? <IconPlayerPause size="1.5rem" /> : <IconPlayerPlay size="1.5rem" />}
                            </ActionIcon>
                            <ActionIcon color="white" variant="transparent" onClick={(e) => { e.stopPropagation(); onClose(); }}>
                                <IconX size="1.5rem" />
                            </ActionIcon>
                        </div>
                        <div style={{ position: 'absolute', bottom: 24, left: 24, color: 'white' }}>
                            <div style={{ fontSize: 14, opacity: 0.8 }}>
                                {currentIndex + 1} / {photos.length}
                            </div>
                            <div style={{ fontSize: 18, fontWeight: 500 }}>
                                {photo?.name}
                            </div>
                        </div>
                    </Overlay>
                )}
            </Transition>
        </div>
    );
};
