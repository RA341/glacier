<script lang="ts">
	import { SearchIcon } from '@lucide/svelte';
	import type { GameSource } from '#lib/gen/search/v1/search_pb.js';

	// let { source = $bindable(null) }: { source: GameSource | null } = $props();
</script>

<div class="space-y-4 rounded-3xl border border-border bg-surface p-6">
	<div class="flex gap-2">
		<div class="rounded-xl border border-border bg-panel px-3 py-2 text-xs font-bold uppercase">
			Indexer
		</div>
		<div class="relative flex-1">
			<SearchIcon size={14} class="absolute top-1/2 left-3 -translate-y-1/2 text-muted" />
			<input
				type="text"
				placeholder="Search for alternative source..."
				class="w-full rounded-xl border border-border bg-panel py-2 pl-9 text-xs outline-none focus:border-frost-500"
			/>
		</div>
	</div>
	<div class="space-y-2">
		{#each [1, 2] as i}
			<div
				class="flex items-center justify-between rounded-xl border border-border bg-panel/50 p-3 text-xs"
			>
				<span>Other.Release.v1.0-Group</span>
				<button class="btn-secondary rounded-2xl px-3 py-1 font-bold hover:border-frost-500">
					Update
				</button>
			</div>
		{/each}
	</div>
</div>
