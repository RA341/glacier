<script lang="ts">
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { type Game, LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		ChevronLeftIcon,
		DatabaseIcon,
		DownloadIcon,
		FileTextIcon,
		LayoutGridIcon,
		LoaderIcon,
		SaveIcon,
		Trash
	} from '@lucide/svelte';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	import TabMetadata from './TabMetadata.svelte';
	import TabSource from './TabSource.svelte';
	import TabDownload from './TabDownload.svelte';
	import TabFile from './TabFile.svelte';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';

	const sm = getSnackbarCtx();
	const libSrv = glacierCli(LibraryService);

	let gameIdStr = page.params.slug;
	let gameId = $derived(BigInt(gameIdStr!));

	let activeTab = $derived(page.url.searchParams.get('tab') || 'metadata');

	let isSaving = $state(false);
	let gameRpc = createRPCRunner(() => libSrv.getGame({ gameId }));
	let editableGame = $state<Game | null>(null);

	onMount(async () => {
		await gameRpc.runner();
		if (gameRpc?.value?.game) {
			editableGame = gameRpc.value.game;
		}
	});

	function handleTabChange(id: string) {
		const url = new URL(page.url.href);
		url.searchParams.set('tab', id);
		goto(url, { replaceState: true, reset: false });
	}

	async function handleSave() {
		if (!editableGame) return;
		isSaving = true;
		const { err } = await callRPC(() => libSrv.edit({ game: editableGame! }));
		if (err) {
			sm.push(`Failed to save: ${err}`, 'error');
		} else {
			sm.push('Game configuration updated', 'success');
			await gameRpc.runner();
		}
		isSaving = false;
	}

	const dm = getDialogCtx();

	async function deleteGame() {
		const ok = await dm.confirm(
			'Delete game',
			'This will permanently delete the game files and metadata from Glacier.',
			'error'
		);
		if (!ok) {
			return;
		}

		const { err } = await callRPC(() => libSrv.delete({ gameId: BigInt(gameIdStr!) }));
		if (err) {
			sm.push(`Error deleting: ${err}`, 'error');
			return;
		}
		await goto('/library');
	}

	const tabs = [
		{ id: 'metadata', label: 'Metadata', icon: LayoutGridIcon },
		{ id: 'source', label: 'Source', icon: DatabaseIcon },
		{ id: 'download', label: 'Download', icon: DownloadIcon },
		{ id: 'files', label: 'Files', icon: FileTextIcon }
	];
</script>

<div class="mx-auto space-y-6 p-6">
	{#if gameRpc.loading && !editableGame}
		<div class="flex h-[60vh] flex-col items-center justify-center gap-4 text-muted">
			<LoaderIcon class="animate-spin text-frost-500" size={40} />
			<p class="animate-pulse text-sm font-bold tracking-[0.3em] uppercase">
				Retrieving Game Data...
			</p>
		</div>
	{:else if gameRpc.error}
		<div class="space-y-4 rounded-4xl border border-red-500/20 bg-red-500/5 p-8 text-center">
			<p class="font-bold text-red-400">{gameRpc.error}</p>
			<button
				onclick={() => gameRpc.runner()}
				class="btn-secondary rounded-2xl px-6 py-2 text-xs font-bold uppercase"
			>
				Retry
			</button>
		</div>
	{:else if editableGame}
		<header
			class="sticky top-0 z-30 -mx-6 flex items-center justify-between border-b border-border bg-background/80 px-6 py-4 backdrop-blur-md"
		>
			<div class="flex items-center gap-4">
				<a href=".." class="rounded-xl p-2 text-muted transition-colors hover:bg-panel">
					<ChevronLeftIcon size={20} />
				</a>
				<div>
					<h1 class="max-w-md truncate text-xl font-bold text-foreground">
						{editableGame.Meta?.Name || 'Manage Game'}
					</h1>
				</div>
			</div>

			<div class="flex flex-row gap-4">
				<button
					onclick={handleSave}
					disabled={isSaving}
					class="btn-frost flex items-center gap-2 rounded-full px-8 py-2.5 text-sm font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-95 disabled:opacity-50"
				>
					{#if isSaving}
						<LoaderIcon size={18} class="animate-spin" />
						Saving...
					{:else}
						<SaveIcon size={18} />
						Save Changes
					{/if}
				</button>

				<button
					onclick={() => deleteGame()}
					class="flex items-center gap-2 rounded-full border border-border bg-red-600 px-6 py-2 text-sm font-bold transition-all hover:border-frost-500"
				>
					<Trash size={16} />
					Delete
				</button>
			</div>
		</header>

		<div class="flex w-fit gap-1 rounded-2xl border border-border bg-panel p-1">
			{#each tabs as tab}
				<button
					onclick={() => handleTabChange(tab.id)}
					class="flex items-center gap-2 rounded-full px-6 py-2 text-sm font-bold transition-all
                    {activeTab === tab.id
						? 'bg-surface text-frost-400 shadow-md'
						: 'text-muted hover:text-foreground'}"
				>
					<tab.icon size={16} />
					{tab.label}
				</button>
			{/each}
		</div>

		<!-- todo refactor game rpc runner and game binding to external state -->
		<main class="min-h-125">
			{#if activeTab === 'metadata'}
				<TabMetadata refresh={gameRpc.runner} bind:game={editableGame} />
			{:else if activeTab === 'source'}
				<TabSource bind:game={editableGame} />
			{:else if activeTab === 'download'}
				<TabDownload refresh={gameRpc.runner} bind:game={editableGame} />
			{:else if activeTab === 'files'}
				<TabFile bind:game={editableGame} />
			{/if}
		</main>
	{:else}
		<div class="p-20 text-center text-muted italic">Game not found.</div>
	{/if}
</div>
