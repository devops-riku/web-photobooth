import { writable } from 'svelte/store';

export type LayoutOption = {
  count: 2 | 3 | 4;
  width: number;   // px
  height: number;  // px
};

export const layoutStore = writable<LayoutOption | null>(null);
