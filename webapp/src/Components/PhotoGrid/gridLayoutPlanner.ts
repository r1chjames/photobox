/**
 * Pure quantum-masonry layout planner (issue #174).
 *
 * Deterministically assigns every photo a display variant and grid spans so
 * the photo grid reads as curated rather than uniform, while guaranteeing:
 *
 * - identical inputs → identical output (no render-time randomness);
 * - photo order is never reshuffled — placement stays CSS dense flow;
 * - appending photos never re-plans the existing prefix (constraints only
 *   look backward);
 * - source aspect ratios are respected — variants only apply when the cover
 *   crop stays within a bounded tolerance.
 *
 * Selection is driven by a stable hash of the photo id, then constrained by
 * a small pattern budget: minimum spacing between large variants and a cap
 * on consecutive identical footprints where an alternate exists.
 */

export type LayoutVariant =
    | 'standard'  // 1-col landscape tile, height jittered
    | 'featured'  // 2-col panorama hero
    | 'portrait'  // 1-col portrait tile, height jittered
    | 'tall'      // 1-col very-tall portrait
    | 'square'    // 1-col square tile, height jittered
    | 'mobile';   // uniform 4:3 tile for small screens

export interface PlannedTileLayout {
    variant: LayoutVariant;
    columnSpan: number;
    rowSpan: number;
}

export interface PlannerPhotoInput {
    id: string;
    width?: number;
    height?: number;
}

export type PlannerMode = 'desktop' | 'mobile';

export interface PlanGridInput {
    photos: PlannerPhotoInput[];
    containerWidth: number;
    columns: number;
    gap: number;
    mode?: PlannerMode;
}

/** Height of one grid auto-row track; every tile edge lands on this rhythm. */
export const GRID_QUANTUM = 16;

// Display aspect (height / width) per variant.
const RATIO: Record<LayoutVariant, number> = {
    standard: 11 / 16,
    featured: 11 / 28,
    portrait: 14 / 10,
    tall: 16 / 9,
    square: 1,
    mobile: 3 / 4,
};

// Deterministic height jitter range per variant, as a fraction of the base
// ratio. Fixed variants (featured/tall/mobile) carry no jitter.
const JITTER: Partial<Record<LayoutVariant, number>> = {
    standard: 0.15,
    portrait: 0.12,
    square: 0.10,
};

// Source-ratio (width / height) classification gates. Variants only apply
// when the cover crop stays within reasonable bounds.
const LANDSCAPE_MIN_RATIO = 1.35;
const PORTRAIT_MAX_RATIO = 0.85;
const FEATURED_MIN_RATIO = 2.0;
const TALL_MAX_RATIO = 0.62;

// Hash thresholds for variant selection.
const FEATURED_CHANCE = 0.25;
const TALL_CHANCE = 0.35;
const SQUARE_STANDARD_CHANCE = 0.45;

// Pattern budget: minimum index spacing between large variants and the max
// length of an identical-variant run (when an alternate exists).
const FEATURED_MIN_GAP = 6;
const TALL_MIN_GAP = 5;
const MAX_CONSECUTIVE = 3;

// FNV-1a with a murmur-style finalizer, mapped to [0, 1). Stable across
// calls, sessions and JS engines for the same id.
const hash01 = (id: string, salt: number): number => {
    let h = (2166136261 ^ Math.imul(salt + 1, 0x9e3779b1)) >>> 0;
    for (let i = 0; i < id.length; i++) {
        h ^= id.charCodeAt(i);
        h = Math.imul(h, 16777619) >>> 0;
    }
    h ^= h >>> 16;
    h = Math.imul(h, 2246822507) >>> 0;
    h ^= h >>> 13;
    h = Math.imul(h, 3266489909) >>> 0;
    h ^= h >>> 16;
    return (h >>> 0) / 4294967296;
};

const lastNIdentical = (recent: LayoutVariant[], variant: LayoutVariant): boolean =>
    recent.length >= MAX_CONSECUTIVE && recent.every(v => v === variant);

