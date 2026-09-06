import { describe, it, expect } from 'vitest';
import {
    planGridLayout,
    GRID_QUANTUM,
    PlannerPhotoInput,
    PlanGridInput,
} from './gridLayoutPlanner';

const photo = (id: string, width?: number, height?: number): PlannerPhotoInput => ({ id, width, height });

// Photo helpers: landscape (4032x3024), portrait (3024x4032), panorama
// (6000x2000), very tall (2000x4000), square (3000x3000), no dimensions.
const landscape = (id: string) => photo(id, 4032, 3024);
const portrait = (id: string) => photo(id, 3024, 4032);
const panorama = (id: string) => photo(id, 6000, 2000);
const veryTall = (id: string) => photo(id, 2000, 4000);
const square = (id: string) => photo(id, 3000, 3000);

const BASE = { containerWidth: 1200, columns: 5, gap: 10 };
const plan = (photos: PlannerPhotoInput[], overrides: Partial<PlanGridInput> = {}) =>
    planGridLayout({ ...BASE, ...overrides, photos });

describe('gridLayoutPlanner', () => {
    it('is deterministic for identical inputs', () => {
        const photos = Array.from({ length: 60 }, (_, i) => landscape(`p${i}`));
        const first = plan(photos);
        const second = plan(photos);
        expect([...first.entries()]).toEqual([...second.entries()]);
    });

    it('returns an empty plan for unknown geometry', () => {
        expect(plan([landscape('p1')], { columns: 0 }).size).toBe(0);
        expect(plan([landscape('p1')], { containerWidth: 0 }).size).toBe(0);
        expect(plan([]).size).toBe(0);
    });

    it('keys the plan by photo id and preserves input order', () => {
        const photos = [landscape('p1'), portrait('p2'), square('p3')];
        const result = plan(photos);
        expect([...result.keys()]).toEqual(['p1', 'p2', 'p3']);
    });

    it('never plans a featured tile wider than the grid', () => {
        const photos = Array.from({ length: 30 }, (_, i) => panorama(`p${i}`));
        const result = plan(photos, { columns: 1 });
        result.forEach(tile => expect(tile.columnSpan).toBe(1));
    });

    it('plans only uniform tiles in mobile mode', () => {
        const photos = [landscape('p1'), portrait('p2'), panorama('p3'), square('p4')];
        const result = plan(photos, { mode: 'mobile' as const });
        result.forEach(tile => {
            expect(tile.variant).toBe('mobile');
            expect(tile.columnSpan).toBe(1);
        });
    });

    it('falls back safely for missing dimensions', () => {
        const result = plan([photo('nodim'), photo('zero', 0, 0)]);
        result.forEach(tile => {
            expect(tile.variant).toBe('standard');
            expect(tile.columnSpan).toBe(1);
            expect(tile.rowSpan).toBeGreaterThanOrEqual(1);
        });
    });

    it('caps consecutive identical variants when an alternate exists', () => {
        // Squares can always alternate square/standard, so no run may exceed 3.
        const photos = Array.from({ length: 60 }, (_, i) => square(`s${i}`));
        const variants = [...plan(photos).values()].map(t => t.variant);
        let run = 1;
        let maxRun = 1;
        for (let i = 1; i < variants.length; i++) {
            run = variants[i] === variants[i - 1] ? run + 1 : 1;
            maxRun = Math.max(maxRun, run);
        }
        expect(maxRun).toBeLessThanOrEqual(3);
    });

    it('spaces featured tiles apart', () => {
        const photos = Array.from({ length: 120 }, (_, i) => panorama(`pan${i}`));
        const variants = [...plan(photos).values()].map(t => t.variant);
        const featured = variants.map((v, i) => (v === 'featured' ? i : -1)).filter(i => i >= 0);
        expect(featured.length).toBeGreaterThan(0);
        for (let i = 1; i < featured.length; i++) {
            expect(featured[i] - featured[i - 1]).toBeGreaterThanOrEqual(6);
        }
    });

    it('spaces tall tiles apart', () => {
        const photos = Array.from({ length: 120 }, (_, i) => veryTall(`t${i}`));
        const variants = [...plan(photos).values()].map(t => t.variant);
        const tall = variants.map((v, i) => (v === 'tall' ? i : -1)).filter(i => i >= 0);
        expect(tall.length).toBeGreaterThan(0);
        for (let i = 1; i < tall.length; i++) {
            expect(tall[i] - tall[i - 1]).toBeGreaterThanOrEqual(5);
        }
    });

    it('introduces variation into a homogeneous landscape run', () => {
        const photos = Array.from({ length: 60 }, (_, i) => landscape(`l${i}`));
        const spans = [...plan(photos).values()].map(t => t.rowSpan);
        const uniqueSpans = new Set(spans);
        // Height jitter must break the single-footprint monotony.
        expect(uniqueSpans.size).toBeGreaterThan(1);
    });

    it('bounds jitter so no standard tile doubles in height', () => {
        const photos = Array.from({ length: 200 }, (_, i) => landscape(`j${i}`));
        const spans = [...plan(photos).values()].filter(t => t.variant === 'standard').map(t => t.rowSpan);
        expect(Math.max(...spans) / Math.min(...spans)).toBeLessThan(1.35);
    });

    it('computes row spans on the 16px quantum from measured geometry', () => {
        // Base display ratio (height/width) and jitter band per variant —
        // mirrors the planner constants.
        const baseRatio: Record<string, number> = {
            standard: 11 / 16, featured: 11 / 28, portrait: 14 / 10,
            tall: 16 / 9, square: 1, mobile: 3 / 4,
        };
        const jitterBand: Record<string, number> = {
            standard: 0.15, portrait: 0.12, square: 0.10,
            featured: 0, tall: 0, mobile: 0,
        };
        const photos = [square('sq1'), landscape('l1'), portrait('p1'), panorama('pan1'), veryTall('vt1')];
        const cardWidth = (BASE.containerWidth - (BASE.columns - 1) * BASE.gap) / BASE.columns;
        plan(photos).forEach(tile => {
            const width = tile.columnSpan * cardWidth + (tile.columnSpan - 1) * BASE.gap;
            const height = tile.rowSpan * (GRID_QUANTUM + BASE.gap) - BASE.gap;
            const ratio = height / width;
            const base = baseRatio[tile.variant];
            const band = jitterBand[tile.variant];
            // Half a quantum track of rounding slack on either side.
            const slack = (GRID_QUANTUM + BASE.gap) / 2 / width;
            expect(ratio).toBeGreaterThanOrEqual(base * (1 - band) - slack);
            expect(ratio).toBeLessThanOrEqual(base * (1 + band) + slack);
        });
    });

    it('keeps the existing prefix stable when photos are appended', () => {
        const photos = Array.from({ length: 60 }, (_, i) =>
            i % 3 === 0 ? landscape(`a${i}`) : i % 3 === 1 ? portrait(`a${i}`) : square(`a${i}`));
        const before = plan(photos);
        const after = plan([...photos, ...Array.from({ length: 60 }, (_, i) => landscape(`b${i + 60}`))]);
        photos.forEach(p => expect(after.get(p.id)).toEqual(before.get(p.id)));
    });

    it('recomputes spans when width, columns or gap change', () => {
        const photos = Array.from({ length: 30 }, (_, i) => portrait(`w${i}`));
        const wide = [...plan(photos, { containerWidth: 1600, columns: 8 }).values()];
        const narrow = [...plan(photos, { containerWidth: 600, columns: 3 }).values()];
        const sameSpans = wide.every((tile, i) => tile.rowSpan === narrow[i].rowSpan);
        expect(sameSpans).toBe(false);
    });

    it('degrades safely across column counts', () => {
        const photos = [landscape('p1'), portrait('p2'), panorama('p3'), veryTall('p4'), square('p5')];
        [1, 2, 5, 12].forEach(columns => {
            plan(photos, { columns }).forEach(tile => {
                expect(tile.columnSpan).toBeLessThanOrEqual(columns);
                expect(tile.rowSpan).toBeGreaterThanOrEqual(1);
            });
        });
    });
});
