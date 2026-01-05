import { writable } from 'svelte/store';

export type LayoutOption = {
    count: 2 | 3 | 4;
    width: number;   // px
    height: number;  // px
    timer: number;   // seconds
};

export const layoutStore = writable<LayoutOption | null>(null);
