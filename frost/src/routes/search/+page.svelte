<script lang="ts">
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {type GameSearchResult, SearchService} from "$lib/gen/search/v1/search_pb";
    import {cli, formatDate} from "$lib/api/api";
    import {
        SearchIcon, LoaderIcon, CircleAlert, SearchXIcon,
        DownloadIcon, FileIcon, CalendarIcon, HardDriveIcon, XIcon
    } from '@lucide/svelte';

    import MatchMetadataModal from "./MatchMetadataModal.svelte";
    import { page } from '$app/state';
    import {goto} from "$app/navigation";
    import {onMount} from "svelte";

    let query = $state(page.url.searchParams.get('q') || "");

    function handleSearch() {
        const url = new URL(page.url);
        if (query) {
            url.searchParams.set('q', query);
        } else {
            url.searchParams.delete('q');
        }
        goto(url, { keepFocus: true });
        rpc.runner();
    }

    onMount(() => {
        if (query) {
            handleSearch()
        }
    })

    const searchClient = cli(SearchService);
    const rpc = createRPCRunner(() => searchClient.search({query}));

    let isModalOpen = $state(false);
    let selectedGame = $state<GameSearchResult | null>(null);

    function handleAdd(gameResult: GameSearchResult) {
        selectedGame = gameResult;
        isModalOpen = true;
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter' && query.trim() !== '') {
            handleSearch()
        }
    }
</script>

<div class="flex flex-col h-full max-w-5xl mx-auto w-full bg-background text-foreground">
    <!-- Header -->
    <header class="p-8 space-y-6">
        <div class="flex flex-col gap-1">
            <h1 class="text-3xl font-bold tracking-tight">Search</h1>
            <p class="text-muted text-sm">Find and download assets from the repository.</p>
        </div>

        <div class="flex gap-3">
            <div class="relative flex-1 group">
                <div class="absolute inset-y-0 left-3 flex items-center pointer-events-none text-muted group-focus-within:text-frost-400 transition-colors">
                    <SearchIcon size={18}/>
                </div>
                <input
                        type="search"
                        bind:value={query}
                        onkeydown={handleKeyDown}
                        placeholder="Search by filename..."
                        class="w-full bg-surface border border-border focus:border-frost-500 rounded-xl p-2.5 pl-10 outline-none transition-all placeholder:text-muted/50"
                />
            </div>

            <button
                    onclick={rpc.runner}
                    disabled={rpc.loading || !query.trim()}
                    class="bg-frost-500 hover:bg-frost-400 disabled:opacity-50 disabled:hover:bg-frost-500 text-background font-semibold px-6 rounded-xl flex items-center gap-2 transition-colors active:scale-95"
            >
                {#if rpc.loading}
                    <LoaderIcon class="animate-spin" size={18}/>
                {:else}
                    Search
                {/if}
            </button>
        </div>
    </header>

    <div class="border-b border-border mx-8"></div>

    <!-- Results Area -->
    <main class="flex-1 overflow-y-auto p-8">
        {#if rpc.error}
            <aside class="flex items-start gap-4 p-4 rounded-xl bg-red-500/5 border border-red-500/20 text-red-400">
                <CircleAlert size={20}/>
                <div>
                    <h3 class="font-bold">Search Failed</h3>
                    <p class="text-sm opacity-80">{rpc.error}</p>
                </div>
            </aside>

        {:else if rpc.loading}
            <div class="flex flex-col items-center justify-center h-64 text-muted gap-4">
                <LoaderIcon class="animate-spin size-10 text-frost-500"/>
                <p class="animate-pulse">Searching database...</p>
            </div>

        {:else if rpc.value}
            {#if rpc.value.results && rpc.value.results.length > 0}
                <div class="grid grid-cols-1 gap-3">
                    {#each rpc.value.results as item}
                        <div class="flex items-center justify-between p-4 rounded-xl border border-border bg-surface hover:border-frost-500/50 transition-all group">
                            <div class="flex items-center gap-4">
                                <div class="bg-panel p-3 rounded-lg text-muted group-hover:text-frost-400 transition-colors">
                                    <FileIcon size={24}/>
                                </div>
                                <div class="flex flex-col">
                                    <span class="font-bold text-lg leading-tight group-hover:text-frost-200 transition-colors">{item.name}</span>
                                    <div class="flex items-center gap-4 mt-1 text-sm text-muted">
                                        <span class="flex items-center gap-1">
                                            <HardDriveIcon size={14}/> {item.size}
                                        </span>
                                        <span class="flex items-center gap-1">
                                            <CalendarIcon size={14}/> {formatDate(item.uploadDate)}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            <button
                                    type="button"
                                    onclick={() => handleAdd(item)}
                                    class="bg-panel hover:bg-frost-500 hover:text-background text-foreground px-4 py-2 rounded-lg flex items-center gap-2 text-sm font-medium transition-all"
                            >
                                <DownloadIcon size={16}/>
                                <span>Add</span>
                            </button>
                        </div>
                    {/each}
                </div>
            {:else}
                <div class="flex flex-col items-center justify-center h-64 text-muted gap-2">
                    <SearchXIcon size={48} strokeWidth={1} class="opacity-20"/>
                    <p class="text-xl font-medium">No files found</p>
                </div>
            {/if}

        {:else}
            <!-- Initial State -->
            <div class="flex flex-col items-center justify-center h-64 border border-dashed border-border rounded-2xl text-muted/30">
                <SearchIcon size={48} class="mb-2"/>
                <p>Enter a query to search the file directory</p>
            </div>
        {/if}
    </main>
</div>

<MatchMetadataModal
        bind:isOpen={isModalOpen}
        localGame={selectedGame}
/>
