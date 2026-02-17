<script lang="ts">
    import {
        CalendarIcon,
        ChevronLeft,
        ChevronRight,
        CircleAlert,
        ImageIcon,
        LoaderIcon,
        MonitorIcon,
        Pen,
        PlayIcon,
        StarIcon
    } from '@lucide/svelte';
    import type {Game} from "$lib/gen/library/v1/library_pb";
    import {getAssetPath} from "$lib/api/assets";
    import AssetImg from "$lib/components/assets/AssetImg.svelte";
    import AssetVideo from "$lib/components/assets/AssetVideo.svelte";
    import GameDownloadButton from "$lib/components/ButtonDownload.svelte";
    import {glacierCli, isFrost} from "$lib/api/api";
    import {AssetService} from "$lib/gen/assets/v1/assets_pb";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {onMount} from "svelte";
    import {fade} from 'svelte/transition';
    import {trimPrefix} from "$lib/api/strings.ts";

    let {game = $bindable(null), onManage}: { game: Game | null, onManage: () => void } = $props();

    let meta = $derived(game?.Meta);
    const thumb = getAssetPath({gameId: game?.ID, assetType: "AssetThumbnail"});
    const trailer = getAssetPath({gameId: game?.ID, assetType: "AssetTrailer"});

    const assetSrv = glacierCli(AssetService)

    const listRpc = createRPCRunner(() => assetSrv.list({
        ID: game?.ID,
        AssetType: assetType
    }))

    const assetType = $state("")
    let activeIndex = $state(0);
    let assets = $derived(listRpc.value?.assets ?? []);
    let currentAsset = $derived(assets[activeIndex]);

    function next() {
        activeIndex = (activeIndex + 1) % assets.length;
    }

    let getImg = (inputLink: string): string => {
        if (!inputLink) return "";

        return getAssetPath({
            gameId: game?.ID,
            assetPath: inputLink
        });
    }

    function prev() {
        activeIndex = (activeIndex - 1 + assets.length) % assets.length;
    }

    function isVideo(type: string) {
        return type.toLowerCase().includes('video') ||
            type.toLowerCase().includes('trailer');
    }

    onMount(async () => {
        if (!game?.ID) {
            return
        }
        await listRpc.runner()
    })
</script>

