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

<div class="h-screen flex flex-col items-center justify-center bg-slate-50 p-6 overflow-hidden">
  <div class="text-center mb-10 shrink-0">
    <h1 class="text-4xl font-black text-slate-800 tracking-tight mb-2">WUBY</h1>
    <p class="text-slate-500 font-medium italic underline decoration-primary-500/30 underline-offset-4">Choose Your Strip Layout</p>
  </div>

  <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-3xl w-full px-4">
    {#each [2, 3, 4] as count}
      <button 
        class="group relative bg-white border-2 border-slate-200 p-5 rounded-3xl transition-all duration-300 hover:border-primary-500 hover:shadow-2xl hover:-translate-y-2 active:scale-95 text-center flex flex-col items-center" 
        on:click={() => select(count as 2|3|4)}
      >
        <div class="w-24 h-48 bg-slate-50 border border-slate-100 rounded-lg flex flex-col overflow-hidden shadow-inner group-hover:bg-primary-50 transition-colors">
          {#each Array(count) as _, i}
            <div class="flex-1 {i < count - 1 ? 'border-b border-slate-100 group-hover:border-primary-100' : ''} flex items-center justify-center p-2 opacity-20 group-hover:opacity-100 transition-opacity">
               <svg class="w-full h-full text-primary-300" fill="currentColor" viewBox="0 0 24 24"><path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm0 16H5V5h14v14zm-5.04-6.71l-2.75 3.54-1.96-2.36L6.5 17h11l-3.54-4.71z"/></svg>
            </div>
          {/each}
        </div>
        
        <div class="mt-4">
          <p class="text-sm font-black text-slate-700 uppercase tracking-widest">{count} Photos</p>
          <div class="h-1 w-0 bg-primary-500 mx-auto mt-1 rounded-full group-hover:w-full transition-all"></div>
        </div>
      </button>
    {/each}
  </div>
</div>

