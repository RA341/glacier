<script lang="ts">
    import GameHero from "./GameHero.svelte";
    import {CircleAlert, LoaderIcon, ServerIcon} from '@lucide/svelte';
    import {LibraryService} from "$lib/gen/library/v1/library_pb";
    import {glacierCli, isFrost} from "$lib/api/api";
    import {createRPCRunner} from "$lib/api/rpc.svelte.js";
    import {onMount} from "svelte";
    import {page} from "$app/state";
    import {goto} from "$app/navigation";
    import TabFrost from "./TabFrost.svelte";
    import TabFrontPage from "./TabFrontPage.svelte";
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
    <div class="max-w-8xl px-5 py-3 mx-auto space-y-8 bg-background text-foreground">
        {#if originalGame}
            <GameHero bind:game={originalGame} onManage={goToEdit}/>

            <div class="flex items-center justify-between border-t border-border py-4">
                <div class="flex gap-1 bg-panel p-1 rounded-2xl border border-border w-fit">
                    <button
                            onclick={() => setTab('details')}
                            class="px-8 py-2 rounded-xl text-sm font-black uppercase tracking-widest transition-all {activeTab === 'details' ? 'bg-surface shadow-md text-frost-400' : 'text-muted hover:text-foreground'}"
                    >
                        Overview
                    </button>
                    {#if isFrost}
                        <button
                                onclick={() => setTab('local')}
                                class="px-8 py-2 rounded-xl text-sm font-black uppercase tracking-widest transition-all {activeTab === 'local' ? 'bg-surface shadow-md text-frost-400' : 'text-muted hover:text-foreground'}"
                        >
                            Local
                        </button>
                    {/if}
                </div>
            </div>

            <main>
                <!-- Main Tab Content -->
                {#if activeTab === 'details'}
                    <TabFrontPage bind:game={originalGame}/>
                {:else if activeTab === 'local'}
                    <TabFrost game={originalGame}/>
                {/if}
            </main>
        {/if}
    </div>
{/if}
