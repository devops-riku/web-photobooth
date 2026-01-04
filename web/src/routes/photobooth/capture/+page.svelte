<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { layoutStore } from '$lib/stores/layout.store';
  import { photoboothStore } from '$lib/stores/photobooth.store';

  let video: HTMLVideoElement;
  let canvas: HTMLCanvasElement;

  let countdown = 3;
  let shooting = false;
  let maxShots = 0;
  let currentShots = 0;

  let layout: any = null;

  // ✅ subscribe without navigation
  layoutStore.subscribe(v => {
    layout = v;
    maxShots = v?.count ?? 0;
  });

  onMount(async () => {
    // ✅ client-side redirect only
    if (!layout) {
      goto('/photobooth');
      return;
    }

    const stream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'user' }
    });
    video.srcObject = stream;
  });

  async function startSequence() {
    if (shooting) return;
    shooting = true;
    currentShots = 0;
    
    // Clear previous shots
    photoboothStore.set({ shots: [] });

    while (currentShots < maxShots) {
      await runCountdown();
      capture();
      currentShots++;
    }

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
    border-radius: 1rem;
  }
  
  video {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    min-width: 100%;
    min-height: 100%;
    object-fit: cover;
  }
</style>

<div class="h-screen flex flex-col overflow-hidden bg-slate-900 text-white">
  <header class="p-4 shrink-0 flex justify-between items-center bg-black/20">
    <h1 class="text-lg font-bold">Capture</h1>
    <div class="text-sm font-mono opacity-60">
      Shot {currentShots + 1} / {maxShots}
    </div>
  </header>

  <main class="flex-1 flex flex-col items-center justify-center p-6 gap-8">
    <div class="relative w-full max-w-[320px] aspect-square">
      <div class="absolute inset-0 video-container rounded-2xl overflow-hidden border-4 border-white/20 shadow-2xl bg-black">
        <video
          bind:this={video}
          autoplay
          playsinline
          class="w-full h-full object-cover"
        ></video>

        {#if shooting && countdown > 0}
          <div class="absolute inset-0 flex items-center justify-center bg-black/40 backdrop-blur-[2px] text-white text-8xl font-black drop-shadow-2xl">
            {countdown}
          </div>
        {/if}
      </div>
      
      <!-- Decorative corners -->
      <div class="absolute -top-2 -left-2 w-8 h-8 border-t-4 border-l-4 border-primary-500 rounded-tl-lg"></div>
      <div class="absolute -top-2 -right-2 w-8 h-8 border-t-4 border-r-4 border-primary-500 rounded-tr-lg"></div>
      <div class="absolute -bottom-2 -left-2 w-8 h-8 border-b-4 border-l-4 border-primary-500 rounded-bl-lg"></div>
      <div class="absolute -bottom-2 -right-2 w-8 h-8 border-b-4 border-r-4 border-primary-500 rounded-br-lg"></div>
    </div>

    <div class="flex flex-col items-center gap-6">
      <button
        class="w-20 h-20 rounded-full border-8 border-white/20 p-1 transition-all active:scale-90 disabled:opacity-50 disabled:scale-100"
        on:click={startSequence}
        disabled={shooting}
      >
        <div class="w-full h-full rounded-full bg-white shadow-inner flex items-center justify-center">
          <div class="w-4 h-4 rounded-full bg-primary-600 animate-pulse {shooting ? 'opacity-0' : ''}"></div>
        </div>
      </button>
      
      <p class="text-sm font-medium text-white/50 animate-pulse">
        {shooting ? 'Smile!' : 'Ready when you are'}
      </p>
    </div>
  </main>
</div>

<canvas bind:this={canvas} class="hidden"></canvas>

