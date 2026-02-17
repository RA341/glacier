<script lang="ts">
    import {ImageIcon} from '@lucide/svelte';
    import IndexerSearch from "$lib/components/IndexerSearch.svelte";

    let {game = $bindable()} = $props();
</script>

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
