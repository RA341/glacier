<script lang="ts">
	import {
		CalendarIcon,
		ChevronLeft,
		ChevronRight,
		CircleAlert,
		ImageIcon,
		LoaderIcon,
		MonitorIcon,
		Pen,
		StarIcon
	} from '@lucide/svelte';
	import type { Game } from '#lib/gen/library/v1/library_pb.js';
	import { getAssetPath } from '#lib/api/assets.js';
	import AssetImg from '#lib/components/assets/AssetImg.svelte';
	import AssetVideo from '#lib/components/assets/AssetVideo.svelte';
	import GameDownloadButton from '#lib/components/ButtonDownload.svelte';
	import { glacierCli, isFrost } from '#lib/api/api.js';
	import { AssetService } from '#lib/gen/assets/v1/assets_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { trimPrefix } from '#lib/api/strings.js';

	let { game = $bindable(null), onManage }: { game: Game | null; onManage: () => void } = $props();

	let meta = $derived(game?.Meta);
	const thumb = getAssetPath({ gameId: game?.ID, assetType: 'AssetThumbnail' });
	const trailer = getAssetPath({ gameId: game?.ID, assetType: 'AssetTrailer' });

	const assetSrv = glacierCli(AssetService);

	const listRpc = createRPCRunner(() =>
		assetSrv.list({
			ID: game?.ID,
			AssetType: [
				'AssetTrailer',
				'AssetGameplayImage',
				'AssetGameplayVideo',
				'AssetArtwork',
				'AssetBanner'
			]
		})
	);

	const assetType = $state('');
	let activeIndex = $state(0);
	let assets = $derived(listRpc.value?.assets ?? []);
	let currentAsset = $derived(assets[activeIndex]);

	function next() {
		activeIndex = (activeIndex + 1) % assets.length;
	}

	let getImg = (inputLink: string): string => {
		if (!inputLink) return '';

		return getAssetPath({
			gameId: game?.ID,
			assetPath: inputLink
		});
	};

	function prev() {
		activeIndex = (activeIndex - 1 + assets.length) % assets.length;
	}

	function isVideo(type: string) {
		return type.toLowerCase().includes('video') || type.toLowerCase().includes('trailer');
	}

	onMount(async () => {
		if (!game?.ID) {
			return;
		}
		await listRpc.runner();
	});
</script>

