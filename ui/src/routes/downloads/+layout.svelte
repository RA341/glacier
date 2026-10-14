<script lang="ts">
	import { DownloadIcon } from '@lucide/svelte';
	import { isFrost } from '#lib/api/api.js';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	let { children } = $props();

	let activeTab = $derived.by(() => {
		const path = page.url.pathname;
		if (path.includes('/downloads/frost')) return 'frost';
		if (path.includes('/downloads/glacier')) return 'server';
		return '';
	});

	function handleTabChange(tab: string) {
		const path = tab === 'frost' ? '/downloads/frost' : '/downloads/glacier';
		goto(path, { replaceState: true, reset: false });
	}
</script>

<div class="mx-auto min-h-screen max-w-7xl space-y-8 bg-background p-6 text-foreground">
	<header class="flex items-center justify-between border-b border-border pb-6">
		<div class="flex items-center gap-4">
			<div class="rounded-2xl border border-frost-500/20 bg-frost-500/10 p-3 text-frost-400">
				<DownloadIcon size={28} />
			</div>
			<div>
				<h1 class="headline-xl">Downloads</h1>
				<p class="text-sm text-muted">Monitor and manage active downloads.</p>
			</div>
		</div>

		<div class="flex gap-1 rounded-xl border border-border bg-panel p-1">
			{#if isFrost}
				<button
					onclick={() => handleTabChange('frost')}
					class="rounded-full px-6 py-2 text-sm font-bold transition-all {activeTab === 'frost'
						? 'bg-surface text-frost-400 '
						: 'text-muted hover:text-foreground'}"
				>
					Frost
				</button>
			{/if}

			<button
				onclick={() => handleTabChange('server')}
				class="rounded-full px-6 py-2 text-sm font-bold transition-all {activeTab === 'server'
					? 'bg-surface text-frost-400 '
					: 'text-muted hover:text-foreground'}"
			>
				Glacier
			</button>
		</div>
	</header>

	<main>
		{@render children()}
	</main>
</div>
