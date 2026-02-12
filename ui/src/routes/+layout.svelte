<script lang="ts">
    import './layout.css';
    import {page} from '$app/state';
    import {appName} from "$lib/api/api";
    import Snackbar from "$lib/components/snackbar/Snackbar.svelte";
    import Sidebar from './sidebar.svelte';
    import User from "$lib/components/user/User.svelte";

    let {children} = $props();
</script>

<svelte:head>
    <title>{appName}</title>
    <link rel="icon" href={"/favicon.svg"}/>
</svelte:head>

<Snackbar>
    <User>
        {#if !page.url.pathname.startsWith('/auth')}
            <Sidebar/>
        {/if}
        <!-- Main Content Area -->
        <main class="relative flex-1 overflow-y-auto">
            {@render children()}
        </main>
    </User>
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

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
