<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { photoboothStore } from '$lib/stores/photobooth.store';
  import { getApiUrl, API_CONFIG } from '$lib/config';
  import { SAVE_SETTINGS } from './settings';

  let finalStrip: string | null = null;
  let isDownloading = false;
  let isSaving = false;
  let token: string | null = null;
  let username: string | null = null;
  let saveMessage = '';
  let title = 'My Memory';
  let uploadedId: string | null = null;

  let savedStripId: string | null = null;
  let guestFileURL: string | null = null;
  let showQR = false;

  photoboothStore.subscribe(v => {
    finalStrip = v.finalStrip;
    uploadedId = v.uploadedId;
  });

  onMount(() => {
    console.log("SAVE PAGE DEBUG: onMount - uploadedId:", uploadedId, "token:", !!token);
    if (!finalStrip && !uploadedId) {
      console.warn("SAVE PAGE DEBUG: No finalStrip or uploadedId in store, redirecting...");
      goto('/photobooth');
      return;
    }
    // Check for auth
    token = localStorage.getItem('sb_token');
    username = localStorage.getItem('sb_user');

    if (uploadedId) {
      if (token) {
        savedStripId = uploadedId;
        saveMessage = 'Successfully Saved!';
        const uid = localStorage.getItem('sb_uid');
        guestFileURL = `${API_CONFIG.DO_SPACES_URL}/strips/${uid}/${uploadedId}.png`;
        console.log("SAVE PAGE DEBUG: ID found, URL set to:", guestFileURL);
      } else {
        guestFileURL = `${API_CONFIG.DO_SPACES_URL}/strips/guest/${uploadedId}.png`;
        showQR = true;
      }
    }
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
    if (!token || !savedStripId) return;
    isSaving = true;
    saveMessage = '';

    try {
      // Update existing strip (Rename)
      const response = await fetch(getApiUrl('STRIP_DETAIL', savedStripId), {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ title: title })
      });

      if (response.ok) {
        saveMessage = 'Memory Updated!';
      } else {
        const data = await response.json();
        saveMessage = data.error || 'Failed to update';
      }
    } catch (e) {
      saveMessage = 'Error connecting to server';
    } finally {
      isSaving = false;
    }
  }


  function copyLink() {
    if (!guestFileURL && !savedStripId) return;
    
    let url = '';
    if (guestFileURL) { // guestFileURL is set for both guest and authenticated users now
      url = guestFileURL;
    } else if (savedStripId) { // Fallback if guestFileURL somehow wasn't set for authenticated user
      const uid = localStorage.getItem('sb_uid');
      url = `${API_CONFIG.DO_SPACES_URL}/strips/${uid}/${savedStripId}.png`;
    }
    
    navigator.clipboard.writeText(url);
    saveMessage = 'Link Copied!';
    setTimeout(() => saveMessage = '', 2000);
  }

  async function shareStrip() {
    if (!finalStrip) return;
    try {
      const shareData: any = {
        title: 'My Wuby Memory',
        text: 'Check out my photobooth memory!',
      };

      if (guestFileURL) {
        shareData.url = guestFileURL;
      }

      if (navigator.share) {
        await navigator.share(shareData);
      } else {
        copyLink();
      }
    } catch (err) {
      console.log('Share failed:', err);
    }
  }
</script>

