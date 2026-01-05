<script lang="ts">
  import { layoutStore, type LayoutOption } from '$lib/stores/layout.store';
  import { goto } from '$app/navigation';
  import { BOOTH_SETTINGS } from './settings';

  let selectedTimer = BOOTH_SETTINGS.DEFAULT_TIMER;
  let currentStep = 1;

  function nextStep() {
    currentStep = 2;
  }

  function prevStep() {
    currentStep = 1;
  }

  function select(count: 2 | 3 | 4) {
    const layout: LayoutOption = {
      count,
      width: BOOTH_SETTINGS.STRIP_WIDTH,
      height: BOOTH_SETTINGS.STRIP_WIDTH * count * BOOTH_SETTINGS.ASPECT_RATIO,
      timer: selectedTimer
    };

    layoutStore.set(layout);
    goto('/photobooth/capture');
  }
</script>

<div class="min-h-screen flex flex-col items-center justify-center p-6 bg-[#f8f2ff] overflow-hidden">
  <div class="text-center mb-12">
    <h1 class="text-4xl font-medium tracking-tight text-purple-900/80 mb-3">{BOOTH_SETTINGS.BRAND_NAME}</h1>
    <p class="text-purple-400/80 text-sm font-medium tracking-wide font-sans">Setup your session</p>
  </div>

  <!-- Step Indicator -->
  <div class="flex items-center gap-3 mb-12">
    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full {currentStep === 1 ? 'bg-purple-500 scale-125' : 'bg-purple-200'} transition-all duration-300"></div>
      <span class="text-[9px] font-bold uppercase tracking-widest {currentStep === 1 ? 'text-purple-600' : 'text-purple-300'}">Countdown</span>
    </div>
    <div class="w-8 h-px bg-purple-100"></div>
    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full {currentStep === 2 ? 'bg-purple-500 scale-125' : 'bg-purple-200'} transition-all duration-300"></div>
      <span class="text-[9px] font-bold uppercase tracking-widest {currentStep === 2 ? 'text-purple-600' : 'text-purple-300'}">Layout</span>
    </div>
  </div>

  <div class="w-full max-w-4xl flex flex-col items-center">
    {#if currentStep === 1}
      <!-- Timer Selection -->
      <div class="flex flex-col items-center gap-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
        <div class="flex flex-col items-center gap-6">
          <h2 class="text-[10px] font-bold uppercase text-purple-300 tracking-[0.3em]">Choose Pacing</h2>
          <div class="flex p-1 bg-purple-100/50 rounded-2xl border border-purple-200/50">
            {#each BOOTH_SETTINGS.TIMER_OPTIONS as t}
              <button
                class="px-6 md:px-10 py-2.5 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all duration-300 {selectedTimer === t ? 'bg-white text-purple-500 shadow-sm' : 'text-purple-400 hover:text-purple-500'}"
                on:click={() => selectedTimer = t}
              >
                {t === 0 ? 'Manual' : t + 's'}
              </button>
            {/each}
          </div>
          <p class="text-[10px] text-purple-300/80 italic max-w-[240px] text-center leading-relaxed">
            {#if selectedTimer === 0}
              Manual mode: Trigger each photo yourself for perfect timing.
            {:else}
              Automatic: We'll count down {selectedTimer}s between each shot.
            {/if}
          </p>
        </div>

        <button 
          class="mt-8 px-10 py-3 rounded-full bg-purple-500 text-white text-[10px] font-bold uppercase tracking-[0.2em] transition-all hover:bg-purple-600 active:scale-95 shadow-lg shadow-purple-200"
          on:click={nextStep}
        >
          Select Layout
        </button>
      </div>
    {:else}
      <!-- Layout Selection -->
      <div class="w-full flex flex-col items-center gap-6 md:gap-10 animate-in fade-in slide-in-from-bottom-8 duration-700">
        <h2 class="text-[10px] font-bold uppercase text-purple-300 tracking-[0.3em]">Choose Frames</h2>
        
        <div class="flex flex-col md:flex-row gap-2 md:gap-12 justify-center w-full px-4">
          {#each BOOTH_SETTINGS.LAYOUT_OPTIONS as count}
            <button 
              class="group flex flex-row md:flex-col items-center gap-4 md:gap-6 transition-all duration-500 active:scale-95 p-2.5 md:p-0 bg-white/60 md:bg-transparent border border-purple-100/50 md:border-none rounded-xl md:rounded-none shadow-sm md:shadow-none" 
              on:click={() => select(count as 2|3|4)}
            >
              <div class="w-10 h-20 md:w-32 md:h-64 bg-white border border-purple-200/60 rounded-lg md:rounded-xl flex flex-col overflow-hidden shadow-sm group-hover:shadow-md transition-all p-1 md:p-2.5 gap-0.5 md:gap-2">
                {#each Array(count) as _, i}
                  <div class="flex-1 bg-purple-100/40 rounded flex items-center justify-center">
                    <svg class="w-2.5 h-2.5 md:w-8 md:h-8 text-purple-300/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                  </div>
                {/each}
              </div>
              
              <div class="text-left md:text-center flex-1">
                <p class="text-[9px] md:text-[11px] font-bold text-purple-500 uppercase tracking-[0.2em] group-hover:text-purple-600 transition-colors">{count} photos</p>
                <p class="md:hidden text-[7px] text-purple-400 font-bold uppercase tracking-[0.1em] mt-0.5 opacity-60">Classic Strip</p>
                <div class="hidden md:block mt-3 h-px w-4 bg-purple-200 mx-auto group-hover:w-8 group-hover:bg-purple-300 transition-all"></div>
              </div>
            </button>
          {/each}
        </div>

        <button 
          class="mt-4 text-[9px] font-bold text-purple-300 uppercase tracking-[0.2em] hover:text-purple-500 transition-colors flex items-center gap-2"
          on:click={prevStep}
        >
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7"/></svg>
          Back to Countdown
        </button>
      </div>
    {/if}
  </div>
</div>

