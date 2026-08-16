import React from 'react';
import {render} from '../../test/test-utils';
import {GridSkeleton} from './GridSkeleton';

describe('GridSkeleton', () => {
    it('renders skeleton cards in a grid', () => {
        render(<GridSkeleton columns={4} rows={2} />);
        // 4 columns × 2 rows = 8 skeleton cards. Mantine Skeleton renders
        // with a mantine-Skeleton-root class.
        expect(document.querySelectorAll('.mantine-Skeleton-root').length).toBe(8);
    });

    it('renders no spinner', () => {
        render(<GridSkeleton columns={4} rows={2} />);
        // Ensure there is no Loader element (spinning arrows) present.
        expect(document.querySelector('.mantine-Loader-root')).toBeNull();
    });

    it('respects custom row/column counts', () => {
        render(<GridSkeleton columns={3} rows={1} />);
        expect(document.querySelectorAll('.mantine-Skeleton-root').length).toBe(3);
    });
});
