<script lang="ts">
    import {
        ChevronDownIcon,
        CircleAlert,
        ExternalLinkIcon,
        FolderOpenIcon,
        HardDriveIcon,
        ImageIcon,
        LinkIcon,
        LoaderIcon,
        SearchIcon,
        StarIcon
    } from '@lucide/svelte';
    import {transferStore} from "../selectedGame.svelte";
    import {glacierCli} from "$lib/api/api";
    import {IndexerService} from "$lib/gen/indexer/v1/indexer_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {type GameSource, GameSourceSchema, SearchService} from "$lib/gen/search/v1/search_pb";
    import {DownloadSchema, GameSchema, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {create,} from "@bufbuild/protobuf";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";


    let game = $derived(transferStore.data);
    let selectedGameType = $state("");

    let searchQuery = $state("");
    let searchService = glacierCli(SearchService)
    let indexerSearchRpc = createRPCRunner(() => searchService.searchIndexers({
        q: {
            indexer: selectedIndexer,
            query: searchQuery,
        }
    }))

    function searchIndexer() {
        if (searchQuery.trim()) {
            indexerSearchRpc.runner()
        }
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') searchIndexer();
    }

    let selectedIndexer = $state("");

    let selectedSource = $state(null as GameSource | null);

    function handleMatch(item: GameSource) {
        if (selectedSource?.Title === item.Title) {
            selectedSource = null
        } else {
            selectedSource = item
        }
    }

    let srvConfig = glacierCli(ServiceConfigService)
    let activeClientsRpc = createRPCRunner(() => srvConfig.getActiveService({serviceType: "Downloader"}))

    let selectedClient = $state("");
    let libraryService = glacierCli(LibraryService)
    let libraryGameTypeRpc = createRPCRunner(() => indexerManagerService.getGameType({}))

    function onAddGame() {
        const {...pureData} = selectedSource;

        libraryService.add({
            game: create(GameSchema, {
                Meta: game ?? {},
                DownloadState: create(DownloadSchema, {
                    DownloadUrl: selectedSource?.DownloadUrl,
                    Client: selectedClient,

                }),
                Source: create(GameSourceSchema, {
                    ...pureData
                })
            })
        })
    }

    let indexerManagerService = glacierCli(IndexerService)
    let activeIndexerRpc = createRPCRunner(() => srvConfig.getActiveService({serviceType: "Indexer"}))
    $effect(() => {
        if (game) {
            searchQuery = game.Name
            activeIndexerRpc.runner()
            activeClientsRpc.runner()
            libraryGameTypeRpc.runner()
        }
    })
</script>

<div class="max-w-7xl mx-auto p-6 space-y-8 bg-background text-foreground">

    <!-- 1. VIDEO IFRAME SECTION (Expanded fully width-wise, reduced height) -->
    <!--    <div class="w-full aspect-[21/9] bg-panel rounded-3xl border border-border flex items-center justify-center overflow-hidden shadow-2xl">-->
    <!--        <div class="flex flex-col items-center gap-4 opacity-20">-->
    <!--            <PlusIcon size={64} class="rotate-45"/>-->
    <!--            <p class="font-bold tracking-widest uppercase">Video Preview Placeholder</p>-->
    <!--        </div>-->
    <!--    </div>-->

    <!-- 2. METADATA & ACTION ROW -->
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-8 items-start">

        <!-- Metadata Content -->
        <div class="lg:col-span-3 bg-surface border border-border rounded-3xl p-8 flex gap-8">
            <div class="flex-1 space-y-4">
                <img src={game?.ThumbnailURL} alt="game thumbnail"/>
            </div>
            <!-- Left: Title & Desc -->
            <div class="flex-1 space-y-4">
                <div class="flex items-center gap-3">
                    <h1 class="text-3xl font-bold">{game?.Name || "Game Title"}</h1>
                    <a href="#" class="p-2 bg-panel rounded-lg text-frost-400 hover:text-frost-300 transition-colors">
                        <ExternalLinkIcon size={18}/>
                    </a>
                </div>
                <div class="h-px bg-border w-full"></div>
                <p class="text-sm text-muted leading-relaxed line-clamp-6">
                    {game?.Summary || "No description available for this title."}
                </p>
            </div>

            <!-- Middle: Provider Specs -->
            <div class="w-px bg-border"></div>
            <div class="flex-1 space-y-4">
                <div class="space-y-3">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Provider Info</p>
                    <div class="space-y-2 text-sm">
                        <div class="flex justify-between"><span class="text-muted">Source:</span> <span
                                class="font-medium">{game?.ProviderType}</span></div>
                        <div class="flex justify-between">
                            <span class="text-muted">Released:</span>
                            <span class="font-medium">{game?.ReleaseDate ? new Date(game?.ReleaseDate).getFullYear() : "unknown"}</span>
                        </div>
                        <div class="flex justify-between">
                            <span class="text-muted">Status:</span>
                            <span class="px-2 bg-green-500/10 text-green-400 rounded text-xs">{game?.ReleaseStatus}</span>
                        </div>
                        <div class="flex items-center gap-2 mt-4 text-frost-400">
                            <StarIcon size={16} fill="currentColor"/>
                            <span class="font-bold">{game?.Rating ?? "0.0"}</span>
                            <span class="text-xs text-muted">({game?.RatingCount})</span>
                        </div>
                    </div>
                </div>
            </div>

            <!-- Right: Tags -->
            <div class="w-px bg-border"></div>
            <div class="flex-1 space-y-6">
                <div class="space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Genres</p>
                    <div class="flex flex-wrap gap-2">
                        {#each game?.Genres ?? [] as ge}
                            <span class="px-2 py-1 bg-panel border border-border rounded text-[11px]">{ge}</span>
                        {/each}
                    </div>
                </div>
                <div class="space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Platforms</p>
                    <div class="flex flex-wrap gap-2">
                        {#each game?.Platforms ?? [] as g}
                            <span class="px-2 py-1 bg-frost-500/10 text-frost-400 border border-frost-500/20 rounded text-[11px]">{g}</span>
                        {/each}
                    </div>
                </div>
            </div>
        </div>

        <!-- Action Card (Shorter, height-fit) -->
        <div class="bg-surface border border-border rounded-3xl p-6 flex flex-col gap-6 shadow-xl h-fit">
            <div class="space-y-4">
                <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Library Config</p>
                <div class="relative">
                    <select bind:value={selectedGameType}
                            class="w-full bg-panel border border-border rounded-xl py-3 px-4 appearance-none outline-none focus:border-frost-500 text-sm font-bold cursor-pointer">
                        {#each libraryGameTypeRpc.value?.gameTypes as gameType}
                            <option value={gameType.Name}>{gameType.Name}</option>
                        {/each}
                    </select>
                    <ChevronDownIcon size={16}
                                     class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-muted"/>
                </div>
            </div>
            <button
                    onclick={onAddGame}
                    class="w-full py-4 bg-frost-500 text-background font-bold rounded-2xl hover:bg-frost-400 active:scale-95 transition-all shadow-lg shadow-frost-500/20">
                Add to Library
            </button>
        </div>
    </div>

    <!-- SEARCH & DOWNLOAD SECTION -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- SEARCH INDEXER SECTION -->
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
                        {#each activeIndexerRpc?.value?.names as ind }
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
                        disabled={indexerSearchRpc.loading || !searchQuery.trim()}
                        class=" right-1.5 top-1.5 bottom-1.5 w-10 flex items-center justify-center rounded-lg bg-surface border border-border text-muted hover:text-frost-500 hover:border-frost-500/50 transition-all disabled:opacity-50"
                >
                    {#if indexerSearchRpc.loading}
                        <LoaderIcon size={18} class="animate-spin text-frost-500"/>
                    {:else}
                        <SearchIcon size={18}/>
                    {/if}
                </button>
            </div>

            <div class="space-y-3 h-100 overflow-y-auto pr-2">

                {#if indexerSearchRpc.error}
                    <aside class="flex items-start gap-4 p-4 rounded-xl bg-red-500/5 border border-red-500/20 text-red-400">
                        <CircleAlert size={20}/>
                        <div>
                            <h3 class="font-bold">Search Failed</h3>
                            <p class="text-sm opacity-80">{indexerSearchRpc.error}</p>
                        </div>
                    </aside>

                {:else if indexerSearchRpc.loading}
                    <div class="flex flex-col items-center justify-center h-96 text-muted gap-4">
                        <LoaderIcon class="animate-spin size-10 text-frost-500"/>
                        <p class="animate-pulse font-medium">Indexing results...</p>
                    </div>

                {:else if indexerSearchRpc.value?.results}
                    {#each indexerSearchRpc?.value?.results as item}
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
                                {selectedSource?.Title === item.Title ? "Selected" : "Select"}
                            </button>
                        </div>
                    {/each}
                {/if}
            </div>
        </div>

        <!-- DOWNLOAD CONFIG SECTION -->
        <div class="bg-surface border border-border rounded-3xl p-6 flex flex-col gap-8">
            <h2 class="text-lg font-bold flex items-center gap-2">
                <HardDriveIcon size={20} class="text-frost-400"/>
                Download Config
            </h2>

            <div class="p-5 bg-panel border border-border rounded-2xl flex items-center justify-between">
                <div class="space-y-1">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Selected Release</p>
                    <h3 class="font-bold text-foreground">{selectedSource ? selectedSource.Title : "Waiting for selection..."}</h3>
                    {#if !selectedSource}
                        <p class="text-xs text-muted italic">No release matched yet</p>
                    {/if}
                </div>
                <div class="w-16 h-20 bg-background rounded border border-border flex items-center justify-center text-muted">
                    <ImageIcon size={24}/>
                </div>
            </div>

            <div class="flex gap-4">
                <div class="w-1/3 space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Download Client</p>
                    <div class="relative">
                        <select
                                bind:value={selectedClient}
                                class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 appearance-none outline-none text-sm font-bold">
                            {#each activeClientsRpc.value?.names as cli}
                                <option>{cli.Name}</option>
                            {/each}
                        </select>
                        <ChevronDownIcon size={16} class="absolute right-3 top-1/2 -translate-y-1/2 text-muted"/>
                    </div>
                </div>
                <div class="flex-1 space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Download URL</p>
                    <div class="relative">
                        <LinkIcon size={14} class="absolute left-4 top-1/2 -translate-y-1/2 text-muted"/>
                        <input type="text" placeholder="Magnet link or URL..."
                               value={selectedSource?.DownloadUrl ?? ""}
                               class="w-full bg-panel border border-border rounded-xl py-2.5 pl-10 pr-4 outline-none text-sm"/>
                    </div>
                </div>
            </div>

            <div class="flex items-end gap-4">
                <div class="flex-1 space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">Download Path</p>
                    <div class="flex gap-2">
                        <div class="flex-1 bg-panel border border-border rounded-xl py-2.5 px-4 text-sm text-muted truncate">
                            /mnt/storage/games/...
                        </div>
                        <button class="p-2.5 bg-panel border border-border rounded-xl hover:text-frost-400 transition-colors">
                            <FolderOpenIcon size={20}/>
                        </button>
                    </div>
                </div>

                <div class="space-y-2 min-w-35">
                    <div class="p-3 bg-panel/50 border border-border rounded-xl space-y-1">
                        <div class="flex justify-between text-[9px] font-bold text-muted uppercase">
                            <span>Available:</span>
                            <span class="text-foreground">1.2 TB</span>
                        </div>
                        <div class="flex justify-between text-[9px] font-bold text-muted uppercase">
                            <span>Remaining:</span>
                            <span class="text-frost-400">1.1 TB</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
