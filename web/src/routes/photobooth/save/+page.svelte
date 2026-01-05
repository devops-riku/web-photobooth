<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { photoboothStore } from '$lib/stores/photobooth.store';
  import { getApiUrl } from '$lib/config';

  let finalStrip: string | null = null;
  let isDownloading = false;
  let isSaving = false;
  let token: string | null = null;
  let username: string | null = null;
  let saveMessage = '';
  let title = 'My Memory';

  let savedStripId: number | null = null;

  photoboothStore.subscribe(v => {
    finalStrip = v.finalStrip;
  });

  onMount(() => {
    if (!finalStrip) {
      goto('/photobooth');
      return;
    }
    // Check for auth
    token = localStorage.getItem('sb_token');
    username = localStorage.getItem('sb_user');
  });

  function downloadStrip() {
    if (!finalStrip) return;
    isDownloading = true;

    const link = document.createElement('a');
    link.href = finalStrip;
    link.download = `wuby-strip-${Date.now()}.png`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    setTimeout(() => {
      isDownloading = false;
    }, 1000);
  }

  async function saveToGallery() {
    if (!token) return;
    isSaving = true;
    saveMessage = '';

    try {
      let response;
      
      if (savedStripId) {
        // Update existing strip (Rename)
        response = await fetch(getApiUrl('STRIP_DETAIL', savedStripId), {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({ title: title })
        });
      } else {
        // Create new strip
        if (!finalStrip) return;
        response = await fetch(getApiUrl('SAVE_STRIP'), {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            image: finalStrip,
            title: title,
            caption: 'Saved from Photobooth'
          })
        });
      }

      const data = await response.json();
      
      if (response.ok) {
        if (!savedStripId && data.id) {
          savedStripId = data.id;
          saveMessage = 'Successfully Saved!';
        } else {
          saveMessage = 'Memory Updated!';
        }
      } else {
        saveMessage = data.error || 'Failed to save';
      }
    } catch (e) {
      saveMessage = 'Error connecting to server';
    } finally {
      isSaving = false;
    }
  }
</script>

