<script lang="ts">
	import {
		ChevronDownIcon,
		CircleAlert,
		ExternalLinkIcon,
		FolderOpenIcon,
		HardDriveIcon,
		ImageIcon,
		LinkIcon,
		LoaderIcon,
		SearchIcon,
		StarIcon
	} from '@lucide/svelte';
	import { transferStore } from '../selectedGame.svelte';
	import { glacierCli } from '#lib/api/api.js';
	import { IndexerService } from '#lib/gen/indexer/v1/indexer_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { type GameSource, SearchService } from '#lib/gen/search/v1/search_pb.js';
	import { ServiceConfigService } from '#lib/gen/service_config/v1/service_config_pb.js';
	import AddButton from './AddButton.svelte';
	import IndexerSearch from './IndexerSearch.svelte';

	const srvConfig = glacierCli(ServiceConfigService);
	let indexerManagerService = glacierCli(IndexerService);

	let game = $derived(transferStore.data);
	$effect(() => {
		if (game) {
			availableDownloadClientsRpc.runner();
			availableGameTypeRpc.runner();
		}
	});

	let availableGameTypeRpc = createRPCRunner(() => indexerManagerService.getGameType({}));
	let selectedGameType = $state('');
	$effect(() => {
		if (selectedGameType === '') {
			selectedGameType = availableGameTypeRpc.value?.gameTypes.at(0)?.Name ?? '';
		}

		if (selectedGameSource) {
			selectedGameSource.GameType = selectedGameType;
		}
	});

	let availableDownloadClientsRpc = createRPCRunner(() =>
		srvConfig.getActiveService({ serviceType: 'Downloader' })
	);
	let selectedDownloadClient = $state('');
	$effect(() => {
		selectedDownloadClient = availableDownloadClientsRpc.value?.names.at(0)?.Name ?? '';
	});

	let selectedGameSource = $state(null as GameSource | null);
</script>

