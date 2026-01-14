<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { API_CONFIG, getApiUrl } from '$lib/config';

  let activeTab: 'users' | 'gallery' = 'users';
  let users: any[] = [];
  let strips: any[] = [];
  let loading = true;
  let error = '';
  let token: string | null = null;

  // Filter state
  let selectedUserId: string | null = null;

  onMount(async () => {
    token = localStorage.getItem('sb_token');
    if (!token) {
      goto('/auth/login');
      return;
    }
    await loadUsers();
    await loadStrips(); // Initially load recent strips
    loading = false;
  });

  async function loadUsers() {
    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_USERS, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (!res.ok) {
        if (res.status === 403) throw new Error('Unauthorized: Admin access required');
        throw new Error('Failed to load users');
      }
      users = await res.json();
    } catch (e: any) {
      error = e.message;
      if (error.includes('Unauthorized')) setTimeout(() => goto('/photobooth'), 2000);
    }
  }

  async function loadStrips(userId?: string) {
    try {
      let url = API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_STRIPS;
      if (userId) {
        url += `?user_id=${userId}`;
      }
      const res = await fetch(url, {
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) {
        strips = await res.json();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deleteUser(id: string) {
    if (!confirm('Are you sure? This will delete the user and ALL their photos permanently.')) return;
    
    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_USER + id, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) {
        users = users.filter(u => u.id !== id);
        // Also refresh strips if we were looking at that user
        if (selectedUserId === id) {
             selectedUserId = null;
             loadStrips();
        } else {
             loadStrips(selectedUserId || undefined);
        }
      } else {
        alert('Failed to delete user');
      }
    } catch (e) {
      alert('Error deleting user');
    }
  }

  // Create User State
  let showCreateModal = false;
  let newUser = { username: '', email: '', password: '', is_admin: false };

  async function createUser() {
    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_USERS, {
        method: 'POST',
        headers: { 
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(newUser)
      });
      
      if (res.ok) {
        alert('User created successfully');
        showCreateModal = false;
        newUser = { username: '', email: '', password: '', is_admin: false };
        loadUsers();
      } else {
        const d = await res.json();
        alert(d.error || 'Failed to create user');
      }
    } catch (e) {
      alert('Error creating user');
    }
  }

  async function deleteStrip(id: string) {
    if (!confirm('Are you sure you want to delete this photo?')) return;
    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_DELETE_STRIP + id, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      });
      if (res.ok) {
        strips = strips.filter(s => s.id !== id);
      } else {
        alert('Failed to delete photo');
      }
    } catch (e) {
      alert('Error deleting photo');
    }
  }

  // Reset Password State
  let showResetPasswordModal = false;
  let resetUserId: string | null = null;
  let resetPasswordStr = '';

  function openResetPasswordModal(id: string) {
    resetUserId = id;
    showResetPasswordModal = true;
    resetPasswordStr = '';
  }

  async function resetPassword() {
    if (!resetUserId || !resetPasswordStr) return;
    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_RESET_PASSWORD + `${resetUserId}/password`, {
        method: 'PATCH',
        headers: { 
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ password: resetPasswordStr })
      });
      
      if (res.ok) {
        alert('Password updated successfully');
        showResetPasswordModal = false;
        resetUserId = null;
        resetPasswordStr = '';
      } else {
        const d = await res.json();
        alert(d.error || 'Failed to update password');
      }
    } catch (e) {
      alert('Error updating password');
    }
  }

  async function toggleAdmin(user: any) {
    const newStatus = !user.is_admin;
    const action = newStatus ? 'promote this user to Admin' : 'demote this user to regular User';
    if (!confirm(`Are you sure you want to ${action}?`)) return;

    try {
      const res = await fetch(API_CONFIG.BASE_URL + API_CONFIG.ENDPOINTS.ADMIN_UPDATE_ROLE + user.id + '/role', {
        method: 'PATCH',
        headers: { 
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ is_admin: newStatus })
      });
      
      if (res.ok) {
        user.is_admin = newStatus;
        users = users; // Trigger reactivity
      } else {
        const d = await res.json();
        alert(d.error || 'Failed to update role');
      }
    } catch (e) {
      alert('Error updating role');
    }
  }

  let viewingStrip: any = null;

  async function downloadImage(url: string, title: string) {
    try {
      const response = await fetch(url);
      const blob = await response.blob();
      const objectUrl = URL.createObjectURL(blob);
      
      const a = document.createElement('a');
      a.href = objectUrl;
      a.download = `${title || 'photo'}.png`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(objectUrl);
    } catch (e) {
      console.error('Download failed:', e);
      window.open(url, '_blank');
    }
  }

  function viewUserGallery(id: string) {
    selectedUserId = id;
    activeTab = 'gallery';
    loadStrips(id);
  }

  function formatDate(str: string) {
    return new Date(str).toLocaleDateString() + ' ' + new Date(str).toLocaleTimeString();
  }
