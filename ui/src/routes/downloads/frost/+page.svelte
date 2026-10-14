<script lang="ts">
	import {
		CircleAlert,
		CloudIcon,
		LoaderIcon,
		RefreshCcwDot,
		TrendingUpIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { callRPC, Frost, frostCli } from '#lib/api/api.js';
	import { FrostLibraryService } from '#lib/gen/frost_library/v1/frost_library_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import DownloadItem from './DownloadItem.svelte';
	import { onMount, untrack } from 'svelte';
	import { handleProcessStatus } from '#lib/api/websockets.js';
	import DownloadThrottle from './DownloadThrottle.svelte';

	const frostLib = frostCli(FrostLibraryService);
	let downloadingRpc = createRPCRunner(() => frostLib.listDownloading({}));

	function refresh() {
		if (downloadingRpc.loading) return;
		downloadingRpc.runner();
	}

	onMount(() => {
		refresh();
		const interval = setInterval(refresh, 2000);
		return () => clearInterval(interval);
	});

	let speed = $state('');
	let isTracking = $state(false);
	let socketCleanup: (() => void) | null = null;

	const hasDownloads = $derived((downloadingRpc.value?.downloads?.length ?? 0) > 0);

	$effect(() => {
		// We only want this effect to react to 'hasDownloads'
		if (hasDownloads) {
			// Use untrack so that setting isTracking doesn't trigger the effect again
			untrack(() => {
				if (!isTracking) {
					socketCleanup = startSpeedTracking();
				}
			});
		} else {
			untrack(() => {
				if (isTracking) {
					socketCleanup?.();
					isTracking = false;
				}
			});
		}

		// Cleanup on component destroy
		return () => {
			socketCleanup?.();
		};
	});

	function startSpeedTracking() {
		isTracking = true;
		const proc = `${Frost.base}/launcher/speed`;

		const ws = handleProcessStatus({
			url: proc,
			onDone: () => {
				speed = '';
				isTracking = false;
			},
			onMessage: (data) => {
				speed = data?.speed ?? '';
			}
		});

		return () => {
			ws?.(); // Call the cleanup returned by your websocket helper
			isTracking = false;
			speed = '';
		};
	}
</script>

<div class="space-y-6">
	<!-- TOOLBAR / HEADER -->
	<header class="flex flex-col justify-between gap-4 px-2 sm:flex-row sm:items-center">
		<div class="flex items-center gap-3">
			<button
				onclick={refresh}
				disabled={downloadingRpc.loading}
				class="group btn-secondary flex items-center gap-2 rounded-2xl p-2 px-4 text-xs font-bold text-muted transition-all hover:border-frost-500/50 hover:text-frost-400 active:scale-95 disabled:opacity-50"
			>
				<RefreshCcwDot
					size={14}
					class={downloadingRpc.loading
						? 'animate-spin text-frost-500'
						: 'transition-transform duration-500 group-hover:rotate-180'}
				/>
				Refresh
			</button>
			<DownloadThrottle />
			<!-- LIVE SPEED BADGE -->
			{#if speed && hasDownloads}
				<div
					in:fade
					class="flex items-center gap-2 rounded-full border border-frost-500/20 bg-frost-500/5 px-3 py-1.5 text-[10px] font-bold tracking-wider text-frost-400 uppercase shadow-[0_0_15px_rgba(130,170,255,0.05)]"
				>
					<span class="relative flex h-2 w-2">
						<span
							class="absolute inline-flex h-full w-full animate-ping rounded-full bg-frost-400 opacity-75"
						></span>
						<span class="relative inline-flex h-2 w-2 rounded-full bg-frost-500"></span>
					</span>
					<TrendingUpIcon size={12} class="opacity-50" />
					<span>{speed}</span>
				</div>
			{/if}
		</div>
	</header>

	<!-- CONTENT AREA -->
	<main class="min-h-100" in:fade={{ duration: 200 }}>
		{#if downloadingRpc.error}
			<!-- ERROR STATE -->
			<div
				class="flex h-96 flex-col items-center justify-center gap-3 rounded-3xl border border-red-500/10 bg-red-500/5 text-red-400"
				in:fly={{ y: 10 }}
			>
				<CircleAlert size={40} strokeWidth={1.5} />
				<div class="text-center">
					<h3 class="text-xs font-bold tracking-widest uppercase">Failed to get downloads</h3>
					<p class="mt-1 text-xs opacity-60">{downloadingRpc.error}</p>
				</div>
				<button
					onclick={refresh}
					class="btn-secondary mt-2 rounded-2xl px-4 py-2 text-[10px] font-black uppercase transition-colors hover:bg-surface"
				>
					Retry Connection
				</button>
			</div>
		{:else if downloadingRpc.loading && !downloadingRpc.value}
			<div class="flex h-96 flex-col items-center justify-center gap-4 text-muted">
				<div class="relative">
					<LoaderIcon class="animate-spin text-frost-500" size={48} strokeWidth={1} />
					<CloudIcon
						size={16}
						class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 opacity-50"
					/>
				</div>
				<p class="animate-pulse text-[10px] font-bold tracking-[0.3em] uppercase">
					Loading Downloads...
				</p>
			</div>
		{:else if !hasDownloads}
			<!-- EMPTY STATE -->
			<div
				class="flex h-96 flex-col items-center justify-center gap-4 rounded-[2.5rem] border-2 border-dashed border-border text-muted/20"
				in:fade
			>
				<div class="rounded-full bg-panel/20 p-6">
					<CloudIcon size={64} strokeWidth={1} />
				</div>
				<div class="text-center">
					<h2 class="text-xl font-bold text-foreground/30">No active downloads</h2>
					<p class="mt-1 text-xs font-medium tracking-widest uppercase">
						Your frost queue is empty
					</p>
				</div>
			</div>
		{:else}
			<!-- LIST STATE -->
			<div class="space-y-4">
				{#each downloadingRpc.value?.downloads ?? [] as down (down.ID)}
					<div in:fly={{ y: 20, duration: 400 }}>
						<DownloadItem detail={down} />
					</div>
				{/each}
			</div>
		{/if}
	</main>
</div>

<style>
	/* Prevent layout shift during refresh */
	main {
		contain: layout;
	}
</style>
