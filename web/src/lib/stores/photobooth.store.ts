import { writable } from 'svelte/store';

const initial = {
  shots: [],
  finalStrip: null,
  uploadedId: null,
  settings: null
};

// Check for existing data in session
let stored = initial;
if (typeof window !== 'undefined') {
  const data = sessionStorage.getItem('photobooth_store');
  if (data) {
    try {
      stored = JSON.parse(data);
    } catch (e) { /* ignore */ }
  }
}

export const photoboothStore = writable<{
  shots: string[];
  finalStrip: string | null;
  uploadedId: string | null;
  settings: {
    filter: string;
    stripColor: string;
    caption: string;
    captionSize: number;
    font: string;
  } | null;
}>(stored);

// Subscribe to changes and save to session (Selective saving to avoid QuotaExceeded)
if (typeof window !== 'undefined') {
  photoboothStore.subscribe(value => {
    try {
      const { finalStrip, ...persistable } = value;
      sessionStorage.setItem('photobooth_store', JSON.stringify(persistable));
    } catch (e) {
      console.warn("Session storage quota exceeded, persistence disabled for this session.");
    }
  });
}