<div class="space-y-4">
    <div class="bg-surface/40 backdrop-blur-xl border border-white/10 rounded-3xl p-6 shadow-sm relative w-full aspect-21/9 lg:aspect-30/9  overflow-hidden group ">
        {#if listRpc.loading && !listRpc.value}
            <div class="flex h-full w-full items-center justify-center bg-background/50 animate-pulse">
                <LoaderIcon class="animate-spin text-frost-500" size={40}/>
            </div>
        {:else if listRpc.error}
            <div class="flex flex-col h-full w-full items-center justify-center text-red-400 gap-2">
                <CircleAlert size={32}/>
                <p class="text-xs font-black uppercase tracking-widest">Failed to load gallery</p>
            </div>
        {:else if assets.length > 0}
            <!-- Asset Label Overlay -->
            <div class="absolute top-6 right-8 px-4 py-1.5 rounded-full bg-background/40 backdrop-blur-md border border-white/10 text-[10px] font-black uppercase tracking-[0.2em] text-white/90 z-20 pointer-events-none">
                {trimPrefix(currentAsset.Type, "Asset")}
            </div>
            <!-- Navigation Arrows -->
            {#if assets.length > 1}
                <button
                        onclick={prev}
                        class="absolute left-4 top-1/2 -translate-y-1/2 z-20 p-3 rounded-full bg-background/20 backdrop-blur-md border border-white/10 opacity-0 group-hover:opacity-100 transition-all text-white hover:bg-frost-500 hover:text-background active:scale-95 shadow-xl"
                >
                    <ChevronLeft size={24} strokeWidth={3}/>
                </button>
                <button
                        onclick={next}
                        class="absolute right-4 top-1/2 -translate-y-1/2 z-20 p-3 rounded-full bg-background/20 backdrop-blur-md border border-white/10 opacity-0 group-hover:opacity-100 transition-all text-white hover:bg-frost-500 hover:text-background active:scale-95 shadow-xl"
                >
                    <ChevronRight size={24} strokeWidth={3}/>
                </button>
            {/if}

            <!-- Active Asset Display -->
            {#key activeIndex}
                <div in:fade={{ duration: 400 }} class="w-full h-full">
                    {#if isVideo(currentAsset.Type)}
                        <AssetVideo
                                src={getImg(currentAsset.LocalPath) || currentAsset.RemoteUrl}
                                autoplay={true}
                                muted={true}
                                controls={true}
                                class="w-full h-full object-cover"
                        >
                            {#snippet loadingSlot()}
                                <div class="flex h-full w-full items-center justify-center bg-black">
                                    <LoaderIcon class="animate-spin text-frost-500"/>
                                </div>
                            {/snippet}
                        </AssetVideo>
                    {:else}
                        <AssetImg
                                src={getImg(currentAsset.LocalPath) || currentAsset.RemoteUrl}
                                class="w-full h-full "
                                imgClass="w-full h-full object-contain"
                        />
                    {/if}
                </div>
            {/key}

            <!-- Indicators (Dots) -->
            <div class="absolute bottom-6 left-1/2 -translate-x-1/2 flex gap-2 z-20">
                {#each assets as _, i}
                    <button
                            title="Indicator"
                            onclick={() => activeIndex = i}
                            class="h-1.5 transition-all rounded-full
                        {activeIndex === i ? 'w-8 bg-frost-500 shadow-[0_0_10px_rgba(130,170,255,0.5)]' : 'w-2 bg-white/20 hover:bg-white/40'}"
                    ></button>
                {/each}
            </div>
        {:else}
            <!-- Empty State -->
            <div class="flex flex-col h-full w-full items-center justify-center text-muted/20 gap-3">
                <ImageIcon size={64} strokeWidth={1}/>
                <p class="text-xs font-black uppercase tracking-[0.3em]">No Media Assets</p>
            </div>
        {/if}
    </div>

    <!-- MIDDLE SECTION: INFO BAR -->
    <div class="bg-surface/40 backdrop-blur-xl border border-white/10 rounded-3xl p-6 shadow-sm">
        <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-6">
            <div class="space-y-4">
                <h1 class="text-3xl font-black tracking-tight text-foreground">{meta?.Name}</h1>

                <div class="flex flex-wrap items-center gap-1 text-sm font-bold">
                    <!-- Rating -->
                    <div class="flex items-center gap-2 px-3 py-1.5 bg-transparent  rounded-xl">
                        <StarIcon size={16} class="text-yellow-500 fill-yellow-500"/>
                        <span>{meta?.Rating} <span class="text-muted font-medium">({meta?.RatingCount})</span></span>
                    </div>

                    <div class="w-px h-4 bg-border hidden lg:block"></div>

                    <!-- Release Year -->
                    <div class="flex items-center gap-2 px-3 py-1.5 bg-transparent  rounded-xl">
                        <CalendarIcon size={16} class="text-frost-400"/>
                        <span>{new Date(meta?.ReleaseDate || "").getFullYear()}</span>
                    </div>

                    <div class="w-px h-4 bg-border hidden lg:block"></div>

                    <!-- Platform Tags -->
                    <div class="flex items-center gap-2 px-3 py-1.5 bg-transparent  rounded-xl">
                        <MonitorIcon size={16} class="text-frost-400"/>
                        {#each meta?.Platforms.slice(0, 2) as Platform}
                            <span class="px-3 py-1 bg-frost-500/5 border border-frost-500/15 text-muted rounded-lg text-[11px] font-black uppercase tracking-wider">
                                {Platform}
                            </span>
                        {/each}
                    </div>
                </div>

                <!-- Genre Tags -->
                <div class="flex flex-wrap gap-2">
                    {#each meta?.Genres.slice(0, 4) ?? [] as genre}
                        <span class="px-3 py-1 bg-frost-500/5 border border-frost-500/15 text-frost-400 rounded-lg text-[11px] font-black uppercase tracking-wider">
                            {genre}
                        </span>
                    {/each}
                    {#if meta?.Genres?.length && meta?.Genres?.length > 4}
                        <span class="px-3 py-1 bg-panel border border-border text-muted rounded-lg text-[11px] font-bold">
                            +{meta.Genres.length - 4} more
                        </span>
                    {/if}
                </div>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-3 shrink-0">
                <button
                        onclick={onManage}
                        class="px-8 py-3 bg-white/5 border border-white/10 rounded-2xl text-sm font-black uppercase tracking-widest hover:border-frost-500/60 hover:bg-white/10 transition-all flex items-center gap-2"
                >
                    <Pen size={18}/>
                    Manage
                </button>

                {#if isFrost && game?.DownloadState?.State === "Complete"}
                    <GameDownloadButton gameId={game.ID}/>
                {/if}
            </div>
        </div>
    </div>
</div>
