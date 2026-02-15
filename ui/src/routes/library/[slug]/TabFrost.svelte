<script lang="ts">
    import {
        RefreshCwIcon, Trash2Icon, AlertTriangleIcon,
        ChevronLeftIcon, ShieldCheckIcon, LoaderIcon
    } from "@lucide/svelte";
    import {fade, fly} from "svelte/transition";
    import {frostCli, callRPC} from "$lib/api/api";
    import {FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import type {Game} from "$lib/gen/library/v1/library_pb";

    let {game}: { game?: Game } = $props();

    const frostLib = frostCli(FrostLibraryService);
    const sm = getSnackbarCtx();

    let isProcessing = $state(false);
    let showDeleteConfirm = $state(false);

    async function handleRecheck() {
        if (!game!.ID || isProcessing) return;

        isProcessing = true;
        const {err} = await callRPC(() => frostLib.download({
            gameId: BigInt(game!.ID),
            recheck: true,
        }));

        if (err) {
            sm.push(`Recheck failed: ${err}`, 'error');
        } else {
            sm.push("Recheck process initiated", 'success');
        }
        isProcessing = false;
    }

    async function handleDeleteRedownload() {
        if (!game!.ID || isProcessing) return;

        isProcessing = true;
        const {err} = await callRPC(() => frostLib.download({
            gameId: BigInt(game!.ID),
            recheck: true,
            force: true,
        }));

        if (err) {
            sm.push(`Operation failed: ${err}`, 'error');
        } else {
            sm.push("Local files purged. Redownload started.", 'success');
        }

        showDeleteConfirm = false;
        isProcessing = false;
    }
</script>

<div class="space-y-8" in:fade>
    <!-- Header Section -->
    <header class="flex flex-col gap-2 px-2">
        <div class="flex items-center gap-4">
            <div class="p-3 bg-panel border border-border rounded-2xl text-frost-400">
                <ShieldCheckIcon size={24}/>
            </div>
            <div>
                <h1 class="text-2xl font-bold tracking-tight text-foreground">
                    Manage local frost files
                </h1>
                <p class="text-sm text-muted font-medium">
                    {game?.Meta?.Name || 'Unknown Title'}
                </p>
            </div>
        </div>
    </header>

    <!-- Operations Card -->
    <main class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Recheck Card -->
        <div class="bg-surface border border-border rounded-3xl p-6 flex flex-col justify-between gap-6 shadow-sm">
            <div class="space-y-2">
                <h3 class="font-bold text-foreground flex items-center gap-2">
                    <RefreshCwIcon size={18} class="text-frost-400"/>
                    Verify Integrity
                </h3>
                <p class="text-xs text-muted leading-relaxed">
                    Scans local files against the manifest to ensure no data is corrupted or missing.
                    This will pause any active downloads for this title.
                </p>
            </div>

            <button
                    onclick={handleRecheck}
                    disabled={isProcessing}
                    class="w-full py-3 bg-panel border border-border rounded-xl text-sm font-bold text-foreground hover:border-frost-500/50 hover:bg-frost-500/5 transition-all flex items-center justify-center gap-2 disabled:opacity-50"
            >
                {#if isProcessing}
                    <LoaderIcon size={16} class="animate-spin"/>
                {:else}
                    <RefreshCwIcon size={16}/>
                {/if}
                Recheck Files
            </button>
        </div>

        <!-- Delete/Redownload Card -->
        <div class="bg-surface border border-red-500/10 rounded-3xl p-6 flex flex-col justify-between gap-6 shadow-sm">
            <div class="space-y-2">
                <h3 class="font-bold text-red-400 flex items-center gap-2">
                    <Trash2Icon size={18}/>
                    Danger Zone
                </h3>
                <p class="text-xs text-muted leading-relaxed">
                    Completely removes all local data for this game and re-queues it for a fresh download.
                    Use this if a recheck fails to fix issues.
                </p>
            </div>

            <button
                    onclick={() => showDeleteConfirm = true}
                    disabled={isProcessing}
                    class="w-full py-3 bg-red-500/10 border border-red-500/20 rounded-xl text-sm font-bold text-red-400 hover:bg-red-500 hover:text-white transition-all flex items-center justify-center gap-2 disabled:opacity-50"
            >
                <Trash2Icon size={16}/>
                Delete and Redownload
            </button>
        </div>
    </main>
</div>

<!-- DELETE CONFIRMATION MODAL -->
{#if showDeleteConfirm}
    <div class="fixed inset-0 z-[200] flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md"
             onclick={() => !isProcessing && (showDeleteConfirm = false)}></div>

        <div
                class="relative w-full max-w-sm bg-surface border border-red-500/20 rounded-[2rem] shadow-2xl p-8 space-y-6"
                transition:fly={{ y: 20, duration: 300 }}
        >
            <div class="flex flex-col items-center text-center gap-4">
                <div class="p-4 bg-red-500/10 text-red-400 rounded-full border border-red-500/20">
                    <AlertTriangleIcon size={40} strokeWidth={1.5}/>
                </div>
                <div class="space-y-2">
                    <h2 class="text-xl font-bold text-foreground">Are you sure?</h2>
                    <p class="text-xs text-muted leading-relaxed">
                        This will permanently delete all local files for <span
                            class="text-foreground font-bold">{game?.Meta?.Name}</span>.
                        A new download will start immediately after.
                    </p>
                </div>
            </div>

            <div class="flex flex-col gap-2">
                <button
                        onclick={handleDeleteRedownload}
                        disabled={isProcessing}
                        class="w-full py-3 bg-red-500 text-white rounded-xl text-sm font-bold hover:bg-red-400 transition-all flex items-center justify-center gap-2"
                >
                    {#if isProcessing}
                        <LoaderIcon size={16} class="animate-spin"/>
                    {:else}
                        Purge and Redownload
                    {/if}
                </button>
                <button
                        onclick={() => showDeleteConfirm = false}
                        disabled={isProcessing}
                        class="w-full py-3 bg-panel border border-border rounded-xl text-sm font-bold text-muted hover:text-foreground transition-all"
                >
                    Cancel
                </button>
            </div>
        </div>
    </div>
{/if}
