<script lang="ts">
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {type GameMetadata, SearchService} from "$lib/gen/search/v1/search_pb";
    import {glacierCli} from "$lib/api/api";
    import {CircleAlert, HardDriveIcon, ImageIcon, LoaderIcon, PlusIcon, SearchIcon, SearchXIcon} from '@lucide/svelte';

    import MatchMetadataModal from "./MatchMetadataModal.svelte";
    import {page} from '$app/state';
    import {goto} from "$app/navigation";
    import {onMount} from "svelte";
    import {MetadataService} from "$lib/gen/metadata/v1/metadata_pb";
    import {transferStore} from "./selectedGame.svelte";

    const searchQueryParam = 'q';
    let query = $state(page.url.searchParams.get(searchQueryParam) || "");

    let metaClients = glacierCli(MetadataService)
    let metadataRpc = createRPCRunner(() => metaClients.getMetadataProviders({}))

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
        if (metadataRpc.value?.providers && !selectedMetadataProvider) {
            selectedMetadataProvider = metadataRpc.value.providers[0].Name;
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

        goto("/search/details", { keepFocus: true })
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
                        {:else if metadataRpc.value?.providers}
                            {#each metadataRpc.value.providers as indexer}
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
                <!-- Card Grid -->
                <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6">
                    {#each searchRpc.value.metadata as item}
                        <button
                                onclick={() => handleAdd(item)}
                                class="group relative flex flex-col bg-surface border border-border rounded-2xl overflow-hidden hover:border-frost-500/50 transition-all hover:-translate-y-1 shadow-sm active:scale-[0.98]"
                        >
                            <!-- Game Photo / Poster -->
                            <div class="aspect-3/4 w-full bg-panel relative overflow-hidden">
                                <!-- Placeholder Logic -->
                                {#if item.ThumbnailURL}
                                    <img
                                            src={item.ThumbnailURL}
                                            alt={item.Name}
                                            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                                    />
                                {:else}
                                    <div class="w-full h-full flex flex-col items-center justify-center text-muted/20 gap-2">
                                        <ImageIcon size={48} strokeWidth={1}/>
                                        <span class="text-[10px] font-bold uppercase tracking-widest">No Preview</span>
                                    </div>
                                {/if}

                                <!-- Hover Overlay -->
                                <div class="absolute inset-0 bg-frost-900/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                                    <div class="bg-frost-500 text-background p-3 rounded-full shadow-xl translate-y-4 group-hover:translate-y-0 transition-transform">
                                        <PlusIcon size={24}/>
                                    </div>
                                </div>
                            </div>

                            <!-- Info Bar -->
                            <div class="p-4 border-t border-border flex flex-col gap-1 bg-surface">
                                <h3 class="font-bold text-sm text-foreground truncate">{item.Name}</h3>
                                <div class="flex items-center justify-between text-[11px] text-muted font-medium">
                                    <span class="flex items-center gap-1">
                                        (2023) <!-- Placeholder Year -->
                                    </span>
                                    <!--                                    <span class="flex items-center gap-1 uppercase tracking-tighter opacity-70">-->
                                    <!--                                        <HardDriveIcon size={12}/> {item.FileSize}-->
                                    <!--                                    </span>-->
                                </div>
                            </div>
                        </button>
                    {/each}
                </div>
            {:else if searchRpc.value}
                <div class="flex flex-col items-center justify-center h-96 text-muted gap-2 opacity-40">
                    <SearchXIcon size={64} strokeWidth={1}/>
                    <p class="text-xl font-medium">No results found</p>
                    <p class="text-sm">Try adjusting your query or indexer</p>
                </div>
            {:else}
                <!-- Initial State -->
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
