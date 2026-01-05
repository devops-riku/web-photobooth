<script lang="ts">
  import { layoutStore, type LayoutOption } from '$lib/stores/layout.store';
  import { goto } from '$app/navigation';
  import { BOOTH_SETTINGS } from './settings';

  let selectedTimer = BOOTH_SETTINGS.DEFAULT_TIMER;

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
  <div class="text-center mb-16">
    <h1 class="text-4xl font-medium tracking-tight text-purple-900/80 mb-3">{BOOTH_SETTINGS.BRAND_NAME}</h1>
    <p class="text-purple-400/80 text-sm font-medium tracking-wide font-sans">Setup your session</p>
  </div>

  <!-- Timer Selection -->
  <div class="flex flex-col items-center gap-6 mb-20 animate-in fade-in slide-in-from-bottom-4 duration-700">
    <h2 class="text-[10px] font-bold uppercase text-purple-300 tracking-[0.3em]">Countdown</h2>
    <div class="flex p-1 bg-purple-100/50 rounded-2xl border border-purple-200/50">
      {#each BOOTH_SETTINGS.TIMER_OPTIONS as t}
        <button
          class="px-8 py-2 rounded-xl text-[11px] font-bold uppercase tracking-widest transition-all duration-300 {selectedTimer === t ? 'bg-white text-purple-500 shadow-sm' : 'text-purple-400 hover:text-purple-500'}"
          on:click={() => selectedTimer = t}
        >
          {t}s
        </button>
      {/each}
    </div>
  </div>

  <div class="w-full max-w-4xl px-4 flex flex-col items-center gap-8">
    <h2 class="text-[10px] font-bold uppercase text-purple-300 tracking-[0.3em]">Select Layout</h2>
    <div class="flex flex-col md:flex-row gap-6 md:gap-12 justify-center w-full animate-in fade-in slide-in-from-bottom-8 duration-1000">
      {#each BOOTH_SETTINGS.LAYOUT_OPTIONS as count}
        <button 
          class="group flex flex-row md:flex-col items-center gap-6 transition-all duration-500 active:scale-95 p-4 md:p-0" 
          on:click={() => select(count as 2|3|4)}
        >
          <div class="w-20 h-40 md:w-32 md:h-64 bg-white border border-purple-200/60 rounded-xl flex flex-col overflow-hidden shadow-[0_4px_12px_rgba(159,122,234,0.06)] group-hover:shadow-[0_20px_40px_rgba(159,122,234,0.12)] group-hover:border-purple-300 transition-all p-1.5 md:p-2.5 gap-1 md:gap-2">
            {#each Array(count) as _, i}
              <div class="flex-1 bg-purple-100/40 rounded-lg flex items-center justify-center">
                 <svg class="w-4 h-4 md:w-8 md:h-8 text-purple-300/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                 </svg>
              </div>
            {/each}
          </div>
          
          <div class="text-left md:text-center flex-1">
            <p class="text-[10px] md:text-[11px] font-bold text-purple-400 uppercase tracking-[0.3em] group-hover:text-purple-600 transition-colors">{count} photos</p>
            <div class="hidden md:block mt-3 h-px w-4 bg-purple-200 mx-auto group-hover:w-8 group-hover:bg-purple-400 transition-all"></div>
          </div>
        </button>
      {/each}
    </div>
  </div>
</div>

