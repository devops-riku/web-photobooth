import { writable } from 'svelte/store';

export const photoboothStore = writable<{
  shots: string[];
  finalStrip: string | null;
}>({
  shots: [],
  finalStrip: null
});
