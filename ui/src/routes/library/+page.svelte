<script lang="ts">
    import {glacierCli} from "$lib/api/api";
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/rpc.svelte.js";
    import {onMount} from "svelte";
    import {goto} from "$app/navigation";
    import {getAssetPath} from "$lib/api/assets";
    import {ImageIcon, LoaderIcon} from "@lucide/svelte";
    import AssetImg from "$lib/components/assets/AssetImg.svelte";

    const libSrv = glacierCli(LibraryService)

    let query = $state("")
    let offset = $state(0)
    let limit = $state(50)

    let libRpc = createRPCRunner(() => libSrv.list({
        limit: limit,
        offset: offset,
        query: query,
    }))

    function search() {
        libRpc.runner()
    }

    function select(game: Game) {
        goto(`/library/${game.ID}`)
    }

    $effect(() => {
        const interval = setInterval(() => {
            search()
        }, 5000);
        return () => clearInterval(interval);
    });

    onMount(() => {
        search()
    })
</script>

<main class="min-h-screen p-6 md:p-10 max-w-7xl mx-auto">
    <!-- Header Area -->
    <header class="flex flex-col md:flex-row md:items-center justify-between gap-6 mb-12">
        <h1 class="text-3xl font-bold tracking-tight">Library</h1>

        <div class="relative flex items-center w-full max-w-2xl group">
            <input
                    type="text"
                    bind:value={query}
                    onkeydown={(e) => e.key === 'Enter' && search()}
                    placeholder="Search games..."
                    class="w-full bg-surface border border-border text-foreground px-4 py-3 rounded-xl focus:outline-none focus:ring-2 focus:ring-frost-500/50 transition-all placeholder:text-muted"
            />
            <button
                    onclick={search}
                    class="absolute right-2 p-2 text-muted hover:text-frost-400 transition-colors"
                    aria-label="Search"
            >
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                     stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="11" cy="11" r="8"/>
                    <path d="m21 21-4.3-4.3"/>
                </svg>
            </button>
        </div>
    </header>

    <!-- Content Area -->
    {#if libRpc.error}
        <div class="flex items-center justify-center h-64 border border-red-900/50 bg-red-950/20 rounded-xl text-red-400">
            <p>Error: {libRpc.error}</p>
        </div>
    {:else if libRpc.loading && !libRpc.value && !libRpc.error}
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
            {#each Array(10) as _}
                <div class="animate-pulse bg-surface border border-border rounded-2xl overflow-hidden h-80">
                    <div class="h-3/4 bg-panel"></div>
                    <div class="p-4 space-y-2">
                        <div class="h-4 bg-panel w-3/4 rounded"></div>
                        <div class="h-3 bg-panel w-1/2 rounded"></div>
                    </div>
                </div>
            {/each}
        </div>
    {:else}
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
            {#each libRpc.value?.gameList ?? [] as ga}
                <div
                        onclick={() => select(ga)}
                        onkeydown={(e) => e.key === 'Enter' && select(ga)}
                        role="button"
                        tabindex="0"
                        class="group cursor-pointer bg-surface border border-border rounded-2xl overflow-hidden transition-all duration-300 hover:border-frost-500/50 hover:shadow-frost focus:outline-none focus:ring-2 focus:ring-frost-500"
                >
                    <div class="aspect-3/4 bg-panel flex items-center justify-center relative overflow-hidden">
                        <AssetImg
                                src={getAssetPath({gameId: ga.ID, assetType: "AssetThumbnail"})}
                                alt={ga.Meta?.Name ?? "Untitled"}
                                class="object-cover w-full h-full"
                        >
                            {#snippet loadingSlot()}
                                <div class="flex h-full w-full items-center justify-center bg-surface animate-pulse">
                                    <LoaderIcon class="animate-spin text-muted"/>
                                </div>
                            {/snippet}

                            {#snippet errorSlot()}
                                <div class="flex h-full w-full flex-col items-center justify-center bg-surface">
                                    <ImageIcon size={48} class="mb-5"/>
                                    <p class="text-sm">Thumbnail Not found</p>
                                </div>
                            {/snippet}
                        </AssetImg>
                    </div>

                    <!-- Info Area -->
                    <div class="p-4 border-t border-border">
                        <h3 class="font-bold truncate text-foreground mb-1">
                            {ga.Meta?.Name}
                            <span class="text-muted font-normal text-sm ml-1">
                                ({new Date(ga?.Meta?.ReleaseDate || 'N/A').getFullYear() || 'N/A'})
                            </span>
                        </h3>
                        <div class="flex items-center gap-2">
                            <div class="w-2 h-2 rounded-full bg-frost-400 animate-pulse"></div>
                            <span class="text-xs text-muted font-medium uppercase tracking-wider">
                                {ga?.DownloadState?.State || "Unknown"}
                            </span>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</main>

<style>
    /* Add smooth scaling effect for the cards */
    .group:hover {
        transform: translateY(-4px);
    }
</style>