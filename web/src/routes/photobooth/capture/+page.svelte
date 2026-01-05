<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';
  import { CAPTURE_SETTINGS } from './settings';

  let video: HTMLVideoElement;
  let canvas: HTMLCanvasElement;

  let countdown = 3; // Default, will be sync'd with layout
  let shooting = false;
  let reviewing = false;
  let flash = false;
  let currentReviewIndex = 0;
  let maxShots = 0;
  let currentShots = 0;
  let swappingIndex: number | null = null;
  let busy = false;

  let layout: any = null;
  let shots: string[] = [];

  // subscribe without navigation
  layoutStore.subscribe(v => {
    layout = v;
    maxShots = v?.count ?? 0;
    if (v) countdown = v.timer;
  });

  photoboothStore.subscribe(v => {
    shots = v.shots;
  });

  let stream: MediaStream | null = null;
  
  // Re-attach stream whenever video element is created/recreated
  $: if (video && stream) {
    video.srcObject = stream;
  }

  onMount(async () => {
    if (!layout) {
      goto('/photobooth');
      return;
    }
    await initCamera();
  });

  async function initCamera() {
    try {
      stream = await navigator.mediaDevices.getUserMedia({
        video: { 
          facingMode: 'user', 
          width: { ideal: CAPTURE_SETTINGS.CAMERA_IDEAL_WIDTH }, 
          height: { ideal: CAPTURE_SETTINGS.CAMERA_IDEAL_HEIGHT } 
        }
      });
    } catch (err) {
      console.error("Camera access error:", err);
    }
  }

  function waitForVideo(): Promise<void> {
    return new Promise((resolve) => {
      if (video && video.readyState >= 2) return resolve();
      const check = setInterval(() => {
        if (video && video.readyState >= 2) {
          clearInterval(check);
          resolve();
        }
      }, 50);
    });
  }

  async function handleTrigger() {
    if (reviewing || busy) return;

    if (swappingIndex !== null) {
      busy = true;
      await runCountdown();
      captureInternal(swappingIndex);
      swappingIndex = null;
      shooting = false;
      reviewing = true;
      busy = false;
      return;
    }

    if (layout?.timer === 0) {
      if (!shooting) {
        shooting = true;
        currentShots = 0;
        photoboothStore.set({ shots: [] });
        await tick();
        await waitForVideo();
      }

      capture();
      currentShots++;

      if (currentShots >= maxShots) {
        shooting = false;
        reviewing = true;
      }
    } else {
      await startSequence();
    }
  }

  async function startSequence() {
    if (shooting || reviewing || busy) return;
    busy = true;
    shooting = true;
    reviewing = false;
    currentShots = 0;
    
    // Clear previous shots
    photoboothStore.set({ shots: [] });
    await tick();
    await waitForVideo();

    while (currentShots < maxShots) {
      await runCountdown();
      capture();
      currentShots++;
    }

    shooting = false;
    reviewing = true;
    busy = false;
  }

  async function retakeSingle(index: number) {
    reviewing = false;
    await tick();
    shooting = true;
    await waitForVideo();
    swappingIndex = index;
  }

  function captureInternal(index: number) {
     // Flash effect
    flash = true;
    setTimeout(() => flash = false, CAPTURE_SETTINGS.FLASH_DURATION_MS);

    const ctx = canvas.getContext('2d')!;
    const size = Math.min(video.videoWidth, video.videoHeight);
    canvas.width = size;
    canvas.height = size;

    ctx.save();
    
    // If mirroring is enabled, flip the context
    if (CAPTURE_SETTINGS.MIRROR_CAPTURE) {
      ctx.translate(size, 0);
      ctx.scale(-1, 1);
    }

    const sx = (video.videoWidth - size) / 2;
    const sy = (video.videoHeight - size) / 2;

    ctx.drawImage(video, sx, sy, size, size, 0, 0, size, size);
    
    ctx.restore();
    
    const image = canvas.toDataURL('image/png');

    photoboothStore.update(v => {
      const newShots = [...v.shots];
      newShots[index] = image;
      return { shots: newShots };
    });
  }

  async function retake() {
    reviewing = false;
    await tick();
    currentShots = 0;
    photoboothStore.set({ shots: [] });
    shooting = false;
    swappingIndex = null;
    busy = false;
  }

  function proceed() {
    goto('/photobooth/preview');
  }


  function runCountdown() {
    return new Promise<void>((resolve) => {
      if (layout.timer === 0) {
        countdown = 0;
        resolve();
        return;
      }

      countdown = layout.timer;
      const interval = setInterval(() => {
        countdown--;
        if (countdown <= 0) {
          clearInterval(interval);
          resolve();
        }
      }, 1000);
    });
  }

  function capture() {
     // Flash effect
    flash = true;
    setTimeout(() => flash = false, CAPTURE_SETTINGS.FLASH_DURATION_MS);

    const ctx = canvas.getContext('2d')!;
    
    // We want a square capture to match the preview
    // Use the smaller dimension to determine the square size
    const size = Math.min(video.videoWidth, video.videoHeight);
    canvas.width = size;
    canvas.height = size;

    ctx.save();

    // If mirroring is enabled, flip the context
    if (CAPTURE_SETTINGS.MIRROR_CAPTURE) {
      ctx.translate(size, 0);
      ctx.scale(-1, 1);
    }

    // Center crop from the video
    const sx = (video.videoWidth - size) / 2;
    const sy = (video.videoHeight - size) / 2;

    ctx.drawImage(
      video,
      sx, sy, size, size, // source (centered square)
      0, 0, size, size    // destination
    );

    ctx.restore();

    const image = canvas.toDataURL('image/png');

    photoboothStore.update(v => ({
      shots: [...v.shots, image]
    }));
  }

  let fileInput: HTMLInputElement;
  let uploadMessage = '';
  let uploadMessageType: 'success' | 'error' = 'success';

  async function handleFileUpload(e: Event) {
    const files = (e.target as HTMLInputElement).files;
    if (!files || files.length === 0) return;

    busy = true;
    let addedCount = 0;
    
    for (let i = 0; i < files.length; i++) {
      if (shots.length >= maxShots) break;

      const file = files[i];
      
      // Strict photo validation
      if (!file.type.startsWith('image/')) {
        uploadMessage = `"${file.name}" is not a photo`;
        uploadMessageType = 'error';
        continue;
      }

      const dataUrl = await new Promise<string>((resolve) => {
        const reader = new FileReader();
        reader.onload = (e) => resolve(e.target?.result as string);
        reader.readAsDataURL(file);
      });

      const img = new Image();
      img.src = dataUrl;
      await new Promise(r => img.onload = r);

      const ctx = canvas.getContext('2d')!;
      const size = Math.min(img.width, img.height);
      canvas.width = size;
      canvas.height = size;

      const sx = (img.width - size) / 2;
      const sy = (img.height - size) / 2;

      ctx.clearRect(0,0,size,size);
      ctx.drawImage(img, sx, sy, size, size, 0, 0, size, size);
      
      const processedShot = canvas.toDataURL('image/png');
      photoboothStore.update(v => {
        const newShots = [...v.shots, processedShot];
        shots = newShots; // Sync local ref for loop check
        currentShots = newShots.length;
        return { shots: newShots };
      });
      addedCount++;
    }

    if (addedCount > 0) {
      uploadMessage = `Successfully added ${addedCount} photo${addedCount > 1 ? 's' : ''}!`;
      uploadMessageType = 'success';
    }

    setTimeout(() => {
      uploadMessage = '';
    }, 3000);

    if (shots.length >= maxShots) {
      shooting = false;
      reviewing = true;
    }
    busy = false;
    // reset input
    (e.target as HTMLInputElement).value = '';
  }
