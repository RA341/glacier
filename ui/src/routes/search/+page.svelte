<script lang="ts">
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {type GameMetadata, SearchService} from "$lib/gen/search/v1/search_pb";
    import {glacierCli} from "$lib/api/api";
    import {CircleAlert, ImageIcon, LoaderIcon, PlusIcon, SearchIcon, SearchXIcon} from '@lucide/svelte';
    import {page} from '$app/state';
    import {goto} from "$app/navigation";
    import {onMount} from "svelte";
    import {transferStore} from "./selectedGame.svelte";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";
    import GameGridItem from "./GameGridItem.svelte";

    const searchQueryParam = 'q';
    let query = $state(page.url.searchParams.get(searchQueryParam) || "");

    let srvConfig = glacierCli(ServiceConfigService)
    let metadataRpc = createRPCRunner(() => srvConfig.getActiveService({serviceType: "Metadata"}))

    async function loadProviders() {
        await metadataRpc.runner()
    }

    onMount(() => {
        loadProviders();
        if (query && selectedMetadataProvider) handleSearch();
    });

    const metadataQueryParam = 'metadata';
    let selectedMetadataProvider = $state(page.url.searchParams.get(metadataQueryParam) || "");

    function handleSearch() {
        const url = new URL(page.url);

        if (query) url.searchParams.set(searchQueryParam, query);
        else url.searchParams.delete(searchQueryParam);

        if (selectedMetadataProvider) url.searchParams.set(metadataQueryParam, selectedMetadataProvider);
        else url.searchParams.delete(metadataQueryParam);

        goto(url, {keepFocus: true});
        searchRpc.runner();
    }

    $effect(() => {
        if (metadataRpc.value?.names && !selectedMetadataProvider) {
            selectedMetadataProvider = metadataRpc.value.names[0].Name;
            if (query) handleSearch();
        }
    });

    const searchClient = glacierCli(SearchService);
    const searchRpc = createRPCRunner(() => searchClient.searchMetadata({
        q: {
            query: query,
            indexer: selectedMetadataProvider
        }
    }));

    let isModalOpen = $state(false);

    function handleAdd(meta: GameMetadata) {
        transferStore.data = meta

        goto("/search/details", {keepFocus: true})
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter' && query.trim() !== '') handleSearch();
    }
</script>

<div class="flex flex-col h-full w-full bg-background text-foreground">
    <header class="p-6 border-b border-border bg-surface/50 backdrop-blur-sm">
        <div class="max-w-7xl mx-auto flex items-center gap-6">
            <h1 class="text-2xl font-bold tracking-tight">Search</h1>

            <div class="flex flex-1 items-center gap-3">
                <div class="relative min-w-40">
                    <select
                            bind:value={selectedMetadataProvider}
                            onchange={handleSearch}
                            class="w-full appearance-none bg-panel border border-border rounded-xl px-4 py-2.5 pr-10 text-sm font-bold outline-none focus:border-frost-500 transition-all cursor-pointer"
                    >
                        {#if metadataRpc.loading}
                            <option disabled>Loading indexers...</option>
                        {:else if metadataRpc.value?.names}
                            {#each metadataRpc.value.names as indexer}
                                <option value={indexer.Name}>{indexer.Name}</option>
                            {/each}
                        {/if}
                    </select>
                </div>

                <!-- Query Input + Search Button -->
                <div class="relative flex-1 group">
                    <input
                            type="text"
                            bind:value={query}
                            onkeydown={handleKeyDown}
                            placeholder="Search for games..."
                            class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 pr-12 outline-none focus:border-frost-500 transition-all text-sm"
                    />
                    <button
                            onclick={handleSearch}
                            disabled={searchRpc.loading || !query.trim()}
                            class="absolute right-1.5 top-1.5 bottom-1.5 w-10 flex items-center justify-center rounded-lg bg-surface border border-border text-muted hover:text-frost-500 hover:border-frost-500/50 transition-all disabled:opacity-50"
                    >
                        {#if searchRpc.loading}
                            <LoaderIcon size={18} class="animate-spin text-frost-500"/>
                        {:else}
                            <SearchIcon size={18}/>
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </header>

    <!-- Results Area -->
    <main class="flex-1 overflow-y-auto p-6">
        <div class="max-w-7xl mx-auto">
            {#if searchRpc.error}
                <aside class="flex items-start gap-4 p-4 rounded-xl bg-red-500/5 border border-red-500/20 text-red-400">
                    <CircleAlert size={20}/>
                    <div>
                        <h3 class="font-bold">Search Failed</h3>
                        <p class="text-sm opacity-80">{searchRpc.error}</p>
                    </div>
                </aside>

            {:else if searchRpc.loading}
                <div class="flex flex-col items-center justify-center h-96 text-muted gap-4">
                    <LoaderIcon class="animate-spin size-10 text-frost-500"/>
                    <p class="animate-pulse font-medium">Indexing results...</p>
                </div>

            {:else if searchRpc.value && searchRpc.value.metadata && searchRpc.value.metadata.length > 0}
                <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                    {#each searchRpc.value.metadata as item (item.ID)}
                        <GameGridItem game={item}/>
                    {/each}
                </div>
            {:else if searchRpc.value}
                <div class="flex flex-col items-center justify-center h-96 text-muted gap-2 opacity-40">
                    <SearchXIcon size={64} strokeWidth={1}/>
                    <p class="text-xl font-medium">No results found</p>
                    <p class="text-sm">Try adjusting your query or indexer</p>
                </div>
            {:else}
                <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-border rounded-3xl text-blue-100">
                    <SearchIcon size={80} strokeWidth={1} class="mb-4"/>
                    <p class="text-lg font-medium">Search for game</p>
                </div>
            {/if}
        </div>
    </main>
</div>

<style>
    main::-webkit-scrollbar {
        width: 8px;
    }

    main::-webkit-scrollbar-track {
        background: transparent;
    }

    main::-webkit-scrollbar-thumb {
        background: var(--border);
        border-radius: 10px;
    }
</style>
