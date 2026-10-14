<script lang="ts">
	import {
		CircleAlert,
		GlobeIcon,
		PowerIcon,
		RefreshCwIcon,
		LoaderIcon,
		SaveIcon,
		SettingsIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';

	import { callRPC, Frost, frostCli } from '#lib/api/api.js';
	import { ConfigService } from '#lib/gen/config/v1/config_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import ConfigField from '#lib/components/ConfigField.svelte';
	import { extractValues } from '#lib/components/util.svelte.js';
	import { getRandomIntInclusive } from '#lib/api/byte-math.js';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';

	const sm = getSnackbarCtx();
	const libConf = frostCli(ConfigService);
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

	let isRestarting = $state(false);
	let isShuttingDown = $state(false);
	let restartStatus = $state('Waiting for server...');

	async function handleRestart() {
		try {
			isRestarting = true;
			restartStatus = 'Restarting client standby...';

			// Trigger restart
			await fetch(`${Frost.base}/restart`);

			// Give the server 2 seconds to actually shut down before we start pinging
			setTimeout(startPingLoop, 2000);
		} catch (e) {
			isRestarting = false;
			sm.push('Restart request failed', 'error');
		}
	}

	function startPingLoop() {
		restartStatus = 'Reconnecting to Glacier...';
		const messages = ['Test test', 'Finding my marbles'];

		const interval = setInterval(async () => {
			try {
				// We use a short timeout so the UI doesn't hang on a stuck socket
				const controller = new AbortController();
				const timeoutId = setTimeout(() => controller.abort(), 900);

				const res = await fetch(`${Frost.base}/hello`, { signal: controller.signal });

				if (res.ok) {
					clearInterval(interval);
					restartStatus = 'Connection established!';

					// Final delay for internal services to warm up
					setTimeout(() => {
						window.location.reload();
					}, 1000);
				}
			} catch (e) {
				// Fetch error is expected while server is down
				restartStatus = messages[getRandomIntInclusive(0, messages.length)];
			}
		}, 1000);
	}

	const dm = getDialogCtx();

	async function handleShutdown() {
		if (!(await dm.confirm('You sure about that', 'Just checkin', 'error'))) {
			return;
		}

		try {
			isShuttingDown = true;
			await fetch(`${Frost.base}/shutdown`);
			sm.push('Shutdown command sent', 'warn');
		} catch (e) {
			isShuttingDown = false;
			sm.push('Shutdown request failed', 'error');
		}
	}
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
		<CircleAlert size={48} class="text-red-400" />
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
			<div class="flex items-center gap-4">
				<div class="rounded-2xl border border-border bg-panel p-3 text-frost-400">
					<SettingsIcon size={24} />
				</div>
				<div>
					<h1 class="text-2xl leading-tight font-black tracking-tight text-foreground uppercase">
						Frost Settings
					</h1>
					<p class="text-xs font-medium tracking-wider text-muted uppercase">
						Configure the frost client
					</p>
				</div>
			</div>

			<div class="flex gap-2">
				<button
					onclick={handleRestart}
					disabled={isSaving || isRestarting}
					class="btn-secondary flex items-center gap-2 rounded-2xl px-5 py-2.5 text-xs font-bold text-muted transition-all hover:text-frost-400 active:scale-95"
				>
					<RefreshCwIcon size={16} class={isRestarting ? 'animate-spin' : ''} />
					Restart
				</button>

				<button
					onclick={handleShutdown}
					disabled={isSaving || isRestarting}
					class="btn-secondary flex items-center gap-2 rounded-2xl px-5 py-2.5 text-xs font-bold text-muted transition-all hover:text-red-400 active:scale-95"
				>
					<PowerIcon size={16} />
					Shutdown
				</button>

				<div class="mx-2 h-8 w-px bg-border"></div>

				<button
					onclick={handleSave}
					disabled={isSaving || isRestarting}
					class="btn-frost flex items-center gap-2 rounded-full px-8 py-2.5 text-sm font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-95 disabled:opacity-50"
				>
					{#if isSaving}
						<LoaderIcon size={18} class="animate-spin" />
					{:else}
						<SaveIcon size={18} />
						Save Changes
					{/if}
				</button>
			</div>
		</header>

		<div class="space-y-12 pb-20">
			<ConfigField schema={configSchema} />
		</div>
	</div>
{/if}

<!-- BLOCKING RESTART MODAL -->
{#if isRestarting}
	<div
		class="fixed inset-0 z-200 flex items-center justify-center p-6"
		transition:fade={{ duration: 200 }}
	>
		<div class="absolute inset-0 bg-background/90 backdrop-blur-xl"></div>

		<div
			class="relative flex w-full max-w-sm flex-col items-center gap-6 rounded-[2.5rem] border border-border bg-surface p-10 text-center shadow-2xl"
			transition:fly={{ y: 20, duration: 400 }}
		>
			<div class="relative">
				<div class="animate-pulse rounded-full bg-frost-500/10 p-5 text-frost-500">
					<RefreshCwIcon size={48} strokeWidth={1.5} class="animate-spin-slow" />
				</div>
				<div
					class="absolute -right-1 -bottom-1 rounded-full border border-border bg-background p-2 text-frost-400"
				>
					<GlobeIcon size={16} />
				</div>
			</div>

			<div class="space-y-2">
				<h2 class="text-xl font-bold text-foreground">Restarting Server</h2>
				<p class="text-sm font-medium text-muted">{restartStatus}</p>
			</div>

			<!-- INDETERMINATE PROGRESS BAR -->
			<div class="h-1.5 w-full overflow-hidden rounded-full border border-border bg-panel">
				<div class="animate-progress-loop h-full w-1/3 rounded-full bg-frost-500"></div>
			</div>

			<p class="text-[10px] font-bold tracking-[0.2em] text-muted uppercase opacity-50">
				So hows it going, hows the family
			</p>
		</div>
	</div>
{/if}

<style>
	:global(.animate-spin-slow) {
		animation: spin 3s linear infinite;
	}

	@keyframes progress-loop {
		0% {
			transform: translateX(-100%);
		}
		100% {
			transform: translateX(300%);
		}
	}

	.animate-progress-loop {
		animation: progress-loop 1.5s infinite linear;
	}
</style>
