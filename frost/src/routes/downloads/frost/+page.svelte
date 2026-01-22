<script lang="ts">
    import {CloudIcon} from '@lucide/svelte';
    import {fade} from 'svelte/transition';
    import {frostCli, glacierCli} from "$lib/api/api";
    import {FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";

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

</script>

<div class="space-y-6" in:fade={{ duration: 200 }}>
    {#if downloadingRpc.error}
        <div>{downloadingRpc.error}</div>
    {:else if downloadingRpc.error}
        <div>Loading...</div>
    {:else}
        <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-frost-500/20 bg-frost-500/5 rounded-3xl text-frost-400/30">
            <CloudIcon size={48} strokeWidth={1} class="mb-4"/>
            <h2 class="text-xl font-bold text-frost-400/50">Frost Instance</h2>

            {#each downloadingRpc.value?.downloads?.files as ff}
                <div>
                    {ff.Name}
                </div>
            {/each}
        </div>
    {/if}
</div>