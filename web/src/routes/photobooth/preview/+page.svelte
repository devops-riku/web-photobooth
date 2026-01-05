<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';

  import { renderStripCanvas } from '$lib/utils/renderStrip';
  import { computeStripLayout } from '$lib/utils/stripLayout';
  import { applyGLFXFilter } from '$lib/utils/glfxFilters';
  import { PREVIEW_SETTINGS } from './settings';
  import { generateUUID } from '$lib/utils/uuid';
  import { API_CONFIG, getApiUrl } from '$lib/config';
  import { SAVE_SETTINGS } from '../save/settings';
  import ColorWheel from './ColorWheel.svelte';


  let layout: { count: 2 | 3 | 4 } | null = null;
  let shots: string[] = [];
  let filter = PREVIEW_SETTINGS.DEFAULT_FILTER;
  let caption = PREVIEW_SETTINGS.DEFAULT_CAPTION;
  let captionSize = PREVIEW_SETTINGS.DEFAULT_CAPTION_SIZE;
  let font = PREVIEW_SETTINGS.DEFAULT_FONT;
  let stripColor = PREVIEW_SETTINGS.STRIP_BG_COLOR;
  let roundedCorners = true;

  // Global State
  let previewUrl: string | null = null;
  let error: string | null = null;
  let logoImg: HTMLImageElement | null = null;
  let qrImg: HTMLImageElement | null = null;
  let qrData = 'https://wuby.link/pending';

  // UI State
  let isConfirming = false;
  let isUploading = false;
  let uploadError = '';
  let showControls = false;
  let showColorWheel = false;

  // Subscribe stores (NO goto here)
  layoutStore.subscribe(v => layout = v);
  photoboothStore.subscribe(v => shots = v.shots);

  const loadImage = (src: string) => new Promise<HTMLImageElement>((resolve, reject) => {
    const img = new Image();
    img.crossOrigin = 'anonymous'; // Critical for external APIs
    img.onload = () => resolve(img);
    img.onerror = (e) => reject(new Error(`Failed to load: ${src}`));
    img.src = src;
  });

  onMount(async () => {
    try {
      logoImg = await loadImage('/wuby_logo.png');
      // Load a placeholder QR for the preview
      qrImg = await loadImage(`https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=${encodeURIComponent(qrData)}`);
    } catch (e) {
      console.error('Failed to pre-load branding:', e);
    }
    await updatePreview();
  });

  function onCaptionInput() {
    debouncedUpdate();
  }

  let updateTimeout: any;
  function debouncedUpdate() {
    clearTimeout(updateTimeout);
    updateTimeout = setTimeout(() => {
      updatePreview();
    }, 50); // Small delay for smooth feel
  }


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
      
      // 2. Use Pre-loaded Brand Assets (Optional for preview)
      const logo = logoImg;
      const qr = qrImg;

      const logoRatio = logo.width / logo.height;
      const qrRatio = qr.width / qr.height;
      const availableBrandWidth = initialLayout.contentWidthPx - (PREVIEW_SETTINGS.BRAND_SIDE_PADDING_PX * 2);
      const logoW_qrW_total = availableBrandWidth - PREVIEW_SETTINGS.BRAND_GAP_PX;

      // 3. Calculate dynamic margins based on paddings in settings
      // Constrain branding height if it's too crazy
      const maxBrandHeight = 300;
      const calculatedBrandHeight = logoW_qrW_total / (logoRatio + qrRatio);
      const brandHeight = Math.min(calculatedBrandHeight, maxBrandHeight);

      const topPx = PREVIEW_SETTINGS.BRAND_TOP_PX + brandHeight + PREVIEW_SETTINGS.BRAND_BOT_PX;
      const topIn = topPx / dpi;

      // Caption and timestamp estimated heights for margin calculation
      const timestampH = 20; // 12px + padding
      const captionH = captionSize * 1.2; 
      const bottomPx = PREVIEW_SETTINGS.TIMESTAMP_TOP_PX + timestampH + PREVIEW_SETTINGS.TIMESTAMP_BOT_PX + captionH + PREVIEW_SETTINGS.CAPTION_BOT_PX;
      const bottomIn = bottomPx / dpi;

      const canvas = await renderStripCanvas(
        shots,
        layout.count,
        dpi,
        filter,
        topIn,
        bottomIn,
        stripColor,
        roundedCorners
      );

      const ctx = canvas.getContext('2d');
      if (!ctx) return;

      // Get precise layout results with our dynamic margins
      const fullLayout = computeStripLayout(layout.count, dpi, topIn, bottomIn);

      // ─────────────────────────
      // TOP BRAND (Logo + QR) - Aligned to Edges
      // ─────────────────────────
      const logoW = logoRatio * brandHeight;
      const qrHeight = brandHeight * 0.6;
      const qrW = qrRatio * qrHeight;

      const logoX = fullLayout.contentX;
      const qrX = fullLayout.contentX + fullLayout.contentWidthPx - qrW;
      const logoY = PREVIEW_SETTINGS.BRAND_TOP_PX;
      const qrY = logoY + (brandHeight - qrHeight) / 2;

      if (logo) ctx.drawImage(logo, logoX, logoY, logoW, brandHeight);
      if (qr) ctx.drawImage(qr, qrX, qrY, qrW, qrHeight);

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

  let processing = false;
  async function handleConfirm() {
    if (processing || isUploading) return;
    processing = true;
    isUploading = true;
    isConfirming = false;
    uploadError = '';

    try {
      const stripId = generateUUID();
      const token = localStorage.getItem('sb_token');
      const uid = localStorage.getItem('sb_uid');
      
      let shareUrl = '';
      if (token && uid) {
        shareUrl = `${API_CONFIG.DO_SPACES_URL}/strips/${uid}/${stripId}.png`;
      } else {
        shareUrl = `${API_CONFIG.DO_SPACES_URL}/strips/guest/${stripId}.png`;
      }
      console.log("PHOTOBOOTH DEBUG: Generating ID:", stripId);
      console.log("PHOTOBOOTH DEBUG: Expected URL:", shareUrl);
      
      // 1. Fetch REAL QR first (with timeout)
      const realQRImg: HTMLImageElement = await Promise.race([
        new Promise<HTMLImageElement>((resolve, reject) => {
          const img = new Image();
          img.crossOrigin = 'anonymous';
          img.onload = () => resolve(img);
          img.onerror = () => reject(new Error("QR Generator failed. Connection issue?"));
          img.src = `${SAVE_SETTINGS.BASE_QR_API}${encodeURIComponent(shareUrl)}`;
        }),
        new Promise<never>((_, reject) => 
          setTimeout(() => reject(new Error("QR Generation timed out. Please try again.")), 12000)
        )
      ]);

      // 2. Re-render with REAL assets for the upload version
      const dpi = PREVIEW_SETTINGS.STRIP_THUMBNAIL_WIDTH;
      // Re-run the exact layout and drawing logic
      if (!logoImg) throw new Error("Logo not loaded");
      
      const logoRatio = logoImg.width / logoImg.height;
      const qrRatio = realQRImg.width / realQRImg.height;
      const initialLayout = computeStripLayout((layout as any).count, dpi);
      const availableBrandWidth = initialLayout.contentWidthPx - (PREVIEW_SETTINGS.BRAND_SIDE_PADDING_PX * 2);
      const logoW_qrW_total = availableBrandWidth - PREVIEW_SETTINGS.BRAND_GAP_PX;
      
      const maxBrandHeight = 200;
      const brandHeight = Math.min(logoW_qrW_total / (logoRatio + qrRatio), maxBrandHeight);

      const topPx = PREVIEW_SETTINGS.BRAND_TOP_PX + brandHeight + PREVIEW_SETTINGS.BRAND_BOT_PX;
      const topIn = topPx / dpi;
      const timestampH = 20;
      const captionH = captionSize * 1.2;
      const bottomPx = PREVIEW_SETTINGS.TIMESTAMP_TOP_PX + timestampH + PREVIEW_SETTINGS.TIMESTAMP_BOT_PX + captionH + PREVIEW_SETTINGS.CAPTION_BOT_PX;
      const bottomIn = bottomPx / dpi;

      const canvas = await renderStripCanvas(shots, (layout as any).count, dpi, filter, topIn, bottomIn, stripColor, roundedCorners);
      const ctx = canvas.getContext('2d')!;
      const fullLayout = computeStripLayout((layout as any).count, dpi, topIn, bottomIn);
      const logoW = logoRatio * brandHeight;
      const qrHeight = brandHeight * 0.6;
      const qrW = qrRatio * qrHeight;
      
      const logoX = fullLayout.contentX;
      const qrX = fullLayout.contentX + fullLayout.contentWidthPx - qrW;
      const logoY = PREVIEW_SETTINGS.BRAND_TOP_PX;
      const qrY = logoY + (brandHeight - qrHeight) / 2;

      if (logoImg) ctx.drawImage(logoImg, logoX, logoY, logoW, brandHeight);
      if (realQRImg) ctx.drawImage(realQRImg, qrX, qrY, qrW, qrHeight);
      
      // Draw Timestamp & Caption (Simplified for re-render)
      const date = new Date();
      const timeStr = date.toLocaleDateString('en-US', { month: 'short', day: '2-digit', year: 'numeric' }).toUpperCase() + " • " + date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit', hour12: true }).toUpperCase();
      const lastPhotoBottom = fullLayout.topCanvasPx + (fullLayout.photoHeightPx * (layout as any).count) + (fullLayout.gapPx * ((layout as any).count - 1));
      const timestampY = lastPhotoBottom + PREVIEW_SETTINGS.TIMESTAMP_TOP_PX;
      ctx.fillStyle = PREVIEW_SETTINGS.TIMESTAMP_COLOR; ctx.textAlign = 'center'; ctx.font = '700 12px Montserrat, sans-serif'; ctx.letterSpacing = '2px';
      ctx.fillText(timeStr, canvas.width / 2, timestampY + 10);
      ctx.fillStyle = PREVIEW_SETTINGS.CAPTION_COLOR; ctx.textAlign = 'center'; ctx.font = `${captionSize}px ${font}, system-ui, sans-serif`;
      const captionY = timestampY + timestampH + PREVIEW_SETTINGS.TIMESTAMP_BOT_PX;
      ctx.fillText(caption || ' ', canvas.width / 2, captionY + (captionSize * 0.8));

      const finalOutput = canvas.toDataURL('image/png');

      // 3. Upload
      const apiEndpoint = token ? getApiUrl('SAVE_STRIP') : getApiUrl('GUEST_SAVE');
      
      const uploadResponse = await fetch(apiEndpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          ...(token ? { 'Authorization': `Bearer ${token}` } : {})
        },
        body: JSON.stringify({
          id: stripId,
          image: finalOutput,
          title: token ? 'My Memory' : 'Guest Memory',
          caption: caption || 'Captured with Wuby'
        })
      });

      if (!uploadResponse.ok) throw new Error("Upload failed");
      const result = await uploadResponse.json();
      const finalId = result.id || stripId; // Use backend ID if available

      // 4. Update Store and Redirect
      photoboothStore.update(s => ({
        ...s,
        finalStrip: finalOutput,
        uploadedId: finalId,
        settings: { filter, stripColor, caption, captionSize, font, roundedCorners }
      }));

      // If guest, we should pass info to save page or handle there
      if (!token) {
        // For guest, save page will need the ID to show the link
        localStorage.setItem('last_guest_id', stripId);
      }
      
      goto('/photobooth/save');
    } catch (e: any) {
      console.error("Critical error during final render/upload:", e);
      uploadError = e.message || "Failed to process final memory";
      // We keep isUploading = true so the error UI remains visible
    }
  }

  function handleContinue() {
    if (previewUrl) {
      isConfirming = true;
    }
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
        class="border border-purple-200 hover:bg-purple-100/50 text-purple-500 px-3 md:px-6 py-1.5 md:py-2 rounded-lg text-[10px] md:text-xs font-bold transition-all active:scale-95 disabled:opacity-30"
        on:click={() => {
          photoboothStore.update(s => ({ ...s, shots: [] }));
          goto('/photobooth/capture');
        }}
        disabled={isUploading}
      >
        Retake
      </button>
      <button
        class="bg-purple-500 hover:bg-purple-600 text-white px-3 md:px-6 py-1.5 md:py-2 rounded-lg text-[10px] md:text-xs font-bold transition-all active:scale-95 shadow-md shadow-purple-200 disabled:opacity-50"
        on:click={handleContinue}
        disabled={!previewUrl || isUploading}
      >
        Continue
      </button>
    </div>
  </header>

  {#if isConfirming}
    <div class="fixed inset-0 z-[100] bg-purple-900/20 backdrop-blur-md flex items-center justify-center p-6 animate-in fade-in duration-300">
      <div class="bg-white w-full max-w-sm rounded-[2.5rem] p-8 shadow-2xl flex flex-col items-center text-center gap-6 animate-in slide-in-from-bottom-4">
        <div class="w-16 h-16 bg-purple-50 rounded-full flex items-center justify-center text-purple-500">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/></svg>
        </div>
        <div class="flex flex-col gap-2">
          <h2 class="text-xl font-medium text-slate-800">Finalize Memory?</h2>
          <p class="text-sm text-slate-400">We'll render your high-resolution strip with a scannable share link.</p>
        </div>
        <div class="flex flex-col w-full gap-3 mt-4">
           <button 
            on:click={handleConfirm}
            disabled={isUploading}
            class="w-full bg-purple-600 text-white py-4 rounded-2xl text-xs font-bold uppercase tracking-widest shadow-lg shadow-purple-200 hover:bg-purple-700 transition-all active:scale-95 disabled:opacity-50"
           >
            {isUploading ? 'Finalizing...' : 'Generate & Continue'}
           </button>
           <button 
            on:click={() => isConfirming = false}
            class="w-full text-slate-400 py-2 text-[10px] font-bold uppercase tracking-widest hover:text-slate-600 transition-colors"
           >
            Wait, I'm not done
           </button>
        </div>
      </div>
    </div>
  {/if}

  {#if isUploading}
    <div class="fixed inset-0 z-[100] bg-white flex flex-col items-center justify-center p-6 animate-in fade-in duration-500">
      <div class="w-16 h-16 border-4 border-purple-100 border-t-purple-500 rounded-full animate-spin mb-8"></div>
      <h2 class="text-2xl font-light text-purple-900 mb-2">Creating your digital memory...</h2>
      <p class="text-xs text-purple-300 font-bold uppercase tracking-[0.3em] animate-pulse">Please stay on this page</p>
      
      {#if uploadError}
        <div class="mt-8 p-6 bg-red-50 text-red-600 rounded-[2rem] text-xs font-bold uppercase tracking-widest max-w-xs text-center border-2 border-red-100 animate-in zoom-in duration-300">
          <p class="mb-4">{uploadError}</p>
          <button 
            on:click={() => { isUploading = false; uploadError = ''; }} 
            class="w-full bg-red-500 text-white py-3 rounded-xl hover:bg-red-600 transition-colors shadow-lg shadow-red-100"
          >
            Go Back
          </button>
        </div>
      {/if}
    </div>
  {/if}

  <main class="flex-1 flex flex-col md:flex-row overflow-hidden relative">
    <!-- Strip Preview Section - Minimalist -->
    <div class="flex-1 flex flex-col items-center justify-center p-4 md:p-12 bg-purple-100/10 overflow-hidden relative">
      {#if error}
        <div class="text-center">
          <p class="text-sm text-red-400 mb-4">{error}</p>
          <button class="text-xs font-bold text-purple-500 underline" on:click={() => location.reload()}>Try again</button>
        </div>
      {:else if previewUrl}
        <div class="h-full w-full max-h-[75vh] md:max-h-full relative flex items-center justify-center p-2">
          <img
            src={previewUrl}
            alt="Photobooth strip"
            class="max-h-full max-w-full object-contain transition-all duration-500 drop-shadow-[0_10px_30px_rgba(0,0,0,0.1)] md:hover:scale-[1.02]"
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
                on:click={() => { filter = f; updatePreview(); }}
              >
                {f.charAt(0).toUpperCase() + f.slice(1)}
              </button>
            {/each}
          </div>
        </section>

        <section>
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-[10px] font-bold uppercase text-purple-400 tracking-[0.2em]">Frame Color</h2>
          </div>
          <div class="flex flex-wrap gap-3">
            {#each PREVIEW_SETTINGS.STRIP_COLORS as color}
              <button
                class="w-8 h-8 rounded-full border-2 transition-all active:scale-90 {stripColor === color.id ? 'border-purple-500 ring-2 ring-purple-100 scale-110' : 'border-slate-200'}"
                style="background-color: {color.id}"
                title={color.label}
                on:click={() => { stripColor = color.id; updatePreview(); }}
              >
              </button>
            {/each}

            <!-- Custom Color Picker Trigger -->
            <button 
              class="w-8 h-8 rounded-full border-2 border-slate-200 flex items-center justify-center cursor-pointer hover:border-purple-300 transition-all active:scale-90 relative overflow-hidden {showColorWheel ? 'border-purple-500 ring-2 ring-purple-100' : ''}"
              style="background: conic-gradient(red, yellow, lime, aqua, blue, magenta, red);"
              on:click={() => showColorWheel = !showColorWheel}
              title="Custom Color Wheel"
            >
              <div class="w-3 h-3 bg-white rounded-full shadow-sm"></div>
            </button>
          </div>
        </section>

        <section>
          <div class="flex justify-between items-center mb-6">
            <h2 class="text-[10px] font-bold uppercase text-purple-400 tracking-[0.2em]">Photo Style</h2>
          </div>
          <button 
            class="w-full flex items-center justify-between p-4 bg-slate-50 rounded-2xl border border-transparent hover:border-purple-200 transition-all group"
            on:click={() => { roundedCorners = !roundedCorners; updatePreview(); }}
          >
            <span class="text-xs font-medium text-slate-600">Rounded Corners</span>
            <div class="w-10 h-5 rounded-full relative transition-colors {roundedCorners ? 'bg-purple-500' : 'bg-slate-200'}">
              <div class="absolute top-1 left-1 w-3 h-3 bg-white rounded-full transition-transform {roundedCorners ? 'translate-x-5' : 'translate-x-0'}"></div>
            </div>
          </button>
        </section>

        {#if showColorWheel}
            <div class="mt-4 p-4 bg-purple-50/50 rounded-3xl border border-purple-100 flex justify-center animate-in slide-in-from-top-2 duration-300">
              <ColorWheel 
                bind:value={stripColor}
                on:change={debouncedUpdate}
              />
            </div>
          {/if}

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
              <div class="grid grid-cols-5 gap-1.5">
                {#each PREVIEW_SETTINGS.AVAILABLE_CAPTION_SIZES as size}
                  <button
                    class="py-2.5 rounded-lg text-[10px] font-bold transition-all border {captionSize === size ? 'bg-purple-100/50 border-purple-300 text-purple-600' : 'bg-slate-50 border-transparent text-slate-400 hover:bg-slate-100 hover:text-slate-500'}"
                    on:click={() => { captionSize = size; updatePreview(); }}
                  >
                    {size}
                  </button>
                {/each}
              </div>
            </div>
          </div>
        </section>
      </div>
    </aside>
  </main>
</div>

<style>
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