</script>

<div class="h-[100dvh] bg-[#f8f2ff] flex flex-col font-sans overflow-hidden fixed inset-0 w-full">
  <header class="bg-white border-b border-purple-100 p-4 sticky top-0 z-10 flex justify-between items-center shadow-sm shrink-0">
    <div class="flex items-center gap-3">
       <div class="w-2.5 h-2.5 md:w-3 md:h-3 rounded-full bg-red-500 animate-pulse"></div>
       <h1 class="text-base md:text-xl font-bold text-slate-800 tracking-tight">Admin Dashboard</h1>
    </div>
    <div class="flex items-center gap-3 md:gap-4">
       <div class="text-[10px] md:text-xs font-bold uppercase tracking-widest text-slate-400 hidden lg:block">
         {users.length} Users • {strips.length} Photos
       </div>
       <button 
         on:click={() => showCreateModal = true}
         class="text-xs font-bold text-white bg-purple-600 hover:bg-purple-700 px-4 py-2 rounded-lg transition-colors shadow-lg shadow-purple-200"
       >
         + New User
       </button>
       <button 
         on:click={() => goto('/home')}
         class="text-xs font-bold text-purple-600 hover:bg-purple-50 px-4 py-2 rounded-lg transition-colors border border-purple-200"
       >
         Back to Home
       </button>
    </div>
  </header>

  <main class="flex-1 p-3 md:p-8 max-w-7xl mx-auto w-full relative overflow-y-auto overflow-x-hidden">
    {#if showCreateModal}
      <!-- ... existing create user modal ... -->
      <div class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in">
        <div class="bg-white rounded-3xl p-8 w-full max-w-md shadow-2xl animate-in zoom-in-95">
          <h2 class="text-xl font-bold text-slate-800 mb-6">Create New Account</h2>
          <div class="flex flex-col gap-4">
            <div>
              <label for="username" class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Username</label>
              <input id="username" bind:value={newUser.username} class="w-full bg-slate-50 border border-slate-200 rounded-xl p-3 text-sm font-bold text-slate-700 outline-none focus:border-purple-400" />
            </div>
            <div>
              <label for="email" class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Email</label>
              <input id="email" bind:value={newUser.email} type="email" class="w-full bg-slate-50 border border-slate-200 rounded-xl p-3 text-sm font-bold text-slate-700 outline-none focus:border-purple-400" />
            </div>
             <div>
              <label for="password" class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Password</label>
              <input id="password" bind:value={newUser.password} type="password" class="w-full bg-slate-50 border border-slate-200 rounded-xl p-3 text-sm font-bold text-slate-700 outline-none focus:border-purple-400" />
            </div>
            <div class="flex items-center justify-between p-2">
              <span id="admin-label" class="text-xs font-bold text-slate-600">Grant Admin Privileges?</span>
              <button 
                aria-labelledby="admin-label"
                class="w-12 h-6 rounded-full transition-colors relative {newUser.is_admin ? 'bg-purple-500' : 'bg-slate-200'}"
                on:click={() => newUser.is_admin = !newUser.is_admin}
              >
                <div class="w-4 h-4 bg-white rounded-full absolute top-1 transition-all {newUser.is_admin ? 'left-7' : 'left-1'}"></div>
              </button>
            </div>
            
            <div class="flex gap-3 mt-4">
              <button on:click={() => showCreateModal = false} class="flex-1 py-3 font-bold text-[10px] text-slate-400 hover:bg-slate-50 rounded-xl transition-colors uppercase tracking-wider">Cancel</button>
              <button on:click={createUser} class="flex-1 py-3 font-bold text-[10px] text-white bg-purple-600 hover:bg-purple-700 rounded-xl shadow-lg shadow-purple-200 transition-colors uppercase tracking-wider">Create</button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    {#if showResetPasswordModal}
      <div class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4 animate-in fade-in">
        <div class="bg-white rounded-3xl p-8 w-full max-w-md shadow-2xl animate-in zoom-in-95">
          <h2 class="text-xl font-bold text-slate-800 mb-6">Reset Password</h2>
          <div class="flex flex-col gap-4">
             <div>
              <label for="new-password" class="text-[10px] font-bold uppercase tracking-widest text-slate-400">New Password</label>
              <input id="new-password" bind:value={resetPasswordStr} type="password" class="w-full bg-slate-50 border border-slate-200 rounded-xl p-3 text-sm font-bold text-slate-700 outline-none focus:border-purple-400" placeholder="Minimum 6 characters" />
            </div>
            
            <div class="flex gap-3 mt-4">
              <button on:click={() => showResetPasswordModal = false} class="flex-1 py-3 font-bold text-[10px] text-slate-400 hover:bg-slate-50 rounded-xl transition-colors uppercase tracking-wider">Cancel</button>
              <button on:click={resetPassword} class="flex-1 py-3 font-bold text-[10px] text-white bg-purple-600 hover:bg-purple-700 rounded-xl shadow-lg shadow-purple-200 transition-colors uppercase tracking-wider">Update</button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    {#if viewingStrip}
      <div 
        class="fixed inset-0 z-50 bg-black/95 backdrop-blur-md flex flex-col items-center justify-center p-4 animate-in fade-in duration-200"
        on:click|self={() => viewingStrip = null}
        role="dialog"
      >
         <button 
           class="absolute top-4 right-4 text-white hover:text-purple-400 transition-colors p-2 z-50"
           on:click={() => viewingStrip = null}
         >
           <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
         </button>

         <img 
           src={viewingStrip.file_url} 
           alt="Fullscreen Strip" 
           class="max-w-full max-h-[85vh] object-contain rounded-lg shadow-2xl"
         />
         
         <div class="mt-6 flex gap-4">
           <button 
             on:click={() => downloadImage(viewingStrip.file_url, viewingStrip.title)}
             class="bg-white text-purple-600 hover:bg-purple-50 px-6 py-2 rounded-full font-bold uppercase tracking-widest text-xs transition-colors shadow-lg"
           >
             Download Photo
           </button>
         </div>
      </div>
    {/if}

    {#if error}
      <div class="bg-red-50 text-red-600 p-4 rounded-xl border border-red-100 mb-6 text-center shadow-sm">
        <p class="font-bold text-sm">{error}</p>
      </div>
    {/if}

    <div class="flex flex-col gap-6">
      <!-- Tabs -->
      <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-hide">
        <button 
          class="whitespace-nowrap px-5 py-2 rounded-full text-[10px] md:text-xs font-bold uppercase tracking-widest transition-all {activeTab === 'users' ? 'bg-purple-600 text-white shadow-md shadow-purple-200' : 'bg-white text-slate-400 hover:text-purple-500'}"
          on:click={() => activeTab = 'users'}
        >
          User Management
        </button>
        <button 
          class="whitespace-nowrap px-5 py-2 rounded-full text-[10px] md:text-xs font-bold uppercase tracking-widest transition-all {activeTab === 'gallery' ? 'bg-purple-600 text-white shadow-md shadow-purple-200' : 'bg-white text-slate-400 hover:text-purple-500'}"
          on:click={() => { activeTab = 'gallery'; selectedUserId = null; loadStrips(); }}
        >
          Gallery Monitor
        </button>
      </div>

      <!-- Content -->
      <div class="bg-white rounded-3xl p-4 md:p-6 shadow-xl shadow-purple-100/50 border border-purple-50 min-h-[500px]">
        {#if loading}
          <div class="flex items-center justify-center h-full py-20">
            <div class="w-8 h-8 border-4 border-purple-100 border-t-purple-500 rounded-full animate-spin"></div>
          </div>
        {:else if activeTab === 'users'}
          <div class="overflow-x-auto">
            <table class="w-full text-left border-collapse">
              <thead>
                <tr class="border-b border-purple-100">
                  <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest">User</th>
                  <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest hidden md:table-cell">Email</th>
                  <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest">Role</th>
                  <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest hidden lg:table-cell">Created</th>
                  <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest text-right">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-slate-50">
                {#each users as user}
                  <tr class="hover:bg-purple-50/50 transition-colors group text-sm">
                    <td class="p-2 md:p-4">
                      <div class="font-bold text-slate-700">{user.username}</div>
                      <div class="text-[11px] text-slate-400 font-mono hidden md:block mt-0.5">{user.id}</div>
                      <!-- Mobile email fallback -->
                      <div class="text-[10px] text-slate-500 md:hidden truncate max-w-[100px]">{user.email}</div>
                    </td>
                    <td class="p-2 md:p-4 text-slate-600 hidden md:table-cell">{user.email}</td>
                    <td class="p-2 md:p-4">
                       {#if user.email === 'wuby@superuser.com'}
                         <span class="bg-purple-100/50 text-purple-400 px-2 py-1 rounded text-[10px] font-bold uppercase tracking-wide cursor-not-allowed" title="Superuser role cannot be changed">Admin</span>
                       {:else if user.is_admin}
                         <button 
                           on:click={() => toggleAdmin(user)}
                           class="bg-purple-100 text-purple-700 hover:bg-purple-200 transition-colors px-2 py-1 rounded text-[10px] font-bold uppercase tracking-wide cursor-pointer"
                           title="Click to demote to User"
                         >
                           Admin
                         </button>
                       {:else}
                         <button 
                           on:click={() => toggleAdmin(user)}
                           class="bg-slate-100 text-slate-500 hover:bg-slate-200 transition-colors px-2 py-1 rounded text-[10px] font-bold uppercase tracking-wide cursor-pointer"
                           title="Click to promote to Admin"
                         >
                           User
                         </button>
                       {/if}
                    </td>
                    <td class="p-2 md:p-4 text-xs text-slate-400 hidden lg:table-cell">{formatDate(user.created_at)}</td>
                    <td class="p-2 md:p-4 text-right flex flex-row justify-end gap-2 items-center">
                      <button 
                        on:click={() => viewUserGallery(user.id)}
                        class="p-2 rounded-lg bg-purple-50 text-purple-600 hover:bg-purple-100 transition-colors"
                        title="View Photos"
                      >
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/></svg>
                      </button>
                      <button 
                        on:click={() => openResetPasswordModal(user.id)}
                        class="p-2 rounded-lg bg-yellow-50 text-yellow-600 hover:bg-yellow-100 transition-colors"
                        title="Reset Password"
                      >
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/></svg>
                      </button>
                      {#if !user.is_admin}
                        <button 
                          on:click={() => deleteUser(user.id)}
                          class="p-2 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors"
                          title="Delete User"
                        >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                        </button>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {:else if activeTab === 'gallery'}
          <div class="flex flex-col gap-3">
            {#if selectedUserId}
              <div class="flex items-center gap-2 mb-2 p-3 bg-purple-50 rounded-xl border border-purple-100">
                <button 
                  aria-label="Clear filter"
                  on:click={() => { selectedUserId = null; loadStrips(); }}
                  class="bg-white p-1 rounded-full hover:bg-purple-200 transition-colors"
                >
                   <svg class="w-3 h-3 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
                </button>
                <div class="flex flex-col">
                  <span class="text-[9px] font-bold uppercase tracking-widest text-purple-400">Filtering User</span>
                  <span class="font-mono text-[10px] md:text-xs text-purple-800 break-all">{selectedUserId}</span>
                </div>
              </div>
            {/if}

            <div class="overflow-x-auto">
              <table class="w-full text-left border-collapse">
                <thead>
                  <tr class="border-b border-purple-100">
                    <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest w-20">Photo</th>
                    <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest">Details</th>
                    <th class="p-2 md:p-4 text-[10px] md:text-xs font-bold text-purple-400 uppercase tracking-widest text-right">Actions</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-slate-50">
                  {#each strips as strip}
                    <tr class="hover:bg-purple-50/50 transition-colors group text-sm">
                      <td class="p-2 md:p-4">
                        <button 
                          on:click={() => viewingStrip = strip}
                          class="w-12 h-16 md:w-16 md:h-20 bg-slate-100 rounded-lg overflow-hidden border border-slate-200 hover:border-purple-300 transition-all hover:scale-105"
                        >
                          <img src={strip.file_url} alt="Strip" class="w-full h-full object-cover" loading="lazy" />
                        </button>
                      </td>
                      <td class="p-2 md:p-4">
                         <div class="flex flex-col gap-0.5">
                            <span class="font-bold text-slate-700 text-xs">{strip.User ? strip.User.username : 'Guest'}</span>
                            <span class="text-[10px] text-slate-400">{formatDate(strip.created_at)}</span>
                            <span class="font-mono text-[9px] text-slate-300 hidden md:inline-block">{strip.id}</span>
                         </div>
                      </td>
                      <td class="p-2 md:p-4 text-right flex flex-row justify-end gap-2 items-center h-full">
                        <button 
                          on:click={() => viewingStrip = strip}
                          class="p-2 rounded-lg bg-purple-50 text-purple-600 hover:bg-purple-100 transition-colors"
                          title="View"
                        >
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/></svg>
                        </button>
                        <button 
                          on:click={() => downloadImage(strip.file_url, strip.title)}
                          class="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
                          title="Download"
                        >
             <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/></svg>
                        </button>
                        <button 
                          on:click={() => deleteStrip(strip.id)}
                          class="p-2 rounded-lg bg-red-50 text-red-600 hover:bg-red-100 transition-colors"
                          title="Delete"
                        >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                        </button>
                      </td>
                    </tr>
                  {/each}
                  {#if strips.length === 0}
                    <tr>
                      <td colspan="3" class="py-10 text-center text-slate-300 font-bold uppercase tracking-widest text-xs">
                        No photos found
                      </td>
                    </tr>
                  {/if}
                </tbody>
              </table>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </main>
</div>
