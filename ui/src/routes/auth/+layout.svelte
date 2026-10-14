<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { fade, fly } from 'svelte/transition';
	import {
		CheckCircle2Icon,
		CircleAlert,
		Globe,
		GlobeIcon,
		LoaderIcon,
		LockIcon,
		PencilIcon,
		RefreshCwIcon,
		UserPlusIcon
	} from '@lucide/svelte';
	import { appName, callRPC, frostCli, isFrost } from '#lib/api/api.js';
	import { onMount } from 'svelte';
	import { ConfigService } from '#lib/gen/config/v1/config_pb.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';

	let { children } = $props();

	const sm = getSnackbarCtx();
	const confSrv = frostCli(ConfigService);

	let activeTab = $derived(page.url.pathname.includes('register') ? 'register' : 'login');

	function navigate(path: string) {
		goto(`/auth/${path}`, { replaceState: true });
	}

	let glacierUrlInput = $state('');
	let isSettingUrl = $state(false);
	let forceEdit = $state(false);
	let glacierUrlCheckRpc = createRPCRunner(() => confSrv.getField({}));
	let isGlacierReady = $derived(
		!isFrost || (!!glacierUrlCheckRpc?.value?.Value && !glacierUrlCheckRpc?.value?.Error)
	);
	let showEditor = $derived(isFrost && (!isGlacierReady || forceEdit));

	onMount(() => {
		if (isFrost) {
			glacierUrlCheckRpc.runner();
		}
	});

	const dm = getDialogCtx();

	async function handleSetUrl() {
		if (!glacierUrlInput) return;
		isSettingUrl = true;
		forceEdit = false;

		console.log('Setting Glacier URL to:', glacierUrlInput);
		const { err } = await callRPC(() => confSrv.setField({ Value: glacierUrlInput }));
		if (err) {
			await dm.alert('Could not connect to glacier', err, 'error');
		} else {
			await glacierUrlCheckRpc.runner();
		}

		isSettingUrl = false;
	}
</script>

<div
	class="relative flex min-h-screen w-full items-center justify-center overflow-hidden bg-background p-6"
>
	<div
		class="absolute top-[-10%] left-[-10%] h-[40%] w-[40%] rounded-full bg-frost-500/10 blur-[120px]"
	></div>
	<div
		class="absolute right-[-10%] bottom-[-10%] h-[40%] w-[40%] rounded-full bg-frost-900/20 blur-[120px]"
	></div>

	<div class="z-10 w-full max-w-110" transition:fly={{ y: 20, duration: 600 }}>
		<div class="mb-8 flex flex-col items-center gap-3">
			<img src="/favicon.svg" alt="logo" height="60" width="60" class="drop-shadow-2xl" />
			<h1 class="headline-xler text-foreground uppercase">{appName}</h1>
			<p class="text-center text-xs font-bold tracking-[0.2em] text-frost-50/50 uppercase">
				<span class="font-black text-frost-500">COLD</span> Storage for your games
			</p>
		</div>

		{#if showEditor}
			{#if isFrost && glacierUrlCheckRpc.loading && !glacierUrlCheckRpc.value}
				<div class="flex flex-col items-center justify-center gap-4 py-20 text-muted" in:fade>
					<LoaderIcon class="animate-spin text-frost-500" size={32} />
					<p class="animate-pulse text-[10px] font-bold tracking-widest uppercase">
						Checking Instance Status...
					</p>
				</div>
			{:else}
				<!-- GLACIER URL ENTRY VIEW -->
				<div
					class="space-y-4 rounded-4xl border border-border bg-surface p-8 shadow-2xl"
					in:fly={{ y: 10 }}
				>
					<div class="space-y-2">
						<h2 class="text-xl font-bold text-foreground">Connect Glacier</h2>
						<p class="text-sm text-muted">Enter url to a Glacier instance</p>
					</div>

					{#if glacierUrlCheckRpc.error}
						<div
							class="flex items-start gap-3 rounded-2xl border border-red-500/20 bg-red-500/5 p-4 text-red-400"
						>
							<CircleAlert size={18} class="shrink-0" />
							<div class="text-xs">
								<p class="font-bold">Connection Failed</p>
								<p class="opacity-80">{glacierUrlCheckRpc.error}</p>
							</div>
						</div>
					{/if}

					<div class="space-y-4">
						<div class="space-y-2">
							<div class="relative">
								<GlobeIcon
									size={18}
									class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50"
								/>
								<input
									type="url"
									bind:value={glacierUrlInput}
									placeholder="glacier.yourdomain.com"
									class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
									onkeydown={(e) => {
										if (e.key === 'Enter' && glacierUrlInput && !isSettingUrl) {
											handleSetUrl();
										}
									}}
								/>
							</div>
						</div>

						<button
							onclick={handleSetUrl}
							disabled={isSettingUrl || !glacierUrlInput}
							class="btn-frost flex w-full items-center justify-center gap-2 rounded-full py-4 font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-[0.98] disabled:opacity-50"
						>
							{#if isSettingUrl}
								<LoaderIcon size={20} class="animate-spin" />
								Connecting...
							{:else}
								<Globe size={18} />
								Connect
							{/if}
						</button>
					</div>
				</div>
			{/if}
		{:else}
			{#if isFrost && glacierUrlCheckRpc.value}
				<div class="mb-6 flex items-center justify-center" in:fade>
					<div
						class="flex items-center gap-3 rounded-full border border-frost-500/10 bg-frost-500/5 px-4 py-2 shadow-sm"
					>
						<CheckCircle2Icon size={14} class="text-green-500" />
						<span class="max-w-50 truncate font-mono text-[10px] text-frost-400/80">
							{glacierUrlCheckRpc.value.Value}
						</span>
						<div class="mx-1 h-3 w-px bg-border"></div>
						<button
							onclick={() => {
								glacierUrlInput = glacierUrlCheckRpc?.value?.Value ?? '';
								forceEdit = true;
							}}
							class="flex items-center gap-1 text-[10px] font-black tracking-widest text-muted uppercase transition-colors hover:text-frost-400"
						>
							<PencilIcon size={10} />
							Edit
						</button>
					</div>
				</div>
			{/if}

			<div
				class="flex flex-col overflow-hidden rounded-4xl border border-border bg-surface shadow-2xl"
				in:fade
			>
				<div class="flex gap-1 border-b border-border bg-panel/30 p-2">
					<button
						onclick={() => navigate('login')}
						class="flex flex-1 items-center justify-center gap-2 rounded-2xl py-3 text-sm font-bold transition-all
                        {activeTab === 'login'
							? 'border border-border bg-surface text-frost-400 shadow-sm'
							: 'text-muted hover:bg-panel hover:text-foreground'}"
					>
						<LockIcon size={16} />
						Login
					</button>
					<button
						onclick={() => navigate('register')}
						class="flex flex-1 items-center justify-center gap-2 rounded-2xl py-3 text-sm font-bold transition-all
                        {activeTab === 'register'
							? 'border border-border bg-surface text-frost-400 shadow-sm'
							: 'text-muted hover:bg-panel hover:text-foreground'}"
					>
						<UserPlusIcon size={16} />
						Register
					</button>
				</div>
				<main class="p-8">
					{@render children()}
				</main>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		background-color: #0a0a0c;
	}
</style>
