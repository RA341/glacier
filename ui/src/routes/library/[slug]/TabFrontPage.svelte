<script lang="ts">
    import type {Game} from "$lib/gen/library/v1/library_pb";
    import {ExternalLinkIcon} from "@lucide/svelte";
    import {trimPrefix} from "$lib/api/strings";
    import {formatReleaseDate} from "$lib/api/date";

    let {game = $bindable(null)}: { game: Game | null } = $props();
    let meta = $derived(game?.Meta);

    const dateInf = $derived(formatReleaseDate(meta?.ReleaseDate ?? ""));
</script>

<div class="grid grid-cols-1 lg:grid-cols-4 gap-8 items-start">
    <div class="lg:col-span-3 space-y-6">
        <div class="p-8 bg-surface border border-border rounded-3xl space-y-3 shadow-sm">
            <h3 class="text-[10px] font-black text-muted uppercase tracking-[0.3em]">
                Summary
            </h3>
            <p class="text-lg font-medium leading-relaxed text-foreground/90">
                {meta?.Summary || "No summary available."}
            </p>
        </div>

        <div class="p-8 bg-surface border border-border rounded-3xl space-y-4 shadow-sm min-h-75">
            <h3 class="text-[10px] font-black text-muted uppercase tracking-[0.3em]">About this game</h3>
            <p class="text-md text-white leading-relaxed whitespace-pre-line">
                {meta?.Description || "No detailed description found."}
            </p>
        </div>
    </div>

    <aside class="space-y-6">
        <div class="bg-surface border border-border rounded-3xl p-6 space-y-6 shadow-sm">
            <div class="space-y-1">
                <span class="text-[10px] font-black text-muted uppercase tracking-widest flex items-center gap-2">
                    Release Date
                </span>
                <p class="text-sm font-bold text-foreground">
                    {dateInf.formatted || 'Unknown'}
                    {#if dateInf.formatted}
                        <span class="text-muted-foreground font-normal ml-1">
                            ({dateInf.relative})
                        </span>
                    {/if}
                </p>
            </div>

            <div class="space-y-1">
                <span class="text-[10px] font-black text-muted uppercase tracking-widest flex items-center gap-2">
                    Category
                </span>
                <p class="text-sm font-bold text-foreground capitalize">{meta?.Category || 'Unknown'}</p>
            </div>

            <div class="space-y-1">
                <span class="text-[10px] font-black text-muted uppercase tracking-widest flex items-center gap-2">
                    Platforms
                </span>
                <p class="text-sm font-bold text-frost-400">{meta?.Platforms.join(", ")}</p>
            </div>

            <div class="h-px bg-border"></div>

            <div class="space-y-1">
                <span class="text-[10px] font-black text-muted uppercase tracking-widest">Metadata Provider</span>
                <div class="flex items-center justify-between bg-panel border border-border rounded-xl p-3 mt-2">
                    <span class="text-xs font-bold">{trimPrefix(meta?.ProviderType || "No Provider Found", "Provider")}</span>
                    <a href={meta?.URL} target="_blank" class="text-frost-400 hover:text-frost-300">
                        <ExternalLinkIcon size={16}/>
                    </a>
                </div>
            </div>

            <!-- Sidebar Genre Tags -->
            <div class="space-y-2 pt-2">
                <span class="text-[10px] font-black text-muted uppercase tracking-widest">Genres</span>
                <div class="flex flex-wrap gap-2 pt-1">
                    {#each meta?.Genres ?? [] as genre}
                        <span class="px-2 py-1 bg-panel border border-border text-[10px] font-bold rounded-lg whitespace-nowrap">
                            {genre}
                        </span>
                    {/each}
                </div>
            </div>
        </div>
    </aside>
</div>
