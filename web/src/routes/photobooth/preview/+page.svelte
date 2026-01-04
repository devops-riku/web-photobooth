<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';

  import { renderStripCanvas } from '$lib/utils/renderStrip';
  import { applyGLFXFilter } from '$lib/utils/glfxFilters';


  let layout: { count: 2 | 3 | 4 } | null = null;
  let shots: string[] = [];

  let filter = 'none';
  let caption = '';
  let captionSize = 60;
  let font = 'Lobster';



  let previewUrl: string | null = null;
  let error: string | null = null;

  // Subscribe stores (NO goto here)
  layoutStore.subscribe(v => layout = v);
  photoboothStore.subscribe(v => shots = v.shots);


  onMount(async () => {
    await updatePreview();
  });


  async function updatePreview() {
    if (!layout || shots.length === 0) {
      if (!layout) goto('/photobooth');
      return;
    }
    error = null;

    try {
      // 1. Build the strip (this now handles square crop + per-photo filter + dynamic height)
      const canvas = await renderStripCanvas(
        shots,
        layout.count,
        300,
        filter
      );

      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      // ─────────────────────────
      // TOP BRAND (inside top canvas)
      // ─────────────────────────
      ctx.fillStyle = '#000';
      ctx.textAlign = 'center';
      ctx.font = 'bold 32px system-ui, sans-serif';
      ctx.fillText('WUBY', canvas.width / 2, 100);

      // ─────────────────────────
      // BOTTOM CAPTION
      // ─────────────────────────
      await document.fonts.ready; // Wait for fonts to load
      ctx.font = `${captionSize}px ${font}, system-ui, sans-serif`;
      ctx.fillText(

        caption || ' ',
        canvas.width / 2,
        canvas.height - 100
      );


      previewUrl = canvas.toDataURL('image/png');
    } catch (e: any) {
      console.error('Update preview failed:', e);
      error = e.message || 'Failed to generate preview';
    }
  }



  async function onFilterChange() {
    await updatePreview();
  }

  async function onCaptionInput() {
    await updatePreview();
  }
</script>

