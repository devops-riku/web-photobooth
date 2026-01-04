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

  <main class="flex-1 flex flex-col items-center justify-center p-6 gap-12">
    <!-- Camera Viewport - Minimalist -->
    <div class="relative w-full max-w-[400px] aspect-square">
      <div class="absolute inset-0 video-container rounded-3xl overflow-hidden border border-purple-50 shadow-sm bg-purple-50/10">
        <video
          bind:this={video}
          autoplay
          playsinline
        ></video>

        {#if shooting && countdown > 0}
          <div class="absolute inset-0 flex items-center justify-center bg-white/60 backdrop-blur-[2px] transition-all duration-300">
            <span class="text-9xl font-light text-purple-900/60 leading-none">
              {countdown}
            </span>
          </div>
        {/if}
      </div>
    </div>

    <!-- Simple Trigger -->
    <div class="flex flex-col items-center gap-6">
      <button
        class="w-16 h-16 rounded-full border border-purple-100 p-1 transition-all active:scale-90 hover:border-purple-200 disabled:opacity-30 disabled:scale-100"
        on:click={startSequence}
        disabled={shooting}
      >
        <div class="w-full h-full rounded-full bg-purple-50 flex items-center justify-center group">
          <div class="w-3 h-3 rounded-full bg-purple-400 {shooting ? 'animate-ping' : ''}"></div>
        </div>
      </button>
      
      <p class="text-[11px] font-medium tracking-widest uppercase text-purple-300">
        {#if shooting}
          Look at the camera
        {:else}
          Tap to capture
        {/if}
      </p>
    </div>
  </main>
</div>

<canvas bind:this={canvas} class="hidden"></canvas>

