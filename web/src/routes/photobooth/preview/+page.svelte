<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';

  import { renderStripCanvas } from '$lib/utils/renderStrip';
  import { computeStripLayout } from '$lib/utils/stripLayout';
  import { applyGLFXFilter } from '$lib/utils/glfxFilters';
  import { PREVIEW_SETTINGS } from './settings';


  let layout: { count: 2 | 3 | 4 } | null = null;
  let shots: string[] = [];

  let filter = PREVIEW_SETTINGS.DEFAULT_FILTER;
  let caption = PREVIEW_SETTINGS.DEFAULT_CAPTION;
  let captionSize = PREVIEW_SETTINGS.DEFAULT_CAPTION_SIZE;
  let font = PREVIEW_SETTINGS.DEFAULT_FONT;



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
      // 1. Prepare layout basics for calculations
      const dpi = PREVIEW_SETTINGS.STRIP_THUMBNAIL_WIDTH;
      const initialLayout = computeStripLayout(layout.count, dpi);
      
      // 2. Load Brand Assets
      const loadImage = (src: string) => new Promise<HTMLImageElement>((resolve, reject) => {
        const img = new Image();
        img.onload = () => resolve(img);
        img.onerror = reject;
        img.src = src;
      });

      const [logo, qr] = await Promise.all([
        loadImage('/wuby_logo.png'),
        loadImage('/sample_qr.png')
      ]);

      const logoRatio = logo.width / logo.height;
      const qrRatio = qr.width / qr.height;
      const availableBrandWidth = initialLayout.contentWidthPx - (PREVIEW_SETTINGS.BRAND_SIDE_PADDING_PX * 2);
      const logoW_qrW_total = availableBrandWidth - PREVIEW_SETTINGS.BRAND_GAP_PX;
      const brandHeight = logoW_qrW_total / (logoRatio + qrRatio);

      // 3. Calculate dynamic margins based on paddings in settings
      const topPx = PREVIEW_SETTINGS.BRAND_TOP_PX + brandHeight + PREVIEW_SETTINGS.BRAND_BOT_PX;
      const topIn = topPx / dpi;

      // Caption and timestamp estimated heights for margin calculation
      const timestampH = 20; // 12px + padding
      const captionH = captionSize * 1.2; 
      const bottomPx = PREVIEW_SETTINGS.TIMESTAMP_TOP_PX + timestampH + PREVIEW_SETTINGS.TIMESTAMP_BOT_PX + captionH + PREVIEW_SETTINGS.CAPTION_BOT_PX;
      const bottomIn = bottomPx / dpi;

      // 4. Build the final strip canvas
      const canvas = await renderStripCanvas(
        shots,
        layout.count,
        dpi,
        filter,
        topIn,
        bottomIn
      );

      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      // Get precise layout results with our dynamic margins
      const fullLayout = computeStripLayout(layout.count, dpi, topIn, bottomIn);

      // ─────────────────────────
      // TOP BRAND (Logo + QR)
      // ─────────────────────────
      const logoW = logoRatio * brandHeight;
      const qrW = qrRatio * brandHeight;
      const brandX = fullLayout.contentX + PREVIEW_SETTINGS.BRAND_SIDE_PADDING_PX;
      const brandY = PREVIEW_SETTINGS.BRAND_TOP_PX;

      ctx.drawImage(logo, brandX, brandY, logoW, brandHeight);
      ctx.drawImage(qr, brandX + logoW + PREVIEW_SETTINGS.BRAND_GAP_PX, brandY, qrW, brandHeight);

      // ─────────────────────────
      // TIMESTAMP (Horizontal)
      // ─────────────────────────
      const date = new Date();
      const timeStr = date.toLocaleDateString('en-US', { 
        month: 'short', 
        day: '2-digit', 
        year: 'numeric' 
      }).toUpperCase() + " • " + date.toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: true
      }).toUpperCase();

      // Position: bottom of last photo + top padding
      const lastPhotoBottom = fullLayout.topCanvasPx + (fullLayout.photoHeightPx * layout.count) + (fullLayout.gapPx * (layout.count - 1));
      const timestampY = lastPhotoBottom + PREVIEW_SETTINGS.TIMESTAMP_TOP_PX;

      ctx.fillStyle = PREVIEW_SETTINGS.TIMESTAMP_COLOR;
      ctx.textAlign = 'center';
      ctx.font = '700 12px Montserrat, sans-serif'; 
      ctx.letterSpacing = '2px';
      ctx.fillText(timeStr, canvas.width / 2, timestampY + 10); // +10 for baseline

      // ─────────────────────────
      // BOTTOM CAPTION
      // ─────────────────────────
      await document.fonts.ready;
      ctx.fillStyle = PREVIEW_SETTINGS.CAPTION_COLOR;
      ctx.textAlign = 'center';
      ctx.font = `${captionSize}px ${font}, system-ui, sans-serif`;
      
      const captionY = timestampY + timestampH + PREVIEW_SETTINGS.TIMESTAMP_BOT_PX;
      
      ctx.fillText(
        caption || ' ',
        canvas.width / 2,
        captionY + (captionSize * 0.8) // Adjust for baseline
      );


      previewUrl = canvas.toDataURL('image/png');
    } catch (e: any) {
      console.error('Update preview failed:', e);
      error = e.message || 'Failed to generate preview';
    }
  }


  let showControls = false;

  async function onFilterChange() {
    await updatePreview();
  }

  async function onCaptionInput() {
    await updatePreview();
  }
