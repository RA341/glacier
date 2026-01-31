<script lang="ts">
    import {CloudIcon, RefreshCcwDot} from '@lucide/svelte';
    import {fade} from 'svelte/transition';
    import {frostCli, glacierCli} from "$lib/api/api";
    import {FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import DownloadItem from "./DownloadItem.svelte";
    import {onMount} from "svelte";

    const frostLib = frostCli(FrostLibraryService)
    let downloadingRpc = createRPCRunner(() => frostLib.listDownloading({}))

    function refresh() {
        downloadingRpc.runner()
    }

    $effect(() => {
        const interval = setInterval(() => {
            refresh()
        }, 5000);
        return () => clearInterval(interval);
    });

    onMount(() => {
        refresh()
    })
</script>

<button
        onclick={refresh}
        disabled={downloadingRpc.loading}
        class="flex items-center gap-2 p-2 px-4 rounded-xl bg-panel border border-border text-xs font-bold text-muted hover:text-frost-400 hover:border-frost-500/50 transition-all active:scale-95 disabled:opacity-50"
>
    <RefreshCcwDot size={14} class={downloadingRpc.loading ? 'animate-spin text-frost-500' : ''}/>
    Refresh
</button>

<div class="space-y-6" in:fade={{ duration: 200 }}>
    {#if downloadingRpc.error}
        <div>{downloadingRpc.error}</div>
    {:else if downloadingRpc.error}
        <div>Loading...</div>
    {:else}
        {#if !downloadingRpc.value?.downloads}
            <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-frost-500/20 bg-frost-500/5 rounded-3xl text-frost-400/30">
                <CloudIcon size={48} strokeWidth={1} class="mb-4"/>
                <h2 class="text-xl font-bold text-frost-400/50">Frost downloads</h2>
            </div>
        {:else}
            {#each Object.entries(downloadingRpc.value?.downloads ?? {}) as [key, ff]}
                <DownloadItem ID={key} detail={ff}/>
            {/each}
        {/if}
    {/if}
</div>