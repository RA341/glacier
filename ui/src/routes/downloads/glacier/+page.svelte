<script lang="ts">
    import {
        CircleAlert,
        CircleCheck,
        DownloadIcon,
        ImageIcon,
        LoaderIcon,
        PauseIcon,
        RefreshCcwDot,
        ServerIcon
    } from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';
    import {glacierCli} from "$lib/api/api";
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {calculatePercent, formatBytes} from "$lib/api/byte-math";

    const libSrv = glacierCli(LibraryService)

    let libRpc = createRPCRunner(() => libSrv.triggerTracker({}))

    function trigger() {
        libRpc.runner()
    }

    let activeDownloadsRpc = createRPCRunner(() => libSrv.listWithState({state: "downloading"}))

    function refresh() {
        activeDownloadsRpc.runner()
    }

    onMount(() => {
        refresh()
    })

    $effect(() => {
        const interval = setInterval(() => {
            refresh()
        }, 5000);
        return () => clearInterval(interval);
    });

    function getStatusTheme(state: string) {
        if (state?.toLowerCase().includes('paused')) return {color: 'bg-amber-500', icon: PauseIcon};
        if (state?.toLowerCase().includes('complete')) return {color: 'bg-green-500', icon: CircleCheck};
        return {color: 'bg-frost-500', icon: DownloadIcon};
    }
</script>
<div class="space-y-6" in:fade={{ duration: 200 }}>
    <!-- Toolbar Section -->
    <div class="flex items-center justify-between px-2">
        <h2 class="text-sm font-bold uppercase tracking-widest text-muted">Active Downloads</h2>
        <button
                onclick={refresh}
                disabled={libRpc.loading}
                class="flex items-center gap-2 p-2 px-4 rounded-xl bg-panel border border-border text-xs font-bold text-muted hover:text-frost-400 hover:border-frost-500/50 transition-all active:scale-95 disabled:opacity-50"
        >
            <RefreshCcwDot size={14} class={libRpc.loading ? 'animate-spin text-frost-500' : ''}/>
            Check Status
        </button>
    </div>

    <div class="min-h-100">
        {#if activeDownloadsRpc.error}
            <!-- Error State -->
            <div class="flex flex-col items-center justify-center h-96 text-red-400 gap-3">
                <CircleAlert size={48} strokeWidth={1}/>
                <h3 class="font-bold">Connection Error</h3>
                <p class="text-xs opacity-80">{activeDownloadsRpc.error}</p>
            </div>
        {:else if activeDownloadsRpc.loading && !activeDownloadsRpc.value}
            <div class="flex flex-col items-center justify-center h-96 text-muted gap-4">
                <LoaderIcon class="animate-spin text-frost-500" size={40}/>
                <p class="animate-pulse text-sm font-medium">Fetching downloads...</p>
            </div>
        {:else if activeDownloadsRpc.value?.game && activeDownloadsRpc.value.game.length === 0}
            <!-- Empty State -->
            <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-border rounded-3xl text-muted/30">
                <ServerIcon size={48} strokeWidth={1} class="mb-4"/>
                <h2 class="text-xl font-bold text-foreground/50">No Active Downloads</h2>
                <p class="text-sm">Manage your library to trigger new transfers.</p>
            </div>
        {:else}
            <!-- Downloads List -->
            <div class="grid grid-cols-1 gap-4">
                {#each activeDownloadsRpc.value?.game ?? [] as download (download.ID)}
                    {@const theme = getStatusTheme(download.DownloadState?.State ?? "Unknown")}
                    <!-- Calculate actual percentage for the progress bar -->
                    {@const progressPercent = calculatePercent(
                        download.DownloadState?.Complete,
                        download.DownloadState?.Done
                    )}
                    {@const
                        totalBytes = Number(download.DownloadState?.Complete || 0) + Number(download.DownloadState?.Done || 0)}

                    <div
                            class="group bg-surface border border-border rounded-3xl p-5 flex items-center gap-6 overflow-hidden hover:border-frost-500/30 transition-all shadow-sm"
                            in:fly={{ y: 10, duration: 300 }}
                    >
                        <div class="w-20 h-28 shrink-0 rounded-xl border border-border bg-panel overflow-hidden shadow-inner flex items-center justify-center">
                            {#if download.Meta?.ThumbnailURL}
                                <img src={download.Meta.ThumbnailURL} alt=""
                                     class="w-full h-full object-cover transition-transform group-hover:scale-105 duration-500"/>
                            {:else}
                                <div class="text-muted/20">
                                    <ImageIcon size={32} strokeWidth={1}/>
                                </div>
                            {/if}
                        </div>

                        <!-- Main Content -->
                        <div class="flex-1 min-w-0 space-y-4">
                            <!-- Title & State -->
                            <div class="flex items-center justify-between gap-4">
                                <h3 class="font-bold text-lg text-foreground truncate">
                                    {download.Meta?.Name || 'Unknown Title'}
                                </h3>
                                <span class="shrink-0 px-3 py-1 rounded-full bg-panel border border-border text-[10px] font-bold uppercase tracking-wider text-muted flex items-center gap-2">
                                <span class="w-1.5 h-1.5 rounded-full {theme.color} {download.DownloadState?.State === 'Downloading' ? 'animate-pulse' : ''}"></span>
                                    {download.DownloadState?.State || 'Connecting'}
                                </span>
                            </div>

                            <!-- Progress Bar Container -->
                            <div class="space-y-2">
                                <div class="flex justify-between items-end text-xs font-mono text-muted mb-1">
                                    <!-- Left side: Speed / ETA string from API -->
                                    <span class="truncate opacity-70">
                                        {download.DownloadState?.Progress || 'Initializing...'}
                                    </span>

                                    <!-- Right side: Bytes Complete / Total -->
                                    <span class="shrink-0 font-bold text-foreground">
                                            {formatBytes(Number(download.DownloadState?.Complete || 0))}
                                        <span class="text-muted/30 mx-1">/</span>
                                        {formatBytes(totalBytes)}
                                    </span>
                                </div>

                                <!-- Visual Progress Bar -->
                                <div class="h-2 w-full bg-panel rounded-full overflow-hidden border border-border/50">
                                    <div class="h-full {theme.color} transition-all duration-700 ease-out shadow-[0_0_8px_rgba(130,170,255,0.2)]"
                                         style="width: {progressPercent}%"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<style>
    /* Add a subtle glass effect to the cards */
    .bg-surface {
        background-color: rgba(var(--surface-rgb), 0.6);
        backdrop-filter: blur(8px);
    }
</style>
