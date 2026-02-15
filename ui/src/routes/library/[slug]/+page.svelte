<script lang="ts">
    import GameHero from "./GameHero.svelte";
    import {CircleAlert, LoaderIcon, Pen, ServerIcon, Trash2} from '@lucide/svelte';
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {callRPC, glacierCli, isFrost} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/rpc.svelte.js";
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import {goto} from "$app/navigation";
    import TabFrost from "./TabFrost.svelte";
    import GameDownloadButton from "./ButtonDownload.svelte";
    import TabFrontPage from "./TabFrontPage.svelte";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {getUserCtx} from "$lib/components/user/provider.svelte";

    let activeTab = $derived(page.url.searchParams.get('tab') || 'details');

    let gameIdStr = page.params.slug;
    let gameId = $derived(BigInt(gameIdStr!));

    function setTab(tab: string) {
        const newUrl = new URL(page.url);
        newUrl.searchParams.set('tab', tab);
        goto(newUrl.toString(), {replaceState: true, noScroll: true});
    }

    let libSrv = glacierCli(LibraryService)
    let gameRpc = createRPCRunner(() => libSrv.getGame({gameId: gameId}))

    function getGame() {
        if (gameIdStr) {
            gameRpc.runner()
        }
    }

    let originalGame = $derived(gameRpc.value?.game)

    onMount(() => {
        getGame()
    })

    const sm = getSnackbarCtx()

    async function deleteGame() {
        const {err} = await callRPC(() => libSrv.delete({gameId: BigInt(gameIdStr!)}))
        if (err) {
            sm.push(`Error deleting: ${err}`, 'error')
            return
        }
        await goto("/library")
    }

    const user = getUserCtx()

    async function goToEdit() {
        await goto("manage")
    }
</script>

{#if !gameIdStr}
    <div class="flex flex-col items-center justify-center h-96 border-2 border-dashed border-border rounded-3xl text-muted/30">
        <ServerIcon size={48} strokeWidth={1} class="mb-4"/>
        <h2 class="text-xl font-bold text-foreground/50">No game Id found</h2>
    </div>
{:else if gameRpc.error}
    <div class="flex flex-col items-center justify-center h-96 text-red-400 gap-3">
        <CircleAlert size={48} strokeWidth={1}/>
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
        {#if originalGame}
            <GameHero bind:game={originalGame}/>
        {/if}

        <div class="flex items-center justify-between border-y border-border py-4 px-2">
            <div class="flex gap-4">
                <div class="flex gap-1 bg-panel p-1 rounded-xl w-fit">
                    <button
                            onclick={() => setTab('details')}
                            class="px-6 py-1.5 rounded-lg text-sm font-bold transition-all {activeTab === 'details' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
                    >
                        Details
                    </button>

                    {#if isFrost}
                        <button
                                onclick={() => setTab('local')}
                                class="px-6 py-1.5 rounded-lg text-sm font-bold transition-all {activeTab === 'local' ? 'bg-surface shadow-sm text-frost-400' : 'text-muted hover:text-foreground'}"
                        >
                            Local
                        </button>
                    {/if}
                </div>
            </div>

            {#if user.isOmni}
                <div class="flex gap-3">
                    <button
                            onclick={goToEdit}
                            class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:border-frost-500 transition-all flex items-center gap-2">
                        <Pen size={16}/>
                        Manage
                    </button>


                    <button onclick={() => deleteGame()}
                            class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:border-frost-500 transition-all flex items-center gap-2">
                        <Trash2 size={16}/>
                        Delete
                    </button>
                </div>
            {/if}

            {#if isFrost}
                <div class="flex gap-3">
                    <GameDownloadButton gameId={gameId}/>
                </div>
            {/if}

        </div>

        <main>
            {#if activeTab === 'details' && originalGame}
                <TabFrontPage bind:game={originalGame}/>
            {:else if activeTab === 'local'}
                <TabFrost game={originalGame}/>
            {/if}
        </main>
    </div>
{/if}
