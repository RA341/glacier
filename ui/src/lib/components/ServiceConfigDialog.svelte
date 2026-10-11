<script lang="ts">
	import {
		AlertTriangleIcon,
		ArrowLeftIcon,
		ChevronRightIcon,
		SaveIcon,
		XIcon
	} from '@lucide/svelte';

	import { fade, fly } from 'svelte/transition';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import {
		type ServiceConfig,
		ServiceConfigService
	} from '#lib/gen/service_config/v1/service_config_pb.js';
	import DynForm from '#lib/components/DynForm.svelte';

	interface Props {
		isOpen: boolean;
		serviceType?: string;
		editConf?: ServiceConfig | null;
	}

	let { isOpen = $bindable(), serviceType = '', editConf = null }: Props = $props();

	let scConfig = glacierCli(ServiceConfigService);

	let formData = $state<Record<string, any>>({});
	let enabled = $state(false);
	let name = $state('');

	$effect(() => {
		if (!isOpen) {
			return;
		}

		if (editConf) {
			let sd = new TextDecoder().decode(editConf.Config);
			formData = JSON.parse(sd);
			serviceType = editConf.ServiceType;
			selectedFlavour = editConf.Flavour;
			name = editConf.Name;
			enabled = editConf.Enabled;
			return;
		}

		// new client get choices
		supportedClientsRpc.runner();
	});

	let supportedClientsRpc = createRPCRunner(() =>
		scConfig.getSupportedValues({
			ServiceType: serviceType
		})
	);
	let configSchemaRpc = createRPCRunner(() =>
		scConfig.getSchema({
			ServiceType: serviceType,
			Flavour: selectedFlavour!
		})
	);

	let errorMessage = $state<string | null>(null);

	let selectedFlavour = $state<string | null>(null);
	$effect(() => {
		if (!selectedFlavour) {
			formData = {};
			return;
		}

		configSchemaRpc.runner();
	});

	async function handleSave() {
		console.log('Submitting:', formData);
		if (!name) {
			console.log('Name is required');
			return;
		}

		const { err } = await callRPC(async () => {
			let conf = {
				ID: BigInt(editConf?.ID ?? 0),
				Enabled: enabled,
				Name: name,
				ServiceType: serviceType,
				Flavour: selectedFlavour!,
				Config: new TextEncoder().encode(JSON.stringify(formData))
			};

			if (editConf) {
				await scConfig.edit({ conf });
			} else {
				await scConfig.new({ conf });
			}
		});
		if (err) {
			errorMessage = err;
			return;
		}

		handleClose();
	}

	function handleClose() {
		isOpen = false;
		selectedFlavour = null;
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-100 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={handleClose}></div>

		<div
			class="relative flex max-h-[90vh] min-h-125 w-full max-w-xl flex-col overflow-hidden rounded-3xl border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<!-- HEADER -->
			<div class="flex items-center justify-between border-b border-border bg-panel/30 p-6">
				<div class="flex items-center gap-3">
					{#if selectedFlavour}
						<button
							onclick={() => (selectedFlavour = null)}
							class="rounded-lg p-2 text-muted transition-colors hover:bg-panel"
						>
							<ArrowLeftIcon size={18} />
						</button>
					{/if}
					<h2 class="text-xl font-bold text-foreground">
						{selectedFlavour ? `Configure ${selectedFlavour}` : 'Select Client Type'}
					</h2>
				</div>
				<button onclick={handleClose} class="p-1 text-muted hover:text-foreground">
					<XIcon size={20} />
				</button>
			</div>

			<!-- BODY -->
			<div class="flex-1 overflow-y-auto p-8">
				{#if !editConf && !selectedFlavour}
					<!-- Type SELECTION -->
					<div class="grid grid-cols-2 gap-4">
						{#if supportedClientsRpc.error}
							<div>Error: {supportedClientsRpc.error}</div>
						{:else if supportedClientsRpc.loading}
							<div>loading</div>
						{:else if !supportedClientsRpc.value}
							<div>No supported flavours found</div>
						{:else}
							{#each supportedClientsRpc.value?.values ?? [] as client}
								<button
									onclick={() => (selectedFlavour = client)}
									class="group flex items-center justify-between rounded-2xl border border-border bg-panel p-4 transition-all hover:border-frost-500/50"
								>
									<span class="text-sm font-bold tracking-wider text-foreground uppercase"
										>{client}</span
									>
									<ChevronRightIcon size={18} class="text-muted/30 group-hover:text-frost-400" />
								</button>
							{/each}
						{/if}
					</div>
				{:else if configSchemaRpc.error}
					<div>Error: {configSchemaRpc.error}</div>
				{:else if configSchemaRpc.loading}
					<div>loading</div>
				{:else if !configSchemaRpc.value}
					<div>No supported flavours found</div>
				{:else}
					<div class="space-y-8" in:fly={{ x: 20, duration: 300 }}>
						{#if configSchemaRpc.value}
							<div class="space-y-8">
								<div class="space-y-3">
									<span class="text-sm text-foreground/80">Name</span>
									<input required bind:value={name} class="..." />

									<div class="flex items-center gap-3">
										<span class="text-sm text-foreground/80">Enable</span>
										<input type="checkbox" bind:checked={enabled} class="..." />
									</div>
								</div>

								<DynForm fields={configSchemaRpc.value.fields} bind:formData />
							</div>
						{/if}
					</div>
				{/if}
			</div>

			<!-- FOOTER -->
			<div class="flex justify-end gap-3 border-t border-border bg-panel/30 p-6">
				<button
					onclick={handleClose}
					class="rounded-xl border border-border px-6 py-2.5 text-sm font-bold text-muted"
					>Cancel
				</button>
				{#if selectedFlavour}
					<button
						onclick={handleSave}
						class="flex items-center gap-2 rounded-xl bg-frost-500 px-8 py-2.5 text-sm font-bold text-background hover:bg-frost-400"
					>
						<SaveIcon size={16} />
						Save Client
					</button>
				{/if}
			</div>
		</div>

		{#if errorMessage}
			<div
				class="fixed inset-0 z-110 flex items-center justify-center p-6"
				transition:fade={{ duration: 100 }}
			>
				<!-- Darker backdrop for the error alert -->
				<!-- svelte-ignore a11y_click_events_have_key_events -->
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					class="absolute inset-0 bg-black/40 backdrop-blur-sm"
					onclick={() => (errorMessage = null)}
				></div>

				<div
					class="relative flex w-full max-w-sm flex-col gap-4 rounded-2xl border border-red-500/30 bg-surface p-6 shadow-2xl"
					transition:fly={{ y: 10, duration: 200 }}
				>
					<div class="flex items-center gap-3 text-red-400">
						<div class="rounded-lg bg-red-500/10 p-2">
							<AlertTriangleIcon size={20} />
						</div>
						<h3 class="font-bold">Action Failed</h3>
					</div>

					<p class="text-sm leading-relaxed text-muted">
						{errorMessage}
					</p>

					<div class="flex justify-end pt-2">
						<button
							onclick={() => (errorMessage = null)}
							class="rounded-xl border border-border bg-panel px-5 py-2 text-sm font-bold transition-all hover:text-foreground"
						>
							Dismiss
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
{/if}
