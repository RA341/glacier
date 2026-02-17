<script lang="ts">
    import {
        GlobeIcon,
        LayoutGridIcon,
        LoaderIcon,
        RefreshCwIcon,
        StarIcon,
        TagIcon,
        TextAlignStart,
        TypeIcon
    } from "@lucide/svelte";
    import {fade} from 'svelte/transition';
    import AssetManager from "./AssetManager.svelte";
    import type {Game} from "$lib/gen/library/v1/library_pb";
    import {callRPC, glacierCli} from "$lib/api/api";
    import {SearchService} from "$lib/gen/search/v1/search_pb";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    let {game = $bindable(), refresh}: { game?: Game; refresh: () => Promise<void>; } = $props();

    function handleArrayInput(e: Event, field: 'Platforms' | 'Genres') {
        const value = (e.currentTarget as HTMLInputElement).value;
        game!.Meta![field] = value.split(',').map(s => s.trim()).filter(Boolean);
    }

    const sm = getSnackbarCtx()

    const searchSrv = glacierCli(SearchService)
    let isRefreshing = $state(false);

    async function refreshMeta() {
        isRefreshing = true;

        const {val, err} = await callRPC(() => searchSrv.getGameMeta({
            provider: game!.Meta!.ProviderType,
            gameDbId: game!.Meta!.ID
        }))
        if (err || !val?.meta) {
            sm.push(`Could not find game metadata ${err}`, 'error');
        } else {
            game!.Meta = val.meta
        }

        isRefreshing = false;
    }

</script>

<div class="space-y-6 pb-20">
    {#if game && game.Meta}
        <div class="bg-surface border border-border rounded-3xl p-8 space-y-6 shadow-sm">
            <div class="flex items-center mb-2">
                <h2 class="text-xs font-black text-frost-400 uppercase tracking-[0.2em] flex items-center gap-2">
                    <LayoutGridIcon size={14}/>
                    Main Content
                </h2>

                <button
                        onclick={refreshMeta}
                        disabled={isRefreshing}
                        class="flex items-center gap-2 px-4 py-1.5 bg-panel border border-border rounded-xl text-[10px] font-bold uppercase tracking-widest text-muted hover:text-frost-400 hover:border-frost-500/50 transition-all active:scale-95 disabled:opacity-50"
                >
                    {#if isRefreshing}
                        <LoaderIcon size={12} class="animate-spin text-frost-500"/>
                        Refreshing...
                    {:else}
                        <RefreshCwIcon size={12}/>
                        Refresh Metadata
                    {/if}
                </button>
            </div>

            <div class="space-y-2">
                <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Display Name</label>
                <div class="relative">
                    <TypeIcon size={18} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted/50"/>
                    <input type="text" bind:value={game.Meta.Name}
                           class="w-full bg-panel border border-border rounded-2xl py-3.5 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm font-bold"/>
                </div>
            </div>


            <div class="space-y-2">
                <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Summary</label>
                <div class="relative">
                    <TextAlignStart size={18} class="absolute left-4 top-4 text-muted/50"/>
                    <textarea bind:value={game.Meta.Summary} rows="2"
                              class="w-full bg-panel border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm"></textarea>
                </div>
            </div>

            <div class="space-y-2">
                <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Full Description</label>
                <div class="relative">
                    <TextAlignStart size={18} class="absolute left-4 top-4 text-muted/50"/>
                    <textarea bind:value={game.Meta.Description} rows="6"
                              class="w-full bg-panel border border-border rounded-2xl py-4 pl-12 pr-4 outline-none focus:border-frost-500 transition-all text-sm leading-relaxed"></textarea>
                </div>
            </div>
        </div>

        <!-- 2. TECHNICAL SPECIFICATIONS GRID -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">

            <!-- Provider Details -->
            <div class="bg-surface border border-border rounded-3xl p-6 space-y-4">
                <h3 class="text-[10px] font-bold text-muted uppercase tracking-widest flex items-center gap-2">
                    <GlobeIcon size={12}/>
                    Source Info
                </h3>

                <div class="space-y-4">
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">Provider</span>
                        <input type="text" bind:value={game.Meta.ProviderType}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                    </div>
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">External ID</span>
                        <input type="text" bind:value={game.Meta.ID}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs font-mono outline-none focus:border-frost-500"/>
                    </div>
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">Original URL</span>
                        <input type="text" bind:value={game.Meta.URL}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                    </div>
                </div>
            </div>

            <!-- Classifications -->
            <div class="bg-surface border border-border rounded-3xl p-6 space-y-4">
                <h3 class="text-[10px] font-bold text-muted uppercase tracking-widest flex items-center gap-2">
                    <TagIcon size={12}/>
                    Classifications
                </h3>

                <div class="space-y-4">
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">Category</span>
                        <input type="text" bind:value={game.Meta.Category}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                    </div>
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">Platforms (Comma separated)</span>
                        <input type="text" value={game.Meta.Platforms?.join(', ')}
                               oninput={(e) => handleArrayInput(e, 'Platforms')}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                    </div>
                    <div class="space-y-1">
                        <span class="text-[9px] text-muted uppercase font-medium">Genres (Comma separated)</span>
                        <input type="text" value={game.Meta.Genres?.join(', ')}
                               oninput={(e) => handleArrayInput(e, 'Genres')}
                               class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                    </div>
                </div>
            </div>

            <!-- Release & Ratings -->
            <div class="bg-surface border border-border rounded-3xl p-6 space-y-4">
                <h3 class="text-[10px] font-bold text-muted uppercase tracking-widest flex items-center gap-2">
                    <StarIcon size={12}/>
                    Release & Rating
                </h3>

                <div class="space-y-4">
                    <div class="grid grid-cols-2 gap-2">
                        <div class="space-y-1">
                            <span class="text-[9px] text-muted uppercase font-medium">Date</span>
                            <input type="text" bind:value={game.Meta.ReleaseDate}
                                   class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        </div>
                        <div class="space-y-1">
                            <span class="text-[9px] text-muted uppercase font-medium">Status</span>
                            <input type="text" bind:value={game.Meta.ReleaseStatus}
                                   class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        </div>
                    </div>
                    <div class="grid grid-cols-2 gap-2">
                        <div class="space-y-1">
                            <span class="text-[9px] text-muted uppercase font-medium">Score</span>
                            <input type="text" bind:value={game.Meta.Rating}
                                   class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500 text-frost-400 font-bold"/>
                        </div>
                        <div class="space-y-1">
                            <span class="text-[9px] text-muted uppercase font-medium">Votes</span>
                            <input type="number" bind:value={game.Meta.RatingCount}
                                   class="w-full bg-panel border border-border rounded-xl px-3 py-2 text-xs outline-none focus:border-frost-500"/>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- 3. ASSETS SECTION -->
        {#if game.ID}
            <AssetManager refresh={refresh} gameId={game.ID} assets={game?.Meta?.assets}/>
        {/if}
    {/if}
</div>

<style>
    textarea {
        resize: vertical;
        min-height: 80px;
    }
</style>
