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

      <!-- Logo/Title -->
      <div class="flex flex-col items-center md:items-start">
        <h1 class="text-xl md:text-2xl font-light text-purple-900 tracking-tight">Wuby</h1>
        <div class="hidden md:flex items-center gap-2 mt-1">
          <div class="w-1.5 h-1.5 rounded-full bg-purple-500 animate-pulse"></div>
          <span class="text-[9px] font-bold uppercase tracking-[0.3em] text-purple-300">Digital Memories</span>
        </div>
      </div>

      <!-- Desktop Nav -->
      <nav class="hidden md:flex items-center gap-8">
        <button on:click={() => goto('/gallery')} class="text-[10px] font-bold uppercase tracking-widest text-purple-400 hover:text-purple-600 transition-colors">Gallery</button>
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

  <!-- Subtle Decorative Accents (Minimalist) -->
  <div class="absolute -top-24 -left-24 w-80 h-80 bg-purple-200/30 rounded-full blur-[120px]"></div>
  <div class="absolute bottom-1/4 -right-24 w-64 h-64 bg-purple-200/20 rounded-full blur-[100px]"></div>

  <!-- Main Content -->
  <main class="relative z-10 w-full max-w-5xl px-8 flex flex-col items-center pt-12 pb-40">
    <!-- Hero Section -->
    <div class="flex flex-col items-center text-center mb-40">
      <!-- Minimalist Mascot -->
      <div class="relative w-48 h-48 md:w-56 md:h-56 mb-12">
        <div class="absolute inset-0 bg-purple-50 rounded-full blur-3xl opacity-60"></div>
        <img 
          src="/booth_mascot_pastel.png" 
          alt="Wuby" 
          class="w-full h-full object-contain drop-shadow-[0_10px_30px_rgba(159,122,234,0.15)] animate-float"
        />
      </div>

      <div class="max-w-2xl px-4">
        <h1 class="text-4xl md:text-5xl font-medium text-purple-900/80 leading-[1.2] mb-6 tracking-tight">
          {HOME_SETTINGS.HERO_TITLE}
        </h1>
        
        <p class="text-[15px] md:text-[16px] text-purple-400 leading-relaxed mb-10 max-w-lg mx-auto">
          {HOME_SETTINGS.HERO_SUBTITLE}
        </p>
        
        <div class="flex flex-wrap justify-center gap-4">
          <button 
            class="px-8 py-2.5 rounded-full bg-purple-500 text-white text-[11px] font-bold uppercase tracking-[0.2em] transition-all hover:bg-purple-600 active:scale-95 shadow-lg shadow-purple-100"
            on:click={() => goto('/photobooth')}
          >
            {HOME_SETTINGS.CTA_TEXT}
          </button>
          
          <button class="px-8 py-2.5 rounded-full border border-purple-100 text-purple-300 text-[11px] font-bold uppercase tracking-[0.2em] hover:bg-purple-50 transition-all">
            See Gallery
          </button>
        </div>
      </div>
    </div>

  </main>
  
  <!-- Minimalist Footer -->
  <footer class="mt-auto w-full py-10 px-8 flex flex-col md:flex-row justify-between items-center gap-6 border-t border-purple-50/50">
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
