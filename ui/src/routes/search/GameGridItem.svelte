<script lang="ts">
    import {ImageIcon, PlusIcon} from '@lucide/svelte';
    import {goto} from "$app/navigation";
    import {transferStore} from "./selectedGame.svelte.ts";
    import {type GameMetadata, SearchService} from "$lib/gen/search/v1/search_pb";


    let {game}: { game: GameMetadata } = $props()

    function handleAdd(meta: GameMetadata) {
        transferStore.data = meta

        goto("/search/details", {keepFocus: true})
    }

</script>

<button
        onclick={() => handleAdd(game)}
        class="group relative flex flex-col bg-surface border border-border rounded-2xl overflow-hidden hover:border-frost-500/50 transition-all hover:-translate-y-1 shadow-sm active:scale-[0.98]"
>
    <div class="aspect-3/4 w-full bg-panel relative overflow-hidden">
        {#if game.ThumbnailURL}
            <img
                    src={game.ThumbnailURL}
                    alt={game.Name}
                    class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
            />
        {:else}
            <div class="w-full h-full flex flex-col items-center justify-center text-muted/20 gap-2">
                <ImageIcon size={48} strokeWidth={1}/>
                <span class="text-[10px] font-bold uppercase tracking-widest">No Preview</span>
            </div>
        {/if}

        <div class="absolute inset-0 bg-frost-900/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
            <div class="bg-frost-500 text-background p-3 rounded-full shadow-xl translate-y-4 group-hover:translate-y-0 transition-transform">
                <PlusIcon size={24}/>
            </div>
        </div>
    </div>

    <div class="p-4 border-t border-border flex flex-col gap-1 bg-surface">
        <h3 class="font-bold text-sm text-foreground truncate">{game.Name}</h3>
        <div class="flex items-center justify-between text-[11px] text-muted font-medium">
            <span class="flex items-center gap-1">
                {new Date(game.ReleaseDate).getFullYear()}
            </span>
        </div>
    </div>
</button>