<div class="space-y-4">
	<div
		class="media-surface group relative aspect-video w-full overflow-hidden rounded-3xl border border-white/10 bg-surface/80 p-6 shadow-sm backdrop-blur-md lg:aspect-20/9"
	>
		{#if listRpc.loading && !listRpc.value}
			<div class="flex h-full w-full animate-pulse items-center justify-center bg-background/50">
				<LoaderIcon class="animate-spin text-frost-500" size={40} />
			</div>
		{:else if listRpc.error}
			<div class="flex h-full w-full flex-col items-center justify-center gap-2 text-red-400">
				<CircleAlert size={32} />
				<p class="text-xs font-black tracking-widest uppercase">Failed to load gallery</p>
			</div>
		{:else if assets.length > 0}
			<!-- Asset Label Overlay -->
			<div
				class="pointer-events-none absolute top-6 right-8 z-20 rounded-full border border-white/10 bg-background/40 px-4 py-1.5 text-[10px] font-black tracking-[0.2em] text-white/90 uppercase backdrop-blur-md"
			>
				{trimPrefix(currentAsset.Type, 'Asset')}
			</div>
			<!-- Navigation Arrows -->
			{#if assets.length > 1}
				<button
					aria-label="Previous media"
					onclick={prev}
					class="absolute top-1/2 left-4 z-20 -translate-y-1/2 rounded-full border border-white/10 bg-background/20 p-3 text-white opacity-0 backdrop-blur-md transition-all group-focus-within:opacity-100 group-hover:opacity-100 hover:bg-frost-500 hover:text-background active:scale-95"
				>
					<ChevronLeft size={24} strokeWidth={3} />
				</button>
				<button
					aria-label="Next media"
					onclick={next}
					class="absolute top-1/2 right-4 z-20 -translate-y-1/2 rounded-full border border-white/10 bg-background/20 p-3 text-white opacity-0 backdrop-blur-md transition-all group-focus-within:opacity-100 group-hover:opacity-100 hover:bg-frost-500 hover:text-background active:scale-95"
				>
					<ChevronRight size={24} strokeWidth={3} />
				</button>
			{/if}

			<!-- Active Asset Display -->
			{#key activeIndex}
				<div in:fade={{ duration: 400 }} class="h-full w-full">
					{#if isVideo(currentAsset.Type)}
						<AssetVideo
							src={getImg(currentAsset.LocalPath) || currentAsset.RemoteUrl}
							autoplay={true}
							muted={true}
							controls={true}
							class="h-full w-full object-cover"
						>
							{#snippet loadingSlot()}
								<div class="flex h-full w-full items-center justify-center bg-black">
									<LoaderIcon class="animate-spin text-frost-500" />
								</div>
							{/snippet}
						</AssetVideo>
					{:else}
						<img
							src={getImg(currentAsset.LocalPath) || currentAsset.RemoteUrl}
							class="absolute inset-0 h-full w-full scale-110 object-cover opacity-50 blur-2xl"
							alt=""
						/>

						<AssetImg
							src={getImg(currentAsset.LocalPath) || currentAsset.RemoteUrl}
							class="h-full w-full "
							imgClass="w-full h-full object-contain"
						/>
					{/if}
				</div>
			{/key}

			<!-- Indicators (Dots) -->
			<div class="absolute bottom-6 left-1/2 z-20 flex -translate-x-1/2 gap-2">
				{#each assets as _, i}
					<button
						title="Indicator"
						onclick={() => (activeIndex = i)}
						class="h-1.5 rounded-full transition-all
 {activeIndex === i ? 'btn-frost w-8 ' : 'w-2 bg-white/20 hover:bg-white/40'}"
					></button>
				{/each}
			</div>
		{:else}
			<!-- Empty State -->
			<div class="flex h-full w-full flex-col items-center justify-center gap-3 text-muted/20">
				<ImageIcon size={64} strokeWidth={1} />
				<p class="text-xs font-black tracking-[0.3em] uppercase">No Media Assets</p>
			</div>
		{/if}
	</div>

	<!-- MIDDLE SECTION: INFO BAR -->
	<div class="rounded-3xl border border-white/10 bg-surface/80 p-6 shadow-sm backdrop-blur-md">
		<div class="flex flex-col justify-between gap-6 lg:flex-row lg:items-center">
			<div class="space-y-4">
				<h1 class="display-hero text-foreground">{meta?.Name}</h1>

				<div class="flex flex-wrap items-center gap-1 text-sm font-bold">
					<!-- Rating -->
					<div class="flex items-center gap-2 rounded-xl bg-transparent px-3 py-1.5">
						<StarIcon size={16} class="fill-yellow-500 text-yellow-500" />
						<span
							>{meta?.Rating}
							<span class="font-medium text-muted">({meta?.RatingCount})</span></span
						>
					</div>

					<div class="hidden h-4 w-px bg-border lg:block"></div>

					<!-- Release Year -->
					<div class="flex items-center gap-2 rounded-xl bg-transparent px-3 py-1.5">
						<CalendarIcon size={16} class="text-frost-400" />
						<span>{new Date(meta?.ReleaseDate || '').getFullYear()}</span>
					</div>

					<div class="hidden h-4 w-px bg-border lg:block"></div>

					<!-- Platform Tags -->
					<div class="flex items-center gap-2 rounded-xl bg-transparent px-3 py-1.5">
						<MonitorIcon size={16} class="text-frost-400" />
						{#each meta?.Platforms.slice(0, 2) as Platform}
							<span
								class="label-hud rounded-full border border-frost-500/15 bg-frost-500/5 px-3 py-1 text-muted uppercase"
							>
								{Platform}
							</span>
						{/each}
					</div>
				</div>

				<!-- Genre Tags -->
				<div class="flex flex-wrap gap-2">
					{#each meta?.Genres.slice(0, 4) ?? [] as genre}
						<span
							class="label-hud rounded-full border border-frost-500/15 bg-frost-500/5 px-3 py-1 text-frost-400 uppercase"
						>
							{genre}
						</span>
					{/each}
					{#if meta?.Genres?.length && meta?.Genres?.length > 4}
						<span
							class="rounded-full border border-border bg-panel px-3 py-1 text-[11px] font-bold text-muted"
						>
							+{meta.Genres.length - 4} more
						</span>
					{/if}
				</div>
			</div>

			<!-- Actions -->
			<div class="flex shrink-0 items-center gap-3">
				<button
					onclick={onManage}
					class="flex items-center gap-2 rounded-full border border-white/10 bg-white/5 px-8 py-3 text-sm font-black tracking-widest uppercase transition-all hover:border-frost-500/60 hover:bg-white/10"
				>
					<Pen size={18} />
					Manage
				</button>

				{#if isFrost && game?.DownloadState?.State === 'Complete'}
					<GameDownloadButton gameId={game.ID} />
				{/if}
			</div>
		</div>
	</div>
</div>
