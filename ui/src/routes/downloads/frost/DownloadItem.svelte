<script lang="ts">
	import {
		ChevronDownIcon,
		CircleCheck,
		ClockIcon,
		HourglassIcon,
		ImageIcon,
		PauseIcon,
		PlayIcon,
		XIcon
	} from '@lucide/svelte';
	import { slide } from 'svelte/transition';
	import {
		type DownloadProgress,
		FrostLibraryService
	} from '#lib/gen/frost_library/v1/frost_library_pb.js';
	import { formatBytes } from '#lib/api/byte-math.js';
	import { callRPC, frostCli } from '#lib/api/api.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { trimPrefix } from '#lib/api/strings.js';

	let { detail }: { detail: DownloadProgress } = $props();

	let isExpanded = $state(false);

	let now = $state(Date.now());

	$effect(() => {
		const interval = setInterval(() => {
			now = Date.now();
		}, 1000);
		return () => clearInterval(interval);
	});

	const elapsedTime = $derived.by(() => {
		if (!detail.download?.Started) return null;
		const start = new Date(detail.download.Started).getTime();
		const diff = Math.max(0, now - start);

		const seconds = Math.floor((diff / 1000) % 60);
		const minutes = Math.floor((diff / (1000 * 60)) % 60);
		const hours = Math.floor(diff / (1000 * 60 * 60));

		const parts = [
			hours > 0 ? hours.toString().padStart(2, '0') : null,
			minutes.toString().padStart(2, '0'),
			seconds.toString().padStart(2, '0')
		].filter(Boolean);

		return parts.join(':');
	});

	const totalComplete = $derived(Number(detail.progress?.Complete ?? 0));
	const totalLeft = $derived(Number(detail.progress?.Left ?? 0));
	const totalSize = $derived(totalComplete + totalLeft);
	const overallProgress = $derived(totalSize > 0 ? (totalComplete / totalSize) * 100 : 0);

	function getStateColor(state: string = '') {
		const s = state.toLowerCase();
		if (s.includes('downloading') || s.includes('done'))
			return 'text-green-400 bg-green-400/10 border-green-400/20';
		if (s.includes('queued')) return 'text-amber-400 bg-amber-400/10 border-amber-400/20';
		return 'text-frost-400 bg-frost-400/10 border-frost-400/20';
	}

	const getFileProgress = (complete: bigint | number, left: bigint | number) => {
		const c = Number(complete);
		const l = Number(left);
		return c + l > 0 ? (c / (c + l)) * 100 : 0;
	};

	const frostLib = frostCli(FrostLibraryService);
	const sm = getSnackbarCtx();

	let isPaused = $derived(detail?.download?.Status === 'StatusPaused');

	async function toggleDownload(ID: bigint) {
		if (isPaused) {
			const { err } = await callRPC(() => frostLib.pause({ id: ID, pause: false }));
			if (err) {
				sm.push(`Could not cancel game: ${err}`, 'error');
			}
			return;
		}

		const { err } = await callRPC(() => frostLib.pause({ id: ID, pause: true }));
		if (err) {
			sm.push(`Could not cancel game: ${err}`, 'error');
		}
	}

	async function cancel(ID: bigint) {
		const { err } = await callRPC(() => frostLib.cancel({ id: ID }));
		if (err) {
			sm.push(`Could not cancel game: ${err}`, 'error');
		}
	}
</script>

<div
	class="overflow-hidden rounded-2xl border border-border bg-panel/30 transition-all duration-300 {isExpanded
		? 'border-frost-500/30 ring-1 ring-frost-500/10'
		: ''}"
