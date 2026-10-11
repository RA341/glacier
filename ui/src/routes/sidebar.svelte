<script lang="ts">
	import { page } from '$app/state';
	import {
		DownloadIcon,
		LibraryIcon,
		LogOutIcon,
		MountainSnow,
		SearchIcon,
		Snowflake,
		UserIcon
	} from '@lucide/svelte';
	import { glacierPubCli, isFrost } from '#lib/api/api.js';
	import { AuthService } from '#lib/gen/auth/v1/auth_pb.js';
	import { goto } from '$app/navigation';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';

	const isActive = (href: string) => page.url.pathname.startsWith(href);

	const links = [
		{ label: 'Library', href: '/library', icon: LibraryIcon },
		{ label: 'Search', href: '/search', icon: SearchIcon },
		{ label: 'Downloads', href: '/downloads', icon: DownloadIcon }
	];

	const user = getUserCtx();

	const footerLinks = [{ label: 'Profile', href: '/settings/user', icon: UserIcon }];

	if (user.isOmni) {
		footerLinks.unshift({
			label: 'Server Settings',
			href: '/settings/glacier',
			icon: MountainSnow
		});
	}

	if (isFrost) {
		footerLinks.unshift({ label: 'Frost Settings', href: '/settings/frost', icon: Snowflake });
	}

	const auth = glacierPubCli(AuthService);

	const dm = getDialogCtx();

	async function logout() {
		const ok = await dm.confirm('Logout', 'Did you grab everything? Phone. Keys. Wallet.', 'info');
		if (!ok) return;

		await auth.logout({});
		await goto('/auth/login', { replaceState: true });
	}
</script>

<aside class="flex h-full w-16 flex-col items-center border-r border-border bg-surface py-6">
	<!-- Logo Area -->
	<div class="mb-8 flex items-center justify-center">
		<img src={'/favicon.svg'} alt="Logo" class="h-8 w-8 rounded-lg shadow-frost" />
	</div>

	<!-- Main Navigation -->
	<nav class="flex flex-1 flex-col gap-4">
		{#each links as link}
			{@const active = isActive(link.href)}
			<a
				href={link.href}
				title={link.label}
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
		{#each footerLinks as link}
			{@const active = isActive(link.href)}

			<a
				href={link.href}
				title={link.label}
				class="group relative flex h-10 w-10 items-center justify-center rounded-xl transition-all
            duration-300
            {active
					? 'bg-frost-500/10 text-frost-400 shadow-frost'
					: 'text-muted hover:bg-panel hover:text-frost-400'}"
			>
				<link.icon size={22} strokeWidth={2} />
			</a>
		{/each}

		<button
			onclick={logout}
			class="group relative flex h-10 w-10 items-center justify-center rounded-xl transition-all duration-300"
		>
			<LogOutIcon size={22} strokeWidth={2} />
		</button>
	</div>
</aside>
