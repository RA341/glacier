<script lang="ts">
    import {
        AlertCircleIcon,
        ChevronRightIcon,
        ClockIcon,
        CornerLeftUpIcon,
        FileIcon,
        FolderIcon,
        HardDriveIcon,
        LoaderIcon,
        RefreshCwIcon,
        SearchXIcon
    } from "@lucide/svelte";
    import {glacierCli} from "$lib/api/api";
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {formatBytes} from "$lib/api/byte-math";

    let {gameId}: { gameId: bigint } = $props();

    const libSrv = glacierCli(LibraryService);

    let base = $state("");
    let downloaded = $state(false);

    let fileRpc = createRPCRunner(() => libSrv.listFiles({
        GameId: gameId,
        BasePath: base,
        Downloaded: downloaded,
    }));

    // Trigger initial load and reacts to tab/path changes
    $effect(() => {
        fileRpc.runner();
    });

    function navigateTo(dirName: string) {
        base = base === "" ? dirName : `${base}/${dirName}`;
    }

    function navigateUp() {
        const parts = base.split('/');
        parts.pop();
        base = parts.join('/');
    }

    const breadcrumbs = $derived(base === "" ? [] : base.split('/'));
</script>

<div class="flex flex-col h-full overflow-hidden">
    <!-- Compact Toolbar -->
    <div class="flex items-center justify-between p-2 border-b border-border bg-panel/50">
        <div class="flex gap-1 bg-surface p-0.5 rounded-lg border border-border">
            <button
                    onclick={() => downloaded = false}
                    class="px-3 py-1 rounded-md text-[10px] font-bold transition-all {!downloaded ? 'bg-panel shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
            >
                Incomplete
            </button>
            <button
                    onclick={() => downloaded = true}
                    class="px-3 py-1 rounded-md text-[10px] font-bold transition-all {downloaded ? 'bg-panel shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
            >
                Complete
            </button>
        </div>

        <button onclick={() => fileRpc.runner()} class="p-1.5 text-muted hover:text-frost-400 transition-colors">
            <RefreshCwIcon size={14} class={fileRpc.loading ? 'animate-spin' : ''}/>
        </button>
    </div>

    <!-- Mini Breadcrumbs -->
    <div class="flex items-center gap-1 px-3 py-1.5 bg-panel/20 border-b border-border text-[10px] overflow-x-auto no-scrollbar whitespace-nowrap">
        <button onclick={() => base = ""} class="text-muted hover:text-frost-400">
            <HardDriveIcon size={12}/>
        </button>
        {#if base !== ""}<span class="text-muted/30">/</span>{/if}
        {#each breadcrumbs as part, i}
            <button onclick={() => base = breadcrumbs.slice(0, i + 1).join('/')}
                    class="hover:text-foreground text-muted font-bold">{part}</button>
            {#if i < breadcrumbs.length - 1}<span class="text-muted/30">/</span>{/if}
        {/each}
    </div>

    <!-- Scrollable File Area -->
    <div class="flex-1 overflow-y-auto custom-scrollbar bg-black/5">
        {#if fileRpc.loading && !fileRpc.value}
            <div class="flex items-center justify-center h-full gap-2 text-muted">
                <LoaderIcon class="animate-spin text-frost-500" size={14}/>
                <span class="text-[10px] font-bold uppercase tracking-widest">Loading...</span>
            </div>
        {:else if fileRpc.error}
            <div class="p-4 text-center text-red-400 text-[10px] font-bold uppercase">
                {fileRpc.error}
            </div>
        {:else}
            <div class="flex flex-col">
                {#if base !== ""}
                    <button onclick={navigateUp}
                            class="flex items-center gap-2 px-4 py-2 hover:bg-frost-500/5 text-frost-400 text-[11px] font-bold border-b border-border/30">
                        <CornerLeftUpIcon size={14}/>
                        ..
                    </button>
                {/if}

                {#if fileRpc.value?.files && fileRpc.value.files.length > 0}
                    {#each fileRpc.value.files as file}
                        <button
                                onclick={() => file.IsDir ? navigateTo(file.RelPath.split('/').pop()!) : null}
                                class="flex items-center justify-between px-4 py-2 border-b border-border/20 last:border-0 hover:bg-panel/50 transition-colors text-left group"
                        >
                            <div class="flex items-center gap-2 min-w-0">
                                <span class="shrink-0 text-muted group-hover:text-frost-400">
                                    {#if file.IsDir}<FolderIcon size={14}/>{:else}<FileIcon size={14}/>{/if}
                                </span>
                                <span class="text-[11px] truncate {file.IsDir ? 'font-bold text-foreground' : 'text-muted'}">
                                    {file.RelPath.split('/').pop()}
                                </span>
                            </div>
                            <span class="text-[9px] font-mono text-muted/40 shrink-0">{formatBytes(file.Size)}</span>
                        </button>
                    {/each}
                {:else}
                    <div class="flex flex-col items-center justify-center py-8 text-muted/20 gap-1">
                        <SearchXIcon size={24}/>
                        <span class="text-[10px] font-bold uppercase">Empty</span>
                    </div>
                {/if}
            </div>
        {/if}
    </div>
</div>

<style>
    .no-scrollbar::-webkit-scrollbar {
        display: none;
    }

    .custom-scrollbar::-webkit-scrollbar {
        width: 4px;
    }

    .custom-scrollbar::-webkit-scrollbar-track {
        background: transparent;
    }

    .custom-scrollbar::-webkit-scrollbar-thumb {
        background: rgba(255, 255, 255, 0.05);
        border-radius: 10px;
    }
</style>
