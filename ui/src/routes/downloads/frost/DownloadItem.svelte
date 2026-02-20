<script lang="ts">
    import {
        CheckCircle2Icon,
        ChevronDownIcon,
        ClockIcon,
        HourglassIcon,
        ImageIcon,
        PauseIcon,
        PlayIcon,
        XIcon
    } from "@lucide/svelte";
    import {slide} from 'svelte/transition';
    import {type DownloadProgress, FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {formatBytes} from "$lib/api/byte-math";
    import {callRPC, frostCli} from "$lib/api/api";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {trimPrefix} from "$lib/api/strings";

    let {detail}: { detail: DownloadProgress } = $props();

    let isExpanded = $state(false);

    let now = $state(Date.now());

    $effect(() => {
        const interval = setInterval(() => {
            now = Date.now();
        }, 1000);
        return () => clearInterval(interval);
    });

    const elapsedTime = $derived.by(() => {
        if (!detail.download?.Started) return null;
        const start = new Date(detail.download.Started).getTime();
        const diff = Math.max(0, now - start);

        const seconds = Math.floor((diff / 1000) % 60);
        const minutes = Math.floor((diff / (1000 * 60)) % 60);
        const hours = Math.floor((diff / (1000 * 60 * 60)));

        const parts = [
            hours > 0 ? hours.toString().padStart(2, '0') : null,
            minutes.toString().padStart(2, '0'),
            seconds.toString().padStart(2, '0')
        ].filter(Boolean);

        return parts.join(':');
    });

    const totalComplete = $derived(Number(detail.progress?.Complete ?? 0));
    const totalLeft = $derived(Number(detail.progress?.Left ?? 0));
    const totalSize = $derived(totalComplete + totalLeft);
    const overallProgress = $derived(totalSize > 0 ? (totalComplete / totalSize) * 100 : 0);

    function getStateColor(state: string = "") {
        const s = state.toLowerCase();
        if (s.includes('downloading') || s.includes('done')) return 'text-green-400 bg-green-400/10 border-green-400/20';
        if (s.includes('queued')) return 'text-amber-400 bg-amber-400/10 border-amber-400/20';
        return 'text-frost-400 bg-frost-400/10 border-frost-400/20';
    }

    const getFileProgress = (complete: bigint | number, left: bigint | number) => {
        const c = Number(complete);
        const l = Number(left);
        return (c + l) > 0 ? (c / (c + l)) * 100 : 0;
    };

    const frostLib = frostCli(FrostLibraryService)
    const sm = getSnackbarCtx()

    let isPaused = $derived(detail?.download?.Status === "StatusPaused")

    async function toggleDownload(ID: bigint) {
        if (isPaused) {
            const {err} = await callRPC(() => frostLib.pause({id: ID, pause: false}))
            if (err) {
                sm.push(`Could not cancel game: ${err}`, 'error')
            }
            return
        }

        const {err} = await callRPC(() => frostLib.pause({id: ID, pause: true}))
        if (err) {
            sm.push(`Could not cancel game: ${err}`, 'error')
        }
    }

    async function cancel(ID: bigint) {
        const {err} = await callRPC(() => frostLib.cancel({id: ID}))
        if (err) {
            sm.push(`Could not cancel game: ${err}`, 'error')
        }
    }
</script>

<div class="bg-panel/30 border border-border rounded-2xl overflow-hidden transition-all duration-300 {isExpanded ? 'border-frost-500/30 ring-1 ring-frost-500/10' : ''}">

    <!-- ACCORDION HEADER -->
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
            onclick={() => isExpanded = !isExpanded}
            class="w-full flex items-center gap-4 p-3 hover:bg-panel/50 transition-colors text-left group cursor-pointer"
    >
        <!-- THUMBNAIL -->
        <div class="w-12 h-16 bg-surface border border-border rounded-xl overflow-hidden shrink-0 flex items-center justify-center text-muted/20">
            {#if detail.Thumbnail}
                <img src={detail.Thumbnail} alt=""
                     class="w-full h-full object-cover transition-transform group-hover:scale-110"/>
            {:else}
                <ImageIcon size={20} strokeWidth={1.5}/>
            {/if}
        </div>

        <!-- MAIN INFO SECTION -->
        <div class="flex-1 min-w-0 space-y-1.5">
            <div class="flex justify-between items-center gap-2">
            <span class="font-bold text-sm text-foreground truncate tracking-tight">
                {detail.Title || 'Unknown Package'}
            </span>
                <span class="px-2 py-0.5 rounded-md bg-panel/50 border border-border text-[10px] font-mono font-bold text-frost-400">
                {overallProgress.toFixed(1)}%
            </span>
            </div>

            <!-- Meta Row: State, Message, Elapsed -->
            <div class="flex items-center gap-3 text-[10px] font-medium">
            <span class="px-1.5 py-0.5 rounded border uppercase tracking-tighter shrink-0 {getStateColor(detail.download?.Status)}">
                {trimPrefix(detail.download?.Status ?? "", "Status") || 'Waiting'}
            </span>
                {#if elapsedTime}
                    <span class="text-white ml-auto whitespace-nowrap flex items-center gap-1 font-mono">
                        <HourglassIcon size={10} class="text-frost-500/50"/>
                        {elapsedTime}
                    </span>
                {/if}
            </div>

            <!-- Overall Progress Bar -->
            <div class="h-1.5 w-full bg-panel rounded-full overflow-hidden border border-white/5 mt-1">
                <div
                        class="h-full bg-frost-500 transition-all duration-700 ease-out shadow-[0_0_8px_rgba(130,170,255,0.3)]"
                        style="width: {overallProgress}%"
                ></div>
            </div>
        </div>

        <!-- BYTE STATS SECTION -->
        <div class="flex items-center gap-3 ml-2 shrink-0">
            <div class="hidden md:flex flex-col items-end text-[10px] font-bold uppercase tracking-widest text-muted/40 whitespace-nowrap">
                <span class="text-foreground/60">{formatBytes(totalComplete)} Done</span>
                <span>{formatBytes(totalLeft)} Left</span>
            </div>

            <!-- Action Buttons: Stop propagation so they don't toggle the accordion -->
            {#if overallProgress < 100}
                <div class="flex items-center gap-1 transition-all translate-x-2 ml-2">
                    <button
                            onclick={(e) => { e.stopPropagation(); toggleDownload(detail.ID); }}
                            class="p-1.5 rounded-lg text-muted hover:text-amber-400 hover:bg-amber-400/10 transition-all"
                            title="Pause Download"
                    >
                        {#if isPaused}
                            <PlayIcon size={14}/>
                        {:else}
                            <PauseIcon size={14}/>
                        {/if}
                    </button>
                    <button
                            onclick={(e) => { e.stopPropagation(); cancel(detail.ID); }}
                            class="p-1.5 rounded-lg text-muted hover:text-red-400 hover:bg-red-400/10 transition-all"
                            title="Cancel Download"
                    >
                        <XIcon size={14}/>
                    </button>
                </div>
            {/if}

            <div class="p-1 rounded-lg hover:bg-panel transition-colors text-muted">
                <ChevronDownIcon
                        size={18}
                        class="transition-transform duration-300 {isExpanded ? 'rotate-180 text-frost-400' : ''}"
                />
            </div>
        </div>
    </div>

    {#if isExpanded}
        <!-- ACCORDION CONTENT (File List) -->
        {#if isExpanded}
            <div transition:slide={{ duration: 300 }} class="border-t border-border bg-black/10">
                <div class="p-2 space-y-1 max-h-80 overflow-y-auto custom-scrollbar">
                    {#each detail.progress?.files ?? [] as file}
                        {@const filePct = getFileProgress(file.BytesComplete, file.BytesLeft)}
                        {@const isDone = Number(file.BytesLeft) === 0}

                        <!-- Chunk Logic -->
                        {@const chunksComplete = Number(file.ChunksComplete ?? 0)}
                        {@const chunksLeft = Number(file.ChunksLeft ?? 0)}
                        {@const chunksTotal = chunksComplete + chunksLeft}

                        <!--
                            Performance Guard:
                            If a file has 10,000 chunks, we don't want 10,000 divs.
                            We cap the visual blocks at 60 for a consistent "bitfield" look.
                        -->
                        {@const visualBlockCount = chunksTotal > 60 ? 60 : chunksTotal}
                        {@const filledBlocks = Math.floor((chunksComplete / chunksTotal) * visualBlockCount)}

                        <div class="flex items-center gap-3 p-3 rounded-xl hover:bg-panel/40 transition-colors group">
                            <div class="text-muted/40 group-hover:text-frost-400/50 transition-colors shrink-0">
                                {#if isDone}
                                    <CheckCircle2Icon size={14} class="text-green-500/60"/>
                                {:else}
                                    <ClockIcon size={14}/>
                                {/if}
                            </div>

                            <div class="flex-1 min-w-0">
                                <div class="flex items-center gap-3 mb-1.5">
                                    <p class="text-[12px] font-medium text-white group-hover:text-foreground transition-colors truncate min-w-0 flex-1">
                                        {file.Name}
                                    </p>

                                    <span class="text-[10px] font-mono text-white whitespace-nowrap shrink-0">
                                        {formatBytes(Number(file.BytesComplete))} / {formatBytes(Number(file.BytesComplete) + Number(file.BytesLeft))}
                                                                        <span class="text-white/60 ml-1">
                                            ({Math.round((Number(file.BytesComplete) / (Number(file.BytesComplete) + Number(file.BytesLeft))) * 100)}%)
                                        </span>
                                    </span>

                                    <span class="text-[10px] font-black uppercase tracking-tighter text-white whitespace-nowrap shrink-0">
                                        Blocks: {chunksComplete}/{chunksTotal}
                                        <span class="text-white/60 ml-1">
                                            ({chunksTotal > 0 ? Math.round((chunksComplete / chunksTotal) * 100) : 0}%)
                                        </span>
                                    </span>
                                </div>


                                <div class="space-y-1">

                                    <!-- 2. Segmented Chunk-level Bitfield Bar -->
                                    {#if chunksTotal > 0}
                                        <div class="flex gap-px h-1 w-full" aria-hidden="true">
                                            {#each Array(visualBlockCount) as _, i}
                                                <div
                                                        class="flex-1 h-full rounded-[1px] transition-colors duration-500
                                            {i < filledBlocks
                                                ? (isDone ? 'bg-green-500/20' : 'bg-frost-500/60 shadow-[0_0_3px_rgba(130,170,255,0.2)]')
                                                : 'bg-white/3'}"
                                                ></div>
                                            {/each}
                                        </div>
                                    {/if}
                                </div>
                            </div>
                        </div>
                    {/each}
                </div>
            </div>
        {/if}
    {/if}
</div>

<style>
    button {
        user-select: none;
    }

    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }

    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }

    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.1);
        border-radius: 10px;
    }
</style>