</script>

<div class="h-screen flex flex-col bg-[#f8f2ff] overflow-hidden">
  <!-- Top Bar -->
  <header class="h-16 border-b border-purple-100 flex justify-between items-center px-4 md:px-8 bg-white shrink-0 z-20">
    <div class="flex items-center gap-2">
      <div class="w-1.5 h-1.5 rounded-full bg-purple-400"></div>
      <span class="text-[10px] font-bold tracking-[0.2em] text-purple-400 uppercase">Preview</span>
    </div>
    
    <div class="flex gap-2 md:gap-3">
      <button 
        class="md:hidden border border-purple-200 bg-purple-100/30 text-purple-500 px-3 py-1.5 rounded-lg text-[10px] font-bold transition-all active:scale-95 flex items-center gap-1.5"
        on:click={() => showControls = !showControls}
      >
        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"/></svg>
        {showControls ? 'Close' : 'Filter'}
      </button>

      <button
        class="border border-purple-200 hover:bg-purple-100/50 text-purple-500 px-4 md:px-6 py-1.5 md:py-2 rounded-lg text-[10px] md:text-xs font-bold transition-all active:scale-95"
        on:click={() => goto('/photobooth/capture')}
      >
        Retake
      </button>
      <button
        class="bg-purple-500 hover:bg-purple-600 text-white px-4 md:px-6 py-1.5 md:py-2 rounded-lg text-[10px] md:text-xs font-bold transition-all active:scale-95 shadow-md shadow-purple-200"
        on:click={() => goto('/photobooth/save')}
      >
        Continue
      </button>
    </div>
  </header>

  <main class="flex-1 flex flex-col md:flex-row overflow-hidden relative">
    <!-- Strip Preview Section - Minimalist -->
    <div class="flex-1 flex items-center justify-center p-6 md:p-12 bg-purple-100/10">
      {#if error}
        <div class="text-center">
          <p class="text-sm text-red-400 mb-4">{error}</p>
          <button class="text-xs font-bold text-purple-500 underline" on:click={() => location.reload()}>Try again</button>
        </div>
      {:else if previewUrl}
        <div class="h-full relative flex items-center justify-center">
          <img
            src={previewUrl}
            alt="Photobooth strip"
            class="h-full w-auto max-w-full object-contain shadow-[0_20px_60px_rgba(159,122,234,0.18)] border-[6px] md:border-[10px] border-white rounded-sm transition-transform duration-500"
          />
        </div>
      {:else}
        <div class="flex flex-col items-center gap-4">
          <div class="w-6 h-6 border-2 border-purple-200 border-t-purple-500 rounded-full animate-spin"></div>
          <p class="text-[10px] font-bold text-purple-400 uppercase tracking-widest">Processing</p>
        </div>
      {/if}
    </div>

    <!-- Controls Section - Responsive Sidebar/Panel -->
    <aside class="
      fixed inset-0 z-10 bg-white/95 backdrop-blur-xl md:static md:bg-white md:backdrop-blur-none
      md:w-80 border-l border-purple-100 flex flex-col overflow-hidden transition-all duration-500 ease-out
      {showControls ? 'translate-y-0 opacity-100' : 'translate-y-full opacity-0 md:translate-y-0 md:opacity-100'}
    ">
      <div class="p-8 pb-32 md:pb-8 flex flex-col gap-10 h-full overflow-y-auto custom-scrollbar mt-16 md:mt-0">
        <section>
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-[10px] font-bold uppercase text-purple-400 tracking-[0.2em]">Filter</h2>
          </div>
          <div class="grid grid-cols-2 md:flex md:flex-wrap gap-2">
            {#each PREVIEW_SETTINGS.AVAILABLE_FILTERS as f}
              <button
                class="px-4 py-2.5 text-[11px] font-medium rounded-xl border transition-all {filter === f ? 'bg-purple-100/50 border-purple-300 text-purple-700' : 'bg-transparent border-slate-100 text-slate-400 hover:border-purple-200 hover:text-purple-500'}"
                on:click={() => { filter = f; onFilterChange(); }}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            {/each}
          </div>
        </section>

        <section>
          <h2 class="text-[10px] font-bold uppercase text-purple-300 tracking-[0.2em] mb-6">Typography</h2>
          
          <div class="flex flex-col gap-6">
            <div class="flex flex-col gap-2">
              <label for="font-select" class="text-[10px] font-medium text-slate-400">Style</label>
              <select
                id="font-select"
                bind:value={font}
                on:change={updatePreview}
                class="w-full bg-slate-50 border-none rounded-xl px-4 py-3 text-xs font-medium text-slate-600 focus:ring-1 focus:ring-purple-200 transition-all outline-none"
              >
                {#each PREVIEW_SETTINGS.FONTS as fontOption}
                  <option value={fontOption.id}>{fontOption.label}</option>
                {/each}
              </select>
            </div>

            <div class="flex flex-col gap-2">
              <label for="message-input" class="text-[10px] font-medium text-slate-400">Caption</label>
              <input
                id="message-input"
                type="text"
                bind:value={caption}
                on:input={onCaptionInput}
                placeholder="Type here..."
                class="w-full bg-slate-50 border-none rounded-xl px-4 py-3 text-xs font-medium text-slate-600 focus:ring-1 focus:ring-purple-200 transition-all outline-none"
              />
            </div>

            <div class="flex flex-col gap-2">
              <div class="flex justify-between">
                <label for="size-range" class="text-[10px] font-medium text-slate-400">Size</label>
                <span class="text-[10px] font-bold text-purple-400">{captionSize}px</span>
              </div>
              <input
                id="size-range"
                type="range"
                min="12"
                max="100"
                step="1"
                bind:value={captionSize}
                on:input={updatePreview}
                class="w-full accent-purple-300"
              />
            </div>
          </div>
        </section>
      </div>
    </aside>
  </main>
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }
  .custom-scrollbar::-webkit-scrollbar-track {
    background: transparent;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.05);
    border-radius: 10px;
  }
  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: rgba(0,0,0,0.1);
  }
  .no-scrollbar::-webkit-scrollbar {
    display: none;
  }
  .no-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
</style>



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

