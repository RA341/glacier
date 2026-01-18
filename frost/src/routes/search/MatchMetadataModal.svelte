<script lang="ts">
    import {
        ChevronDownIcon,
        FolderOpenIcon,
        ImageIcon,
        LinkIcon,
        LoaderIcon,
        SearchIcon,
        SearchXIcon,
        TriangleAlert,
        XIcon,
    } from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';
    import {callRPC, cli} from "$lib/api/api";
    import {type GameMetadata, type GameSearchResult, SearchService} from "$lib/gen/search/v1/search_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {GameSchema, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {create} from "@bufbuild/protobuf";
    import {DownloaderService} from "$lib/gen/downloader/v1/downloader_pb";

    type Props = {
        isOpen?: boolean;
        localGame: GameSearchResult | null;
    };

    let {
        isOpen = $bindable(false),
        localGame,
    }: Props = $props();

    const searchClient = cli(SearchService);

    let searchQuery = $state("");
    let matchedMetadata = $state<GameMetadata>();

    $effect(() => {
        if (isOpen && localGame?.name) {
            searchQuery = localGame.name;
            clientRpc.runner()
        } else {
            searchQuery = ""
        }
    })

    const rpc = createRPCRunner(() => searchClient.match({query: searchQuery.trim()}));

    function search() {
        if (searchQuery.trim() !== '') {
            rpc.runner();
        }
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            search()
        }
    }

    const downClientManger = cli(DownloaderService)
    const clientRpc = createRPCRunner(() => downClientManger.getActiveClients({}))


    const libClient = cli(LibraryService)

    let downloadClient = $state("")
    let downloadUrl = $derived(localGame?.downloadUrl ?? "")
    let downloadPath = $state("./games/")
    let gameType = $state("")

    $effect(() => {
        if (clientRpc?.value?.clients && clientRpc?.value?.clients?.length > 0 && !downloadClient) {
            downloadClient = clientRpc.value.clients[0].name;
        }
    });

    async function handleAdd() {
        if (!downloadUrl || !downloadPath) {
            console.log("empty url and path")
            return
        }

        isOpen = false;

        const game = create(GameSchema, {
            DownloadState: {
                DownloadPath: downloadPath,
                DownloadUrl: downloadUrl,
                Client: downloadClient,
            },
            Meta: matchedMetadata,
            GameType: gameType
        });

        const {err} = await callRPC(() => libClient.add({game}))
        if (err) {
            console.error("failed to add game", err);
        }
    }

    function setMetadata(meta: GameMetadata, isMatched: boolean) {
        matchedMetadata = isMatched ? undefined : meta
    }
</script>

