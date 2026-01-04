<script lang="ts">
  import { layoutStore, type LayoutOption } from '$lib/stores/layout.store';
  import { goto } from '$app/navigation';
  import { BOOTH_SETTINGS } from './settings';

  function select(count: 2 | 3 | 4) {
    const layout: LayoutOption = {
      count,
      width: BOOTH_SETTINGS.STRIP_WIDTH,
      height: BOOTH_SETTINGS.STRIP_WIDTH * count * BOOTH_SETTINGS.ASPECT_RATIO
    };

    layoutStore.set(layout);
    goto('/photobooth/capture');
  }
</script>

<div class="min-h-screen flex flex-col items-center justify-start md:justify-center p-6 bg-white">
  <div class="text-center mt-8 mb-12 md:mt-0">
    <h1 class="text-4xl font-medium tracking-tight text-purple-900/80 mb-3">{BOOTH_SETTINGS.BRAND_NAME}</h1>
    <p class="text-purple-400/80 text-sm font-medium tracking-wide">Select your layout</p>
  </div>

  <div class="flex flex-col md:flex-row gap-6 md:gap-8 max-w-4xl w-full justify-center px-4 mb-12">
    {#each BOOTH_SETTINGS.LAYOUT_OPTIONS as count}
      <button 
        class="group flex flex-row md:flex-col items-center gap-6 transition-all duration-400 active:scale-95 bg-purple-50/20 md:bg-transparent p-4 md:p-0 rounded-3xl" 
        on:click={() => select(count as 2|3|4)}
      >
        <div class="w-20 h-40 md:w-32 md:h-64 bg-white border border-purple-100 rounded-xl flex flex-col overflow-hidden shadow-[0_4px_12px_rgba(159,122,234,0.05)] group-hover:shadow-[0_12px_32px_rgba(159,122,234,0.12)] group-hover:border-purple-200 transition-all p-1.5 md:p-2 gap-1 md:gap-1.5">
          {#each Array(count) as _, i}
            <div class="flex-1 bg-purple-50/50 rounded-md flex items-center justify-center">
               <svg class="w-4 h-4 md:w-8 md:h-8 text-purple-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
               </svg>
            </div>
          {/each}
        </div>
        
        <div class="text-left md:text-center flex-1">
          <p class="text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest">{count} photos</p>
          <p class="text-[9px] text-purple-300/60 uppercase tracking-wider mt-1 md:hidden">Classic Strip</p>
          <div class="hidden md:block mt-2 h-0.5 w-4 bg-purple-200 mx-auto rounded-full group-hover:w-8 transition-all"></div>
        </div>
      </button>
    {/each}
  </div>
</div>

