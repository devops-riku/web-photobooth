<script lang="ts">
  import { goto } from '$app/navigation';
  import { getApiUrl, BRAND_CONFIG } from '$lib/config';

  import { page } from '$app/stores';

  let identifier = '';
  let password = '';
  let isLoading = false;
  let message = '';
  let messageType: 'success' | 'error' = 'success';

  async function handleLogin() {
    isLoading = true;
    message = '';
    
    try {
      const response = await fetch(getApiUrl('LOGIN'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ identifier, password })
      });

      const data = await response.json();

      if (!response.ok) {
        message = data.error || 'Login failed';
        messageType = 'error';
      } else {
        message = 'Welcome back! Redirecting...';
        messageType = 'success';
        
        if (data.access_token) {
          localStorage.setItem('sb_token', data.access_token);
        }
        if (data.user) {
          if (data.user.username) localStorage.setItem('sb_user', data.user.username);
          if (data.user.id) localStorage.setItem('sb_uid', data.user.id);
        }
        
        const redirect = $page.url.searchParams.get('redirect') || '/gallery';
        setTimeout(() => {
          goto(redirect);
        }, 800);
      }
    } catch (e) {
      message = 'An unexpected error occurred';
      messageType = 'error';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="min-h-screen bg-[#fcf9ff] flex items-center justify-center p-6 relative overflow-hidden">
  <!-- Decorative Pastel Orbs -->
  <div class="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-purple-100/40 rounded-full blur-[120px]"></div>
  <div class="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-purple-200/30 rounded-full blur-[120px]"></div>

  <div class="w-full max-w-sm flex flex-col items-center relative z-10">
    <!-- Header -->
    <div class="flex flex-col items-center gap-6 mb-16 text-center animate-in">
      <div class="w-14 h-14 bg-white rounded-[2rem] flex items-center justify-center shadow-xl shadow-purple-200/50">
        <svg class="w-7 h-7 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 11c0 3.517-1.009 6.799-2.753 9.571m-3.44-2.04l.054-.09A10.003 10.003 0 0011.203 3c-2.123 0-4.047.665-5.625 1.799M9 15h.01M9 19h.01M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
      </div>
      <div>
        <h1 class="text-3xl font-light text-purple-900 tracking-tight">Sign In</h1>
        <p class="text-[10px] font-bold uppercase tracking-[0.3em] text-purple-300 mt-2">{BRAND_CONFIG.NAME} Photobooth</p>
      </div>
    </div>

    <form on:submit|preventDefault={handleLogin} class="w-full flex flex-col gap-8 animate-in delay-100">
      <div class="space-y-4">
        <div class="flex flex-col gap-2">
          <input 
            id="identifier"
            type="text" 
            bind:value={identifier}
            placeholder="Email or Username"
            class="w-full bg-white/50 backdrop-blur-sm border-2 border-purple-50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 placeholder:text-purple-200 focus:border-purple-200 focus:bg-white focus:outline-none transition-all"
            required
          />
        </div>

        <div class="flex flex-col gap-2 relative">
          <input 
            id="password"
            type="password" 
            bind:value={password}
            placeholder="Password"
            class="w-full bg-white/50 backdrop-blur-sm border-2 border-purple-50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 placeholder:text-purple-200 focus:border-purple-200 focus:bg-white focus:outline-none transition-all"
            required
          />
          <button type="button" class="absolute right-6 bottom-4 text-[10px] font-bold text-purple-300 hover:text-purple-500 transition-colors uppercase tracking-widest">Forgot?</button>
        </div>
      </div>

      {#if message}
        <div class="px-4 py-3 rounded-xl text-center text-[10px] font-bold uppercase tracking-widest {messageType === 'error' ? 'bg-red-50 text-red-400' : 'bg-purple-50 text-purple-500'} animate-in active-message">
          {message}
        </div>
      {/if}

      <button 
        type="submit"
        disabled={isLoading}
        class="w-full bg-purple-500 hover:bg-purple-600 text-white text-xs font-bold uppercase tracking-[0.4em] py-5 rounded-2xl transition-all active:scale-[0.98] shadow-lg shadow-purple-100 disabled:opacity-50"
      >
        {#if isLoading}
          <div class="w-4 h-4 border-2 border-white/20 border-t-white rounded-full animate-spin m-auto"></div>
        {:else}
          Sign In
        {/if}
      </button>
    </form>

    <div class="mt-12 flex flex-col items-center gap-6 animate-in delay-200">
      <p class="text-[10px] font-bold text-purple-300 uppercase tracking-widest">
        No account? 
        <button on:click={() => { 
          const redirect = $page.url.searchParams.get('redirect');
          goto(redirect ? `/auth/signup?redirect=${redirect}` : '/auth/signup');
        }} class="text-purple-500 hover:text-purple-700 underline underline-offset-4 ml-1">Join {BRAND_CONFIG.NAME}</button>
      </p>
      
      <button 
        on:click={() => goto('/home')}
        class="text-[9px] font-bold text-purple-200 uppercase tracking-[0.2em] hover:text-purple-400 transition-colors"
      >
        Back to Home
      </button>
    </div>
  </div>
</div>

<style>
  :global(body) {
    background-color: #fcf9ff;
  }

  .animate-in {
    animation: fade-in 0.8s cubic-bezier(0.16, 1, 0.3, 1) both;
  }

  .active-message {
    animation: slide-up 0.4s ease-out;
  }

  .delay-100 { animation-delay: 100ms; }
  .delay-200 { animation-delay: 200ms; }

  @keyframes fade-in {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @keyframes slide-up {
    from { opacity: 0; transform: scale(0.95); }
    to { opacity: 1; transform: scale(1); }
  }
</style>
