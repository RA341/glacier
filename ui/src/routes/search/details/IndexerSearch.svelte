<script lang="ts">
	import {
		ChevronDownIcon,
		CircleAlert,
		HardDriveIcon,
		ImageIcon,
		LoaderIcon,
		SearchIcon
	} from '@lucide/svelte';
	import { glacierCli } from '#lib/api/api.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import {
		type GameMetadata,
		type GameSource,
		SearchService
	} from '#lib/gen/search/v1/search_pb.js';
	import { ServiceConfigService } from '#lib/gen/service_config/v1/service_config_pb.js';

	const searchService = glacierCli(SearchService);
	const srvConfig = glacierCli(ServiceConfigService);

	let {
		selectedGameSource = $bindable(),
		game
	}: {
		game: GameMetadata | null;
		selectedGameSource: GameSource | null;
	} = $props();

	function searchIndexer() {
		if (searchQuery.trim()) {
			searchIndexerRpc.runner();
		}
	}

	function handleMatch(item: GameSource) {
		if (selectedGameSource?.Title === item.Title) {
			selectedGameSource = null;
		} else {
			selectedGameSource = item;
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter') searchIndexer();
	}

	let searchQuery = $state('');
	$effect(() => {
		if (game) {
			searchQuery = game.Name;
		}
	});

	let searchIndexerRpc = createRPCRunner(() =>
		searchService.searchIndexers({
			q: {
				indexer: selectedIndexer,
				query: searchQuery
			}
		})
	);
	let availableIndexersRpc = createRPCRunner(() =>
		srvConfig.getActiveService({ serviceType: 'Indexer' })
	);
	let selectedIndexer = $state('');

	$effect(() => {
		availableIndexersRpc.runner();
	});
	$effect(() => {
		if (selectedIndexer === '') {
			selectedIndexer = availableIndexersRpc.value?.names.at(0)?.Name ?? '';
		}
	});
</script>

<div class="flex flex-col gap-6 rounded-3xl border border-border bg-surface p-6">
	<h2 class="flex items-center gap-2 text-lg font-bold">
		<SearchIcon size={20} class="text-frost-400" />
		Search Indexers
	</h2>

	<div class="flex gap-3">
		<div class="relative min-w-35">
			<select
				bind:value={selectedIndexer}
				class="w-full appearance-none rounded-xl border border-border bg-panel px-4 py-2.5 text-sm font-bold outline-none"
			>
				{#each availableIndexersRpc?.value?.names as ind}
					<option>{ind.Name}</option>
				{/each}
			</select>
			<ChevronDownIcon
				size={16}
				class="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-muted"
			/>
		</div>
		<div class="relative flex-1">
			<input
				type="text"
				placeholder="Search for release..."
				onkeydown={handleKeyDown}
				bind:value={searchQuery}
				class="w-full rounded-xl border border-border bg-panel px-4 py-2.5 text-sm outline-none focus:border-frost-500"
			/>
		</div>
		<button
			onclick={searchIndexer}
			disabled={searchIndexerRpc.loading || !searchQuery.trim()}
			class=" top-1.5 right-1.5 bottom-1.5 flex w-10 items-center justify-center rounded-full border border-border bg-surface text-muted transition-all hover:border-frost-500/50 hover:text-frost-500 disabled:opacity-50"
		>
			{#if searchIndexerRpc.loading}
				<LoaderIcon size={18} class="animate-spin text-frost-500" />
			{:else}
				<SearchIcon size={18} />
			{/if}
		</button>
	</div>

	<div class="h-100 space-y-3 overflow-y-auto pr-2">
		{#if searchIndexerRpc.error}
			<aside
				class="flex items-start gap-4 rounded-xl border border-red-500/20 bg-red-500/5 p-4 text-red-400"
			>
				<CircleAlert size={20} />
				<div>
					<h3 class="font-bold">Search Failed</h3>
					<p class="text-sm opacity-80">{searchIndexerRpc.error}</p>
				</div>
			</aside>
		{:else if searchIndexerRpc.loading}
			<div class="flex h-96 flex-col items-center justify-center gap-4 text-muted">
				<LoaderIcon class="size-10 animate-spin text-frost-500" />
				<p class="animate-pulse font-medium">Indexing results...</p>
			</div>
		{:else if searchIndexerRpc.value?.results}
			{#each searchIndexerRpc?.value?.results as item}
				<div
					class="flex items-center justify-between rounded-2xl border border-border bg-panel/30 p-4 transition-colors hover:border-muted"
				>
					<div class="flex items-center gap-4">
						<div
							class="flex h-16 w-12 items-center justify-center rounded-lg border border-border bg-panel text-muted"
						>
							<ImageIcon size={20} />
						</div>
						<div class="flex flex-col">
							<span class="text-sm font-bold">{item.Title}</span>
							<div class="mt-1 flex items-center gap-3 text-[10px] font-bold text-muted">
								<span class="flex items-center gap-1 uppercase">
									<HardDriveIcon size={10} />
									{item.FileSize}
								</span>
								<span>•</span>
								<span>{new Date(item.CreatedISO).getFullYear()}</span>
							</div>
						</div>
					</div>
					<button
						onclick={() => handleMatch(item)}
						class="btn-secondary rounded-2xl px-4 py-1.5 text-xs font-bold transition-colors hover:border-frost-500"
					>
						{selectedGameSource?.Title === item.Title ? 'Selected' : 'Select'}
					</button>
				</div>
			{/each}
		{/if}
	</div>
</div>
