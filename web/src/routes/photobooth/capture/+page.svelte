<script lang="ts">
  import { onMount, tick } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';

  let video: HTMLVideoElement;
  let canvas: HTMLCanvasElement;

  let countdown = 3;
  let shooting = false;
  let reviewing = false;
  let flash = false;
  let currentReviewIndex = 0;
  let maxShots = 0;
  let currentShots = 0;

  let layout: any = null;
  let shots: string[] = [];

  // ✅ subscribe without navigation
  layoutStore.subscribe(v => {
    layout = v;
    maxShots = v?.count ?? 0;
  });

  photoboothStore.subscribe(v => {
    shots = v.shots;
  });

  let stream: MediaStream | null = null;
  
  // ✅ Re-attach stream whenever video element is created/recreated
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
        video: { facingMode: 'user', width: { ideal: 1080 }, height: { ideal: 1080 } }
      });
    } catch (err) {
      console.error("Camera access error:", err);
    }
  }

  async function startSequence() {
    if (shooting || reviewing) return;
    shooting = true;
    reviewing = false;
    currentShots = 0;
    
    // Clear previous shots
    photoboothStore.set({ shots: [] });

    while (currentShots < maxShots) {
      await runCountdown();
      capture();
      currentShots++;
    }

    shooting = false;
    reviewing = true;
  }

  async function retakeSingle(index: number) {
    reviewing = false;
    await tick();
    shooting = true;
    
    await runCountdown();
    
    // Flash effect
    flash = true;
    setTimeout(() => flash = false, 100);

    const ctx = canvas.getContext('2d')!;
    const size = Math.min(video.videoWidth, video.videoHeight);
    canvas.width = size;
    canvas.height = size;
    const sx = (video.videoWidth - size) / 2;
    const sy = (video.videoHeight - size) / 2;

    ctx.drawImage(video, sx, sy, size, size, 0, 0, size, size);
    const image = canvas.toDataURL('image/png');

    photoboothStore.update(v => {
      const newShots = [...v.shots];
      newShots[index] = image;
      return { shots: newShots };
    });

    shooting = false;
    reviewing = true;
  }

  async function retake() {
    reviewing = false;
    await tick();
    currentShots = 0;
    photoboothStore.set({ shots: [] });
  }

  function proceed() {
    goto('/photobooth/preview');
  }


  function runCountdown() {
    return new Promise<void>((resolve) => {
      countdown = 3;
      const interval = setInterval(() => {
        countdown--;
        if (countdown === 0) {
          clearInterval(interval);
          resolve();
        }
      }, 1000);
    });
  }

  function capture() {
    // Flash effect
    flash = true;
    setTimeout(() => flash = false, 100);

    const ctx = canvas.getContext('2d')!;
    
    // We want a square capture to match the preview
    // Use the smaller dimension to determine the square size
    const size = Math.min(video.videoWidth, video.videoHeight);
    canvas.width = size;
    canvas.height = size;

    // Center crop from the video
    const sx = (video.videoWidth - size) / 2;
    const sy = (video.videoHeight - size) / 2;

    ctx.drawImage(
      video,
      sx, sy, size, size, // source (centered square)
      0, 0, size, size    // destination
    );

    const image = canvas.toDataURL('image/png');

    photoboothStore.update(v => ({
      shots: [...v.shots, image]
    }));
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

<div class="h-screen flex flex-col bg-white overflow-hidden">
  <!-- Status Bar -->
  <header class="p-6 flex justify-between items-center">
    <div class="flex items-center gap-2">
      <div class="w-1.5 h-1.5 rounded-full bg-purple-400"></div>
      <span class="text-[10px] font-bold tracking-[0.2em] text-purple-400 uppercase">Booth Session</span>
    </div>
    
    <div class="flex gap-1.5">
      {#each Array(maxShots) as _, i}
        <div class="w-8 h-1 rounded-full transition-all duration-700 {i < currentShots ? 'bg-purple-300' : 'bg-purple-50'}"></div>
      {/each}
    </div>
  </header>

  <main class="flex-1 flex flex-col items-center justify-center p-4 md:p-6 gap-8 md:gap-12 overflow-hidden bg-white">
    {#if !reviewing}
      <!-- Camera Viewport - Only visible when not reviewing -->
      <div class="relative w-full max-w-[320px] md:max-w-[420px] aspect-square animate-in zoom-in-95 duration-500">
        <div class="absolute inset-0 video-container rounded-[2rem] md:rounded-[2.5rem] overflow-hidden border border-purple-50 shadow-[0_8px_30px_rgb(159,122,234,0.05)] bg-purple-50/10">
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
      <div class="w-full flex flex-col items-center gap-6 md:gap-10 animate-in fade-in slide-in-from-bottom-8 duration-1000">
        <div class="flex gap-4 md:gap-6 px-6 md:px-12 py-4 overflow-x-auto no-scrollbar snap-x max-w-full">
          {#each shots as shot, i}
            <div class="flex-shrink-0 w-48 md:w-64 aspect-square relative group/shot snap-center">
              <img 
                src={shot} 
                alt="Capture {i + 1}" 
                class="w-full h-full object-cover rounded-2xl md:rounded-3xl shadow-[0_15px_40px_rgba(159,122,234,0.1)] border-4 border-white transition-all duration-500 md:group-hover/shot:scale-[1.05] md:group-hover/shot:-rotate-1" 
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
                <span class="text-[8px] md:text-[9px] font-bold text-purple-300/60 uppercase tracking-[0.3em]">Photo 0{i + 1}</span>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    <!-- Simple Trigger / Decision Controls -->
    <div class="flex flex-col items-center gap-6 md:gap-8">
      {#if !reviewing}
        <div class="flex flex-col items-center gap-4">
          <button
            class="w-16 h-16 md:w-20 md:h-20 rounded-full border-2 border-purple-50 p-1 md:p-1.5 transition-all active:scale-90 hover:border-purple-100 disabled:opacity-30 disabled:scale-100 shadow-sm"
            on:click={startSequence}
            disabled={shooting}
          >
            <div class="w-full h-full rounded-full bg-purple-50/50 flex items-center justify-center group border border-purple-100">
              <div class="w-3 h-3 md:w-4 md:h-4 rounded-full bg-purple-400 {shooting ? 'animate-ping' : ''}"></div>
            </div>
          </button>
          
          <p class="text-[9px] md:text-[10px] font-bold tracking-[0.25em] uppercase text-purple-200">
            {#if shooting}
              Capture in progress
            {:else}
              Start Sequence
            {/if}
          </p>
        </div>
      {:else}
        <div class="flex flex-col items-center gap-6 md:gap-8 animate-in fade-in slide-in-from-bottom-4 duration-1000 delay-300">
          <div class="flex flex-col md:flex-row gap-3 md:gap-4 px-6 w-full md:w-auto">
            <button
              class="w-full md:w-auto px-10 py-3.5 md:py-4 rounded-full border border-purple-100 text-purple-300 text-[10px] md:text-[11px] font-bold uppercase tracking-[0.2em] hover:bg-purple-50 hover:text-purple-500 transition-all active:scale-95"
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
          <p class="text-[9px] md:text-[10px] font-bold text-purple-200 uppercase tracking-[0.3em]">Tap a photo to swap or proceed</p>
        </div>
      {/if}
    </div>
  </main>
</div>

<canvas bind:this={canvas} class="hidden"></canvas>

