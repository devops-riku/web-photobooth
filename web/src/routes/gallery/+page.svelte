<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { getApiUrl } from '$lib/config';

  interface Strip {
    id: number;
    title: string;
    caption: string;
    file_url: string;
    created_at: string;
  }

  let strips: Strip[] = [];
  let isLoading = true;
  let token: string | null = null;
  let error = '';

  // View Modal State
  let viewingStrip: Strip | null = null;

  // Edit Modal State
  let isEditing = false;
  let editingStrip: Strip | null = null;
  let newTitle = '';
  let isUpdating = false;

  let isMenuOpen = false;

  onMount(async () => {
    token = localStorage.getItem('sb_token');
    if (!token) {
      goto('/auth/login');
      return;
    }
    await fetchStrips();
  });

  async function fetchStrips() {
    isLoading = true;
    try {
      const response = await fetch(getApiUrl('GET_STRIPS'), {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      const data = await response.json();
      if (response.ok) {
        strips = data.strips || [];
      } else {
        error = data.error || 'Failed to load gallery';
      }
    } catch (e) {
      error = 'Server connection failed';
    } finally {
      isLoading = false;
    }
  }

  async function deleteStrip(id: number) {
    if (!confirm('Are you sure you want to delete this memory?')) return;

    try {
      const response = await fetch(getApiUrl('STRIP_DETAIL', id), {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });

      if (response.ok) {
        strips = strips.filter(s => s.id !== id);
        viewingStrip = null; // Close view modal if deleted
      } else {
        alert('Delete failed');
      }
    } catch (e) {
      alert('Error connecting to server');
    }
  }

  function openViewModal(strip: Strip) {
    viewingStrip = strip;
  }

  function openEditModal(strip: Strip) {
    editingStrip = strip;
    newTitle = strip.title;
    isEditing = true;
    // Keep view modal open if it called this, or close it? 
    // Usually standard to edit "on top" or close view. Let's keep view open in background or close it.
    // For simplicity, let's just use the edit modal as a standalone interaction.
    viewingStrip = null; 
  }

  async function updateStrip() {
    if (!editingStrip || !newTitle.trim()) return;
    isUpdating = true;

    try {
      const response = await fetch(getApiUrl('STRIP_DETAIL', editingStrip.id), {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ title: newTitle })
      });

      if (response.ok) {
        const data = await response.json();
        strips = strips.map(s => s.id === editingStrip!.id ? data.strip : s);
        isEditing = false;
        editingStrip = null;
        // Re-open view modal with updated strip?
        viewingStrip = data.strip;
      } else {
        alert('Update failed');
      }
    } catch (e) {
      alert('Error updating');
    } finally {
      isUpdating = false;
    }
  }
</script>

