<script lang="ts">
	import type { Game } from '#lib/gen/library/v1/library_pb.js';
	import { ExternalLinkIcon } from '@lucide/svelte';
	import { trimPrefix } from '#lib/api/strings.js';
	import { formatReleaseDate } from '#lib/api/date.js';

	let { game = $bindable(null) }: { game: Game | null } = $props();
	let meta = $derived(game?.Meta);

	const dateInf = $derived(formatReleaseDate(meta?.ReleaseDate ?? ''));
</script>

<div class="grid grid-cols-1 items-start gap-8 lg:grid-cols-4">
	<div class="space-y-6 lg:col-span-3">
		<!-- Summary: Shorter, tighter padding, no fixed aspect ratio -->
		<div
			class="group relative w-full space-y-4 overflow-hidden rounded-2xl border border-frost-500/20 bg-linear-to-br from-surface/40 to-frost-500/10 p-4 shadow-sm backdrop-blur-xl"
		>
			<h3 class="text-[14px] font-black text-frost-400 uppercase">Summary</h3>
			<p class="text-base leading-snug font-medium text-foreground/90">
				{meta?.Summary || 'No summary available.'}
			</p>
		</div>

		<!-- Description: Kept large and spacious -->
		<div
			class="group relative min-h-75 w-full space-y-4 overflow-hidden rounded-3xl border border-frost-500/20 bg-linear-to-br from-surface/40 to-frost-500/10 p-6 shadow-sm backdrop-blur-xl"
		>
			<h3 class="text-[14px] font-black text-frost-400 uppercase">About this game</h3>
			<p class="text-md leading-relaxed whitespace-pre-line text-white">
				{meta?.Description || 'No detailed description found.'}
			</p>
		</div>
	</div>

	<aside class="space-y-6">
		<div
			class="space-y-6 rounded-3xl border border-border border-white/10 bg-surface/80 p-6 shadow-sm backdrop-blur-md"
		>
			<div class="space-y-1">
				<span
					class="flex items-center gap-2 text-[10px] font-black tracking-widest text-muted uppercase"
				>
					Release Date
				</span>
				<p class="text-sm font-bold text-foreground">
					{dateInf.formatted || 'Unknown'}
					{#if dateInf.formatted}
						<span class="text-muted-foreground ml-1 font-normal">
							({dateInf.relative})
						</span>
					{/if}
				</p>
			</div>

			<div class="space-y-1">
				<span
					class="flex items-center gap-2 text-[10px] font-black tracking-widest text-muted uppercase"
				>
					Category
				</span>
				<p class="text-sm font-bold text-foreground capitalize">{meta?.Category || 'Unknown'}</p>
			</div>

			<div class="space-y-1">
				<span
					class="flex items-center gap-2 text-[10px] font-black tracking-widest text-muted uppercase"
				>
					Platforms
				</span>
				<p class="text-sm font-bold text-frost-400">{meta?.Platforms.join(', ')}</p>
			</div>

			<div class="h-px bg-border"></div>

			<div class="space-y-1">
				<span class="text-[10px] font-black tracking-widest text-muted uppercase"
					>Metadata Provider</span
				>
				<div
					class="mt-2 flex items-center justify-between rounded-xl border border-border bg-panel p-3"
				>
					<span class="text-xs font-bold"
						>{trimPrefix(meta?.ProviderType || 'No Provider Found', 'Provider')}</span
					>
					<a href={meta?.URL} target="_blank" class="text-frost-400 hover:text-frost-300">
						<ExternalLinkIcon size={16} />
					</a>
				</div>
			</div>

			<!-- Sidebar Genre Tags -->
			<div class="space-y-2 pt-2">
				<span class="text-[10px] font-black tracking-widest text-muted uppercase">Genres</span>
				<div class="flex flex-wrap gap-2 pt-1">
					{#each meta?.Genres ?? [] as genre}
						<span
							class="rounded-lg border border-border bg-panel px-2 py-1 text-[10px] font-bold whitespace-nowrap"
						>
							{genre}
						</span>
					{/each}
				</div>
			</div>
		</div>
	</aside>
</div>
