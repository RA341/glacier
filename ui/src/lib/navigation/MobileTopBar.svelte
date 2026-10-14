<script lang="ts">
	import { appName } from '#lib/api/api.js';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';
	import { createNavigationStore } from './store.svelte.js';
	import LogoutButton from './LogoutButton.svelte';

	const navigation = createNavigationStore(getUserCtx());
</script>

<header
	class="mobile-top-bar flex shrink-0 items-center justify-between gap-2 border-b border-border bg-background/80 px-3 py-2 backdrop-blur-md md:hidden"
>
	<div class="flex min-w-0 items-center gap-2">
		<img src="/favicon.svg" alt="" class="h-7 w-7 rounded-lg" />
		<span class="truncate text-sm font-bold text-foreground">{appName}</span>
	</div>
	<nav aria-label="Account and settings" class="flex shrink-0 items-center gap-1">
		{#each navigation.footer as link}
			{@const active = navigation.isActive(link.href)}
			<a
				href={link.href}
				title={link.label}
				aria-label={link.label}
				aria-current={active ? 'page' : undefined}
				class="flex h-11 w-11 items-center justify-center rounded-full transition-colors {active
					? 'border border-frost-500/40 bg-frost-500/15 text-frost-300 shadow-frost-glow-sm'
					: 'text-muted hover:bg-panel hover:text-foreground'}"
			>
				<link.icon size={21} />
			</a>
		{/each}
		<LogoutButton />
	</nav>
</header>

<style>
	.mobile-top-bar {
		padding-top: calc(0.5rem + env(safe-area-inset-top));
	}
</style>
