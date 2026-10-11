<script lang="ts">
	import {
		ChevronDownIcon,
		FolderOpenIcon,
		ImageIcon,
		LinkIcon,
		LoaderIcon,
		SearchIcon,
		SearchXIcon,
		TriangleAlert,
		XIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import {
		type GameMetadata,
		type GameSource,
		SearchService
	} from '#lib/gen/search/v1/search_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { GameSchema, LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { create } from '@bufbuild/protobuf';
	import { DownloaderService } from '#lib/gen/downloader/v1/downloader_pb.js';

	type Props = {
		isOpen?: boolean;
		localGame: GameSource | null;
	};

	let { isOpen = $bindable(false), localGame }: Props = $props();

	const searchClient = glacierCli(SearchService);

	let searchQuery = $state('');
	let matchedMetadata = $state<GameMetadata>();

	$effect(() => {
		if (isOpen && localGame?.name) {
			searchQuery = localGame.name;
			clientRpc.runner();
		} else {
			searchQuery = '';
		}
	});

	const rpc = createRPCRunner(() =>
		searchClient.searchMetadata({
			q: {
				query: searchQuery,
				indexer: ''
			}
		})
	);

	function search() {
		if (searchQuery.trim() !== '') {
			rpc.runner();
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			search();
		}
	}

	const downClientManger = glacierCli(DownloaderService);
	const clientRpc = createRPCRunner(() => downClientManger.getActiveClients({}));

	const libClient = glacierCli(LibraryService);

	let downloadClient = $state('');
	let downloadUrl = $derived(localGame?.downloadUrl ?? '');
	let downloadPath = $state('./games/');
	let gameType = $state('');

	$effect(() => {
		if (clientRpc?.value?.clients && clientRpc?.value?.clients?.length > 0 && !downloadClient) {
			downloadClient = clientRpc.value.clients[0].name;
		}
	});

	async function handleAdd() {
		if (!downloadUrl || !downloadPath) {
			console.log('empty url and path');
			return;
		}

		isOpen = false;

		const game = create(GameSchema, {
			DownloadState: {
				DownloadPath: downloadPath,
				DownloadUrl: downloadUrl,
				Client: downloadClient
			},
			Meta: matchedMetadata,
			GameType: gameType
		});

		const { err } = await callRPC(() => libClient.add({ game }));
		if (err) {
			console.error('failed to add game', err);
		}
	}

	function setMetadata(meta: GameMetadata, isMatched: boolean) {
		matchedMetadata = isMatched ? undefined : meta;
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<!-- Backdrop -->
		<div
			class="absolute inset-0 bg-background/90 backdrop-blur-md"
			onclick={() => (isOpen = false)}
		></div>

		<!-- Content Container -->
		<div
			class="relative flex h-212.5 w-full max-w-7xl flex-col overflow-hidden rounded-2xl border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<!-- HEADER -->
			<div class="flex items-center justify-between border-b border-border p-5">
				<h2 class="text-2xl font-bold tracking-tight text-foreground">Add to library</h2>
				<button onclick={() => (isOpen = false)} class="p-1 text-muted hover:text-foreground">
					<XIcon size={24} />
				</button>
			</div>

			<div class="flex flex-1 overflow-hidden">
				<div class="flex flex-[1.2] flex-col border-r border-border bg-background/30">
					<!-- Provider Toolbar -->
					<div class="flex items-center gap-3 border-b border-border bg-surface p-4">
						<!-- Provider Selector -->
						<div class="relative">
							<select
								class="min-w-[120px] cursor-pointer appearance-none rounded-xl border border-border bg-panel px-4 py-2.5 pr-10 text-sm font-bold tracking-widest uppercase transition-colors outline-none focus:border-frost-500"
							>
								<option>IGDB</option>
								<option>Steam</option>
								<option>TMDB</option>
							</select>
							<ChevronDownIcon
								size={16}
								class="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-muted"
							/>
						</div>

						<div class="group relative flex-1">
							<input
								type="text"
								bind:value={searchQuery}
								onkeydown={(e) => e.key === 'Enter' && rpc.runner()}
								placeholder="Search for metadata..."
								class="w-full rounded-xl border border-border bg-panel px-4 py-2.5 pr-14 text-sm transition-all outline-none placeholder:text-muted/50 focus:border-frost-500"
							/>

							<button
								onclick={search}
								onkeydown={handleKeyDown}
								disabled={rpc.loading}
								class="absolute top-1.5 right-1.5 bottom-1.5 flex w-10 items-center justify-center rounded-lg border border-border bg-frost-500 text-muted shadow-sm transition-all hover:border-frost-500/50 hover:text-frost-500 active:scale-95 disabled:opacity-50"
								title="Search"
							>
								{#if rpc.loading}
									<LoaderIcon size={18} class="animate-spin text-frost-500" />
								{:else}
									<SearchIcon size={18} />
								{/if}
							</button>
						</div>
					</div>

					<!-- Provider Results List -->
					<div class="flex-1 space-y-3 overflow-y-auto p-4">
						{#if rpc.loading}
							<!-- Loading State -->
							<div class="flex h-full flex-col items-center justify-center text-muted">
								<LoaderIcon class="mb-2 animate-spin" size={32} />
								<p class="animate-pulse text-sm font-medium">Fetching metadata...</p>
							</div>
						{:else if rpc.error}
							<!-- Error State -->
							<div
								class="flex items-start gap-3 rounded-xl border border-red-500/20 bg-red-500/5 p-4 text-red-400"
							>
								<TriangleAlert size={18} />
								<div class="text-sm">
									<p class="font-bold">Search failed</p>
									<p class="opacity-80">{rpc.error}</p>
								</div>
							</div>
						{:else if rpc.value && rpc.value.metadata.length > 0}
							<!-- Results State -->
							{#each rpc.value.metadata as item}
								{@const isMatched = matchedMetadata?.ID === item.ID}
								<div
									class="flex gap-4 rounded-2xl border p-4 transition-all {isMatched
										? 'border-frost-500 bg-frost-500/10 ring-1 ring-frost-500/20'
										: 'border-border bg-surface hover:bg-panel/50'}"
								>
									<!-- Poster -->
									<div
										class="h-28 w-20 shrink-0 overflow-hidden rounded-lg border border-border/50 bg-panel"
									>
										{#if item.ThumbnailURL}
											<img src={item.ThumbnailURL} alt="" class="h-full w-full object-cover" />
										{:else}
											<div class="flex h-full w-full items-center justify-center text-muted">
												<ImageIcon size={24} />
											</div>
										{/if}
									</div>

									<!-- Meta Info -->
									<div class="flex min-w-0 flex-1 flex-col">
										<div class="mb-1 flex items-center gap-2">
											<h3 class="truncate font-bold text-foreground">{item.Name}</h3>
											<span
												class="rounded-md border border-border bg-panel px-2 py-0.5 text-[10px] font-bold text-muted uppercase"
												>2023</span
											>
										</div>
										<p class="mb-4 line-clamp-3 text-xs leading-relaxed text-muted">
											{item.Description || 'No description available.'}
										</p>

										<div class="mt-auto flex justify-end">
											<button
												onclick={() => setMetadata(item, isMatched)}
												class="rounded-lg px-5 py-1.5 text-xs font-bold transition-all {isMatched
													? 'bg-frost-500 text-background'
													: 'border border-border bg-panel hover:border-frost-500'}"
											>
												{isMatched ? 'Matched' : 'Match'}
											</button>
										</div>
									</div>
								</div>
							{/each}
						{:else if searchQuery.trim() !== ''}
							<!-- EMPTY STATE: No Results Found (Query exists but no metadata) -->
							<div
								class="flex h-full flex-col items-center justify-center py-20 text-muted"
								in:fade
							>
								<div
									class="mb-4 flex h-16 w-16 items-center justify-center rounded-full border border-border bg-panel"
								>
									<SearchXIcon size={32} strokeWidth={1.5} />
								</div>
								<h3 class="font-bold text-foreground">No results found</h3>
								<p class="mt-1 px-8 text-center text-xs opacity-60">
									We couldn't find any games matching "{searchQuery}" on the selected provider.
								</p>
							</div>
						{:else}
							<!-- EMPTY STATE: No Query Found (Initial state / Search bar is empty) -->
							<div
								class="flex h-full flex-col items-center justify-center py-20 text-muted"
								in:fade
							>
								<div
									class="mb-4 flex h-16 w-16 items-center justify-center rounded-full border border-border bg-panel"
								>
									<SearchIcon size={32} strokeWidth={1.5} />
								</div>
								<h3 class="font-bold text-foreground">No query found</h3>
								<p class="mt-1 px-8 text-center text-xs opacity-60">
									Type a game title in the search bar above to start matching metadata.
								</p>
							</div>
						{/if}
					</div>
				</div>

				<!-- RIGHT SIDE: INDEXER DATA -->
				<div class="flex flex-1 flex-col overflow-y-auto bg-surface p-6">
					<div class="space-y-8">
						<!-- Block 1: Basic File Info -->
						<section class="space-y-4">
							<div class="rounded-xl border border-border bg-panel p-4">
								<p class="mb-2 text-[10px] font-bold tracking-[0.2em] text-muted uppercase">
									Indexer Title
								</p>
								<h3 class="truncate text-lg font-bold text-foreground">
									{localGame?.name || 'Waiting for selection...'}
								</h3>
							</div>

							<div class="grid grid-cols-3 gap-3">
								<div class="flex flex-col rounded-xl border border-border bg-panel/50 p-3">
									<span class="mb-1 text-[9px] font-bold text-muted uppercase">Size</span>
									<span class="text-sm font-medium">{localGame?.size || '--'}</span>
								</div>

								<div class="flex flex-col rounded-xl border border-border bg-panel/50 p-3">
									<span class="mb-1 text-[9px] font-bold text-muted uppercase">Date</span>
									<span class="text-sm font-medium">TODO</span>
								</div>

								<div
									class="group flex cursor-pointer items-center justify-between rounded-xl border border-border bg-panel/50 p-3 transition-colors hover:border-frost-500"
								>
									<div class="flex flex-col truncate">
										<span class="mb-1 text-[9px] font-bold text-muted uppercase">Download Path</span
										>
										<span class="truncate text-sm font-medium">{downloadPath}</span>
									</div>
									<FolderOpenIcon size={16} class="ml-2 text-muted group-hover:text-frost-400" />
								</div>
							</div>
						</section>

						<hr class="border-border/50" />

						<!-- Block 2: Download Config -->
						<section class="flex gap-4">
							<div class="min-w-40">
								<p class="mb-2 text-[10px] font-bold tracking-[0.2em] text-muted uppercase">
									Clients
								</p>
								<div class="relative">
									{#if clientRpc.loading}
										<!-- Styled as a placeholder to prevent layout shift -->
										<div
											class="flex w-full animate-pulse items-center gap-2 rounded-xl border border-border bg-panel/50 px-3 py-2.5 text-sm text-muted"
										>
											<LoaderIcon size={14} class="animate-spin" />
											<span>Loading...</span>
										</div>
									{:else if clientRpc.error}
										<!-- Error state with same dimensions -->
										<div
											class="w-full truncate rounded-xl border border-red-500/20 bg-red-500/5 px-3 py-2.5 text-xs text-red-400"
											title={clientRpc.error}
										>
											Error: {clientRpc.error}
										</div>
									{:else}
										<!-- Actual Select -->
										<select
											bind:value={downloadClient}
											class="w-full cursor-pointer appearance-none rounded-xl border border-border bg-panel px-3 py-2.5 text-sm font-bold transition-colors outline-none focus:border-frost-500"
										>
											{#if !clientRpc.value || clientRpc.value.clients.length === 0}
												<option value="" disabled selected>No clients found</option>
											{:else}
												{#each clientRpc.value.clients as client}
													<!-- Assuming client might be a string or object; adjust accordingly -->
													<option value={client.name}>{client.name}</option>
												{/each}
											{/if}
										</select>

										<!-- Chevron only shows when the select is active -->
										<ChevronDownIcon
											size={16}
											class="pointer-events-none absolute top-1/2 right-3 -translate-y-1/2 text-muted"
										/>
									{/if}
								</div>
							</div>
							<div class="flex-1">
								<p class="mb-2 text-[10px] font-bold tracking-[0.2em] text-muted uppercase">
									Download URL
								</p>
								<div class="relative">
									<LinkIcon size={16} class="absolute top-1/2 left-3 -translate-y-1/2 text-muted" />
									<input
										type="text"
										bind:value={downloadUrl}
										placeholder="magnet:?xt=..."
										class="w-full rounded-xl border border-border bg-panel py-2.5 pr-4 pl-10 text-sm outline-none focus:border-frost-500"
									/>
								</div>
							</div>
						</section>

						<hr class="border-border/50" />

						<!-- Block 3: Matched Meta Preview -->
						<section>
							<p class="mb-3 text-[10px] font-bold tracking-[0.2em] text-muted uppercase">
								Matched Meta
							</p>
							{#if matchedMetadata}
								<div
									class="flex gap-4 rounded-2xl border border-frost-500 bg-frost-500/5 p-4"
									in:fade
								>
									<div
										class="h-20 w-16 shrink-0 overflow-hidden rounded border border-border bg-panel"
									>
										<img
											src={matchedMetadata.ThumbnailURL}
											class="h-full w-full object-cover"
											alt=""
										/>
									</div>
									<div class="min-w-0 flex-1">
										<div class="mb-1 flex items-center gap-2">
											<h4 class="truncate text-sm font-bold">{matchedMetadata.Name}</h4>
											<span class="font-mono text-[9px] text-frost-400">2023</span>
										</div>
										<p class="line-clamp-2 text-[11px] text-muted">{matchedMetadata.Description}</p>
									</div>
									<button
										onclick={() => (matchedMetadata = undefined)}
										class="self-start p-1.5 text-muted transition-colors hover:text-red-400"
									>
										<XIcon size={16} />
									</button>
								</div>
							{:else}
								<div
									class="flex h-24 items-center justify-center rounded-2xl border-2 border-dashed border-border text-sm text-muted italic"
								>
									No matched metadata
								</div>
							{/if}
						</section>
					</div>
				</div>
			</div>

			<!-- FOOTER -->
			<div class="flex items-center justify-end gap-3 border-t border-border bg-panel/20 p-5">
				<button
					onclick={() => (isOpen = false)}
					class="rounded-xl border border-border px-8 py-2.5 text-sm font-bold transition-all hover:bg-panel"
				>
					Cancel
				</button>
				<button
					onclick={handleAdd}
					class={`rounded-xl px-10 py-2.5 ${downloadUrl ? 'bg-frost-500' : 'bg-frost-600'} text-sm
                            font-bold text-background
                            shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-400 active:scale-95`}
				>
					Add Game
				</button>
			</div>
		</div>
	</div>
{/if}
