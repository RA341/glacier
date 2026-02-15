<script lang="ts">
    import {callRPC, glacierCli} from "$lib/api/api";
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {page} from "$app/state";
    import {goto} from "$app/navigation"; // Added for URL updates
    import {onMount} from "svelte";
    import {
        ChevronLeftIcon, DatabaseIcon, DownloadIcon,
        FileTextIcon, LayoutGridIcon, LoaderIcon, SaveIcon
    } from "@lucide/svelte";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";

    import TabMetadata from "./TabMetadata.svelte";
    import TabSource from "./TabSource.svelte";
    import TabDownload from "./TabDownload.svelte";
    import TabFile from "./TabFile.svelte";

    const sm = getSnackbarCtx();
    const libSrv = glacierCli(LibraryService);

    let gameIdStr = page.params.slug;
    let gameId = $derived(BigInt(gameIdStr!));

    let activeTab = $derived(page.url.searchParams.get('tab') || "metadata");

    let isSaving = $state(false);
    let gameRpc = createRPCRunner(() => libSrv.getGame({gameId}));
    let editableGame = $state<Game | null>(null);

    onMount(async () => {
        await gameRpc.runner();
        if (gameRpc?.value?.game) {
            editableGame = gameRpc.value.game;
        }
    });

    function handleTabChange(id: string) {
        const url = new URL(page.url);
        url.searchParams.set('tab', id);
        goto(url, {replaceState: true, noScroll: true, keepFocus: true});
    }

    async function handleSave() {
        if (!editableGame) return;
        isSaving = true;
        const {err} = await callRPC(() => libSrv.edit({game: editableGame!}));
        if (err) {
            sm.push(`Failed to save: ${err}`, "error");
        } else {
            sm.push("Game configuration updated", "success");
            await gameRpc.runner();
        }
        isSaving = false;
    }

    const tabs = [
        {id: 'metadata', label: 'Metadata', icon: LayoutGridIcon},
        {id: 'source', label: 'Source', icon: DatabaseIcon},
        {id: 'download', label: 'Download', icon: DownloadIcon},
        {id: 'files', label: 'Files', icon: FileTextIcon},
    ];
</script>

<div class="mx-auto p-6 space-y-6">
    {#if gameRpc.loading && !editableGame}
        <div class="flex flex-col items-center justify-center h-[60vh] gap-4 text-muted">
            <LoaderIcon class="animate-spin text-frost-500" size={40}/>
            <p class="text-sm font-bold uppercase tracking-[0.3em] animate-pulse">Retrieving Game Data...</p>
        </div>
    {:else if gameRpc.error}
        <div class="p-8 bg-red-500/5 border border-red-500/20 rounded-4xl text-center space-y-4">
            <p class="text-red-400 font-bold">{gameRpc.error}</p>
            <button onclick={() => gameRpc.runner()}
                    class="px-6 py-2 bg-panel border border-border rounded-xl text-xs font-bold uppercase">
                Retry
            </button>
        </div>
    {:else if editableGame}
        <header class="sticky top-0 z-30 bg-background/80 backdrop-blur-md border-b border-border -mx-6 px-6 py-4 flex justify-between items-center">
            <div class="flex items-center gap-4">
                <a href=".." class="p-2 hover:bg-panel rounded-xl text-muted transition-colors">
                    <ChevronLeftIcon size={20}/>
                </a>
                <div>
                    <h1 class="text-xl font-bold text-foreground truncate max-w-md">
                        {editableGame.Meta?.Name || 'Manage Game'}
                    </h1>
                </div>
            </div>

            <button
                    onclick={handleSave}
                    disabled={isSaving}
                    class="flex items-center gap-2 px-8 py-2.5 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 active:scale-95 disabled:opacity-50"
            >
                {#if isSaving}
                    <LoaderIcon size={18} class="animate-spin"/>
                    Saving...
                {:else}
                    <SaveIcon size={18}/>
                    Save Changes
                {/if}
            </button>
        </header>

        <div class="flex gap-1 bg-panel p-1 rounded-2xl border border-border w-fit">
            {#each tabs as tab}
                <button
                        onclick={() => handleTabChange(tab.id)}
                        class="flex items-center gap-2 px-6 py-2 rounded-xl text-sm font-bold transition-all
                    {activeTab === tab.id ? 'bg-surface shadow-md text-frost-400' : 'text-muted hover:text-foreground'}"
                >
                    <tab.icon size={16}/>
                    {tab.label}
                </button>
            {/each}
        </div>

        <main class="min-h-125">
            {#if activeTab === 'metadata'}
                <TabMetadata refresh={gameRpc.runner} bind:game={editableGame}/>
            {:else if activeTab === 'source'}
                <TabSource bind:game={editableGame}/>
            {:else if activeTab === 'download'}
                <TabDownload bind:game={editableGame}/>
            {:else if activeTab === 'files'}
                <TabFile bind:game={editableGame}/>
            {/if}
        </main>
    {:else}
        <div class="p-20 text-center text-muted italic">Game not found.</div>
    {/if}
</div>
