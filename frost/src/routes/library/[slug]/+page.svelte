<script lang="ts">
    import GameHero from "./GameHero.svelte";
    import FrontPageTab from "./FrontPageTab.svelte";
    import ManageTab from "./ManageTab.svelte";
    import {DownloadIcon, RefreshCwIcon} from '@lucide/svelte';
    import {type Game, LibraryService} from "$lib/gen/library/v1/library_pb";
    import {cli} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/svelte-api.svelte";
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import {toJson} from "@bufbuild/protobuf";
    import {GameSchema} from "$lib/gen/library/v1/library_pb";

    let activeTab = $state('details');
    let gameId = page.params.slug
    let libSrv = cli(LibraryService)
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

    let isModified = $derived(
        editGame && originalGame &&
        toJson(GameSchema, editGame) !== toJson(GameSchema, originalGame)
    )

    onMount(() => {
        getGame()
    })

</script>

{#if !gameId}
    <div>No game id found</div>
{:else if gameRpc.error}
    <div>Error</div>
    <div>{gameRpc.error}</div>
{:else if gameRpc.loading}
    <div>Loading</div>
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
                <button class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:border-frost-500 transition-all flex items-center gap-2">
                    <RefreshCwIcon size={16}/>
                    Update
                </button>
                <button class="px-8 py-2 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all flex items-center gap-2 shadow-lg shadow-frost-500/20">
                    <DownloadIcon size={16}/>
                    Download
                </button>
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