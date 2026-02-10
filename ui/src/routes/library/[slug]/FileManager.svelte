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
        SearchXIcon,
        Trash2Icon,        // Added
        AlertTriangleIcon  // Added
    } from "@lucide/svelte";
    import {fade, fly} from 'svelte/transition'; // Added
    import {glacierCli} from "$lib/api/api";
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {formatBytes} from "$lib/api/byte-math";
    import {callRPC} from "$lib/api/api";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    let {gameId}: { gameId: bigint } = $props();

    const libSrv = glacierCli(LibraryService);

    let base = $state("");
    let downloaded = $state(false);
    let pendingDeletePath = $state<string | null>(null); // Track file to delete

    let fileRpc = createRPCRunner(() => libSrv.listFiles({
        GameId: gameId,
        BasePath: base,
        Downloaded: downloaded,
    }));

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
    const sm = getSnackbarCtx();

    async function confirmDelete() {
        if (!pendingDeletePath) return;

        const {err} = await callRPC(() => libSrv.deleteFile({
            Path: `${base}/${pendingDeletePath}`,
            Downloaded: downloaded,
            GameId: gameId,
        }));

        if (err) {
            sm.push(`Error deleting file: ${err}`, 'error');
        } else {
            sm.push("File deleted successfully", "success");
        }

        pendingDeletePath = null;
        await fileRpc.runner();
    }
</script>

<div class="flex flex-col h-full overflow-hidden relative">
    <!-- Compact Toolbar -->
    <div class="flex items-center justify-between p-2 border-b border-border bg-panel/50">
        <div class="flex gap-1 bg-surface p-0.5 rounded-lg border border-border">
            <button
                    onclick={() => downloaded = false}
                    class="px-3 py-1 rounded-md text-[10px] font-bold transition-all {!downloaded ? 'bg-panel shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
            >
                Complete
            </button>
            <button
                    onclick={() => downloaded = true}
                    class="px-3 py-1 rounded-md text-[10px] font-bold transition-all {downloaded ? 'bg-panel shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
            >
                Incomplete
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
                        <div class="flex items-center justify-between px-4 py-2 border-b border-border/20 last:border-0 hover:bg-panel/50 transition-colors group">
                            <button
                                    onclick={() => file.IsDir ? navigateTo(file.RelPath.split('/').pop()!) : null}
                                    class="flex items-center gap-2 min-w-0 flex-1 text-left"
                            >
                                <span class="shrink-0 text-muted group-hover:text-frost-400">
                                    {#if file.IsDir}<FolderIcon size={14}/>{:else}<FileIcon size={14}/>{/if}
                                </span>
                                <span class="text-[11px] truncate {file.IsDir ? 'font-bold text-foreground' : 'text-muted'}">
                                    {file.RelPath.split('/').pop()}
                                </span>
                            </button>

                            <span class="text-[9px] font-mono text-muted/40 shrink-0 mr-2">{formatBytes(file.Size)}</span>

                            <!-- Trash Icon Button -->
                            <button
                                    onclick={() => pendingDeletePath = file.RelPath}
                                    class="p-1.5 text-muted hover:text-red-400 hover:bg-red-400/10 rounded-lg transition-all opacity-0 group-hover:opacity-100"
                                    title="Delete File"
                            >
                                <Trash2Icon size={14}/>
                            </button>
                        </div>
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

    <!-- DELETE CONFIRMATION DIALOG -->
    {#if pendingDeletePath}
        <div class="absolute inset-0 z-50 flex items-center justify-center p-4" transition:fade={{ duration: 100 }}>
            <!-- Local Backdrop -->
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_static_element_interactions -->
            <div class="absolute inset-0 bg-background/60 backdrop-blur-sm"
                 onclick={() => pendingDeletePath = null}></div>

            <div
                    class="relative w-full max-w-xs bg-surface border border-red-500/20 rounded-2xl shadow-2xl p-5 space-y-4"
                    transition:fly={{ y: 10, duration: 200 }}
            >
                <div class="flex items-center gap-3 text-red-400">
                    <div class="p-2 bg-red-500/10 rounded-lg">
                        <AlertTriangleIcon size={18}/>
                    </div>
                    <span class="font-bold text-sm">Confirm Deletion</span>
                </div>

                <p class="text-xs text-muted leading-relaxed">
                    Are you sure you want to delete <span class="text-foreground font-mono">{pendingDeletePath}</span>?
                    This action cannot be undone.
                </p>

                <div class="flex gap-2 justify-end">
                    <button
                            onclick={() => pendingDeletePath = null}
                            class="px-3 py-1.5 rounded-lg text-[11px] font-bold bg-panel border border-border text-muted hover:text-foreground transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                            onclick={confirmDelete}
                            class="px-4 py-1.5 rounded-lg text-[11px] font-bold bg-red-500 text-white hover:bg-red-400 transition-colors shadow-lg shadow-red-500/20"
                    >
                        Delete
                    </button>
                </div>
            </div>
        </div>
    {/if}
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
