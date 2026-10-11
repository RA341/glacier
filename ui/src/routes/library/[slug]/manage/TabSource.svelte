<script lang="ts">
	import { ImageIcon } from '@lucide/svelte';
	import IndexerSearch from '#lib/components/IndexerSearch.svelte';

	let { game = $bindable() } = $props();
</script>

<!-- Source Section -->
<section class="space-y-4">
	<h2 class="px-2 text-sm font-bold tracking-widest text-muted uppercase">Source</h2>
	<div class="grid grid-cols-1 gap-6 lg:grid-cols-2">
		<!-- Current Source Info -->
		<div class="flex gap-6 rounded-3xl border border-border bg-surface p-6">
			<div
				class="flex h-32 w-24 shrink-0 items-center justify-center rounded-xl border border-border bg-panel text-muted"
			>
				<ImageIcon size={32} />
			</div>
			<div class="grid flex-1 grid-cols-2 gap-4 text-xs">
				<div class="col-span-2">
					<p class="text-[9px] font-bold text-muted uppercase">Title</p>
					<p class="text-sm font-bold">Game.Name.2023-REPACK</p>
				</div>
				<div>
					<p class="text-[9px] font-bold text-muted uppercase">Indexer</p>
					<p>1337x</p>
				</div>
				<div>
					<p class="text-[9px] font-bold text-muted uppercase">Size</p>
					<p>45.2 GB</p>
				</div>
				<div>
					<p class="text-[9px] font-bold text-muted uppercase">Created</p>
					<p>2024-01-15</p>
				</div>
			</div>
		</div>

		<IndexerSearch />
	</div>
</section>