<div class="min-h-screen bg-[#f8f2ff] flex flex-col items-center overflow-x-hidden selection:bg-purple-200">
  <!-- Decorative background elements -->
  <div class="fixed top-0 left-0 w-full h-full pointer-events-none overflow-hidden -z-10">
    <div class="absolute -top-[10%] -left-[10%] w-[40%] h-[40%] bg-purple-200/20 blur-[120px] rounded-full"></div>
    <div class="absolute top-[20%] -right-[5%] w-[30%] h-[30%] bg-pink-200/10 blur-[100px] rounded-full"></div>
  </div>

  <header class="w-full max-w-5xl flex justify-between items-center p-6 md:p-8 sticky top-0 bg-[#f8f2ff]/80 backdrop-blur-md shrink-0 z-50 border-b border-purple-100/50">
    <button 
      on:click={() => goto('/photobooth/preview')}
      class="text-purple-400 hover:text-purple-600 transition-all flex items-center gap-2 group active:scale-90"
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

  <main class="w-full max-w-5xl flex-1 flex flex-col md:flex-row gap-8 md:gap-12 items-center md:items-start justify-center p-6 md:p-12 pb-32 md:pb-24">
    <!-- Strip Display -->
    <div class="w-full md:w-1/2 flex justify-center shrink-0">
      {#if finalStrip}
        <div class="flex flex-col items-center gap-6">
          <div class="relative group">
            <img 
              src={finalStrip || guestFileURL} 
              alt="Final Photobooth Strip" 
              class="max-h-[50vh] md:max-h-[60vh] w-auto shadow-[0_40px_100px_rgba(159,122,234,0.4)] rounded-lg transition-transform duration-700 md:group-hover:scale-[1.02]"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-purple-900/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity rounded-lg pointer-events-none"></div>
          </div>
          <p class="text-[10px] font-bold text-purple-300 uppercase tracking-widest md:hidden animate-pulse">
            Tip: Long press image to save
          </p>
        </div>
      {/if}
    </div>

    <!-- Actions Section -->
    <div class="w-full md:w-80 flex flex-col gap-6">
      <section class="flex flex-col gap-4">
        <h1 class="text-2xl font-light text-purple-900 leading-tight">Your memory is ready.</h1>
        <p class="text-sm text-purple-500/70">Download your strip instantly or save it to your personal gallery.</p>
      </section>

      <section class="flex flex-col gap-3">
        <div class="flex flex-col gap-2">
          <button 
            on:click={downloadStrip}
            disabled={isDownloading}
            class="w-full bg-white border-2 border-purple-100 hover:border-purple-300 text-purple-600 font-bold py-4 rounded-2xl transition-all active:scale-[0.98] flex items-center justify-center gap-3 shadow-sm hover:shadow-md disabled:opacity-50"
          >
            {#if isDownloading}
              <div class="w-5 h-5 border-2 border-purple-200 border-t-purple-500 rounded-full animate-spin"></div>
              <span>Preparing...</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
              </svg>
              <span>Download PNG</span>
            {/if}
          </button>

          <button 
            on:click={shareStrip}
            class="w-full bg-purple-100/50 hover:bg-purple-100 text-purple-600 font-bold py-3 rounded-2xl transition-all active:scale-[0.98] flex items-center justify-center gap-3"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6a3 3 0 100-2.684m0 2.684l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"/>
            </svg>
            <span class="text-xs uppercase tracking-widest">Send to Friends</span>
          </button>
        </div>

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
              <span>Updating...</span>
            {:else}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
              </svg>
              <span>{savedStripId ? 'Update Memory Name' : 'Save to My Gallery'}</span>
            {/if}
          </button>

            <div class="flex flex-col items-center gap-2 mt-2 animate-in fade-in">
              <p class="text-[10px] text-center font-bold uppercase tracking-widest {saveMessage.includes('Successfully') || saveMessage.includes('Updated') ? 'text-green-600' : 'text-red-500'}">
                {saveMessage}
              </p>
              {#if saveMessage.includes('Successfully') || saveMessage.includes('Updated')}
                <div class="flex flex-col gap-3 mt-2">
                  <button 
                    on:click={copyLink}
                    class="text-[10px] font-bold text-purple-600 hover:bg-purple-100 px-4 py-2 rounded-full transition-all uppercase tracking-widest border border-purple-200"
                  >
                    Copy Share Link
                  </button>
                  <button 
                    on:click={() => goto('/gallery')}
                    class="text-[10px] font-bold text-purple-400 hover:text-purple-600 transition-colors uppercase tracking-widest text-center"
                  >
                    Go to Gallery →
                  </button>
                </div>
              {/if}
            </div>

          <button 
            on:click={() => { localStorage.removeItem('sb_token'); localStorage.removeItem('sb_user'); location.reload(); }}
            class="w-full mt-4 text-[10px] font-bold text-purple-300 hover:text-purple-500 uppercase tracking-widest transition-colors py-2"
          >
            Sign Out
          </button>

        {:else}
          <!-- Guest State -->
          <div class="flex flex-col gap-4">
            {#if !guestFileURL}
              <div class="bg-white border-2 border-purple-100 p-8 rounded-[2rem] text-center flex flex-col gap-4">
                <p class="text-xs font-bold text-purple-400 uppercase tracking-widest">Upload failed or missing</p>
                <button 
                  on:click={() => goto('/photobooth/preview')}
                  class="bg-purple-600 text-white px-8 py-3 rounded-xl text-[10px] font-bold uppercase tracking-widest active:scale-95"
                >
                  Go Back to Preview
                </button>
              </div>
            {:else}
              <div class="bg-white p-6 rounded-[2rem] border-2 border-purple-100 flex flex-col items-center gap-6 animate-in">
                <div class="bg-purple-50 p-3 rounded-2xl">
                  <img 
                    src={`${SAVE_SETTINGS.BASE_QR_API}${encodeURIComponent(guestFileURL)}`} 
                    alt="Scan to view on your phone" 
                    class="w-40 h-40"
                  />
                </div>
                
                <div class="flex flex-col items-center gap-2">
                  <p class="text-[10px] font-bold uppercase tracking-widest text-purple-900">Your Share Link</p>
                  <div class="flex gap-2 w-full">
                    <input 
                      readonly 
                      value={guestFileURL} 
                      class="flex-1 bg-purple-50 text-[10px] font-medium text-purple-400 px-4 py-3 rounded-xl border border-purple-100 outline-none"
                    />
                    <button 
                      on:click={copyLink}
                      class="bg-purple-600 p-3 rounded-xl text-white hover:bg-purple-700 transition-colors"
                      aria-label="Copy share link"
                      title="Copy share link"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3"/></svg>
                    </button>
                  </div>
                  <p class="text-[9px] text-purple-300 font-medium mt-2">✨ {SAVE_SETTINGS.GUEST_EXPIRY_MESSAGE}</p>
                </div>
              </div>
            {/if}

            <div class="grid grid-cols-2 gap-3 mt-2">
               <button 
                on:click={() => goto('/auth/login?redirect=/photobooth/save')}
                class="w-full bg-purple-100 hover:bg-purple-200 text-purple-700 font-bold py-4 rounded-2xl transition-all active:scale-[0.98] text-xs uppercase tracking-wider"
              >
                Log In
              </button>
              <button 
                on:click={() => goto('/auth/signup?redirect=/photobooth/save')}
                class="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-4 rounded-2xl transition-all active:scale-[0.98] shadow-lg shadow-purple-200 text-xs uppercase tracking-wider"
              >
                Sign Up
              </button>
            </div>
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
