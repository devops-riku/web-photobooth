<script lang="ts">
  import { HOME_SETTINGS } from './settings';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  let isMenuOpen = false;
  let isLoggedIn = false;

  onMount(() => {
    const token = localStorage.getItem('sb_token');
    isLoggedIn = !!token;
  });

  function handleSignOut() {
    localStorage.removeItem('sb_token');
    localStorage.removeItem('sb_user');
    localStorage.removeItem('sb_uid');
    isLoggedIn = false;
    goto('/auth/login');
  }
</script>

<div class="min-h-screen bg-[#f8f2ff] flex flex-col relative overflow-hidden selection:bg-purple-200">
  <!-- Responsive Header -->
  <header class="w-full sticky top-0 z-40 bg-[#f8f2ff]/80 backdrop-blur-md px-4 py-4 md:px-8 md:py-6 flex justify-center">
    <div class="w-full max-w-5xl flex justify-between items-center relative">
      <!-- Mobile Burger -->
      <button 
        class="md:hidden p-2 -ml-2 text-purple-900"
        on:click={() => isMenuOpen = true}
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/></svg>
      </button>

      <!-- Logo/Title (Centered on Mobile) -->
      <div class="flex flex-col items-center md:items-start flex-1 md:flex-none">
        <h1 class="text-xl md:text-2xl font-semibold text-purple-900 leading-none" style="font-family: 'Quicksand', sans-serif">Wuby</h1>
        <div class="flex items-center gap-2 mt-1">
          <div class="w-1 h-1 rounded-full bg-purple-500 animate-pulse"></div>
          <span class="text-[8px] font-bold uppercase tracking-[0.3em] text-purple-300">Digital Memories</span>
        </div>
      </div>

      <!-- Desktop Nav -->
      <nav class="hidden md:flex items-center gap-8 flex-1 justify-end">
        <button on:click={() => goto('/gallery')} class="text-[10px] font-bold uppercase tracking-widest text-purple-400 hover:text-purple-600 transition-colors">Your Gallery</button>
        {#if isLoggedIn}
          <button 
            on:click={handleSignOut}
            class="px-6 py-2 rounded-full border-2 border-purple-100 text-purple-500 text-[10px] font-bold uppercase tracking-widest hover:bg-purple-500 hover:text-white hover:border-purple-500 transition-all"
          >
            Sign Out
          </button>
        {:else}
          <button 
            on:click={() => goto('/auth/login')}
            class="px-6 py-2 rounded-full border-2 border-purple-100 text-purple-500 text-[10px] font-bold uppercase tracking-widest hover:bg-purple-500 hover:text-white hover:border-purple-500 transition-all"
          >
            Login
          </button>
        {/if}
      </nav>

      <!-- Mobile Spacer -->
      <div class="md:hidden w-8"></div>
    </div>
  </header>

  <!-- Mobile Menu Overlay -->
  {#if isMenuOpen}
    <div 
      class="fixed inset-0 z-50 bg-[#f8f2ff] flex flex-col p-6 animate-in slide-in-from-top-4"
    >
      <div class="flex justify-end mb-4 relative z-50">
        <button 
          on:click={() => isMenuOpen = false}
          class="p-4 -mr-4 text-purple-900 tap-highlight-transparent active:scale-95 transition-transform"
          aria-label="Close menu"
        >
          <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <nav class="flex flex-col gap-8 items-center flex-grow justify-center -mt-20">
        <button 
          on:click={() => { isMenuOpen = false; goto('/photobooth'); }}
          class="text-3xl font-light text-purple-900"
        >
          Start Session
        </button>
        <button 
          on:click={() => { isMenuOpen = false; goto('/gallery'); }}
          class="text-3xl font-light text-purple-900"
        >
          Gallery
        </button>
        <div class="w-16 h-px bg-purple-100"></div>
        {#if isLoggedIn}
          <button 
            on:click={() => { isMenuOpen = false; handleSignOut(); }}
            class="text-3xl font-light text-purple-900"
          >
            Sign Out
          </button>
        {:else}
          <button 
            on:click={() => { isMenuOpen = false; goto('/auth/login'); }}
            class="text-3xl font-light text-purple-900"
          >
            Login
          </button>
        {/if}
      </nav>
    </div>
  {/if}

  <!-- Subtle Decorative Accents (More playful blobs) -->
  <div class="absolute -top-24 -left-24 w-96 h-96 bg-purple-200/40 rounded-full blur-[120px] animate-pulse"></div>
  <div class="absolute top-1/2 -right-32 w-80 h-80 bg-pink-100/30 rounded-full blur-[100px] animate-pulse duration-[7000ms]"></div>
  <div class="absolute bottom-1/4 -left-24 w-64 h-64 bg-purple-100/30 rounded-full blur-[100px]"></div>

  <!-- Main Content -->
  <main class="relative z-10 w-full max-w-5xl mx-auto px-8 flex-1 flex flex-col items-center justify-center py-8">
    <!-- Hero Section -->
    <div class="flex flex-col items-center text-center">
      <!-- Minimalist Mascot with Square Frame -->
      <div class="relative w-48 md:w-64 aspect-square mb-12">
        <div class="absolute inset-0 bg-purple-100/40 rounded-full blur-3xl opacity-60 scale-125"></div>
        
        <!-- Floating Sparkles -->
        <span class="absolute -top-6 -right-6 text-2xl animate-float duration-[3s] drop-shadow-sm z-20">✨</span>
        <span class="absolute top-1/2 -left-12 text-xl animate-float duration-[4s] delay-700 drop-shadow-sm z-20">⭐</span>
        <span class="absolute -bottom-6 right-10 text-2xl animate-float duration-[3.5s] delay-300 drop-shadow-sm z-20">💖</span>
        
        <!-- The Perfect Square Frame -->
        <div class="w-full h-full bg-white rounded-3xl p-8 shadow-xl shadow-purple-200/50 border border-purple-50 relative z-10 overflow-hidden flex items-center justify-center animate-float aspect-square">
          <img 
            src="/booth_mascot_pastel.png" 
            alt="Wuby" 
            class="w-full h-full object-contain"
          />
        </div>
      </div>

      <div class="max-w-2xl px-4">
        <h1 class="text-3xl md:text-5xl font-bold text-purple-900/80 leading-[1.2] mb-6" style="font-family: 'Quicksand', sans-serif">
          {HOME_SETTINGS.HERO_TITLE}
        </h1>
        
        <p class="text-[14px] md:text-[16px] text-purple-400 leading-relaxed mb-8 max-w-lg mx-auto">
          {HOME_SETTINGS.HERO_SUBTITLE}
        </p>
        
        <div class="flex flex-col items-center gap-6">
          <p class="text-[9px] font-bold text-purple-300 uppercase tracking-[0.4em] animate-pulse">Ready to smile? ✨</p>
          <button 
            class="px-12 py-4 rounded-full bg-purple-500 text-white text-[11px] font-bold uppercase tracking-[0.2em] transition-all hover:bg-purple-600 hover:scale-105 active:scale-95 shadow-xl shadow-purple-200 group relative"
            on:click={() => goto('/photobooth')}
          >
            <span class="relative z-10">{HOME_SETTINGS.CTA_TEXT}</span>
            <div class="absolute inset-0 rounded-full bg-white/20 opacity-0 group-hover:opacity-100 transition-opacity"></div>
          </button>
        </div>
      </div>
    </div>
  </main>
  
  <!-- Minimalist Footer -->
  <footer class="w-full border-t border-purple-50/50 flex-shrink-0">
    <div class="max-w-5xl mx-auto py-10 px-8 flex flex-col md:flex-row justify-between items-center gap-6">
      <div class="flex items-center gap-2">
        <div class="w-1.5 h-1.5 rounded-full bg-purple-200"></div>
        <span class="text-[9px] font-bold text-purple-200 uppercase tracking-widest leading-none">&copy; 2026 {HOME_SETTINGS.HERO_TITLE}</span>
      </div>
      
      <div class="flex gap-10">
        {#each ['Instagram', 'Documentation', 'Contact'] as link}
          <a href="#" class="text-[9px] font-bold text-purple-200 uppercase tracking-widest hover:text-purple-400 transition-colors">
            {link}
          </a>
        {/each}
      </div>
    </div>
  </footer>
</div>

<style>
  @keyframes float {
    0%, 100% { transform: translateY(0px); }
    50% { transform: translateY(-12px); }
  }
  
  .animate-float {
    animation: float 5s ease-in-out infinite;
  }

  /* Standardizing line height for better control */
  :global(body) {
    line-height: 1.5;
  }
</style>