>
	<!-- ACCORDION HEADER -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div
		onclick={() => (isExpanded = !isExpanded)}
		class="group flex w-full cursor-pointer items-center gap-4 p-3 text-left transition-colors hover:bg-panel/50"
	>
		<!-- THUMBNAIL -->
		<div
			class="flex h-16 w-12 shrink-0 items-center justify-center overflow-hidden rounded-xl border border-border bg-surface text-muted/20"
		>
			{#if detail.Thumbnail}
				<img
					src={detail.Thumbnail}
					alt=""
					class="h-full w-full object-cover transition-transform group-hover:scale-110"
				/>
			{:else}
				<ImageIcon size={20} strokeWidth={1.5} />
			{/if}
		</div>

		<!-- MAIN INFO SECTION -->
		<div class="min-w-0 flex-1 space-y-1.5">
			<div class="flex items-center justify-between gap-2">
				<span class="truncate text-sm font-bold tracking-tight text-foreground">
					{detail.Title || 'Unknown Package'}
				</span>
				<span
					class="rounded-md border border-border bg-panel/50 px-2 py-0.5 font-hud text-[10px] font-bold text-frost-400 tabular-nums"
				>
					{overallProgress.toFixed(1)}%
				</span>
			</div>

			<!-- Meta Row: State, Message, Elapsed -->
			<div class="flex items-center gap-3 text-[10px] font-medium">
				<span
					class="shrink-0 rounded border px-1.5 py-0.5 tracking-tighter uppercase {getStateColor(
						detail.download?.Status
					)}"
				>
					{trimPrefix(detail.download?.Status ?? '', 'Status') || 'Waiting'}
				</span>
				{#if elapsedTime}
					<span
						class="ml-auto flex items-center gap-1 font-hud whitespace-nowrap text-white tabular-nums"
					>
						<HourglassIcon size={10} class="text-frost-500/50" />
						{elapsedTime}
					</span>
				{/if}
			</div>

			<!-- Overall Progress Bar -->
			<div class="mt-1 h-1.5 w-full overflow-hidden rounded-full border border-white/5 bg-panel">
				<div
					class="h-full bg-frost-500 shadow-[0_0_8px_rgba(56,189,248,0.3)] transition-all duration-700 ease-out"
					style="width: {overallProgress}%"
				></div>
			</div>
		</div>

		<!-- BYTE STATS SECTION -->
		<div class="ml-2 flex shrink-0 items-center gap-3">
			<div
				class="hidden flex-col items-end text-[10px] font-bold tracking-widest whitespace-nowrap text-muted/40 uppercase md:flex"
			>
				<span class="text-foreground/60">{formatBytes(totalComplete)} Done</span>
				<span>{formatBytes(totalLeft)} Left</span>
			</div>

			<!-- Action Buttons: Stop propagation so they don't toggle the accordion -->
			{#if overallProgress < 100}
				<div class="ml-2 flex translate-x-2 items-center gap-1 transition-all">
					<button
						onclick={(e) => {
							e.stopPropagation();
							toggleDownload(detail.ID);
						}}
						class="rounded-full p-1.5 text-muted transition-all hover:bg-amber-400/10 hover:text-amber-400"
						title="Pause Download"
					>
						{#if isPaused}
							<PlayIcon size={14} />
						{:else}
							<PauseIcon size={14} />
						{/if}
					</button>
					<button
						onclick={(e) => {
							e.stopPropagation();
							cancel(detail.ID);
						}}
						class="rounded-full p-1.5 text-muted transition-all hover:bg-red-400/10 hover:text-red-400"
						title="Cancel Download"
					>
						<XIcon size={14} />
					</button>
				</div>
			{/if}

			<div class="rounded-lg p-1 text-muted transition-colors hover:bg-panel">
				<ChevronDownIcon
					size={18}
					class="transition-transform duration-300 {isExpanded ? 'rotate-180 text-frost-400' : ''}"
				/>
			</div>
		</div>
	</div>

	{#if isExpanded}
		<!-- ACCORDION CONTENT (File List) -->
		{#if isExpanded}
			<div transition:slide={{ duration: 300 }} class="border-t border-border bg-black/10">
				<div class="custom-scrollbar max-h-80 space-y-1 overflow-y-auto p-2">
					{#each detail.progress?.files ?? [] as file}
						{@const filePct = getFileProgress(file.BytesComplete, file.BytesLeft)}
						{@const isDone = Number(file.BytesLeft) === 0}

						<!-- Chunk Logic -->
						{@const chunksComplete = Number(file.ChunksComplete ?? 0)}
						{@const chunksLeft = Number(file.ChunksLeft ?? 0)}
						{@const chunksTotal = chunksComplete + chunksLeft}

						<!--
                            Performance Guard:
                            If a file has 10,000 chunks, we don't want 10,000 divs.
                            We cap the visual blocks at 60 for a consistent "bitfield" look.
                        -->
						{@const visualBlockCount = chunksTotal > 60 ? 60 : chunksTotal}
						{@const filledBlocks = Math.floor((chunksComplete / chunksTotal) * visualBlockCount)}

						<div
							class="group flex items-center gap-3 rounded-xl p-3 transition-colors hover:bg-panel/40"
						>
							<div class="shrink-0 text-muted/40 transition-colors group-hover:text-frost-400/50">
								{#if isDone}
									<CircleCheck size={14} class="text-green-500/60" />
								{:else}
									<ClockIcon size={14} />
								{/if}
							</div>

							<div class="min-w-0 flex-1">
								<div class="mb-1.5 flex items-center gap-3">
									<p
										class="min-w-0 flex-1 truncate text-[12px] font-medium text-white transition-colors group-hover:text-foreground"
									>
										{file.Name}
									</p>

									<span
										class="shrink-0 font-hud text-[10px] whitespace-nowrap text-white tabular-nums"
									>
										{formatBytes(Number(file.BytesComplete))} / {formatBytes(
											Number(file.BytesComplete) + Number(file.BytesLeft)
										)}
										<span class="ml-1 text-white/60">
											({Math.round(
												(Number(file.BytesComplete) /
													(Number(file.BytesComplete) + Number(file.BytesLeft))) *
													100
											)}%)
										</span>
									</span>

									<span
										class="shrink-0 text-[10px] font-black tracking-tighter whitespace-nowrap text-white uppercase"
									>
										Blocks: {chunksComplete}/{chunksTotal}
										<span class="ml-1 text-white/60">
											({chunksTotal > 0 ? Math.round((chunksComplete / chunksTotal) * 100) : 0}%)
										</span>
									</span>
								</div>

								<div class="space-y-1">
									<!-- 2. Segmented Chunk-level Bitfield Bar -->
									{#if chunksTotal > 0}
										<div class="flex h-1 w-full gap-px" aria-hidden="true">
											{#each Array(visualBlockCount) as _, i}
												<div
													class="h-full flex-1 rounded-[1px] transition-colors duration-500
                                            {i < filledBlocks
														? isDone
															? 'bg-green-500/20'
															: 'bg-frost-500/60 shadow-[0_0_3px_rgba(56,189,248,0.2)]'
														: 'bg-white/3'}"
												></div>
											{/each}
										</div>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	{/if}
</div>

<style>
	button {
		user-select: none;
	}

	.custom-scrollbar::-webkit-scrollbar {
		width: 4px;
	}

	.custom-scrollbar::-webkit-scrollbar-track {
		background: transparent;
	}

	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.1);
		border-radius: 10px;
	}
</style>
