<script lang="ts">
    import GameHero from "./GameHero.svelte";
    import FrontPageTab from "./FrontPageTab.svelte";
    import ManageTab from "./ManageTab.svelte";
    import {AlertCircleIcon, DownloadIcon, LoaderIcon, Save, ServerIcon, Trash2} from '@lucide/svelte';
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {frostCli, glacierCli, isFrost} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import {toJson} from "@bufbuild/protobuf";
    import {GameSchema} from "$lib/gen/library/v1/library_pb";
    import {FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {goto} from "$app/navigation";

    let activeTab = $state('details');
    let gameId = page.params.slug
    let libSrv = glacierCli(LibraryService)
    let gameRpc = createRPCRunner(() => libSrv.getGame({gameId: BigInt(gameId!)}))

    function getGame() {
        console.log(`Fetching game... ${gameId}`);
        if (gameId) {
            gameRpc.runner()
        }
    }

    let originalGame = $derived(gameRpc.value?.game)

    let editGame = $state<Game | null>(null)

    $effect(() => {
        if (originalGame) {
            editGame = originalGame
        }
    })

    const llervice = frostCli(FrostLibraryService)

    function download() {
        llervice.download({
            gameId: BigInt(gameId!),
            downloadFolder: "./downloads"
        })
    }

    let isModified = $derived(
        editGame && originalGame &&
        toJson(GameSchema, editGame) !== toJson(GameSchema, originalGame)
    )

    onMount(() => {
        getGame()
    })

    function deleteGame() {
        libSrv.delete({gameId: BigInt(gameId!)})
        goto("/library")
    }

</script>

{#if !gameId}
    <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-border rounded-3xl text-muted/30">
        <ServerIcon size={48} strokeWidth={1} class="mb-4"/>
        <h2 class="text-xl font-bold text-foreground/50">No game Id found</h2>
    </div>
{:else if gameRpc.error}
    <div class="flex flex-col items-center justify-center h-96 text-red-400 gap-3">
        <AlertCircleIcon size={48} strokeWidth={1}/>
        <h3 class="font-bold">The game you are looking for does not exist</h3>
        <p class="text-xs opacity-80">{gameRpc.error}</p>
    </div>
{:else if gameRpc.loading}
    <div class="flex flex-col items-center justify-center h-96 text-muted gap-4">
        <LoaderIcon class="animate-spin text-frost-500" size={40}/>
        <p class="animate-pulse text-sm font-medium">Fetching game...</p>
    </div>
{:else}
    <div class="max-w-7xl mx-auto p-6 space-y-8 bg-background text-foreground">

        <GameHero bind:game={editGame}/>

        <div class="flex items-center justify-between border-y border-border py-4 px-2">
            <div class="flex gap-4">
                <div class="flex gap-1 bg-panel p-1 rounded-xl w-fit">
                    <button
                            onclick={() => activeTab = 'details'}
                            class="px-6 py-1.5 rounded-lg text-sm font-bold transition-all {activeTab === 'details' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
                    >
                        Details
                    </button>
                    <button
                            onclick={() => activeTab = 'manage'}
                            class="px-6 py-1.5 rounded-lg text-sm font-bold transition-all {activeTab === 'manage' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
                    >
                        Manage
                    </button>
                </div>
            </div>
            <div class="flex gap-3">
                <button
                        onclick={() => deleteGame()}
                        class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:border-frost-500 transition-all flex items-center gap-2"
                >
                    <Trash2 size={16}/>

                    Delete
                </button>
                <button
                        disabled={isModified}
                        onclick={() => isModified ? deleteGame(): {} }
                        class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:border-frost-500 transition-all flex items-center gap-2">
                    <Save size={16}/>
                    Save
                </button>
                {#if isFrost}
                    <button
                            onclick={download}
                            class="px-8 py-2 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all flex items-center gap-2 shadow-lg shadow-frost-500/20"
                    >
                        <DownloadIcon size={16}/>
                        Download
                    </button>
                {/if}
            </div>
        </div>

        <main>
            {#if activeTab === 'details'}
                <FrontPageTab bind:game={editGame}/>
            {:else}
                <ManageTab bind:game={editGame}/>
            {/if}
        </main>
    </div>
{/if}