</script>

<style>
  .video-container {
    aspect-ratio: 1 / 1;
    overflow: hidden;
    position: relative;
  }
  
  video {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
    transform: scaleX(-1);
  }
</style>

<div class="h-screen flex flex-col bg-[#f8f2ff] overflow-hidden">
  <!-- Status Bar -->
  <header class="p-6 flex justify-between items-center">
    <!-- Back Button -->
    <button 
      on:click={() => goto('/photobooth')}
      class="text-purple-400 hover:text-purple-600 transition-colors flex items-center gap-2 group"
    >
      <svg class="w-5 h-5 transition-transform group-hover:-translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
      </svg>
      <span class="text-[10px] font-bold uppercase tracking-widest hidden md:inline">Exit</span>
    </button>
    
    <div class="flex gap-1.5">
      {#each Array(maxShots) as _, i}
        <div class="w-8 h-1 rounded-full transition-all duration-700 {i < currentShots ? 'bg-purple-400' : 'bg-purple-100'}"></div>
      {/each}
    </div>

    <div class="flex items-center gap-2">
      <div class="w-1.5 h-1.5 rounded-full bg-purple-400"></div>
      <span class="text-[10px] font-bold tracking-[0.2em] text-purple-400 uppercase hidden md:inline">Booth Session</span>
    </div>
  </header>

  <main class="flex-1 flex flex-col items-center justify-center p-4 md:p-6 gap-8 md:gap-12 overflow-hidden bg-transparent">
    {#if !reviewing}
      <!-- Camera Viewport - Only visible when not reviewing -->
      <div class="relative w-full max-w-[320px] md:max-w-[420px] aspect-square">
        <div class="absolute inset-0 video-container rounded-[2rem] md:rounded-[2.5rem] overflow-hidden border border-purple-200/50 shadow-[0_8px_30px_rgb(159,122,234,0.08)] bg-purple-100/20">
          {#if flash}
            <div class="absolute inset-0 bg-white z-[60] pointer-events-none transition-opacity duration-300"></div>
          {/if}

          <video
            bind:this={video}
            autoplay
            playsinline
            class="w-full h-full object-cover"
          ></video>

          {#if shooting && countdown > 0}
            <div class="absolute inset-0 flex items-center justify-center bg-white/60 backdrop-blur-[2px] z-50 transition-all duration-300">
              <span class="text-7xl md:text-9xl font-light text-purple-900/60 leading-none">
                {countdown}
              </span>
            </div>
          {/if}
        </div>
      </div>
    {:else}
      <!-- Sequence Review - Horizontal Lineup (not trapped in a card) -->
      <div class="w-full flex flex-col items-center gap-6 md:gap-10">
        <div class="flex gap-4 md:gap-6 px-6 md:px-12 py-4 overflow-x-auto no-scrollbar snap-x max-w-full">
          {#each shots as shot, i}
            <div class="flex-shrink-0 w-48 md:w-64 aspect-square relative group/shot snap-center">
              <img 
                src={shot} 
                alt="Capture {i + 1}" 
                class="w-full h-full object-cover rounded-2xl md:rounded-3xl shadow-[0_15px_40px_rgba(159,122,234,0.12)] border-4 border-white transition-all duration-500 md:group-hover/shot:scale-[1.05] md:group-hover/shot:-rotate-1" 
              />
              
              <!-- Retake overlay - Simplified for mobile (visible by default or on tap) -->
              <div class="absolute inset-0 bg-purple-900/40 opacity-0 group-hover/shot:opacity-100 md:opacity-0 transition-all duration-300 rounded-2xl md:rounded-3xl flex flex-col items-center justify-center gap-3">
                <button 
                  on:click={() => retakeSingle(i)}
                  class="bg-white text-purple-600 px-4 md:px-5 py-2 md:py-2.5 rounded-full text-[9px] md:text-[10px] font-bold uppercase tracking-widest shadow-xl transform active:scale-95 transition-all"
                >
                  Swap
                </button>
              </div>
              
              <div class="absolute -top-6 left-0 right-0 text-center">
                <span class="text-[8px] md:text-[9px] font-bold text-purple-300 uppercase tracking-[0.3em]">Photo 0{i + 1}</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Simple Trigger / Decision Controls -->
    <div class="flex flex-col items-center gap-6 md:gap-8">
      {#if !reviewing}
        <div class="flex flex-col items-center gap-6">
          <div class="flex items-center gap-12">
            <!-- Upload Button (New) -->
            <div class="flex flex-col items-center gap-2">
              <button 
                on:click={() => fileInput.click()}
                disabled={busy}
                class="w-12 h-12 md:w-14 md:h-14 rounded-full border border-purple-200 flex items-center justify-center text-purple-400 hover:bg-purple-50 transition-all active:scale-95 disabled:opacity-30"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
                </svg>
              </button>
              <span class="text-[8px] font-bold uppercase tracking-widest text-purple-300">Upload</span>
            </div>

            <!-- Main Trigger -->
            <button
              class="w-16 h-16 md:w-20 md:h-20 rounded-full border-2 border-purple-100 p-1 md:p-1.5 transition-all active:scale-90 hover:border-purple-200 disabled:opacity-30 disabled:scale-100 shadow-sm"
              on:click={handleTrigger}
              disabled={busy}
            >
              <div class="w-full h-full rounded-full bg-purple-200/20 flex items-center justify-center group border border-purple-200/50">
                <div class="w-3 h-3 md:w-4 md:h-4 rounded-full bg-purple-500 {busy && layout?.timer !== 0 ? 'animate-ping' : ''}"></div>
              </div>
            </button>

            <!-- Spacer for symmetry -->
            <div class="w-12 h-12 md:w-14 md:h-14 opacity-0 pointer-events-none"></div>
          </div>
          
          <input 
            type="file" 
            bind:this={fileInput} 
            on:change={handleFileUpload} 
            accept="image/*" 
            multiple 
            class="hidden" 
          />

          <p class="text-[9px] md:text-[10px] font-bold tracking-[0.25em] uppercase text-purple-300 transition-all">
            {#if uploadMessage}
              <span class={uploadMessageType === 'error' ? 'text-red-400' : 'text-purple-500'}>{uploadMessage}</span>
            {:else if busy && layout?.timer !== 0}
              Capture in progress
            {:else if swappingIndex !== null}
              {layout?.timer > 0 ? 'Start Countdown' : 'Capture Swap'}
            {:else if shooting && layout?.timer === 0}
              Capture {currentShots + 1} of {maxShots}
            {:else}
              {layout?.timer === 0 ? 'Capture First Shot' : 'Start Sequence'}
            {/if}
          </p>
        </div>
      {:else}
        <div class="flex flex-col items-center gap-6 md:gap-8">
          <div class="flex flex-col md:flex-row gap-3 md:gap-4 px-6 w-full md:w-auto">
            <button
              class="w-full md:w-auto px-10 py-3.5 md:py-4 rounded-full border border-purple-200 text-purple-400 text-[10px] md:text-[11px] font-bold uppercase tracking-[0.2em] hover:bg-purple-100/50 hover:text-purple-600 transition-all active:scale-95"
              on:click={retake}
            >
              Retake All
            </button>
            <button
              class="w-full md:w-auto px-12 py-3.5 md:py-4 rounded-full bg-purple-500 text-white text-[10px] md:text-[11px] font-bold uppercase tracking-[0.2em] hover:bg-purple-600 shadow-[0_10px_30px_rgba(168,85,247,0.3)] transition-all active:scale-95"
              on:click={proceed}
            >
              Continue to Edit
            </button>
          </div>
          <p class="text-[9px] md:text-[10px] font-bold text-purple-300 uppercase tracking-[0.3em]">Tap a photo to swap or proceed</p>
        </div>
      {/if}
    </div>
  </main>
</div>

<canvas bind:this={canvas} class="hidden"></canvas>