<div class="h-screen overflow-y-auto bg-[#f8f2ff] flex flex-col relative w-full">
  <!-- Responsive Header -->
  <header class="w-full sticky top-0 z-40 bg-[#f8f2ff]/80 backdrop-blur-md border-b border-purple-100/50 px-4 py-4 md:px-8 md:py-6 flex-shrink-0">
    <div class="max-w-6xl mx-auto flex justify-between items-center relative">
      <!-- Mobile Burger -->
      <button 
        class="md:hidden p-2 -ml-2 text-purple-900"
        on:click={() => isMenuOpen = true}
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/></svg>
      </button>

      <!-- Desktop Left: Back Button -->
      <div class="hidden md:block w-32">
        <button 
          on:click={() => goto('/photobooth')}
          class="text-purple-400 hover:text-purple-600 transition-colors flex items-center gap-2 group"
        >
          <svg class="w-5 h-5 transition-transform group-hover:-translate-x-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
          </svg>
          <span class="text-xs font-bold uppercase tracking-widest">Photobooth</span>
        </button>
      </div>

      <!-- Center Title -->
      <div class="flex flex-col items-center">
        <h1 class="text-xl md:text-3xl font-light text-purple-900 tracking-tight">Your Gallery</h1>
        <div class="flex items-center gap-2 mt-1 md:mt-2">
          <div class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse"></div>
          <span class="text-[9px] md:text-[10px] font-bold uppercase tracking-[0.3em] text-purple-300">Member Memories</span>
        </div>
      </div>
      
      <!-- Desktop Right: Sign Out -->
      <div class="hidden md:flex w-32 justify-end">
        <button 
          on:click={() => { localStorage.removeItem('sb_token'); goto('/auth/login'); }}
          class="text-xs font-bold uppercase tracking-widest text-purple-400 hover:text-purple-700 transition-colors"
        >
          Sign Out
        </button>
      </div>

      <!-- Mobile Spacer (balances layout) -->
      <div class="md:hidden w-8"></div>
    </div>
  </header>

  <!-- Mobile Menu Overlay -->
  {#if isMenuOpen}
    <div 
      class="fixed inset-0 z-50 bg-[#f8f2ff] flex flex-col p-6 animate-in slide-in-from-top-4"
    >
      <div class="flex justify-end mb-8">
        <button 
          on:click={() => isMenuOpen = false}
          class="p-2 -mr-2 text-purple-900"
        >
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>
      </div>

      <nav class="flex flex-col gap-8 items-center flex-grow justify-center -mt-20">
        <button 
          on:click={() => goto('/photobooth')}
          class="text-3xl font-light text-purple-900 hover:text-purple-600 transition-colors tracking-tight"
        >
          Photobooth
        </button>
        <div class="w-16 h-px bg-purple-100"></div>
        <button 
          on:click={() => { localStorage.removeItem('sb_token'); goto('/auth/login'); }}
          class="text-3xl font-light text-purple-900 hover:text-purple-600 transition-colors tracking-tight"
        >
          Sign Out
        </button>
      </nav>

      <div class="text-center pb-8 opacity-50">
        <p class="text-[10px] font-bold uppercase tracking-[0.3em] text-purple-400">Wuby Photobooth</p>
      </div>
    </div>
  {/if}

  <main class="w-full max-w-6xl mx-auto flex-grow p-4 md:p-8">
    {#if isLoading}
      <div class="flex flex-col items-center justify-center min-h-[50vh] gap-4">
        <div class="w-12 h-12 border-4 border-purple-100 border-t-purple-500 rounded-full animate-spin"></div>
        <p class="text-[10px] font-bold uppercase tracking-widest text-purple-300">Developing your memories...</p>
      </div>
    {:else if error}
      <div class="flex flex-col items-center justify-center min-h-[50vh] gap-4">
        <p class="text-red-400 font-medium">{error}</p>
        <button on:click={fetchStrips} class="text-purple-500 underline text-sm">Try Again</button>
      </div>
    {:else if strips.length === 0}
      <div class="flex flex-col items-center justify-center min-h-[50vh] gap-6 text-center">
        <div class="w-20 h-20 bg-purple-50 rounded-full flex items-center justify-center">
          <svg class="w-10 h-10 text-purple-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
          </svg>
        </div>
        <p class="text-purple-400/70 font-light text-lg">Your gallery is empty.<br><span class="text-sm">Start taking some photos!</span></p>
        <button 
          on:click={() => goto('/photobooth')}
          class="bg-purple-600 text-white px-8 py-3 rounded-full font-bold uppercase tracking-[0.2em] text-xs hover:bg-purple-700 transition-all shadow-lg shadow-purple-200"
        >
          Open Photobooth
        </button>
      </div>
    {:else}
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 md:gap-8">
        {#each strips as strip (strip.id)}
          <button 
            class="group flex flex-col gap-3 animate-in text-left w-full"
            on:click={() => openViewModal(strip)}
          >
            <!-- Square container for all thumbnails -->
            <div class="relative w-full aspect-square overflow-hidden rounded-2xl bg-white shadow-lg shadow-purple-100/50 transition-all duration-300 group-hover:scale-[1.03] group-hover:shadow-xl group-hover:shadow-purple-200/50 ring-1 ring-purple-50">
              <img src={strip.file_url} alt={strip.title} class="w-full h-full object-cover" />
              <!-- Subtle overlay on hover -->
              <div class="absolute inset-0 bg-purple-900/10 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                <svg class="w-8 h-8 text-white drop-shadow-lg transform scale-90 group-hover:scale-100 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                   <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
              </div>
            </div>
            <div class="px-1">
              <h3 class="text-xs font-bold text-purple-900 truncate">{strip.title || 'Untitled'}</h3>
              <p class="text-[9px] text-purple-400 font-medium uppercase tracking-wider mt-0.5">
                {new Date(strip.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
              </p>
            </div>
          </button>
        {/each}
      </div>
    {/if}
  </main>

  <!-- View Modal (Mobile Responsive) -->
  {#if viewingStrip}
    <div 
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/90 backdrop-blur-xl animate-in fade-in p-4 md:p-8"
      on:click|self={() => viewingStrip = null}
    >
      <div class="relative w-full h-full md:w-auto md:h-auto md:max-h-[90vh] flex flex-col items-center gap-4 md:gap-4">
        <!-- Close Button -->
        <button 
          on:click={() => viewingStrip = null}
          class="absolute top-4 right-4 md:top-2 md:right-2 lg:-right-12 lg:-top-0 z-20 p-2 text-white/50 hover:text-white transition-colors bg-black/20 hover:bg-black/40 rounded-full md:bg-transparent"
        >
          <svg class="w-6 h-6 md:w-8 md:h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
        </button>

        <!-- Image Container -->
        <div class="relative rounded-lg overflow-hidden shadow-2xl shadow-black/50 w-full flex-1 md:flex-none min-h-0 md:h-auto flex items-center justify-center bg-black/50">
          <img 
            src={viewingStrip.file_url} 
            alt={viewingStrip.title} 
            class="w-full h-full object-contain md:w-auto md:h-auto md:max-h-[70vh]" 
          />
        </div>

        <!-- Info & Actions -->
        <div class="w-full md:w-[400px] flex flex-col items-center gap-4 text-center">
            <div>
              <h2 class="text-xl font-medium text-white tracking-wide">{viewingStrip.title || 'Untitled'}</h2>
              <p class="text-xs uppercase tracking-widest text-white/50 mt-1">
                {new Date(viewingStrip.created_at).toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}
              </p>
            </div>
            
            <div class="flex items-center gap-3 w-full justify-center flex-wrap">
               <button 
                on:click={() => openEditModal(viewingStrip!)}
                class="px-6 py-3 bg-white/10 hover:bg-white/20 text-white rounded-full text-[10px] font-bold uppercase tracking-widest backdrop-blur-md transition-all flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/></svg>
                Rename
              </button>
              <button 
                on:click={() => deleteStrip(viewingStrip!.id)}
                class="px-6 py-3 bg-red-500/20 hover:bg-red-500/40 text-red-200 hover:text-white rounded-full text-[10px] font-bold uppercase tracking-widest backdrop-blur-md transition-all flex items-center gap-2"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                Delete
              </button>
              <a 
                href={viewingStrip.file_url} 
                download
                target="_blank"
                class="px-6 py-3 bg-white text-purple-900 hover:bg-purple-50 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all flex items-center gap-2"
              >
                 <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
                 Download
              </a>
            </div>
        </div>
      </div>
    </div>
  {/if}

  <!-- Edit Modal -->
  {#if isEditing}
    <div class="fixed inset-0 z-[60] flex items-center justify-center p-6 bg-purple-900/40 backdrop-blur-md animate-in fade-in">
      <div class="w-full max-w-sm bg-white rounded-[2.5rem] shadow-2xl p-8 border border-purple-50">
        <h2 class="text-xl font-light text-purple-900 mb-6">Rename Memory</h2>
        <div class="space-y-6">
          <div class="flex flex-col gap-2">
            <label class="text-[10px] font-bold uppercase tracking-widest text-purple-400 px-2 opacity-70">New Name</label>
            <input 
              type="text" 
              bind:value={newTitle}
              class="w-full bg-purple-50/50 border border-purple-100/50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 focus:ring-4 focus:ring-purple-100/50 focus:bg-white transition-all outline-none"
              placeholder="E.g. Summer Night 🌙"
            />
          </div>
          
          <div class="flex items-center gap-3">
            <button 
              on:click={() => { isEditing = false; if (editingStrip) viewingStrip = editingStrip; }}
              class="flex-grow bg-purple-50 text-purple-400 font-bold py-4 rounded-2xl text-[10px] uppercase tracking-widest hover:bg-purple-100 transition-all"
            >
              Cancel
            </button>
            <button 
              on:click={updateStrip}
              disabled={isUpdating}
              class="flex-grow bg-purple-600 text-white font-bold py-4 rounded-2xl text-[10px] uppercase tracking-widest hover:bg-purple-700 transition-all shadow-lg shadow-purple-200 disabled:opacity-50"
            >
              {isUpdating ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

<style>
  :global(body) {
    background-color: #f8f2ff;
  }

  .animate-in {
    animation: fade-in 0.4s cubic-bezier(0.16, 1, 0.3, 1) both;
  }

  @keyframes fade-in {
    from { opacity: 0; }
    to { opacity: 1; }
  }
</style>
