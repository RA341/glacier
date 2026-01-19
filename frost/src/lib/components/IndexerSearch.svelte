<script lang="ts">
    import {SearchIcon} from '@lucide/svelte';
    import type {GameSource} from "$lib/gen/search/v1/search_pb";

    // let { source = $bindable(null) }: { source: GameSource | null } = $props();

</script>


<div class="p-6 bg-surface border border-border rounded-3xl space-y-4">
    <div class="flex gap-2">
        <div class="bg-panel border border-border px-3 py-2 rounded-xl text-xs font-bold uppercase">
            Indexer
        </div>
        <div class="relative flex-1">
            <SearchIcon size={14} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted"/>
            <input type="text" placeholder="Search for alternative source..."
                   class="w-full bg-panel border border-border rounded-xl py-2 pl-9 text-xs outline-none focus:border-frost-500"/>
        </div>
    </div>
    <div class="space-y-2">
        {#each [1, 2] as i}
            <div class="flex items-center justify-between p-3 bg-panel/50 border border-border rounded-xl text-xs">
                <span>Other.Release.v1.0-Group</span>
                <button class="px-3 py-1 bg-panel border border-border rounded-md hover:border-frost-500 font-bold">
                    Update
                </button>
            </div>
        {/each}
    </div>
</div>