/**
 * Plan the grid layout for an ordered photo list. Returns a map keyed by
 * photo id with the variant, column span and row span each tile should
 * occupy. Empty when the geometry is unknown (columns/width unset) — the
 * caller should render its skeleton instead.
 */
export function planGridLayout(input: PlanGridInput): Map<string, PlannedTileLayout> {
    const {photos, containerWidth, columns, gap, mode = 'desktop'} = input;
    const plan = new Map<string, PlannedTileLayout>();
    if (columns < 1 || containerWidth <= 0 || photos.length === 0) return plan;

    const cardWidth = Math.max(80, (containerWidth - (columns - 1) * gap) / columns);

    let lastFeaturedIndex = Number.NEGATIVE_INFINITY;
    let lastTallIndex = Number.NEGATIVE_INFINITY;
    const recentVariants: LayoutVariant[] = [];

    photos.forEach((photo, index) => {
        const w = photo.width ?? 0;
        const h = photo.height ?? 0;
        const sourceRatio = w > 0 && h > 0 ? w / h : 0;

        // Candidate footprints for this photo's class, primary first.
        let candidates: LayoutVariant[];
        if (mode === 'mobile') {
            candidates = ['mobile'];
        } else if (sourceRatio <= 0) {
            // Missing/zero dimensions: deterministic safe fallback.
            candidates = ['standard'];
        } else if (sourceRatio > LANDSCAPE_MIN_RATIO) {
            candidates = ['standard'];
            if (sourceRatio >= FEATURED_MIN_RATIO && columns >= 2) {
                candidates.push('featured');
            }
        } else if (sourceRatio < PORTRAIT_MAX_RATIO) {
            candidates = ['portrait'];
            if (sourceRatio <= TALL_MAX_RATIO) {
                candidates.push('tall');
            }
        } else {
            candidates = ['square', 'standard'];
        }

        // Hash-driven primary pick within the class.
        const pick = hash01(photo.id, 0);
        let primary: LayoutVariant;
        if (candidates.length === 1) {
            primary = candidates[0];
        } else if (candidates[0] === 'standard' && candidates[1] === 'featured') {
            primary = pick < FEATURED_CHANCE ? 'featured' : 'standard';
        } else if (candidates[0] === 'portrait' && candidates[1] === 'tall') {
            primary = pick < TALL_CHANCE ? 'tall' : 'portrait';
        } else {
            primary = pick < SQUARE_STANDARD_CHANCE ? 'standard' : 'square';
        }

        const satisfies = (variant: LayoutVariant): boolean => {
            if (variant === 'featured') return index - lastFeaturedIndex >= FEATURED_MIN_GAP;
            if (variant === 'tall') return index - lastTallIndex >= TALL_MIN_GAP;
            return true;
        };

        let variant = primary;
        if (!satisfies(variant)) {
            // Primary violates spacing — fall back to the base footprint.
            variant = candidates.find(c => c !== primary && satisfies(c)) ?? candidates[0];
        } else if (candidates.length > 1 && lastNIdentical(recentVariants, variant)) {
            // Break the run when the class admits an alternate that fits.
            const alternate = candidates.find(c => c !== variant && satisfies(c));
            if (alternate) variant = alternate;
        }

        if (variant === 'featured') lastFeaturedIndex = index;
        if (variant === 'tall') lastTallIndex = index;
        recentVariants.push(variant);
        if (recentVariants.length > MAX_CONSECUTIVE) recentVariants.shift();

        const jitterRange = JITTER[variant];
        const ratio = jitterRange === undefined
            ? RATIO[variant]
            : RATIO[variant] * (1 + (hash01(photo.id, 2) - 0.5) * 2 * jitterRange);

        const columnSpan = variant === 'featured' ? 2 : 1;
        const tileWidth = columnSpan * cardWidth + (columnSpan - 1) * gap;
        const rowSpan = Math.max(1, Math.round((tileWidth * ratio + gap) / (GRID_QUANTUM + gap)));

        plan.set(photo.id, {variant, columnSpan, rowSpan});
    });

    return plan;
}
