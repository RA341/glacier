<script lang="ts">
    import {FileTextIcon, ImageIcon, SearchIcon, ShieldCheckIcon, XIcon} from '@lucide/svelte';
    import {fade} from "svelte/transition";
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import IndexerSearch from "$lib/components/IndexerSearch.svelte";
    import FileManager from "./FileManager.svelte";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    let {game = $bindable(null)}: { game: Game | null; refresh: () => Promise<void> } = $props();

    const lib = glacierCli(LibraryService)
    const sm = getSnackbarCtx()

    async function handleRedownload() {
        if (!game?.ID) {
            return
        }

        const {err} = await callRPC(() => lib.redownload({gameId: game.ID}))
        if (err) {
            sm.push(`Could not redownload ${err}`, 'error')
        }
        await refresh()

    }

</script>

<div class="space-y-8" in:fade>
    <!-- Download Section -->
    <section class="space-y-4">
        <h2 class="text-sm font-bold uppercase tracking-widest text-muted px-2">Download</h2>
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 p-6 bg-surface border border-border rounded-3xl">
            <!-- State Info -->
            <div class="lg:col-span-5 grid grid-cols-2 gap-y-6 gap-x-4 border-r border-border pr-6">
                <div>
                    <button onclick={handleRedownload} class="text-frost-400 hover:text-foreground p-1">
                        Retry
                    </button>
                </div>
                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">Client Type</p>
                    <p class="text-sm font-bold">{game?.DownloadState?.Client}</p>
                </div>
                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">Download ID</p>
                    <p class="text-sm font-mono opacity-60 truncate">{game?.DownloadState?.DownloadId}</p>
                </div>
                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">State</p>
                    <p class="text-sm font-bold text-frost-400">{game?.DownloadState?.State}</p>
                </div>
                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">Progress</p>
                    <p class="text-sm font-bold">{game?.DownloadState?.Progress}</p>
                </div>

                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">Download Url</p>
                    <p class="text-sm font-bold text-frost-400 overflow-hidden">{game?.DownloadState?.DownloadUrl}</p>
                </div>
                <div>
                    <p class="text-[9px] font-bold text-muted uppercase">Path</p>
                    <p class="text-sm font-bold wrap-break-word">{game?.DownloadState?.DownloadPath}</p>
                </div>

                <div class="col-span-2 pt-4 border-t border-border">
                    <p class="text-[9px] font-bold text-muted uppercase mb-2 flex items-center gap-1">
                        <ShieldCheckIcon size={12}/>
                        Virus Scan
                    </p>
                    <p class="text-xs text-green-400 font-medium">TODO</p>
                </div>
            </div>
            <!-- File List -->
            <div class="lg:col-span-7 flex flex-col h-95 bg-panel/30 border border-border rounded-2xl overflow-hidden">
                {#if game}
                    <FileManager game={game}/>
                {/if}
            </div>
        </div>
    </section>

    <!-- Source Section -->
    <section class="space-y-4">
        <h2 class="text-sm font-bold uppercase tracking-widest text-muted px-2">Source</h2>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Current Source Info -->
            <div class="p-6 bg-surface border border-border rounded-3xl flex gap-6">
                <div class="w-24 h-32 bg-panel rounded-xl border border-border shrink-0 flex items-center justify-center text-muted">
                    <ImageIcon size={32}/>
                </div>
                <div class="flex-1 grid grid-cols-2 gap-4 text-xs">
                    <div class="col-span-2"><p class="text-[9px] text-muted uppercase font-bold">Title</p>
                        <p class="font-bold text-sm">Game.Name.2023-REPACK</p></div>
                    <div><p class="text-[9px] text-muted uppercase font-bold">Indexer</p>
                        <p>1337x</p></div>
                    <div><p class="text-[9px] text-muted uppercase font-bold">Size</p>
                        <p>45.2 GB</p></div>
                    <div><p class="text-[9px] text-muted uppercase font-bold">Created</p>
                        <p>2024-01-15</p></div>
                </div>
            </div>

            <IndexerSearch/>
        </div>
    </section>
</div>