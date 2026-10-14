<script lang="ts">
	import { AlertCircleIcon, LoaderIcon, SaveIcon } from '@lucide/svelte';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { ConfigService } from '#lib/gen/config/v1/config_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import ConfigField from '#lib/components/ConfigField.svelte';
	import { extractValues } from '#lib/components/util.svelte.js';

	const sm = getSnackbarCtx();
	const libConf = glacierCli(ConfigService);
	let configSchema = $state({});

	let isSaving = $state(false);

	async function handleSave() {
		isSaving = true;

		const payload = extractValues(configSchema);

		const { err } = await callRPC(() =>
			libConf.set({
				configSchema: JSON.stringify(payload)
			})
		);

		if (err) {
			sm.push(`Could not update config: ${err}`, 'error');
		}

		await configRpc.runner();

		isSaving = false;
	}

	const configRpc = createRPCRunner(() => libConf.get({}));

	onMount(() => {
		configRpc.runner();
	});

	$effect(() => {
		if (configRpc.value) {
			configSchema = JSON.parse(configRpc.value.configSchema);
		}
	});
</script>

{#if configRpc.loading && !configRpc.value}
	<!-- THEMED LOADING STATE -->
	<div class="flex h-[60vh] flex-col items-center justify-center gap-4 text-muted">
		<LoaderIcon class="animate-spin text-frost-500" size={40} />
		<p class="animate-pulse text-sm font-bold tracking-widest uppercase">
			Loading Configuration...
		</p>
	</div>
{:else if configRpc.error}
	<!-- THEMED ERROR STATE -->
	<div
		class="mx-auto mt-20 flex max-w-md flex-col items-center gap-4 rounded-3xl border border-red-500/20 bg-red-500/5 p-6 text-center"
	>
		<AlertCircleIcon size={48} class="text-red-400" />
		<div>
			<h3 class="font-bold text-foreground">Failed to load config</h3>
			<p class="mt-1 text-sm text-muted">{configRpc.error}</p>
		</div>
		<button
			onclick={() => configRpc.runner()}
			class="btn-secondary rounded-2xl px-6 py-2 text-xs font-bold transition-all hover:bg-surface"
		>
			Retry Connection
		</button>
	</div>
{:else}
	<div class="mx-auto min-h-screen max-w-4xl p-8">
		<!-- STICKY HEADER -->
		<header
			class="sticky top-0 z-30 -mx-8 mb-10 flex items-center justify-between border-b border-border bg-background/80 px-8 py-6 backdrop-blur-md"
		>
			<button
				onclick={handleSave}
				disabled={isSaving}
				class="btn-frost flex items-center gap-2 rounded-full px-8 py-3 text-sm font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-95 disabled:opacity-50"
			>
				{#if isSaving}
					<LoaderIcon size={18} class="animate-spin" />
					Saving...
				{:else}
					<SaveIcon size={18} />
					Save Changes
				{/if}
			</button>
		</header>

		<div class="space-y-12 pb-20">
			<ConfigField schema={configSchema} />
		</div>
	</div>
{/if}
