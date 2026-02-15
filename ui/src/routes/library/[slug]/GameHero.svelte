<script lang="ts">
    import {CalendarIcon, LoaderIcon, ImageIcon, PlusIcon, StarIcon, SwatchBookIcon} from '@lucide/svelte';
    import type {Game} from "$lib/gen/library/v1/library_pb";
    import {getAssetPath} from "$lib/api/assets";
    import AssetImg from "$lib/components/assets/AssetImg.svelte";
    import AssetVideo from "$lib/components/assets/AssetVideo.svelte";


    let {game = $bindable(null)}: { game: Game | null } = $props();

    let meta = $derived(game?.Meta)

    const thumb = getAssetPath({gameId: game?.ID, assetType: "AssetThumbnail"})

    const trailer = getAssetPath({gameId: game?.ID, assetType: "AssetTrailer"})
    let videoError = $state(false);

</script>

<div class="space-y-6">
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 h-80">
        <!-- Main Poster -->
        <div class="lg:col-span-2 bg-panel rounded-3xl border border-border overflow-hidden relative">
            <AssetImg
                    src={thumb}
                    class="w-full h-full"
            >
                {#snippet loadingSlot()}
                    <div class="flex h-full w-full items-center justify-center bg-surface animate-pulse">
                        <LoaderIcon class="animate-spin text-muted"/>
                    </div>
                {/snippet}

                {#snippet errorSlot()}
                    <div class="flex h-full w-full items-center justify-center bg-surface text-muted/20">
                        <ImageIcon size={48}/>
                    </div>
                {/snippet}
            </AssetImg>
        </div>

        <!-- Video Trailer -->
        <div class="lg:col-span-8 bg-panel rounded-3xl border border-border flex items-center justify-center overflow-hidden relative group">
            <AssetVideo
                    src={trailer}
                    autoplay={true}
                    muted={true}
                    class="w-full h-full object-contain"
            >
                {#snippet loadingSlot()}
                    <div class="flex h-full w-full items-center justify-center bg-black">
                        <LoaderIcon class="animate-spin text-white/50"/>
                    </div>
                {/snippet}

                {#snippet errorSlot()}
                    <div class="flex flex-col h-full w-full items-center justify-center bg-surface text-muted-foreground">
                        <PlusIcon size={48} class="rotate-45 opacity-50"/>
                        <p class="text-xs font-bold uppercase tracking-widest">Video Unavailable</p>
                    </div>
                {/snippet}
            </AssetVideo>
        </div>

        <!-- Quick Meta Tags -->
        <div class="lg:col-span-2 flex flex-col gap-3">
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Release Date</span>
                <span class="text-sm font-bold flex items-center gap-2">
                    <CalendarIcon size={14} class="text-frost-400"/>
                    {new Date(meta?.ReleaseDate || "N/A").getFullYear()}
                </span>
            </div>
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Genre</span>
                <span class="text-sm font-bold flex items-center gap-2">
                    <SwatchBookIcon size={14} class="text-frost-400"/> {meta?.Genres.join(", ")}</span>
            </div>
            <div class="p-4 bg-surface border border-border rounded-2xl flex flex-col gap-1">
                <span class="text-[10px] font-bold text-muted uppercase">Rating</span>
                <span class="text-sm font-bold flex items-center gap-2">
                    <StarIcon size={14} class="text-yellow-500 fill-yellow-500"/>
                    {meta?.Rating} ({meta?.RatingCount ?? 0})
                </span>
            </div>
        </div>
    </div>
</div>