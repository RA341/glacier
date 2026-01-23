<script lang="ts">
    import {CalendarIcon, ImageIcon, PlusIcon, StarIcon, SwatchBookIcon} from '@lucide/svelte';
    import type {Game} from "$lib/gen/library/v1/library_pb";

    let {game = $bindable(null)}: { game: Game | null } = $props();

    let meta = $derived(game?.Meta)
</script>

<div class="space-y-6">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 h-75">
        <!-- Main Poster -->
        <div class="lg:col-span-2 bg-panel rounded-3xl border border-border overflow-hidden relative">
            {#if meta?.ThumbnailURL}
                <img src={meta?.ThumbnailURL} alt="" class="w-full h-full object-cover"/>
            {:else}
                <div class="w-full h-full flex items-center justify-center text-muted/20">
                    <ImageIcon size={48}/>
                </div>
            {/if}
        </div>

        <!-- Video Trailer -->
        <div class="lg:col-span-8 bg-panel rounded-3xl border border-border flex items-center justify-center overflow-hidden relative group">
            <div class="flex flex-col items-center gap-2">
                {#if meta?.Videos}
                    <iframe
                            width="560"
                            height="315"
                            src="https://www.youtube.com/embed/{meta.Videos[0]}"
                            title="YouTube video player"
                            allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                            allowfullscreen>
                    </iframe>
                {:else}
                    <PlusIcon size={48} class="rotate-45"/>
                    <p class="text-xs font-bold uppercase tracking-widest">No video found</p>
                {/if}
            </div>
        </div>
        <!-- Quick Meta Tags -->
        <div class="lg:col-span-2 flex flex-col gap-3">
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Release Date</span>
                <span class="text-sm font-bold flex items-center gap-2">
                    <CalendarIcon size={14} class="text-frost-400"/>
                    {new Date(meta?.ReleaseDate || "N/A").getFullYear()}
                </span>
            </div>
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Genre</span>
                <span class="text-sm font-bold flex items-center gap-2"><SwatchBookIcon size={14}
                                                                                        class="text-frost-400"/> {meta?.Genres.join(", ")}</span>
            </div>
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Rating</span>
                <span class="text-sm font-bold flex items-center gap-2">
                    <StarIcon size={14} class="text-yellow-500 fill-yellow-500"/>
                    {meta?.Rating} ({meta?.RatingCount ?? 0})
                </span>
            </div>
        </div>
    </div>
</div>