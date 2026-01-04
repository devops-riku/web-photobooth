import { writable } from 'svelte/store';

export const photoboothStore = writable<{
  shots: string[];
}>({
  shots: []
});
