<script lang="ts">
    import {page} from '$app/state';
    import {goto} from '$app/navigation';
    import {fly} from 'svelte/transition';
    import {LockIcon, ShieldCheckIcon, UserPlusIcon} from '@lucide/svelte';
    import {appName} from "$lib/api/api";

    let {children} = $props();

    let activeTab = $derived(page.url.pathname.includes('register') ? 'register' : 'login');

    function navigate(path: string) {
        goto(`/auth/${path}`, {replaceState: true});
    }
</script>

<div class="min-h-screen w-full bg-background flex items-center justify-center p-6 relative overflow-hidden">
    <!-- Decorative background elements -->
    <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-frost-500/10 blur-[120px] rounded-full"></div>
    <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-frost-900/20 blur-[120px] rounded-full"></div>

    <div
            class="w-full max-w-110 z-10"
            transition:fly={{ y: 20, duration: 600 }}
    >
        <!-- Logo / Brand Section -->
        <div class="flex flex-col items-center gap-3 mb-8">
            <img src="/favicon.svg" alt="logo" height="60" width="60"/>
            <h1 class="text-3xl font-black tracking-tighter text-foreground uppercase">{appName}</h1>
        </div>

        <!-- Auth Card -->
        <div class="bg-surface border border-border rounded-4xl shadow-2xl overflow-hidden flex flex-col">

            <!-- Tab Switcher -->
            <div class="p-2 bg-panel/30 border-b border-border flex gap-1">
                <button
                        onclick={() => navigate('login')}
                        class="flex-1 flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-bold transition-all
                    {activeTab === 'login' ? 'bg-surface border border-border text-frost-400 shadow-sm' : 'text-muted hover:text-foreground hover:bg-panel'}"
                >
                    <LockIcon size={16}/>
                    Login
                </button>
                <button
                        onclick={() => navigate('register')}
                        class="flex-1 flex items-center justify-center gap-2 py-3 rounded-2xl text-sm font-bold transition-all
                    {activeTab === 'register' ? 'bg-surface border border-border text-frost-400 shadow-sm' : 'text-muted hover:text-foreground hover:bg-panel'}"
                >
                    <UserPlusIcon size={16}/>
                    Register
                </button>
            </div>

            <!-- Page Content -->
            <main class="p-8">
                {@render children()}
            </main>
        </div>

        <p class="mt-8 text-center text-xs text-muted/50 font-medium uppercase tracking-[0.2em]">
            TODO FOOTER MESSAGE LIKE JELLYFIN
        </p>
    </div>
</div>

<style>
    :global(body) {
        background-color: #0a0a0c;
    }
</style>
