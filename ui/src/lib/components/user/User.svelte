<script lang="ts">
    import {onMount} from 'svelte';
    import {setUserCtx} from "$lib/components/user/provider.svelte";

    let {children} = $props();

    const userCtx = setUserCtx();
    onMount(() => {
        userCtx.check();
    });
</script>

<div class="flex h-screen w-full overflow-hidden bg-background text-foreground">
    {#if userCtx.isLoading}
        <!-- Loading Spinner -->
        <div class="flex items-center justify-center w-full h-full">
            <div class="flex flex-col items-center gap-4">
                <div class="spinner"></div>
                <p class="text-muted text-sm">Checking authentication...</p>
            </div>
        </div>
    {:else if userCtx.error}
        <div class="flex items-center justify-center w-full h-full">
            <div class="flex flex-col items-center gap-4">
                <p class="text-muted text-sm">Could not checked auth status, is the glacier running ?...</p>
                <p class="text-muted text-sm">{userCtx.error}</p>
            </div>
        </div>
    {:else}
        {@render children()}
    {/if}
</div>

<style>

    /* Loading Spinner */
    .spinner {
        width: 48px;
        height: 48px;
        border: 4px solid var(--color-border);
        border-top-color: var(--color-frost-400, #60a5fa);
        border-radius: 50%;
        animation: spin 0.8s linear infinite;
    }
</style>
