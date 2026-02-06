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
    import {type GameSource, SearchService} from "$lib/gen/search/v1/search_pb";
    import {ServiceConfigService} from "$lib/gen/service_config/v1/service_config_pb";
    import AddButton from "./AddButton.svelte";
    import IndexerSearch from "./IndexerSearch.svelte";

    const srvConfig = glacierCli(ServiceConfigService)
    let indexerManagerService = glacierCli(IndexerService)

    let game = $derived(transferStore.data);
    $effect(() => {
        if (game) {
            availableDownloadClientsRpc.runner()
            availableGameTypeRpc.runner()
        }
    })

    let availableGameTypeRpc = createRPCRunner(() => indexerManagerService.getGameType({}))
    let selectedGameType = $state("");
    $effect(() => {
        if (selectedGameType === "") {
            selectedGameType = availableGameTypeRpc.value?.gameTypes.at(0)?.Name ?? ""
        }

        if (selectedGameSource) {
            selectedGameSource.GameType = selectedGameType
        }
    })

    let availableDownloadClientsRpc = createRPCRunner(() => srvConfig.getActiveService({serviceType: "Downloader"}))
    let selectedDownloadClient = $state("");
    $effect(() => {
        selectedDownloadClient = availableDownloadClientsRpc.value?.names.at(0)?.Name ?? ""
    })

    let selectedGameSource = $state(null as GameSource | null);


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
                    <!--                    <a href="#" class="p-2 bg-panel rounded-lg text-frost-400 hover:text-frost-300 transition-colors">-->
                    <!--                        <ExternalLinkIcon size={18}/>-->
                    <!--                    </a>-->
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
            <AddButton
                    gameSrc={selectedGameSource}
                    gameMetadata={game}
                    downloadClient={selectedDownloadClient}
            />
        </div>
    </div>

    <!-- SEARCH & DOWNLOAD SECTION -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- SEARCH INDEXER SECTION -->
        <IndexerSearch game={game} bind:selectedGameSource={selectedGameSource}/>
        <!-- DOWNLOAD CONFIG SECTION -->
        <div class="bg-surface border border-border rounded-3xl p-6 flex flex-col gap-8">
            <h2 class="text-lg font-bold flex items-center gap-2">
                <HardDriveIcon size={20} class="text-frost-400"/>
                Download Config
            </h2>

            <div class="p-5 bg-panel border border-border rounded-2xl flex items-center start">
                <div class="w-16 h-20 bg-background rounded border border-border flex items-center justify-center text-muted">
                    <ImageIcon size={24}/>
                </div>
                <div class="px-5 space-y-1">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">
                        Selected Release
                    </p>
                    <h3 class="font-bold text-foreground">
                        {selectedGameSource ? selectedGameSource.Title : "Waiting for selection..."}
                    </h3>
                    {#if !selectedGameSource}
                        <p class="text-xs text-muted italic">
                            No release matched yet
                        </p>
                    {/if}
                </div>
            </div>

            <div class="flex items-center gap-4">
                <p class="text-[10px] font-bold text-muted uppercase tracking-widest whitespace-nowrap">
                    Game Type
                </p>

                <div class="relative flex-1">
                    <select bind:value={selectedGameType}
                            class="w-full bg-panel border border-border rounded-xl py-3 px-4 appearance-none outline-none focus:border-frost-500 text-sm font-bold cursor-pointer">
                        {#each availableGameTypeRpc.value?.gameTypes as gameType}
                            <option value={gameType.Name}>{gameType.Name}</option>
                        {/each}
                    </select>
                    <ChevronDownIcon
                            size={16}
                            class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-muted"/>
                </div>
            </div>

            <div class="flex gap-4">
                <div class="w-1/3 space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">
                        Download Client
                    </p>
                    <div class="relative">
                        <select bind:value={selectedDownloadClient}
                                class="w-full bg-panel border border-border rounded-xl py-2.5 px-4 appearance-none outline-none text-sm font-bold"
                        >
                            {#each availableDownloadClientsRpc.value?.names as cli}
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
                        <input type="text" placeholder="link or URL..."
                               value={selectedGameSource?.DownloadUrl ?? ""}
                               class="w-full bg-panel border border-border rounded-xl py-2.5 pl-10 pr-4 outline-none text-sm"
                        />
                    </div>
                </div>
            </div>

            <div class="flex items-end gap-4">
                <div class="flex-1 space-y-2">
                    <p class="text-[10px] font-bold text-muted uppercase tracking-widest">
                        Download Path
                    </p>
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
