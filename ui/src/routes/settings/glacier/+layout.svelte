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

<div class="flex h-full min-h-0 flex-col bg-background">
	<header class="shrink-0 pt-4 sm:pt-6">
		<div class="mb-4 h-14 px-4 sm:mb-6 sm:px-8">
			<nav aria-label="Breadcrumb" class="breadcrumb-scroll overflow-x-auto">
				<ol
					class="flex min-w-max items-center gap-2 text-base font-semibold whitespace-nowrap sm:text-lg"
				>
					<li>
						<a href="/settings" class="text-muted transition-colors hover:text-frost-400">
							Settings
						</a>
					</li>
					<li aria-hidden="true" class="text-muted/40">/</li>
					<li>
						<a href={prefix} class="text-muted transition-colors hover:text-frost-400"> Glacier </a>
					</li>
					{#if activeTab}
						<li aria-hidden="true" class="text-muted/40">/</li>
						<li>
							<a
								href={withPrefix(activeTab.href)}
								aria-current="page"
								class="text-foreground transition-colors hover:text-frost-400"
							>
								{activeTab.label}
							</a>
						</li>
					{/if}
				</ol>
			</nav>

			<p class="mt-1 h-5 truncate text-xs text-muted sm:text-sm">
				{activeTab?.desc ?? 'Configure Glacier'}
			</p>
		</div>

		<div class="tab-scroll overflow-x-auto border-b border-border px-4 sm:px-8">
			<nav aria-label="Glacier settings" class="flex min-w-max gap-2 pb-2">
				{#each tabs as tab}
					{@const active = currentPath.startsWith(withPrefix(tab.href))}
					<a
						href={withPrefix(tab.href)}
						aria-current={active ? 'page' : undefined}
						class="group relative shrink-0 rounded-full px-4 py-2 text-xs font-medium transition-colors sm:text-sm
                        {active
							? 'border border-frost-500/40 bg-frost-500/10 text-frost-400'
							: 'border border-transparent text-muted hover:bg-panel hover:text-foreground'}"
					>
						{tab.label}

						{#if active}
							<span
								aria-hidden="true"
								class="absolute right-2 -bottom-2.25 left-2 h-0.5 rounded-full bg-frost-500"
							></span>
						{/if}
					</a>
				{/each}
			</nav>
		</div>
	</header>

	<main class="min-h-0 flex-1 overflow-y-auto p-4 sm:p-8">
		{@render children()}
	</main>
</div>

<style>
	.breadcrumb-scroll,
	.tab-scroll {
		scrollbar-width: none;
	}

	.breadcrumb-scroll::-webkit-scrollbar,
	.tab-scroll::-webkit-scrollbar {
		display: none;
	}
</style>
