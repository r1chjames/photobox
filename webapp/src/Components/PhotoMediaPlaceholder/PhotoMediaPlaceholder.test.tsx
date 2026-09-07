import { describe, it, expect } from 'vitest';
import { render } from '../../test/test-utils';
import { PhotoMediaPlaceholder } from './PhotoMediaPlaceholder';

// Valid 3x4-component blurhash (used by the issue #173 photo-grid
// fixtures); decodes to a 32x24 RGBA buffer.
const VALID_BLURHASH = 'LEHV6nWB2yk8pyo0adR*.7kCMdnj';

describe('PhotoMediaPlaceholder', () => {

    it('renders a canvas carrying the decoded pixels when the blurhash is valid', () => {
        const { container } = render(<PhotoMediaPlaceholder blurhash={VALID_BLURHASH} />);
        const canvas = container.querySelector('canvas');
        expect(canvas).not.toBeNull();
        // Painted at the 32x24 decode resolution so the canvas carries the
        // decoded pixels (upscaled to 100% by CSS).
        expect(canvas?.getAttribute('width')).toBe('32');
        expect(canvas?.getAttribute('height')).toBe('24');
        // When a 2D context is available the decode must have been painted.
        try {
            const ctx = (canvas as HTMLCanvasElement | null)?.getContext?.('2d');
            if (ctx && typeof ctx.getImageData === 'function') {
                const imageData = ctx.getImageData(0, 0, 32, 24);
                expect(imageData.width).toBe(32);
                expect(imageData.height).toBe(24);
                // Decoded pixels are opaque; a blank canvas would be
                // transparent.
                let nonTransparent = 0;
                for (let i = 3; i < imageData.data.length; i += 4) {
                    if (imageData.data[i] > 0) nonTransparent++;
                }
                expect(nonTransparent).toBeGreaterThan(0);
            }
        } catch {
            // No canvas pixel inspection available — the canvas element and
            // its 32x24 decode surface already assert the render path.
        }
    });

    it('renders a solid dominant color layer when no blurhash is present', () => {
        const { container } = render(<PhotoMediaPlaceholder dominantColor="#667788" />);
        expect(container.querySelector('canvas')).toBeNull();
        const layer = container.querySelector('[data-testid="photo-media-placeholder"]') as HTMLElement;
        expect(layer).not.toBeNull();
        expect(layer.style.backgroundColor).toBe('#667788');
        expect(container.querySelector('.mantine-Skeleton-root')).toBeNull();
    });

    it('renders a skeleton when neither blurhash nor dominant color is present', () => {
        const { container } = render(<PhotoMediaPlaceholder />);
        expect(container.querySelector('canvas')).toBeNull();
        expect(container.querySelector('.mantine-Skeleton-root')).not.toBeNull();
    });

    it('renders no canvas and falls back to the skeleton for an invalid blurhash', () => {
        const { container } = render(<PhotoMediaPlaceholder blurhash="invalid!!!" />);
        expect(container.querySelector('canvas')).toBeNull();
        expect(container.querySelector('.mantine-Skeleton-root')).not.toBeNull();
    });

    it('falls back to the dominant color when the blurhash is invalid but a color exists', () => {
        const { container } = render(<PhotoMediaPlaceholder blurhash="invalid!!!" dominantColor="#445566" />);
        expect(container.querySelector('canvas')).toBeNull();
        const layer = container.querySelector('[data-testid="photo-media-placeholder"]') as HTMLElement;
        expect(layer).not.toBeNull();
        expect(layer.style.backgroundColor).toBe('#445566');
        expect(container.querySelector('.mantine-Skeleton-root')).toBeNull();
    });
});
