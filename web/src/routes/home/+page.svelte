<script lang="ts">
  import { HOME_SETTINGS } from './settings';
  import { goto } from '$app/navigation';
</script>

<div class="min-h-screen bg-[#f8f2ff] flex flex-col items-center relative overflow-hidden selection:bg-purple-200">
  <!-- Subtle Decorative Accents (Minimalist) -->
  <div class="absolute -top-24 -left-24 w-80 h-80 bg-purple-200/30 rounded-full blur-[120px]"></div>
  <div class="absolute bottom-1/4 -right-24 w-64 h-64 bg-purple-200/20 rounded-full blur-[100px]"></div>

  <!-- Main Content -->
  <main class="relative z-10 w-full max-w-5xl px-8 flex flex-col items-center pt-24 pb-40">
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
