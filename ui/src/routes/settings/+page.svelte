<script lang="ts">
	import { ArrowUpRightIcon, MountainSnow, SettingsIcon, Snowflake } from '@lucide/svelte';
	import { isFrost } from '#lib/api/api.js';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';

	const user = getUserCtx();
	const configs = [
		{
			label: 'Glacier',
			href: '/settings/glacier',
			icon: MountainSnow,
			description: 'Manage server settings, users, metadata, indexers, and downloads.',
			visible: () => user.isOmni
		},
		{
			label: 'Frost',
			href: '/settings/frost',
			icon: Snowflake,
			description: 'Configure the Frost client and manage its connection and services.',
			visible: () => isFrost
		}
	];
	let visibleConfigs = $derived(configs.filter((config) => config.visible()));
</script>

<div class="mx-auto w-full max-w-5xl space-y-8 p-4 sm:p-8">
	<header class="flex items-center gap-4 border-b border-border pb-6">
		<div class="rounded-2xl border border-frost-500/20 bg-frost-500/10 p-3 text-frost-400">
			<SettingsIcon size={28} />
		</div>
		<div>
			<h1 class="text-3xl font-bold tracking-tight text-foreground">Config</h1>
			<p class="mt-1 text-sm text-muted">Choose what you want to configure.</p>
		</div>
	</header>

	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
		{#each visibleConfigs as config}
			<a
				href={config.href}
				class="group flex flex-col gap-5 rounded-2xl border border-border bg-surface p-6 transition-colors hover:border-frost-500/40 hover:bg-panel focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-frost-400"
			>
				<div class="flex items-center justify-between">
					<div class="rounded-xl bg-frost-500/10 p-3 text-frost-400">
						<config.icon size={28} />
					</div>
					<ArrowUpRightIcon
						size={20}
						class="text-muted transition-colors group-hover:text-frost-400"
					/>
				</div>
				<div>
					<h2 class="text-xl font-bold text-foreground">{config.label}</h2>
					<p class="mt-2 text-sm leading-relaxed text-muted">{config.description}</p>
				</div>
			</a>
		{:else}
			<p class="col-span-full rounded-2xl border border-border bg-surface p-6 text-sm text-muted">
				No configuration options are available for your account.
			</p>
		{/each}
	</div>
</div>
