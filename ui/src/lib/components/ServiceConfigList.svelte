<script lang="ts">
	import ServiceConfigDialog from '#lib/components/ServiceConfigDialog.svelte';
	import {
		CircleAlertIcon,
		PenLineIcon,
		InboxIcon,
		LoaderIcon,
		PlusIcon,
		RefreshCcw,
		Settings2Icon,
		TrashIcon
	} from '@lucide/svelte';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import {
		type ServiceConfig,
		ServiceConfigService
	} from '#lib/gen/service_config/v1/service_config_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';

	let { ServiceType }: { ServiceType: string } = $props();
	let scConfig = glacierCli(ServiceConfigService);
	let isOpen = $state(false);

	let listRpc = createRPCRunner(() => scConfig.list({ serviceType: ServiceType }));

	function refresh() {
		listRpc.runner();
	}

	$effect(() => {
		if (!isOpen) {
			refresh();
		}
	});

	onMount(() => {
		refresh();
	});

	let editConf = $state<ServiceConfig | null>(null);

	function openEdit(config: ServiceConfig) {
		editConf = config;
		isOpen = true;
	}

	async function deleteClient(config: ServiceConfig) {
		const { err } = await callRPC(() => scConfig.delete({ id: config.ID }));
		if (err) {
			console.error(err);
		}

		refresh();
	}
</script>

<div class="space-y-6">
	<!-- TOOLBAR -->
	<div class="flex items-center justify-between px-2">
		<h2 class="text-sm font-bold tracking-widest text-muted uppercase">
			{ServiceType} Configurations
		</h2>
		<div class="flex items-center gap-3">
			<button
				onclick={refresh}
				class="btn-secondary rounded-2xl p-2.5 text-muted transition-all hover:text-frost-400 active:scale-95"
				title="Refresh List"
			>
				<RefreshCcw size={18} class={listRpc.loading ? 'animate-spin' : ''} />
			</button>

			<button
				onclick={() => (isOpen = true)}
				class="btn-frost flex items-center gap-2 rounded-full px-6 py-2.5 text-sm font-bold text-foreground transition-all hover:bg-frost-500 hover:text-background active:scale-95"
			>
				<PlusIcon size={18} />
				Add {ServiceType}
			</button>
		</div>
	</div>

	<!-- MAIN LIST -->
	<div class="min-h-100">
		{#if listRpc.error}
			<div
				class="flex h-64 flex-col items-center justify-center gap-3 rounded-3xl border border-red-500/10 bg-red-500/5 text-red-400"
			>
				<CircleAlertIcon size={32} />
				<p class="text-sm font-medium">{listRpc.error}</p>
			</div>
		{:else if listRpc.loading && !listRpc.value}
			<div class="flex h-64 flex-col items-center justify-center gap-4 text-muted">
				<LoaderIcon class="animate-spin text-frost-500" size={32} />
				<p class="animate-pulse text-sm">Fetching configurations...</p>
			</div>
		{:else if !listRpc.value?.conf || listRpc.value.conf.length === 0}
			<div
				class="flex h-64 flex-col items-center justify-center rounded-3xl border-2 border-dashed border-border text-muted/30"
			>
				<InboxIcon size={48} strokeWidth={1} class="mb-2" />
				<p class="text-sm font-medium">No {ServiceType} configurations found</p>
			</div>
		{:else}
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<div class="grid grid-cols-[repeat(auto-fill,minmax(320px,1fr))] gap-4">
					{#each listRpc.value.conf as client (client.ID)}
						<div
							class="group relative flex flex-col rounded-3xl border border-border bg-panel/30 p-5 shadow-sm transition-all hover:border-frost-500/30 hover:bg-panel/50"
						>
							<!-- Top Section: Status and Actions -->
							<div class="mb-4 flex items-start justify-between">
								<div class="relative">
									<div
										class="rounded-2xl border border-border bg-surface p-3 text-muted transition-colors group-hover:text-frost-400"
									>
										<Settings2Icon size={24} />
									</div>
									<div
										class="absolute -top-1 -right-1 h-3.5 w-3.5 rounded-full border-2 border-background {client.Enabled
											? 'bg-green-500'
											: 'bg-red-500'}"
									></div>
								</div>

								<div class="flex gap-1 opacity-0 transition-opacity group-hover:opacity-100">
									<button
										onclick={() => openEdit(client)}
										class="rounded-full p-2 text-muted transition-all hover:bg-frost-500/10 hover:text-frost-400"
									>
										<PenLineIcon size={16} />
									</button>
									<button
										onclick={() => deleteClient(client)}
										class="rounded-full p-2 text-muted transition-all hover:bg-red-500/10 hover:text-red-400"
									>
										<TrashIcon size={16} />
									</button>
								</div>
							</div>

							<!-- Bottom Section: Info -->
							<div class="space-y-1">
								<h3 class="truncate text-lg font-bold text-foreground">{client.Name}</h3>
								<div class="flex items-center gap-2">
									<span
										class="rounded-md border border-border bg-panel px-2 py-0.5 text-[10px] font-bold tracking-widest text-muted uppercase"
									>
										{client.Flavour}
									</span>
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{/if}
	</div>
</div>

<ServiceConfigDialog {editConf} serviceType={ServiceType} bind:isOpen />

<style>
	/* Subtle pulsing for the list while refreshing */
	.animate-spin {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		from {
			transform: rotate(0deg);
		}
		to {
			transform: rotate(360deg);
		}
	}
</style>
