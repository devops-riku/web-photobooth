<script lang="ts">
  import { goto } from '$app/navigation';
  import { getApiUrl } from '$lib/config';

  import { page } from '$app/stores';

  let username = '';
  let email = '';
  let password = '';
  let isLoading = false;
  let message = '';
  let messageType: 'success' | 'error' = 'success';

  async function handleSignup() {
    isLoading = true;
    message = '';
    
    try {
      const response = await fetch(getApiUrl('SIGNUP'), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, email, password })
      });

      const data = await response.json();

      if (!response.ok) {
        message = data.error || 'Signup failed';
        messageType = 'error';
      } else {
        message = 'Account created! Redirecting to login...';
        messageType = 'success';
        
        const redirect = $page.url.searchParams.get('redirect');
        setTimeout(() => {
          goto(redirect ? `/auth/login?redirect=${redirect}` : '/auth/login');
        }, 1000);
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
  <div class="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-purple-100/40 rounded-full blur-[120px]"></div>
  <div class="absolute bottom-[-10%] left-[-10%] w-[50%] h-[50%] bg-purple-200/30 rounded-full blur-[120px]"></div>

  <div class="w-full max-w-sm flex flex-col items-center relative z-10">
    <!-- Header -->
    <div class="flex flex-col items-center gap-6 mb-16 text-center animate-in">
      <div class="w-14 h-14 bg-white rounded-[2rem] flex items-center justify-center shadow-xl shadow-purple-200/50">
        <svg class="w-7 h-7 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
        </svg>
      </div>
      <div>
        <h1 class="text-3xl font-light text-purple-900 tracking-tight">Create Account</h1>
        <p class="text-[10px] font-bold uppercase tracking-[0.3em] text-purple-300 mt-2">Join the Wuby community</p>
      </div>
    </div>

    <form on:submit|preventDefault={handleSignup} class="w-full flex flex-col gap-8 animate-in delay-100">
      <div class="space-y-4">
        <div class="flex flex-col gap-2">
          <input 
            id="username"
            type="text" 
            bind:value={username}
            placeholder="Username"
            class="w-full bg-white/50 backdrop-blur-sm border-2 border-purple-50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 placeholder:text-purple-200 focus:border-purple-200 focus:bg-white focus:outline-none transition-all"
            required
          />
        </div>

        <div class="flex flex-col gap-2">
          <input 
            id="email"
            type="email" 
            bind:value={email}
            placeholder="Email Address"
            class="w-full bg-white/50 backdrop-blur-sm border-2 border-purple-50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 placeholder:text-purple-200 focus:border-purple-200 focus:bg-white focus:outline-none transition-all"
            required
          />
        </div>

        <div class="flex flex-col gap-2">
          <input 
            id="password"
            type="password" 
            bind:value={password}
            placeholder="Choose Password"
            class="w-full bg-white/50 backdrop-blur-sm border-2 border-purple-50 rounded-2xl px-6 py-4 text-sm font-medium text-purple-900 placeholder:text-purple-200 focus:border-purple-200 focus:bg-white focus:outline-none transition-all"
            required
            minlength="6"
          />
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
          Register
        {/if}
      </button>
    </form>

    <div class="mt-12 flex flex-col items-center gap-6 animate-in delay-200">
      <p class="text-[10px] font-bold text-purple-300 uppercase tracking-widest">
        Have an account? 
        <button on:click={() => {
          const redirect = $page.url.searchParams.get('redirect');
          goto(redirect ? `/auth/login?redirect=${redirect}` : '/auth/login');
        }} class="text-purple-500 hover:text-purple-700 underline underline-offset-4 ml-1">Sign In</button>
      </p>
      
      <button 
        on:click={() => goto('/photobooth')}
        class="text-[9px] font-bold text-purple-200 uppercase tracking-[0.2em] hover:text-purple-400 transition-colors"
      >
        Continue as Guest
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
