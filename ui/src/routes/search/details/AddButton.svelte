<script lang="ts">
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { create } from '@bufbuild/protobuf';
	import { DownloadSchema, GameSchema, LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { type GameMetadata, type GameSource } from '#lib/gen/search/v1/search_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import {
		CircleAlertIcon,
		ArrowRightIcon,
		LibraryIcon,
		LoaderIcon,
		PlusIcon
	} from '@lucide/svelte';
	import { fade, scale } from 'svelte/transition';

	const sm = getSnackbarCtx();
	let libraryService = glacierCli(LibraryService);

	let {
		gameMetadata,
		gameSrc,
		downloadClient
	}: {
		gameMetadata: GameMetadata | null;
		gameSrc: GameSource | null;
		downloadClient: string;
	} = $props();

	const existsRpc = createRPCRunner(() =>
		libraryService.exists({
			MetadataGameId: gameMetadata!.ID,
			MetadataType: gameMetadata!.ProviderType
		})
	);

	onMount(() => {
		if (!gameMetadata) {
			sm.push('game metadata not found', 'warn');
			return;
		}
		existsRpc.runner();
	});

	let isAdding = $state(false);

	async function onAddGame() {
		if (!existsRpc.value?.gameId) {
			sm.push('game already in lib', 'warn');
		}

		if (!gameMetadata) {
			sm.push('No game metadata found', 'warn');
			return;
		}

		if (!gameSrc) {
			sm.push('Game from the indexer must be selected', 'warn');
			return;
		}

		isAdding = true;
		const { err } = await callRPC(() =>
			libraryService.add({
				game: create(GameSchema, {
					Meta: gameMetadata!,
					DownloadState: create(DownloadSchema, {
						DownloadUrl: gameSrc?.DownloadUrl,
						Client: downloadClient
					}),
					Source: gameSrc
				})
			})
		);

		if (err) {
			sm.push(`Unable to add game ${err}`, 'error');
		}

		isAdding = false;
	}

	async function goToGame() {
		await goto(`/library/${existsRpc.value!.gameId!}`, { replaceState: true });
	}
</script>

<div class="w-full">
	{#if existsRpc.error}
		<!-- ERROR STATE -->
		<div
			transition:fade
			class="flex items-center gap-3 rounded-2xl border border-red-500/20 bg-red-500/10 p-4 text-sm text-red-400"
		>
			<CircleAlertIcon size={18} />
			<span class="flex-1 truncate">Status Check Failed: {existsRpc.error}</span>
			<button onclick={() => existsRpc.runner()} class="text-[10px] font-bold uppercase underline"
				>Retry</button
			>
		</div>
	{:else if existsRpc.loading}
		<!-- LOADING STATE -->
		<div
			transition:fade
			class="flex w-full items-center justify-center gap-3 rounded-2xl border border-border bg-panel py-4 text-muted"
		>
			<LoaderIcon size={18} class="animate-spin text-frost-500" />
			<span class="text-sm font-bold tracking-widest uppercase">Checking Library...</span>
		</div>
	{:else}
		<!-- ACTION BUTTONS -->
		<div class="group relative overflow-hidden rounded-2xl">
			{#if existsRpc.value?.gameId}
				<!-- STATE: ALREADY IN LIBRARY -->
				<button
					transition:scale={{ start: 0.95, duration: 200 }}
					onclick={goToGame}
					class="flex w-full items-center justify-center gap-2 rounded-full border border-frost-500/30 bg-surface py-4 font-bold text-frost-400 shadow-lg transition-all hover:bg-frost-500/10 active:scale-[0.98]"
				>
					<LibraryIcon size={18} />
					<span>View in Library</span>
					<ArrowRightIcon
						size={16}
						class="opacity-50 transition-transform group-hover:translate-x-1"
					/>
				</button>
			{:else}
				<!-- STATE: NOT IN LIBRARY -->
				<button
					transition:scale={{ start: 0.95, duration: 200 }}
					onclick={onAddGame}
					disabled={!gameSrc || !downloadClient || isAdding}
					class="btn-frost flex w-full items-center justify-center gap-2 rounded-full py-4 font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-[0.98] disabled:opacity-50 disabled:grayscale"
				>
					{#if isAdding}
						<LoaderIcon size={18} class="animate-spin" />
						<span>Adding...</span>
					{:else if !gameSrc}
						<PlusIcon size={18} />
						<span>Select a game source</span>
					{:else if !downloadClient}
						<PlusIcon size={18} />
						<span>Select a download client</span>
					{:else}
						<PlusIcon size={18} />
						<span>Add to Library</span>
					{/if}
				</button>
			{/if}
		</div>
	{/if}
</div>

<style>
	button {
		user-select: none;
	}
</style>
