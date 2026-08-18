import {Photo} from './Photo';

// MemoryGroup is a set of photos from one prior year for "On This Day".
export interface MemoryGroup {
    year: number;
    yearsAgo: number;
    photos: Photo[];
}
