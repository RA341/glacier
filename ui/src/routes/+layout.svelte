<script lang="ts">
    import './layout.css';
    import {page} from '$app/state';
    import {appName, Glacier} from "$lib/api/api";
    import Snackbar from "$lib/components/snackbar/Snackbar.svelte";
    import {goto} from "$app/navigation";
    import {onMount} from 'svelte';
    import Sidebar from './sidebar.svelte';

    let {children} = $props();

    let isCheckingAuth = $state(true);

    onMount(async () => {
        await checkAuth();
    });

    async function checkAuth() {
        const isAuthPage = page.url.pathname.startsWith('/auth');

        if (isAuthPage) {
            isCheckingAuth = false;
            return;
        }

        const isRoot = page.url.pathname === '/';

        try {
            const res = await fetch(`${Glacier.base}/ping`, {
                credentials: 'include'
            });

            if (!res.ok) {
                await goto('/auth/login', {replaceState: true});

            } else if (isRoot) {
                await goto('/library', {replaceState: true});
            }
        } catch (e: any) {
            console.log('error fetching ping', e);
            await goto('/auth/login', {replaceState: true});
        }

        isCheckingAuth = false
    }
</script>

<svelte:head>
    <title>{appName}</title>
    <link rel="icon" href={"/favicon.svg"}/>
</svelte:head>

<Snackbar>
    <div class="flex h-screen w-full overflow-hidden bg-background text-foreground">
        {#if isCheckingAuth}
            <!-- Loading Spinner -->
            <div class="flex items-center justify-center w-full h-full">
                <div class="flex flex-col items-center gap-4">
                    <div class="spinner"></div>
                    <p class="text-muted text-sm">Checking authentication...</p>
                </div>
            </div>
        {:else}
            {#if !page.url.pathname.startsWith('/auth')}
                <Sidebar/>
            {/if}
            <!-- Main Content Area -->
            <main class="relative flex-1 overflow-y-auto">
                {@render children()}
            </main>
        {/if}
    </div>
</Snackbar>

<style>
    ::-webkit-scrollbar {
        width: 6px;
    }

    ::-webkit-scrollbar-track {
        background: transparent;
    }

    ::-webkit-scrollbar-thumb {
        background: var(--border);
        border-radius: 10px;
    }

    ::-webkit-scrollbar-thumb:hover {
        background: #444;
    }

    /* Loading Spinner */
    .spinner {
        width: 48px;
        height: 48px;
        border: 4px solid var(--color-border);
        border-top-color: var(--color-frost-400, #60a5fa);
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
