<script lang="ts">
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';
	import { createNavigationStore } from './store.svelte.js';

	const navigation = createNavigationStore(getUserCtx());
</script>

<nav
	aria-label="Main navigation"
	class="mobile-bottom-nav flex shrink-0 items-center gap-2 border-t border-border bg-surface px-3 py-2 md:hidden"
>
	{#each navigation.main as link}
		{@const active = navigation.isActive(link.href)}
		<a
			href={link.href}
			aria-current={active ? 'page' : undefined}
			class="flex min-h-12 flex-1 flex-col items-center justify-center gap-1 rounded-xl py-1 text-xs font-medium transition-colors {active
				? 'bg-frost-500/10 text-frost-400'
				: 'text-muted hover:bg-panel hover:text-foreground'}"
		>
			<link.icon size={22} strokeWidth={active ? 2.5 : 2} />
			<span>{link.label}</span>
		</a>
	{/each}
</nav>

<style>
	.mobile-bottom-nav {
		padding-bottom: calc(0.5rem + env(safe-area-inset-bottom));
	}
</style>
