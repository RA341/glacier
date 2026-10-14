<script lang="ts">
	import { glacierCli } from '#lib/api/api.js';
	import { LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { getAssetPath } from '#lib/api/assets.js';
	import { ImageIcon, LoaderIcon } from '@lucide/svelte';
	import AssetImg from '#lib/components/assets/AssetImg.svelte';

	const libSrv = glacierCli(LibraryService);

	let query = $state('');
	let offset = $state(0);
	let limit = $state(50);

	let libRpc = createRPCRunner(() =>
		libSrv.list({
			limit: limit,
			offset: offset,
			query: query
		})
	);

	function search() {
		libRpc.runner();
	}

	$effect(() => {
		const interval = setInterval(() => {
			search();
		}, 5000);
		return () => clearInterval(interval);
	});

	onMount(() => {
		search();
	});
</script>

<main class="console-page">
	<!-- Header Area -->
	<header class="mb-8 flex flex-col justify-between gap-6 md:flex-row md:items-center">
		<h1 class="headline-xl">Library</h1>

		<div class="group relative flex w-full max-w-2xl items-center">
			<input
				type="search"
				aria-label="Search library"
				bind:value={query}
				onkeydown={(e) => e.key === 'Enter' && search()}
				placeholder="Search games..."
				class="w-full rounded-full border border-border bg-surface px-4 py-3 pr-14 text-foreground transition-all placeholder:text-muted focus:ring-2 focus:ring-frost-500/50 focus:outline-none"
			/>
			<button
				onclick={search}
				class="absolute right-2 p-2 text-muted transition-colors hover:text-frost-400"
				aria-label="Search"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="11" cy="11" r="8" />
					<path d="m21 21-4.3-4.3" />
				</svg>
			</button>
		</div>
	</header>

	<!-- Content Area -->
	{#if libRpc.error}
		<div
			class="flex h-64 items-center justify-center rounded-xl border border-red-900/50 bg-red-950/20 text-red-400"
		>
			<p>Error: {libRpc.error}</p>
		</div>
	{:else if libRpc.loading && !libRpc.value && !libRpc.error}
		<div class="poster-grid">
			{#each Array(10) as _}
				<div class="h-80 animate-pulse overflow-hidden rounded-2xl border border-border bg-surface">
					<div class="h-3/4 bg-panel"></div>
					<div class="space-y-2 p-4">
						<div class="h-4 w-3/4 rounded bg-panel"></div>
						<div class="h-3 w-1/2 rounded bg-panel"></div>
					</div>
				</div>
			{/each}
		</div>
	{:else}
		<div class="poster-grid">
			{#each libRpc.value?.gameList ?? [] as ga}
				<a
					href={`/library/${ga.ID}`}
					class="poster-card group"
					aria-label={ga.Meta?.Name ?? 'Untitled game'}
				>
					<div
						class="relative flex aspect-2/3 items-center justify-center overflow-hidden bg-panel"
					>
						<AssetImg
							src={getAssetPath({ gameId: ga.ID, assetType: 'AssetThumbnail' })}
							alt={ga.Meta?.Name ?? 'Untitled'}
							class="h-full w-full object-cover"
						>
							{#snippet loadingSlot()}
								<div
									class="flex h-full w-full animate-pulse items-center justify-center bg-surface"
								>
									<LoaderIcon class="animate-spin text-muted" />
								</div>
							{/snippet}

							{#snippet errorSlot()}
								<div class="flex h-full w-full flex-col items-center justify-center bg-surface">
									<ImageIcon size={48} class="mb-5" />
									<p class="text-sm">Thumbnail Not found</p>
								</div>
							{/snippet}
						</AssetImg>
					</div>

					<!-- Info Area -->
					<div class="poster-metadata">
						<h3 class="mb-1 truncate font-bold text-foreground">
							{ga.Meta?.Name}
						</h3>
						<div class="flex items-center gap-2">
							<span class="telemetry-label text-muted">
								{new Date(ga?.Meta?.ReleaseDate || 'N/A').getFullYear() || 'N/A'}
							</span>
							{#if ga?.DownloadState?.State !== 'Complete'}
								<span
									class="label-hud rounded-full border border-frost-500/30 bg-surface/80 px-2 py-1 text-frost-300 uppercase"
								>
									{ga?.DownloadState?.State || 'Unknown'}
								</span>
							{/if}
						</div>
					</div>
				</a>
			{/each}
		</div>
	{/if}
</main>
