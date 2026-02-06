<script lang="ts">
    import {ChevronDownIcon, CircleAlert, HardDriveIcon, ImageIcon, LoaderIcon, SearchIcon} from '@lucide/svelte';
    import {glacierCli} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {type GameMetadata, type GameSource, SearchService} from "$lib/gen/search/v1/search_pb";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";

    const searchService = glacierCli(SearchService)
    const srvConfig = glacierCli(ServiceConfigService)

    let {selectedGameSource = $bindable(), game}: {
        game: GameMetadata | null;
        selectedGameSource: GameSource | null
    } = $props()

    function searchIndexer() {
        if (searchQuery.trim()) {
            searchIndexerRpc.runner()
        }
    }

    function handleMatch(item: GameSource) {
        if (selectedGameSource?.Title === item.Title) {
            selectedGameSource = null
        } else {
            selectedGameSource = item
        }
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') searchIndexer();
    }


    let searchQuery = $state("");
    $effect(() => {
        if (game) {
            searchQuery = game.Name
        }
    })

    let searchIndexerRpc = createRPCRunner(() => searchService.searchIndexers(
        {
            q: {
                indexer: selectedIndexer,
                query: searchQuery,
            }
        }
    ))
    let availableIndexersRpc = createRPCRunner(() => srvConfig.getActiveService({serviceType: "Indexer"}))
    let selectedIndexer = $state("");

    $effect(() => {
        availableIndexersRpc.runner()
    })
    $effect(() => {
        if (selectedIndexer === "") {
            selectedIndexer = availableIndexersRpc.value?.names.at(0)?.Name ?? ""
        }
    })

</script>

<div class="bg-surface border border-border rounded-3xl p-6 flex flex-col gap-6">
    <h2 class="text-lg font-bold flex items-center gap-2">
        <SearchIcon size={20} class="text-frost-400"/>
        Search Indexers
    </h2>

    <div class="flex gap-3">
        <div class="relative min-w-35">
            <select
                    bind:value={selectedIndexer}
                    class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 appearance-none outline-none text-sm font-bold"
            >
                {#each availableIndexersRpc?.value?.names as ind }
                    <option>{ind.Name}</option>
                {/each}
            </select>
            <ChevronDownIcon
                    size={16}
                    class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted"
            />
        </div>
        <div class="relative flex-1">
            <input type="text" placeholder="Search for release..."
                   onkeydown={handleKeyDown}
                   bind:value={searchQuery}
                   class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 outline-none focus:border-frost-500 text-sm"/>
        </div>
        <button
                onclick={searchIndexer}
                disabled={searchIndexerRpc.loading || !searchQuery.trim()}
                class=" right-1.5 top-1.5 bottom-1.5 w-10 flex items-center justify-center rounded-lg bg-surface border border-border text-muted hover:text-frost-500 hover:border-frost-500/50 transition-all disabled:opacity-50"
        >
            {#if searchIndexerRpc.loading}
                <LoaderIcon size={18} class="animate-spin text-frost-500"/>
            {:else}
                <SearchIcon size={18}/>
            {/if}
        </button>
    </div>

    <div class="space-y-3 h-100 overflow-y-auto pr-2">

        {#if searchIndexerRpc.error}
            <aside class="flex items-start gap-4 p-4 rounded-xl bg-red-500/5 border border-red-500/20 text-red-400">
                <CircleAlert size={20}/>
                <div>
                    <h3 class="font-bold">Search Failed</h3>
                    <p class="text-sm opacity-80">{searchIndexerRpc.error}</p>
                </div>
            </aside>

        {:else if searchIndexerRpc.loading}
            <div class="flex flex-col items-center justify-center h-96 text-muted gap-4">
                <LoaderIcon class="animate-spin size-10 text-frost-500"/>
                <p class="animate-pulse font-medium">Indexing results...</p>
            </div>

        {:else if searchIndexerRpc.value?.results}
            {#each searchIndexerRpc?.value?.results as item}
                <div class="flex items-center justify-between p-4 bg-panel/30 border border-border rounded-2xl hover:border-muted transition-colors">
                    <div class="flex items-center gap-4">
                        <div class="w-12 h-16 bg-panel rounded-lg flex items-center justify-center border border-border text-muted">
                            <ImageIcon size={20}/>
                        </div>
                        <div class="flex flex-col">
                            <span class="font-bold text-sm">{item.Title}</span>
                            <div class="flex items-center gap-3 mt-1 text-[10px] text-muted font-bold">
                                    <span class="flex items-center gap-1 uppercase">
                                        <HardDriveIcon size={10}/>
                                        {item.FileSize}
                                    </span>
                                <span>•</span>
                                <span>{new Date(item.CreatedISO).getFullYear()}</span>
                            </div>
                        </div>
                    </div>
                    <button
                            onclick={() => handleMatch(item)}
                            class="px-4 py-1.5 bg-panel border border-border rounded-lg text-xs font-bold hover:border-frost-500 transition-colors">
                        {selectedGameSource?.Title === item.Title ? "Selected" : "Select"}
                    </button>
                </div>
            {/each}
        {/if}
    </div>
</div>
