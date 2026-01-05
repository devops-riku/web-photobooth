<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { getApiUrl, API_CONFIG } from '$lib/config';
  import { goto } from '$app/navigation';

  let strip: any = null;
  let loading = true;
  let error: string | null = null;

  const id = $page.params.id;

  onMount(async () => {
    try {
      const response = await fetch(`${API_CONFIG.BASE_URL}/api/strips/public/${id}`);
      if (!response.ok) {
        if (response.status === 410) {
          throw new Error('This memory has expired and is no longer available.');
        }
        throw new Error('Memory not found');
      }
      strip = await response.json();
    } catch (e: any) {
      error = e.message;
    } finally {
      loading = false;
    }
  });

  function downloadImage() {
    if (!strip) return;
    const link = document.createElement('a');
    link.href = strip.file_url;
    link.download = `wuby-memory-${id}.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<svelte:head>
  <title>Wuby Photobooth - A Digital Memory</title>
  <meta name="description" content="View this digital photobooth memory captured with Wuby." />
</svelte:head>

<div class="min-h-screen bg-[#f8f2ff] flex flex-col items-center py-12 px-6">
  <!-- Small Logo Header -->
  <header class="mb-12 flex flex-col items-center">
    <button on:click={() => goto('/')} class="flex flex-col items-center group">
        <h1 class="text-3xl font-light text-purple-900 tracking-tight group-hover:scale-105 transition-transform">Wuby</h1>
        <div class="flex items-center gap-2 mt-1">
          <div class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse"></div>
          <span class="text-[9px] font-bold uppercase tracking-[0.3em] text-purple-300">Digital Memories</span>
        </div>
    </button>
  </header>

  <main class="w-full max-w-lg flex flex-col items-center">
    {#if loading}
      <div class="flex flex-col items-center gap-4 py-20">
        <div class="w-10 h-10 border-4 border-purple-100 border-t-purple-500 rounded-full animate-spin"></div>
        <p class="text-[10px] font-bold uppercase tracking-widest text-purple-400">Loading your memory...</p>
      </div>
    {:else if error}
      <div class="bg-white p-12 rounded-[2.5rem] shadow-xl shadow-purple-900/5 text-center flex flex-col items-center gap-6">
        <div class="w-16 h-16 bg-red-50 rounded-full flex items-center justify-center text-red-400">
           <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/></svg>
        </div>
        <h2 class="text-xl font-medium text-slate-800">Oops!</h2>
        <p class="text-sm text-slate-400 max-w-[200px]">{error}</p>
        <button 
          on:click={() => goto('/')}
          class="mt-4 px-8 py-3 bg-purple-600 text-white text-xs font-bold uppercase tracking-widest rounded-full hover:bg-purple-700 transition-all active:scale-95"
        >
          Back to Home
        </button>
      </div>
    {:else}
      <div class="flex flex-col items-center gap-10 w-full">
        <!-- The Strip -->
        <div class="relative animate-in fade-in slide-in-from-bottom-8 duration-1000">
          <img 
            src={strip.file_url} 
            alt={strip.title} 
            class="max-h-[75vh] w-auto rounded-lg shadow-[0_30px_80px_rgba(159,122,234,0.2)] border-4 border-white"
          />
          <div class="absolute -bottom-4 -right-4 w-12 h-12 bg-purple-500 rounded-full flex items-center justify-center text-white shadow-lg border-2 border-white rotate-12">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/></svg>
          </div>
        </div>

        <!-- Info & Actions -->
        <div class="bg-white w-full p-8 md:p-10 rounded-[2.5rem] shadow-xl shadow-purple-900/5 flex flex-col items-center text-center gap-6">
           <div class="flex flex-col gap-2">
             <h2 class="text-2xl font-light text-purple-900">{strip.title || 'Untitled Memory'}</h2>
             <p class="text-xs text-purple-300 font-bold uppercase tracking-widest">
                {new Date(strip.created_at).toLocaleDateString('en-US', { month: 'long', day: 'numeric', year: 'numeric' })}
             </p>
           </div>

           {#if strip.caption}
             <p class="text-sm text-slate-500 max-w-[280px] leading-relaxed italic">"{strip.caption}"</p>
           {/if}

           <div class="w-full grid grid-cols-1 gap-3 mt-4">
              <button 
                on:click={downloadImage}
                class="w-full py-4 rounded-2xl bg-purple-50 text-purple-600 text-xs font-bold uppercase tracking-widest hover:bg-purple-100 transition-all active:scale-95 flex items-center justify-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
                Download Photo
              </button>
              
              <div class="h-px bg-purple-50 my-2"></div>
              
              <button 
                on:click={() => goto('/')}
                class="w-full py-4 rounded-2xl bg-purple-600 text-white text-xs font-bold uppercase tracking-widest shadow-lg shadow-purple-100 hover:bg-purple-700 transition-all active:scale-95 flex items-center justify-center gap-2"
              >
                Capture Your Own Memory
              </button>
           </div>
        </div>
      </div>
    {/if}
  </main>
  
  <footer class="mt-20">
    <p class="text-[9px] font-bold uppercase tracking-widest text-purple-200">Wuby Photobooth &copy; 2026</p>
  </footer>
</div>

<style>
  .animate-in {
    animation-fill-mode: both;
  }
</style>
