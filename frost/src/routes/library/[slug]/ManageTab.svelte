<script lang="ts">
    import {FileTextIcon, ImageIcon, SearchIcon, ShieldCheckIcon} from '@lucide/svelte';
    import {fade} from "svelte/transition";
    import type {Game} from "$lib/gen/library/v1/library_pb";
    import IndexerSearch from "$lib/components/IndexerSearch.svelte";

    let {game = $bindable(null)}: { game: Game | null } = $props();

</script>

<div class="space-y-8" in:fade>
    <!-- Download Section -->
    <section class="space-y-4">
        <h2 class="text-sm font-bold uppercase tracking-widest text-muted px-2">Download</h2>
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 p-6 bg-surface border border-border rounded-3xl">
            <!-- State Info -->
            <div class="lg:col-span-5 grid grid-cols-2 gap-y-6 gap-x-4 border-r border-border pr-6">
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
            <div class="lg:col-span-7 bg-panel/30 border border-border rounded-2xl flex flex-col h-50">
                <div class="p-3 border-b border-border text-[10px] font-bold uppercase text-muted">Downloaded Files
                </div>
                <div class="flex-1 overflow-y-auto p-4 flex flex-col items-center justify-center opacity-20">
                    <FileTextIcon size={48} strokeWidth={1}/>
                    <p class="text-xs mt-2">File listing not available</p>
                </div>
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