{#if isOpen}
    <div class="fixed inset-0 z-50 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-background/90 backdrop-blur-md" onclick={() => isOpen = false}></div>

        <!-- Content Container -->
        <div class="relative w-full max-w-7xl h-212.5 bg-surface border border-border rounded-2xl flex flex-col overflow-hidden shadow-2xl"
             transition:fly={{ y: 20, duration: 300 }}>

            <!-- HEADER -->
            <div class="flex items-center justify-between p-5 border-b border-border">
                <h2 class="text-2xl font-bold tracking-tight text-foreground">Add to library</h2>
                <button onclick={() => isOpen = false} class="text-muted hover:text-foreground p-1">
                    <XIcon size={24}/>
                </button>
            </div>

            <div class="flex-1 flex overflow-hidden">
                <div class="flex-[1.2] flex flex-col border-r border-border bg-background/30">
                    <!-- Provider Toolbar -->
                    <div class="p-4 border-b border-border bg-surface flex items-center gap-3">
                        <!-- Provider Selector -->
                        <div class="relative">
                            <select
                                    class="appearance-none bg-panel border border-border rounded-xl px-4 py-2.5 pr-10 text-sm font-bold uppercase tracking-widest outline-none focus:border-frost-500 transition-colors cursor-pointer min-w-[120px]"
                            >
                                <option>IGDB</option>
                                <option>Steam</option>
                                <option>TMDB</option>
                            </select>
                            <ChevronDownIcon
                                    size={16}
                                    class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted"
                            />
                        </div>

                        <div class="flex-1 relative group">
                            <input
                                    type="text"
                                    bind:value={searchQuery}
                                    onkeydown={(e) => e.key === 'Enter' && rpc.runner()}
                                    placeholder="Search for metadata..."
                                    class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 pr-14 outline-none focus:border-frost-500 transition-all text-sm placeholder:text-muted/50"
                            />

                            <button
                                    onclick={search}
                                    onkeydown={handleKeyDown}
                                    disabled={rpc.loading}
                                    class="absolute  right-1.5 top-1.5 bottom-1.5 w-10 flex items-center justify-center rounded-lg bg-frost-500 border border-border text-muted hover:text-frost-500 hover:border-frost-500/50 transition-all disabled:opacity-50 active:scale-95 shadow-sm"
                                    title="Search"
                            >
                                {#if rpc.loading}
                                    <LoaderIcon size={18} class="animate-spin text-frost-500"/>
                                {:else}
                                    <SearchIcon size={18}/>
                                {/if}
                            </button>
                        </div>
                    </div>

                    <!-- Provider Results List -->
                    <div class="flex-1 overflow-y-auto p-4 space-y-3">
                        {#if rpc.loading}
                            <!-- Loading State -->
                            <div class="flex flex-col items-center justify-center h-full text-muted">
                                <LoaderIcon class="animate-spin mb-2" size={32}/>
                                <p class="text-sm font-medium animate-pulse">Fetching metadata...</p>
                            </div>

                        {:else if rpc.error}
                            <!-- Error State -->
                            <div class="p-4 rounded-xl border border-red-500/20 bg-red-500/5 text-red-400 flex gap-3 items-start">
                                <TriangleAlert size={18}/>
                                <div class="text-sm">
                                    <p class="font-bold">Search failed</p>
                                    <p class="opacity-80">{rpc.error}</p>
                                </div>
                            </div>

                        {:else if rpc.value && rpc.value.metadata.length > 0}
                            <!-- Results State -->
                            {#each rpc.value.metadata as item}
                                {@const isMatched = matchedMetadata?.ID === item.ID}
                                <div class="flex gap-4 p-4 rounded-2xl border transition-all {isMatched ? 'border-frost-500 bg-frost-500/10 ring-1 ring-frost-500/20' : 'border-border bg-surface hover:bg-panel/50'}">
                                    <!-- Poster -->
                                    <div class="w-20 h-28 shrink-0 bg-panel rounded-lg overflow-hidden border border-border/50">
                                        {#if item.ThumbnailURL}
                                            <img src={item.ThumbnailURL} alt="" class="w-full h-full object-cover"/>
                                        {:else}
                                            <div class="w-full h-full flex items-center justify-center text-muted">
                                                <ImageIcon size={24}/>
                                            </div>
                                        {/if}
                                    </div>

                                    <!-- Meta Info -->
                                    <div class="flex-1 flex flex-col min-w-0">
                                        <div class="flex items-center gap-2 mb-1">
                                            <h3 class="font-bold text-foreground truncate">{item.Name}</h3>
                                            <span class="px-2 py-0.5 rounded-md bg-panel border border-border text-[10px] text-muted font-bold uppercase">2023</span>
                                        </div>
                                        <p class="text-xs text-muted line-clamp-3 leading-relaxed mb-4">{item.Description || 'No description available.'}</p>

                                        <div class="mt-auto flex justify-end">
                                            <button
                                                    onclick={() => setMetadata(item, isMatched)}
                                                    class="px-5 py-1.5 rounded-lg text-xs font-bold transition-all {isMatched ? 'bg-frost-500 text-background' : 'bg-panel border border-border hover:border-frost-500'}">
                                                {isMatched ? 'Matched' : 'Match'}
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            {/each}

                        {:else if searchQuery.trim() !== ''}
                            <!-- EMPTY STATE: No Results Found (Query exists but no metadata) -->
                            <div class="flex flex-col items-center justify-center h-full text-muted py-20" in:fade>
                                <div class="w-16 h-16 rounded-full bg-panel flex items-center justify-center mb-4 border border-border">
                                    <SearchXIcon size={32} strokeWidth={1.5}/>
                                </div>
                                <h3 class="font-bold text-foreground">No results found</h3>
                                <p class="text-xs text-center px-8 mt-1 opacity-60">
                                    We couldn't find any games matching "{searchQuery}" on the selected provider.
                                </p>
                            </div>

                        {:else}
                            <!-- EMPTY STATE: No Query Found (Initial state / Search bar is empty) -->
                            <div class="flex flex-col items-center justify-center h-full text-muted py-20" in:fade>
                                <div class="w-16 h-16 rounded-full bg-panel flex items-center justify-center mb-4 border border-border">
                                    <SearchIcon size={32} strokeWidth={1.5}/>
                                </div>
                                <h3 class="font-bold text-foreground">No query found</h3>
                                <p class="text-xs text-center px-8 mt-1 opacity-60">
                                    Type a game title in the search bar above to start matching metadata.
                                </p>
                            </div>
                        {/if}
                    </div>
                </div>

                <!-- RIGHT SIDE: INDEXER DATA -->
                <div class="flex-1 flex flex-col bg-surface p-6 overflow-y-auto">
                    <div class="space-y-8">

                        <!-- Block 1: Basic File Info -->
                        <section class="space-y-4">
                            <div class="p-4 bg-panel border border-border rounded-xl">
                                <p class="text-[10px] font-bold text-muted uppercase tracking-[0.2em] mb-2">Indexer
                                    Title</p>
                                <h3 class="text-lg font-bold text-foreground truncate">{localGame?.name || 'Waiting for selection...'}</h3>
                            </div>

                            <div class="grid grid-cols-3 gap-3">
                                <div class="bg-panel/50 border border-border p-3 rounded-xl flex flex-col">
                                    <span class="text-[9px] font-bold text-muted uppercase mb-1">Size</span>
                                    <span class="text-sm font-medium">{localGame?.size || '--'}</span>
                                </div>

                                <div class="bg-panel/50 border border-border p-3 rounded-xl flex flex-col">
                                    <span class="text-[9px] font-bold text-muted uppercase mb-1">Date</span>
                                    <span class="text-sm font-medium">TODO</span>
                                </div>

                                <div class="bg-panel/50 border border-border p-3 rounded-xl flex items-center justify-between group cursor-pointer hover:border-frost-500 transition-colors">
                                    <div class="flex flex-col truncate">
                                        <span class="text-[9px] font-bold text-muted uppercase mb-1">Download Path</span>
                                        <span class="text-sm font-medium truncate">{downloadPath}</span>
                                    </div>
                                    <FolderOpenIcon size={16} class="text-muted group-hover:text-frost-400 ml-2"/>
                                </div>
                            </div>
                        </section>

                        <hr class="border-border/50"/>

                        <!-- Block 2: Download Config -->
                        <section class="flex gap-4">
                            <div class="min-w-40">
                                <p class="text-[10px] font-bold text-muted uppercase tracking-[0.2em] mb-2">Clients</p>
                                <div class="relative">
                                    {#if clientRpc.loading}
                                        <!-- Styled as a placeholder to prevent layout shift -->
                                        <div class="w-full bg-panel/50 border border-border rounded-xl py-2.5 px-3 text-sm text-muted flex items-center gap-2 animate-pulse">
                                            <LoaderIcon size={14} class="animate-spin"/>
                                            <span>Loading...</span>
                                        </div>
                                    {:else if clientRpc.error}
                                        <!-- Error state with same dimensions -->
                                        <div class="w-full bg-red-500/5 border border-red-500/20 rounded-xl py-2.5 px-3 text-xs text-red-400 truncate"
                                             title={clientRpc.error}>
                                            Error: {clientRpc.error}
                                        </div>
                                    {:else}
                                        <!-- Actual Select -->
                                        <select
                                                bind:value={downloadClient}

                                                class="w-full bg-panel border border-border rounded-xl py-2.5 px-3 outline-none focus:border-frost-500 appearance-none text-sm font-bold cursor-pointer transition-colors"
                                        >
                                            {#if !clientRpc.value || clientRpc.value.clients.length === 0}
                                                <option value="" disabled selected>No clients found</option>
                                            {:else}
                                                {#each clientRpc.value.clients as client}
                                                    <!-- Assuming client might be a string or object; adjust accordingly -->
                                                    <option value={client.name}>{client.name}</option>
                                                {/each}
                                            {/if}
                                        </select>

                                        <!-- Chevron only shows when the select is active -->
                                        <ChevronDownIcon
                                                size={16}
                                                class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-muted"
                                        />
                                    {/if}
                                </div>
                            </div>
                            <div class="flex-1">
                                <p class="text-[10px] font-bold text-muted uppercase tracking-[0.2em] mb-2">Download
                                    URL</p>
                                <div class="relative">
                                    <LinkIcon size={16} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted"/>
                                    <input type="text" bind:value={downloadUrl} placeholder="magnet:?xt=..."
                                           class="w-full bg-panel border border-border rounded-xl py-2.5 pl-10 pr-4 outline-none focus:border-frost-500 text-sm"/>
                                </div>
                            </div>
                        </section>

                        <hr class="border-border/50"/>

                        <!-- Block 3: Matched Meta Preview -->
                        <section>
                            <p class="text-[10px] font-bold text-muted uppercase tracking-[0.2em] mb-3">Matched Meta</p>
                            {#if matchedMetadata}
                                <div class="flex gap-4 p-4 rounded-2xl border border-frost-500 bg-frost-500/5" in:fade>
                                    <div class="w-16 h-20 shrink-0 bg-panel rounded border border-border overflow-hidden">
                                        <img src={matchedMetadata.ThumbnailURL} class="w-full h-full object-cover"
                                             alt=""/>
                                    </div>
                                    <div class="flex-1 min-w-0">
                                        <div class="flex items-center gap-2 mb-1">
                                            <h4 class="font-bold text-sm truncate">{matchedMetadata.Name}</h4>
                                            <span class="text-[9px] font-mono text-frost-400">2023</span>
                                        </div>
                                        <p class="text-[11px] text-muted line-clamp-2">{matchedMetadata.Description}</p>
                                    </div>
                                    <button onclick={() => matchedMetadata = undefined}
                                            class="self-start p-1.5 text-muted hover:text-red-400 transition-colors">
                                        <XIcon size={16}/>
                                    </button>
                                </div>
                            {:else}
                                <div class="h-24 border-2 border-dashed border-border rounded-2xl flex items-center justify-center text-muted italic text-sm">
                                    No matched metadata
                                </div>
                            {/if}
                        </section>
                    </div>
                </div>
            </div>

            <!-- FOOTER -->
            <div class="p-5 border-t border-border flex justify-end items-center gap-3 bg-panel/20">
                <button onclick={() => isOpen = false}
                        class="px-8 py-2.5 rounded-xl border border-border hover:bg-panel transition-all font-bold text-sm">
                    Cancel
                </button>
                <button
                        onclick={handleAdd}
                        class={`px-10 py-2.5 rounded-xl ${downloadUrl ? "bg-frost-500" : "bg-frost-600"} text-background
                            hover:bg-frost-400 active:scale-95
                            transition-all font-bold text-sm shadow-lg shadow-frost-500/20`}
                >

                    Add Game
                </button>
            </div>
        </div>
    </div>
{/if}
