<script lang="ts">
  import { layoutStore, type LayoutOption } from '$lib/stores/layout.store';
  import { goto } from '$app/navigation';

  const STRIP_WIDTH = 600;

  function select(count: 2 | 3 | 4) {
    const layout: LayoutOption = {
      count,
      width: STRIP_WIDTH,
      height: STRIP_WIDTH * count * 1.5 // classic booth ratio
    };

    layoutStore.set(layout);
    goto('/photobooth/capture');
  }
</script>

<div class="h-screen flex flex-col items-center justify-center p-6 overflow-hidden">
  <div class="text-center mb-12">
    <h1 class="text-3xl font-medium tracking-tight text-purple-900/80 mb-2">Wuby</h1>
    <p class="text-purple-400/80 text-sm font-medium">Select a layout to begin</p>
  </div>

  <div class="flex gap-8 max-w-4xl w-full justify-center px-4">
    {#each [2, 3, 4] as count}
      <button 
        class="group flex flex-col items-center gap-6 transition-all duration-400 active:scale-95" 
        on:click={() => select(count as 2|3|4)}
      >
        <div class="w-32 h-64 bg-white border border-purple-100 rounded-xl flex flex-col overflow-hidden shadow-[0_4px_12px_rgba(159,122,234,0.05)] group-hover:shadow-[0_8px_24px_rgba(159,122,234,0.1)] group-hover:border-purple-200 transition-all p-2 gap-1.5">
          {#each Array(count) as _, i}
            <div class="flex-1 bg-purple-50/50 rounded-md flex items-center justify-center">
               <svg class="w-8 h-8 text-purple-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                 <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
               </svg>
            </div>
          {/each}
        </div>
        
        <div class="text-center">
          <p class="text-xs font-semibold text-purple-400 uppercase tracking-widest">{count} photos</p>
          <div class="mt-2 h-0.5 w-4 bg-purple-200 mx-auto rounded-full group-hover:w-8 transition-all"></div>
        </div>
      </button>
    {/each}
  </div>
</div>

