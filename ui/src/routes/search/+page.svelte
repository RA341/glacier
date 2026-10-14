<script lang="ts">
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { type GameMetadata, SearchService } from '#lib/gen/search/v1/search_pb.js';
	import { glacierCli } from '#lib/api/api.js';
	import { CircleAlert, Search, LoaderIcon, SearchIcon, SearchXIcon } from '@lucide/svelte';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { transferStore } from './selectedGame.svelte';
	import { ServiceConfigService } from '#lib/gen/service_config/v1/service_config_pb.js';
	import GameGridItem from './GameGridItem.svelte';
	import NoMetadataProviders from './NoMetadataProviders.svelte';

	const searchQueryParam = 'q';
	let query = $state(page.url.searchParams.get(searchQueryParam) || '');

	let srvConfig = glacierCli(ServiceConfigService);
	let metadataRpc = createRPCRunner(() => srvConfig.getActiveService({ serviceType: 'Metadata' }));

	let providersLoaded = $state(false);
	let noMetadataProviders = $derived(
		providersLoaded &&
			!metadataRpc.loading &&
			!metadataRpc.error &&
			!metadataRpc.value?.names?.length
	);

	async function loadProviders() {
		await metadataRpc.runner();
		providersLoaded = true;
	}

	onMount(() => {
		loadProviders();
		if (query && selectedMetadataProvider) handleSearch();
	});

	const metadataQueryParam = 'metadata';
	let selectedMetadataProvider = $state(page.url.searchParams.get(metadataQueryParam) || '');

	function handleSearch() {
		const url = new URL(page.url.href);

		if (query) url.searchParams.set(searchQueryParam, query);
		else url.searchParams.delete(searchQueryParam);

		if (selectedMetadataProvider)
			url.searchParams.set(metadataQueryParam, selectedMetadataProvider);
		else url.searchParams.delete(metadataQueryParam);

		goto(url, { reset: false });
		searchRpc.runner();
	}

	$effect(() => {
		if (metadataRpc.value && metadataRpc.value?.names.length !== 0 && !selectedMetadataProvider) {
			selectedMetadataProvider = metadataRpc.value!.names[0].Name;
			if (query) handleSearch();
		}
	});

	const searchClient = glacierCli(SearchService);
	const searchRpc = createRPCRunner(() =>
		searchClient.searchMetadata({
			q: {
				query: query,
				indexer: selectedMetadataProvider
			}
		})
	);

	let isModalOpen = $state(false);

	function handleAdd(meta: GameMetadata) {
		transferStore.data = meta;

		goto('/search/details', { reset: false });
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter' && query.trim() !== '') handleSearch();
	}
</script>

<div class="console-page space-y-8">
	<header class="flex flex-wrap items-center justify-between gap-4 border-b border-border pb-6">
		<div class="flex w-full flex-wrap items-center gap-4">
			<div class="rounded-2xl border border-frost-500/20 bg-frost-500/10 p-3 text-frost-400">
				<Search size={28} />
			</div>
			<div>
				<h1 class="headline-xl">Search</h1>
				<p class="text-sm text-muted">Search and add New games</p>
			</div>

			{#if !noMetadataProviders}
				<div class="flex min-w-0 flex-1 flex-wrap items-center gap-3">
					<div class="relative min-w-40">
						<select
							aria-label="Metadata provider"
							bind:value={selectedMetadataProvider}
							onchange={handleSearch}
							class="w-full cursor-pointer appearance-none rounded-xl border border-border bg-panel px-4 py-2.5 pr-10 text-sm font-bold transition-all outline-none focus:border-frost-500"
						>
							{#if metadataRpc.loading}
								<option disabled>Loading indexers...</option>
							{:else if metadataRpc.value?.names}
								{#each metadataRpc.value.names as indexer}
									<option value={indexer.Name}>{indexer.Name}</option>
								{/each}
							{/if}
						</select>
					</div>

					<!-- Query Input + Search Button -->
					<div class="group relative flex-1">
						<input
							type="search"
							aria-label="Search games"
							bind:value={query}
							onkeydown={handleKeyDown}
							placeholder="Search for games..."
							class="w-full rounded-full border border-border bg-panel px-4 py-2.5 pr-12 text-sm transition-all outline-none focus:border-frost-500"
						/>
						<button
							onclick={handleSearch}
							aria-label="Search"
							disabled={searchRpc.loading || !query.trim()}
							class="absolute top-1.5 right-1.5 bottom-1.5 flex w-10 items-center justify-center rounded-full border border-border bg-surface text-muted transition-all hover:border-frost-500/50 hover:text-frost-500 disabled:opacity-50"
						>
							{#if searchRpc.loading}
								<LoaderIcon size={18} class="animate-spin text-frost-500" />
							{:else}
								<SearchIcon size={18} />
							{/if}
						</button>
					</div>
				</div>
			{/if}
		</div>
	</header>

	<!-- Results Area -->
	<main class="min-w-0 flex-1">
		<div class="w-full">
			{#if noMetadataProviders}
				<NoMetadataProviders />
			{:else if searchRpc.error}
				<aside
					class="flex items-start gap-4 rounded-xl border border-red-500/20 bg-red-500/5 p-4 text-red-400"
				>
					<CircleAlert size={20} />
					<div>
						<h3 class="font-bold">Search Failed</h3>
						<p class="text-sm opacity-80">{searchRpc.error}</p>
					</div>
				</aside>
			{:else if searchRpc.loading}
				<div class="flex h-96 flex-col items-center justify-center gap-4 text-muted">
					<LoaderIcon class="size-10 animate-spin text-frost-500" />
					<p class="animate-pulse font-medium">Indexing results...</p>
				</div>
			{:else if searchRpc.value && searchRpc.value.metadata && searchRpc.value.metadata.length > 0}
				<div class="poster-grid">
					{#each searchRpc.value.metadata as item (item.ID)}
						<GameGridItem game={item} />
					{/each}
				</div>
			{:else if searchRpc.value}
				<div class="flex h-96 flex-col items-center justify-center gap-2 text-muted opacity-40">
					<SearchXIcon size={64} strokeWidth={1} />
					<p class="text-xl font-medium">No results found</p>
					<p class="text-sm">Try adjusting your query or indexer</p>
				</div>
			{:else}
				<div
					class="flex h-96 flex-col items-center justify-center rounded-3xl border-2 border-dashed border-border text-blue-100"
				>
					<SearchIcon size={80} strokeWidth={1} class="mb-4" />
					<p class="text-lg font-medium">Search for game</p>
				</div>
			{/if}
		</div>
	</main>
</div>

<style>
	main::-webkit-scrollbar {
		width: 8px;
	}

	main::-webkit-scrollbar-track {
		background: transparent;
	}

	main::-webkit-scrollbar-thumb {
		background: var(--color-border);
		border-radius: 10px;
	}
</style>