<div class="min-h-screen bg-[#f8f2ff] flex flex-col p-4 md:p-8 items-center">
  <header class="w-full max-w-4xl flex justify-between items-center mb-8 md:mb-12">
    <button 
      on:click={() => goto('/photobooth/preview')}
      class="text-purple-400 hover:text-purple-600 transition-colors flex items-center gap-2 group"
    >
      <svg class="w-5 h-5 transition-transform group-hover:-translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
      </svg>
      <span class="text-xs font-bold uppercase tracking-widest">Back</span>
    </button>

    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full bg-purple-500 animate-pulse"></div>
      <span class="text-xs font-bold uppercase tracking-[0.3em] text-purple-900">Final Result</span>
    </div>
    
    <div class="w-20"></div> <!-- Spacer -->
  </header>

  <main class="w-full max-w-6xl flex flex-col md:flex-row gap-8 md:gap-16 items-start justify-center">
    <!-- Strip Display -->
    <div class="w-full md:w-1/2 flex justify-center">
      {#if finalStrip}
        <div class="relative group">
          <img 
            src={finalStrip} 
            alt="Final Photobooth Strip" 
            class="max-h-[70vh] w-auto shadow-[0_30px_100px_rgba(159,122,234,0.3)] rounded-lg transition-transform duration-700 group-hover:scale-[1.02]"
          />
          <div class="absolute inset-0 bg-gradient-to-t from-purple-900/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-lg pointer-events-none"></div>
        </div>
      {/if}
    </div>

    <!-- Actions Section -->
    <div class="w-full md:w-80 flex flex-col gap-10">
      <section class="flex flex-col gap-4">
        <h1 class="text-2xl font-light text-purple-900 leading-tight">Your memory is ready.</h1>
        <p class="text-sm text-purple-500/70">Download your strip instantly or save it to your personal gallery.</p>
      </section>

      <section class="flex flex-col gap-3">
        <!-- Download Button -->
        <button 
          on:click={downloadStrip}
          disabled={isDownloading}
          class="w-full bg-white border-2 border-purple-100 hover:border-purple-300 text-purple-600 font-bold py-4 rounded-2xl transition-all active:scale-[0.98] flex items-center justify-center gap-3 shadow-sm hover:shadow-md disabled:opacity-50"
        >
          {#if isDownloading}
            <div class="w-5 h-5 border-2 border-purple-200 border-t-purple-500 rounded-full animate-spin"></div>
            <span>Downloading...</span>
          {:else}
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
            </svg>
            <span>Download PNG</span>
          {/if}
        </button>

        <div class="relative py-4 flex items-center">
          <div class="flex-grow border-t border-purple-100"></div>
          <span class="flex-shrink mx-4 text-[10px] font-bold uppercase tracking-widest text-purple-300">
            {#if token}Save to Collection{:else}Login to Save{/if}
          </span>
          <div class="flex-grow border-t border-purple-100"></div>
        </div>

        {#if token}
          <!-- Logged In State -->
          <div class="bg-purple-50 p-4 rounded-2xl border border-purple-100/50 mb-2">
            <div class="flex flex-col gap-2">
              <label for="title" class="text-[9px] font-bold uppercase tracking-widest text-purple-400 px-1 opacity-70">Name this memory</label>
              <input 
                id="title"
                type="text" 
                bind:value={title}
                placeholder="E.g. Fun times! ✨"
                class="w-full bg-white border border-purple-100 rounded-xl px-4 py-3 text-sm font-medium text-purple-900 focus:ring-2 focus:ring-purple-200 focus:outline-none transition-all"
              />
            </div>
          </div>

          <button 
            on:click={saveToGallery}
            disabled={isSaving}
            class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-4 rounded-2xl transition-all active:scale-[0.98] shadow-lg shadow-purple-200 flex items-center justify-center gap-3 disabled:opacity-50"
          >
            {#if isSaving}
              <div class="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin"></div>
              <span>Saving...</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>Save to My Gallery</span>
            {/if}
          </button>

          {#if saveMessage}
            <div class="flex flex-col items-center gap-2 mt-2 animate-in fade-in">
              <p class="text-[10px] text-center font-bold uppercase tracking-widest {saveMessage.includes('Successfully') ? 'text-green-600' : 'text-red-500'}">
                {saveMessage}
              </p>
              {#if saveMessage.includes('Successfully')}
                <button 
                  on:click={() => goto('/gallery')}
                  class="text-[10px] font-bold text-purple-600 hover:text-purple-800 underline underline-offset-4 uppercase tracking-widest"
                >
                  Go to Gallery →
                </button>
              {/if}
            </div>
          {/if}

        {:else}
          <!-- Guest State -->
          <div class="grid grid-cols-2 gap-3">
             <button 
              on:click={() => goto('/auth/login')}
              class="w-full bg-purple-100 hover:bg-purple-200 text-purple-700 font-bold py-4 rounded-2xl transition-all active:scale-[0.98] text-xs uppercase tracking-wider"
            >
              Log In
            </button>
            <button 
              on:click={() => goto('/auth/signup')}
              class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-4 rounded-2xl transition-all active:scale-[0.98] shadow-lg shadow-purple-200 text-xs uppercase tracking-wider"
            >
              Sign Up
            </button>
          </div>
        {/if}
      </section>

      <p class="text-[10px] text-center text-purple-400 font-medium">
        By continuing, you agree to our <a href="/terms" class="underline hover:text-purple-600">Terms</a> and <a href="/privacy" class="underline hover:text-purple-600">Privacy Policy</a>.
      </p>
    </div>
  </main>
</div>

<style>
  :global(body) {
    background-color: #f8f2ff;
  }
  
  .animate-in {
    animation: fade-in 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
  }

  @keyframes fade-in {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
