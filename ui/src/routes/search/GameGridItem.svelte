<script lang="ts">
	import { ImageIcon, PlusIcon } from '@lucide/svelte';
	import { goto } from '$app/navigation';
	import { transferStore } from './selectedGame.svelte.ts';
	import { type GameMetadata, SearchService } from '#lib/gen/search/v1/search_pb.js';

	let { game }: { game: GameMetadata } = $props();

	function handleAdd(meta: GameMetadata) {
		transferStore.data = meta;

		goto('/search/details', { reset: false });
	}
</script>

<button onclick={() => handleAdd(game)} class="poster-card group flex flex-col">
	<div class="relative aspect-2/3 w-full overflow-hidden bg-panel">
		{#if game.assets.at(0)?.RemoteUrl}
			<img
				src={game.assets.at(0)?.RemoteUrl}
				alt={game.Name}
				class="h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
			/>
		{:else}
			<div class="flex h-full w-full flex-col items-center justify-center gap-2 text-muted/20">
				<ImageIcon size={48} strokeWidth={1} />
				<span class="text-[10px] font-bold tracking-widest uppercase">No Preview</span>
			</div>
		{/if}

		<div
			class="absolute inset-0 flex items-center justify-center bg-frost-900/40 opacity-0 transition-opacity group-hover:opacity-100 group-focus-visible:opacity-100"
		>
			<div
				class="translate-y-4 rounded-full bg-frost-500 p-3 text-background transition-transform group-hover:translate-y-0 group-focus-visible:translate-y-0"
			>
				<PlusIcon size={24} />
			</div>
		</div>
	</div>

	<div class="poster-metadata flex flex-col gap-1">
		<h3 class="truncate text-sm font-bold text-foreground">{game.Name}</h3>
		<div class="telemetry-label flex items-center justify-between text-muted">
			<span class="flex items-center gap-1">
				{new Date(game.ReleaseDate).getFullYear()}
			</span>
		</div>
	</div>
</button>
