<script lang="ts">
	import GameHero from './GameHero.svelte';
	import { CircleAlert, LoaderIcon, ServerIcon } from '@lucide/svelte';
	import { LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { glacierCli, isFrost } from '#lib/api/api.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import TabFrost from './TabFrost.svelte';
	import TabFrontPage from './TabFrontPage.svelte';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';
	import { getAssetPath } from '#lib/api/assets.js';

	let activeTab = $derived(page.url.searchParams.get('tab') || 'details');

	let gameIdStr = page.params.slug;
	let gameId = $derived(BigInt(gameIdStr!));

	function setTab(tab: string) {
		const newUrl = new URL(page.url.href);
		newUrl.searchParams.set('tab', tab);
		goto(newUrl.toString(), { replaceState: true, reset: false });
	}

	let libSrv = glacierCli(LibraryService);
	let gameRpc = createRPCRunner(() => libSrv.getGame({ gameId: gameId }));

	function getGame() {
		if (gameIdStr) {
			gameRpc.runner();
		}
	}

	let originalGame = $derived(gameRpc.value?.game);

	onMount(() => {
		getGame();
	});

	async function goToEdit() {
		await goto(`/library/${gameIdStr}/manage`);
	}
</script>

<main class="mx-auto min-h-screen max-w-7xl p-6 md:p-10">
	{#if !gameIdStr}
		<div
			class="flex h-96 flex-col items-center justify-center rounded-3xl border-2 border-dashed border-border text-muted/30"
		>
			<ServerIcon size={48} strokeWidth={1} class="mb-4" />
			<h2 class="text-xl font-bold text-foreground/50">No game Id found</h2>
		</div>
	{:else if gameRpc.error}
		<div class="flex h-96 flex-col items-center justify-center gap-3 text-red-400">
			<CircleAlert size={48} strokeWidth={1} />
			<h3 class="font-bold">The game you are looking for does not exist</h3>
			<p class="text-xs opacity-80">{gameRpc.error}</p>
		</div>
	{:else if gameRpc.loading}
		<div class="flex h-96 flex-col items-center justify-center gap-4 text-muted">
			<LoaderIcon class="animate-spin text-frost-500" size={40} />
			<p class="animate-pulse text-sm font-medium">Fetching game...</p>
		</div>
	{:else}
		<div class="max-w-8xl mx-auto space-y-8 bg-background px-5 py-3 text-foreground">
			{#if originalGame}
				<GameHero bind:game={originalGame} onManage={goToEdit} />

				<div class="flex items-center justify-between border-t border-border py-4">
					<div class="flex w-fit gap-1 rounded-2xl border border-border bg-panel p-1">
						<button
							onclick={() => setTab('details')}
							class="rounded-full px-8 py-2 text-sm font-black tracking-widest uppercase transition-all {activeTab ===
							'details'
								? 'bg-surface text-frost-400 shadow-md'
								: 'text-muted hover:text-foreground'}"
						>
							Overview
						</button>
						{#if isFrost}
							<button
								onclick={() => setTab('local')}
								class="rounded-full px-8 py-2 text-sm font-black tracking-widest uppercase transition-all {activeTab ===
								'local'
									? 'bg-surface text-frost-400 shadow-md'
									: 'text-muted hover:text-foreground'}"
							>
								Local
							</button>
						{/if}
					</div>
				</div>

				<main>
					<!-- Main Tab Content -->
					{#if activeTab === 'details'}
						<TabFrontPage bind:game={originalGame} />
					{:else if activeTab === 'local'}
						<TabFrost game={originalGame} />
					{/if}
				</main>
			{/if}
		</div>
	{/if}
</main>