<div class="h-screen flex flex-col overflow-hidden bg-slate-50">
  <header class="p-3 border-b bg-white flex justify-between items-center shrink-0">
    <h1 class="text-lg font-bold">Preview</h1>
    <button
      class="btn variant-filled-primary px-6 py-1.5 text-sm font-bold shadow-sm"
      on:click={() => goto('/photobooth/save')}
    >
      Continue
    </button>
  </header>

  <main class="flex-1 flex flex-col md:flex-row overflow-hidden p-4 gap-6">
    <!-- Strip Preview Section -->
    <div class="flex-1 flex items-center justify-center min-h-0">
      {#if error}
        <div class="bg-red-50 border border-red-200 text-red-700 px-6 py-8 rounded-xl text-center max-w-sm shadow-sm">
          <p class="font-bold mb-2">Something went wrong</p>
          <p class="text-sm opacity-80 mb-4">{error}</p>
          <button class="btn variant-ghost-surface py-2 px-4 text-xs" on:click={() => location.reload()}>
            Try Again
          </button>
        </div>
      {:else if previewUrl}
        <div class="h-full flex items-center justify-center p-2">
          <img
            src={previewUrl}
            alt="Photobooth strip"
            class="h-full w-auto object-contain rounded shadow-2xl border-[6px] border-white"
          />
        </div>
      {:else}
        <div class="flex flex-col items-center gap-4">
          <div class="w-12 h-12 border-4 border-primary-500/20 border-t-primary-500 rounded-full animate-spin"></div>
          <p class="text-slate-400 font-medium animate-pulse">Building your strip...</p>
        </div>
      {/if}
    </div>

    <!-- Controls Section -->
    <aside class="w-full md:w-80 shrink-0 overflow-y-auto bg-white rounded-2xl border border-slate-200 shadow-xl flex flex-col">
      <div class="p-5 flex flex-col gap-5">
        <h2 class="text-xs font-black uppercase text-slate-400 tracking-widest">Customization</h2>
        
        <!-- Filter -->
        <label class="flex flex-col gap-1.5">
          <span class="text-sm font-bold text-slate-700">Effect</span>
          <div class="grid grid-cols-3 gap-2">
            {#each ['none', 'cinematic', 'film', 'warm', 'bw'] as f}
              <button
                class="py-2 text-xs rounded-lg border transition-all {filter === f ? 'bg-primary-500 text-white border-primary-600 shadow-md ring-2 ring-primary-500/20' : 'bg-slate-50 text-slate-600 border-slate-200 hover:border-slate-300'}"
                on:click={() => { filter = f; onFilterChange(); }}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            {/each}
          </div>
        </label>

        <!-- Font -->
        <label class="flex flex-col gap-1.5">
          <span class="text-sm font-bold text-slate-700">Font Style</span>
          <select
            bind:value={font}
            on:change={updatePreview}
            class="border border-slate-200 rounded-xl px-3 py-2 text-sm bg-slate-50 focus:ring-2 focus:ring-primary-500 focus:border-primary-500 outline-none transition-all"
          >
            {#each [
              { id: 'Lobster', label: 'Lobster' },
              { id: 'Pacifico', label: 'Pacifico' },
              { id: 'Caveat', label: 'Caveat' },
              { id: 'Dancing Script', label: 'Dancing' },
              { id: 'Bebas Neue', label: 'Bebas' },
              { id: 'Righteous', label: 'Righteous' },
              { id: 'Abril Fatface', label: 'Abril' },
              { id: 'Cormorant Garamond', label: 'Classic' },
              { id: 'Permanent Marker', label: 'Marker' },
              { id: 'Special Elite', label: 'Typewriter' },
              { id: 'Monoton', label: 'Retro' },
              { id: 'Montserrat', label: 'Clean' }
            ] as fontOption}
              <option value={fontOption.id}>{fontOption.label}</option>
            {/each}
          </select>
        </label>

        <!-- Caption -->
        <label class="flex flex-col gap-1.5">
          <span class="text-sm font-bold text-slate-700">Message</span>
          <input
            type="text"
            bind:value={caption}
            on:input={onCaptionInput}
            placeholder="Write something..."
            class="border border-slate-200 rounded-xl px-3 py-2 text-sm bg-slate-50 focus:ring-2 focus:ring-primary-500 focus:border-primary-500 outline-none transition-all"
          />
        </label>

        <!-- Caption Size -->
        <label class="flex flex-col gap-1.5 pb-4">
          <div class="flex justify-between items-center">
            <span class="text-sm font-bold text-slate-700">Font Size</span>
            <span class="text-[10px] font-mono bg-slate-100 px-2 py-0.5 rounded text-slate-500">{captionSize}px</span>
          </div>
          <input
            type="range"
            min="12"
            max="60"
            step="1"
            bind:value={captionSize}
            on:input={updatePreview}
            class="w-full h-1.5 bg-slate-200 rounded-lg appearance-none cursor-pointer accent-primary-500"
          />
        </label>
      </div>
    </aside>
  </main>
</div>



<!-- Invisible fonts to keep them loaded for Canvas -->
<div style="visibility: hidden; position: absolute; height: 0; width: 0; overflow: hidden;">
  <span style="font-family: 'Lobster';">.</span>
  <span style="font-family: 'Pacifico';">.</span>
  <span style="font-family: 'Caveat';">.</span>
  <span style="font-family: 'Dancing Script';">.</span>
  <span style="font-family: 'Bebas Neue';">.</span>
  <span style="font-family: 'Righteous';">.</span>
  <span style="font-family: 'Abril Fatface';">.</span>
  <span style="font-family: 'Cormorant Garamond';">.</span>
  <span style="font-family: 'Permanent Marker';">.</span>
  <span style="font-family: 'Special Elite';">.</span>
  <span style="font-family: 'Monoton';">.</span>
  <span style="font-family: 'Montserrat';">.</span>
</div>

