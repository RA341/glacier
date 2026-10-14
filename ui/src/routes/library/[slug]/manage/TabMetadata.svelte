<script lang="ts">
	import {
		GlobeIcon,
		LayoutGridIcon,
		LoaderIcon,
		RefreshCwIcon,
		StarIcon,
		TagIcon,
		TextAlignStart,
		TypeIcon
	} from '@lucide/svelte';
	import { fade } from 'svelte/transition';
	import AssetManager from './AssetManager.svelte';
	import type { Game } from '#lib/gen/library/v1/library_pb.js';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { SearchService } from '#lib/gen/search/v1/search_pb.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	let { game = $bindable(), refresh }: { game?: Game; refresh: () => Promise<void> } = $props();

	function handleArrayInput(e: Event, field: 'Platforms' | 'Genres') {
		const value = (e.currentTarget as HTMLInputElement).value;
		game!.Meta![field] = value
			.split(',')
			.map((s) => s.trim())
			.filter(Boolean);
	}

	const sm = getSnackbarCtx();

	const searchSrv = glacierCli(SearchService);
	let isRefreshing = $state(false);

	async function refreshMeta() {
		isRefreshing = true;

		const { val, err } = await callRPC(() =>
			searchSrv.getGameMeta({
				provider: game!.Meta!.ProviderType,
				gameDbId: game!.Meta!.ID
			})
		);
		if (err || !val?.meta) {
			sm.push(`Could not find game metadata ${err}`, 'error');
		} else {
			game!.Meta = val.meta;
		}

		isRefreshing = false;
	}
</script>

<div class="space-y-6 pb-20">
	{#if game && game.Meta}
		<div class="space-y-6 rounded-3xl border border-border bg-surface p-8 shadow-sm">
			<div class="mb-2 flex items-center">
				<h2
					class="flex items-center gap-2 text-xs font-black tracking-[0.2em] text-frost-400 uppercase"
				>
					<LayoutGridIcon size={14} />
					Main Content
				</h2>

				<button
					onclick={refreshMeta}
					disabled={isRefreshing}
					class="btn-secondary flex items-center gap-2 rounded-2xl px-4 py-1.5 text-[10px] font-bold tracking-widest text-muted uppercase transition-all hover:border-frost-500/50 hover:text-frost-400 active:scale-95 disabled:opacity-50"
				>
					{#if isRefreshing}
						<LoaderIcon size={12} class="animate-spin text-frost-500" />
						Refreshing...
					{:else}
						<RefreshCwIcon size={12} />
						Refresh Metadata
					{/if}
				</button>
			</div>

			<div class="space-y-2">
				<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
					>Display Name</label
				>
				<div class="relative">
					<TypeIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
					<input
						type="text"
						bind:value={game.Meta.Name}
						class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm font-bold transition-all outline-none focus:border-frost-500"
					/>
				</div>
			</div>

			<div class="space-y-2">
				<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
					>Summary</label
				>
				<div class="relative">
					<TextAlignStart size={18} class="absolute top-4 left-4 text-muted/50" />
					<textarea
						bind:value={game.Meta.Summary}
						rows="2"
						class="w-full rounded-[0.75rem] border border-border bg-panel py-4 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
					></textarea>
				</div>
			</div>

			<div class="space-y-2">
				<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
					>Full Description</label
				>
				<div class="relative">
					<TextAlignStart size={18} class="absolute top-4 left-4 text-muted/50" />
					<textarea
						bind:value={game.Meta.Description}
						rows="6"
						class="w-full rounded-[0.75rem] border border-border bg-panel py-4 pr-4 pl-12 text-sm leading-relaxed transition-all outline-none focus:border-frost-500"
					></textarea>
				</div>
			</div>
		</div>

		<!-- 2. TECHNICAL SPECIFICATIONS GRID -->
		<div class="grid grid-cols-1 gap-6 md:grid-cols-3">
			<!-- Provider Details -->
			<div class="space-y-4 rounded-3xl border border-border bg-surface p-6">
				<h3
					class="flex items-center gap-2 text-[10px] font-bold tracking-widest text-muted uppercase"
				>
					<GlobeIcon size={12} />
					Source Info
				</h3>

				<div class="space-y-4">
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase">Provider</span>
						<input
							type="text"
							bind:value={game.Meta.ProviderType}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
					</div>
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase">External ID</span>
						<input
							type="text"
							bind:value={game.Meta.ID}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 font-mono text-xs outline-none focus:border-frost-500"
						/>
					</div>
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase">Original URL</span>
						<input
							type="text"
							bind:value={game.Meta.URL}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
					</div>
				</div>
			</div>

			<!-- Classifications -->
			<div class="space-y-4 rounded-3xl border border-border bg-surface p-6">
				<h3
					class="flex items-center gap-2 text-[10px] font-bold tracking-widest text-muted uppercase"
				>
					<TagIcon size={12} />
					Classifications
				</h3>

				<div class="space-y-4">
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase">Category</span>
						<input
							type="text"
							bind:value={game.Meta.Category}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
					</div>
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase"
							>Platforms (Comma separated)</span
						>
						<input
							type="text"
							value={game.Meta.Platforms?.join(', ')}
							oninput={(e) => handleArrayInput(e, 'Platforms')}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
					</div>
					<div class="space-y-1">
						<span class="text-[9px] font-medium text-muted uppercase">Genres (Comma separated)</span
						>
						<input
							type="text"
							value={game.Meta.Genres?.join(', ')}
							oninput={(e) => handleArrayInput(e, 'Genres')}
							class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
						/>
					</div>
				</div>
			</div>

			<!-- Release & Ratings -->
			<div class="space-y-4 rounded-3xl border border-border bg-surface p-6">
				<h3
					class="flex items-center gap-2 text-[10px] font-bold tracking-widest text-muted uppercase"
				>
					<StarIcon size={12} />
					Release & Rating
				</h3>

				<div class="space-y-4">
					<div class="grid grid-cols-2 gap-2">
						<div class="space-y-1">
							<span class="text-[9px] font-medium text-muted uppercase">Date</span>
							<input
								type="text"
								bind:value={game.Meta.ReleaseDate}
								class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
							/>
						</div>
						<div class="space-y-1">
							<span class="text-[9px] font-medium text-muted uppercase">Status</span>
							<input
								type="text"
								bind:value={game.Meta.ReleaseStatus}
								class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
							/>
						</div>
					</div>
					<div class="grid grid-cols-2 gap-2">
						<div class="space-y-1">
							<span class="text-[9px] font-medium text-muted uppercase">Score</span>
							<input
								type="text"
								bind:value={game.Meta.Rating}
								class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs font-bold text-frost-400 outline-none focus:border-frost-500"
							/>
						</div>
						<div class="space-y-1">
							<span class="text-[9px] font-medium text-muted uppercase">Votes</span>
							<input
								type="number"
								bind:value={game.Meta.RatingCount}
								class="w-full rounded-xl border border-border bg-panel px-3 py-2 text-xs outline-none focus:border-frost-500"
							/>
						</div>
					</div>
				</div>
			</div>
		</div>
		<!-- 3. ASSETS SECTION -->
		{#if game.ID}
			<AssetManager {refresh} gameId={game.ID} assets={game?.Meta?.assets} />
		{/if}
	{/if}
</div>

<style>
	textarea {
		resize: vertical;
		min-height: 80px;
	}
</style>
