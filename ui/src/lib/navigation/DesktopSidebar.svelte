<script lang="ts">
	import { createNavigationStore } from './store.svelte.js';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';
	import LogoutButton from './LogoutButton.svelte';

	const navigation = createNavigationStore(getUserCtx());
</script>

<aside
	class="hidden h-full w-16 shrink-0 flex-col items-center border-r border-border bg-surface py-6 md:flex"
>
	<!-- Logo Area -->
	<div class="mb-8 flex items-center justify-center">
		<img src={'/favicon.svg'} alt="Logo" class="h-8 w-8 rounded-lg shadow-frost" />
	</div>

	<!-- Main Navigation -->
	<nav class="flex flex-1 flex-col gap-4">
		{#each navigation.main as link}
			{@const active = navigation.isActive(link.href)}
			<a
				href={link.href}
				title={link.label}
				aria-label={link.label}
				aria-current={active ? 'page' : undefined}
				class="group relative flex h-10 w-10 items-center justify-center rounded-xl transition-all
            duration-300
            {active
					? 'bg-frost-500/10 text-frost-400 shadow-frost'
					: 'text-muted hover:bg-panel hover:text-frost-400'}"
			>
				<link.icon
					size={22}
					strokeWidth={active ? 2.5 : 2}
					class="transition-transform duration-200 {active ? 'scale-110' : 'group-active:scale-90'}"
				/>
			</a>
		{/each}
	</nav>

	<!-- Footer Navigation -->
	<div class="flex flex-col gap-4">
		{#each navigation.footer as link}
			{@const active = navigation.isActive(link.href)}

			<a
				href={link.href}
				title={link.label}
				aria-label={link.label}
				aria-current={active ? 'page' : undefined}
				class="group relative flex h-10 w-10 items-center justify-center rounded-xl transition-all
            duration-300
            {active
					? 'bg-frost-500/10 text-frost-400 shadow-frost'
					: 'text-muted hover:bg-panel hover:text-frost-400'}"
			>
				<link.icon size={22} strokeWidth={2} />
			</a>
		{/each}

		<LogoutButton />
	</div>
</aside>
