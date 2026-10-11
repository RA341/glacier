<script lang="ts">
	import {
		CircleAlert,
		CircleCheck,
		DownloadIcon,
		ImageIcon,
		LoaderIcon,
		PauseIcon,
		RefreshCcwDot,
		ServerIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { glacierCli } from '#lib/api/api.js';
	import { LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { calculatePercent, formatBytes } from '#lib/api/byte-math.js';

	const libSrv = glacierCli(LibraryService);

	let libRpc = createRPCRunner(() => libSrv.triggerTracker({}));

	function trigger() {
		libRpc.runner();
		refresh();
	}

	let activeDownloadsRpc = createRPCRunner(() => libSrv.listWithState({ state: 'downloading' }));

	function refresh() {
		activeDownloadsRpc.runner();
	}

	onMount(() => {
		refresh();
	});

	$effect(() => {
		const interval = setInterval(() => {
			refresh();
		}, 5000);
		return () => clearInterval(interval);
	});

	function getStatusTheme(state: string) {
		if (state?.toLowerCase().includes('paused')) return { color: 'bg-amber-500', icon: PauseIcon };
		if (state?.toLowerCase().includes('complete'))
			return { color: 'bg-green-500', icon: CircleCheck };
		return { color: 'bg-frost-500', icon: DownloadIcon };
	}
</script>

<div class="space-y-6" in:fade={{ duration: 200 }}>
	<!-- Toolbar Section -->
	<div class="flex items-center justify-between px-2">
		<h2 class="text-sm font-bold tracking-widest text-muted uppercase">Active Downloads</h2>
		<button
			onclick={trigger}
			disabled={libRpc.loading}
			class="flex items-center gap-2 rounded-xl border border-border bg-panel p-2 px-4 text-xs font-bold text-muted transition-all hover:border-frost-500/50 hover:text-frost-400 active:scale-95 disabled:opacity-50"
		>
			<RefreshCcwDot size={14} class={libRpc.loading ? 'animate-spin text-frost-500' : ''} />
			Check Status
		</button>
	</div>

	<div class="min-h-100">
		{#if activeDownloadsRpc.error}
			<!-- Error State -->
			<div class="flex h-96 flex-col items-center justify-center gap-3 text-red-400">
				<CircleAlert size={48} strokeWidth={1} />
				<h3 class="font-bold">Connection Error</h3>
				<p class="text-xs opacity-80">{activeDownloadsRpc.error}</p>
			</div>
		{:else if activeDownloadsRpc.loading && !activeDownloadsRpc.value}
			<div class="flex h-96 flex-col items-center justify-center gap-4 text-muted">
				<LoaderIcon class="animate-spin text-frost-500" size={40} />
				<p class="animate-pulse text-sm font-medium">Fetching downloads...</p>
			</div>
		{:else if activeDownloadsRpc.value?.game && activeDownloadsRpc.value.game.length === 0}
			<!-- Empty State -->
			<div
				class="flex h-96 flex-col items-center justify-center rounded-3xl border-2 border-dashed border-border text-muted/30"
			>
				<ServerIcon size={48} strokeWidth={1} class="mb-4" />
				<h2 class="text-xl font-bold text-foreground/50">No Active Downloads</h2>
				<p class="text-sm">Manage your library to trigger new transfers.</p>
			</div>
		{:else}
			<!-- Downloads List -->
			<div class="grid grid-cols-1 gap-4">
				{#each activeDownloadsRpc.value?.game ?? [] as download (download.ID)}
					{@const theme = getStatusTheme(download.DownloadState?.State ?? 'Unknown')}
					{@const progressPercent = calculatePercent(
						download.DownloadState?.Complete,
						download.DownloadState?.Left
					)}
					{@const totalBytes =
						Number(download.DownloadState?.Complete || 0) +
						Number(download.DownloadState?.Left || 0)}

					<div
						class="group flex items-center gap-6 overflow-hidden rounded-3xl border border-border bg-surface p-5 shadow-sm transition-all hover:border-frost-500/30"
						in:fly={{ y: 10, duration: 300 }}
					>
						<div
							class="flex h-28 w-20 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-border bg-panel shadow-inner"
						>
							{#if download.Meta?.ThumbnailURL}
								<img
									src={download.Meta.ThumbnailURL}
									alt=""
									class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
								/>
							{:else}
								<div class="text-muted/20">
									<ImageIcon size={32} strokeWidth={1} />
								</div>
							{/if}
						</div>

						<!-- Main Content -->
						<div class="min-w-0 flex-1 space-y-4">
							<!-- Title & State -->
							<div class="flex items-center justify-between gap-4">
								<h3 class="truncate text-lg font-bold text-foreground">
									{download.Meta?.Name || 'Unknown Title'}
								</h3>
								<span
									class="flex shrink-0 items-center gap-2 rounded-full border border-border bg-panel px-3 py-1 text-[10px] font-bold tracking-wider text-muted uppercase"
								>
									<span
										class="h-1.5 w-1.5 rounded-full {theme.color} {download.DownloadState?.State ===
										'Downloading'
											? 'animate-pulse'
											: ''}"
									></span>
									{download.DownloadState?.State || 'Connecting'}
								</span>
							</div>

							<!-- Progress Bar Container -->
							<div class="space-y-2">
								<div class="mb-1 flex items-end justify-between font-mono text-xs text-muted">
									<!-- Left side: Speed / ETA string from API -->
									<span class="truncate opacity-70">
										{download.DownloadState?.Progress || 'Initializing...'}
									</span>

									<!-- Right side: Bytes Complete / Total -->
									<span class="shrink-0 font-bold text-foreground">
										{formatBytes(Number(download.DownloadState?.Complete || 0))}
										<span class="mx-1 text-muted/30">/</span>
										{formatBytes(totalBytes)}
									</span>
								</div>

								<!-- Visual Progress Bar -->
								<div
									class="h-2 w-full overflow-hidden rounded-full border border-border/50 bg-panel"
								>
									<div
										class="h-full {theme.color} shadow-[0_0_8px_rgba(130,170,255,0.2)] transition-all duration-700 ease-out"
										style="width: {progressPercent}%"
									></div>
								</div>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<style>
	/* Add a subtle glass effect to the cards */
	.bg-surface {
		background-color: rgba(var(--surface-rgb), 0.6);
		backdrop-filter: blur(8px);
	}
</style>
