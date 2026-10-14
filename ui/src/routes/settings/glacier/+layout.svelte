<script lang="ts">
	import { page } from '$app/state';

	let { children } = $props();

	const prefix = '/settings/glacier';

	const tabs = [
		{ label: 'General', href: 'general', desc: 'Configure basic server settings' },
		{ label: 'Users', href: 'users', desc: 'Manage user accounts' },
		{ label: 'Metadata', href: 'metadata', desc: 'Where to get game information' },
		{ label: 'Indexer', href: 'indexer', desc: 'Where to find and download games' },
		{ label: 'Download', href: 'downloads', desc: 'How downloads are handled' }
	];

	let currentPath = $derived(page.url.pathname);

	function withPrefix(ref: string) {
		return `${prefix}/${ref}`;
	}

	let activeTab = $derived(tabs.find((tab) => currentPath.startsWith(withPrefix(tab.href))));
</script>

<div class="flex h-full flex-col bg-background">
	<header class="px-8 pt-8 pb-0">
		<div class="mb-6 flex items-baseline gap-4">
			<h1 class="headline-xl text-foreground">Glacier</h1>

			{#if activeTab}
				<span class="text-2xl font-light text-muted/30">/</span>
				<p class="text-lg font-medium text-muted transition-all duration-300">
					{activeTab.label}
				</p>
				<span class="text-2xl font-light text-muted/30">/</span>
				<p class="text-lg font-medium text-muted transition-all duration-300">
					{activeTab.desc}
				</p>
			{/if}
		</div>

		<nav class="relative flex gap-2 border-b border-border pb-2">
			{#each tabs as tab}
				{@const active = currentPath.startsWith(withPrefix(tab.href))}
				<a
					href={withPrefix(tab.href)}
					class="group relative rounded-lg px-4 py-2 text-sm font-medium transition-all duration-300
                    {active
						? 'bg-frost-500/10 text-frost-400'
						: 'text-muted hover:bg-panel hover:text-foreground'}"
				>
					{tab.label}

					{#if active}
						<div class="absolute right-0 -bottom-2.25 left-0 h-0.5 rounded-full bg-frost-500"></div>
					{/if}
				</a>
			{/each}
		</nav>
	</header>

	<main class="flex-1 overflow-y-auto p-8">
		{@render children()}
	</main>
</div>