<div class="mx-auto max-w-7xl space-y-8 bg-background p-6 text-foreground">
	<!-- 1. VIDEO IFRAME SECTION (Expanded fully width-wise, reduced height) -->
	<!--    <div class="w-full aspect-[21/9] bg-panel rounded-3xl border border-border flex items-center justify-center overflow-hidden shadow-2xl">-->
	<!--        <div class="flex flex-col items-center gap-4 opacity-20">-->
	<!--            <PlusIcon size={64} class="rotate-45"/>-->
	<!--            <p class="font-bold tracking-widest uppercase">Video Preview Placeholder</p>-->
	<!--        </div>-->
	<!--    </div>-->

	<!-- 2. METADATA & ACTION ROW -->
	<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-4">
		<!-- Metadata Content -->
		<div class="flex gap-8 rounded-3xl border border-border bg-surface p-8 lg:col-span-3">
			<div class="flex-1 space-y-4">
				<img src={game?.assets.at(0)?.RemoteUrl ?? ''} alt="game thumbnail" />
			</div>
			<!-- Left: Title & Desc -->
			<div class="flex-1 space-y-4">
				<div class="flex items-center gap-3">
					<h1 class="text-3xl font-bold">{game?.Name || 'Game Title'}</h1>
					<!--                    <a href="#" class="p-2 bg-panel rounded-lg text-frost-400 hover:text-frost-300 transition-colors">-->
					<!--                        <ExternalLinkIcon size={18}/>-->
					<!--                    </a>-->
				</div>
				<div class="h-px w-full bg-border"></div>
				<p class="line-clamp-6 text-sm leading-relaxed text-muted">
					{game?.Summary || 'No description available for this title.'}
				</p>
			</div>

			<!-- Middle: Provider Specs -->
			<div class="w-px bg-border"></div>
			<div class="flex-1 space-y-4">
				<div class="space-y-3">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Provider Info</p>
					<div class="space-y-2 text-sm">
						<div class="flex justify-between">
							<span class="text-muted">Source:</span>
							<span class="font-medium">{game?.ProviderType}</span>
						</div>
						<div class="flex justify-between">
							<span class="text-muted">Released:</span>
							<span class="font-medium"
								>{game?.ReleaseDate ? new Date(game?.ReleaseDate).getFullYear() : 'unknown'}</span
							>
						</div>
						<div class="flex justify-between">
							<span class="text-muted">Status:</span>
							<span class="rounded bg-green-500/10 px-2 text-xs text-green-400"
								>{game?.ReleaseStatus}</span
							>
						</div>
						<div class="mt-4 flex items-center gap-2 text-frost-400">
							<StarIcon size={16} fill="currentColor" />
							<span class="font-bold">{game?.Rating ?? '0.0'}</span>
							<span class="text-xs text-muted">({game?.RatingCount})</span>
						</div>
					</div>
				</div>
			</div>

			<!-- Right: Tags -->
			<div class="w-px bg-border"></div>
			<div class="flex-1 space-y-6">
				<div class="space-y-2">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Genres</p>
					<div class="flex flex-wrap gap-2">
						{#each game?.Genres ?? [] as ge}
							<span class="rounded border border-border bg-panel px-2 py-1 text-[11px]">{ge}</span>
						{/each}
					</div>
				</div>
				<div class="space-y-2">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Platforms</p>
					<div class="flex flex-wrap gap-2">
						{#each game?.Platforms ?? [] as g}
							<span
								class="rounded border border-frost-500/20 bg-frost-500/10 px-2 py-1 text-[11px] text-frost-400"
								>{g}</span
							>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<!-- Action Card (Shorter, height-fit) -->
		<div
			class="flex h-fit flex-col gap-6 rounded-3xl border border-border bg-surface p-6 shadow-xl"
		>
			<AddButton
				gameSrc={selectedGameSource}
				gameMetadata={game}
				downloadClient={selectedDownloadClient}
			/>
		</div>
	</div>

	<!-- SEARCH & DOWNLOAD SECTION -->
	<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
		<!-- SEARCH INDEXER SECTION -->
		<IndexerSearch {game} bind:selectedGameSource />
		<!-- DOWNLOAD CONFIG SECTION -->
		<div class="flex flex-col gap-8 rounded-3xl border border-border bg-surface p-6">
			<h2 class="flex items-center gap-2 text-lg font-bold">
				<HardDriveIcon size={20} class="text-frost-400" />
				Download Config
			</h2>

			<div class="start flex items-center rounded-2xl border border-border bg-panel p-5">
				<div
					class="flex h-20 w-16 items-center justify-center rounded border border-border bg-background text-muted"
				>
					<ImageIcon size={24} />
				</div>
				<div class="space-y-1 px-5">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Selected Release</p>
					<h3 class="font-bold text-foreground">
						{selectedGameSource ? selectedGameSource.Title : 'Waiting for selection...'}
					</h3>
					{#if !selectedGameSource}
						<p class="text-xs text-muted italic">No release matched yet</p>
					{/if}
				</div>
			</div>

			<div class="flex items-center gap-4">
				<p class="text-[10px] font-bold tracking-widest whitespace-nowrap text-muted uppercase">
					Game Type
				</p>

				<div class="relative flex-1">
					<select
						bind:value={selectedGameType}
						class="w-full cursor-pointer appearance-none rounded-xl border border-border bg-panel px-4 py-3 text-sm font-bold outline-none focus:border-frost-500"
					>
						{#each availableGameTypeRpc.value?.gameTypes as gameType}
							<option value={gameType.Name}>{gameType.Name}</option>
						{/each}
					</select>
					<ChevronDownIcon
						size={16}
						class="pointer-events-none absolute top-1/2 right-4 -translate-y-1/2 text-muted"
					/>
				</div>
			</div>

			<div class="flex gap-4">
				<div class="w-1/3 space-y-2">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Download Client</p>
					<div class="relative">
						<select
							bind:value={selectedDownloadClient}
							class="w-full appearance-none rounded-xl border border-border bg-panel px-4 py-2.5 text-sm font-bold outline-none"
						>
							{#each availableDownloadClientsRpc.value?.names as cli}
								<option>{cli.Name}</option>
							{/each}
						</select>
						<ChevronDownIcon
							size={16}
							class="absolute top-1/2 right-3 -translate-y-1/2 text-muted"
						/>
					</div>
				</div>
				<div class="flex-1 space-y-2">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Download URL</p>
					<div class="relative">
						<LinkIcon size={14} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted" />
						<input
							type="text"
							placeholder="link or URL..."
							value={selectedGameSource?.DownloadUrl ?? ''}
							class="w-full rounded-xl border border-border bg-panel py-2.5 pr-4 pl-10 text-sm outline-none"
						/>
					</div>
				</div>
			</div>

			<div class="flex items-end gap-4">
				<div class="flex-1 space-y-2">
					<p class="text-[10px] font-bold tracking-widest text-muted uppercase">Download Path</p>
					<div class="flex gap-2">
						<div
							class="flex-1 truncate rounded-xl border border-border bg-panel px-4 py-2.5 text-sm text-muted"
						>
							/mnt/storage/games/...
						</div>
						<button
							class="rounded-xl border border-border bg-panel p-2.5 transition-colors hover:text-frost-400"
						>
							<FolderOpenIcon size={20} />
						</button>
					</div>
				</div>

				<div class="min-w-35 space-y-2">
					<div class="space-y-1 rounded-xl border border-border bg-panel/50 p-3">
						<div class="flex justify-between text-[9px] font-bold text-muted uppercase">
							<span>Available:</span>
							<span class="text-foreground">1.2 TB</span>
						</div>
						<div class="flex justify-between text-[9px] font-bold text-muted uppercase">
							<span>Remaining:</span>
							<span class="text-frost-400">1.1 TB</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
