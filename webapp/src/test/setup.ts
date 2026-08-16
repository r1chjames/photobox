import '@testing-library/jest-dom';
import { expect, afterEach, beforeEach, vi } from 'vitest';
import { cleanup } from '@testing-library/react';
import * as matchers from '@testing-library/jest-dom/matchers';

// Extend Vitest's expect with @testing-library/jest-dom matchers
expect.extend(matchers);

// ResizeObserver does not exist in happy-dom; the PhotoGrid measures its
// container with it to derive card width/height. Provide a stub that fires
// the callback with a plausible content width so layout-dependent renders
// (virtualized grid, skeleton gating) behave like a real browser. Registered
// in beforeEach because afterEach calls unstubAllGlobals.
class ResizeObserverStub {
    private cb: ResizeObserverCallback;
    constructor(cb: ResizeObserverCallback) {
        this.cb = cb;
    }
    observe(el: Element) {
        this.cb([{ target: el, contentRect: { width: 1000 } } as unknown as ResizeObserverEntry], this as unknown as ResizeObserver);
    }
    unobserve() {}
    disconnect() {}
}

beforeEach(() => {
    vi.stubGlobal('ResizeObserver', ResizeObserverStub);
});

// Cleanup after each test
afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
});
