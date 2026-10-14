<script lang="ts">
	import {
		RefreshCwIcon,
		Trash2Icon,
		AlertTriangleIcon,
		ChevronLeftIcon,
		ShieldCheckIcon,
		LoaderIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { frostCli, callRPC } from '#lib/api/api.js';
	import { FrostLibraryService } from '#lib/gen/frost_library/v1/frost_library_pb.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import type { Game } from '#lib/gen/library/v1/library_pb.js';

	let { game }: { game?: Game } = $props();

	const frostLib = frostCli(FrostLibraryService);
	const sm = getSnackbarCtx();

	let isProcessing = $state(false);
	let showDeleteConfirm = $state(false);

	async function handleRecheck() {
		if (!game!.ID || isProcessing) return;

		isProcessing = true;
		const { err } = await callRPC(() =>
			frostLib.download({
				gameId: BigInt(game!.ID),
				recheck: true
			})
		);

		if (err) {
			sm.push(`Recheck failed: ${err}`, 'error');
		} else {
			sm.push('Recheck process initiated', 'success');
		}
		isProcessing = false;
	}

	async function handleDeleteRedownload() {
		if (!game!.ID || isProcessing) return;

		isProcessing = true;
		const { err } = await callRPC(() =>
			frostLib.download({
				gameId: BigInt(game!.ID),
				recheck: true,
				force: true
			})
		);

		if (err) {
			sm.push(`Operation failed: ${err}`, 'error');
		} else {
			sm.push('Local files purged. Redownload started.', 'success');
		}

		showDeleteConfirm = false;
		isProcessing = false;
	}
</script>

<div class="space-y-8">
	<!-- Header Section -->
	<header class="flex flex-col gap-2 px-2">
		<div class="flex items-center gap-4">
			<div class="rounded-2xl border border-border bg-panel p-3 text-frost-400">
				<ShieldCheckIcon size={24} />
			</div>
			<div>
				<h1 class="text-2xl font-bold tracking-tight text-foreground">Manage local frost files</h1>
				<p class="text-sm font-medium text-muted">
					{game?.Meta?.Name || 'Unknown Title'}
				</p>
			</div>
		</div>
	</header>

	<!-- Operations Card -->
	<main class="grid grid-cols-1 gap-6 md:grid-cols-2">
		<!-- Recheck Card -->
		<div
			class="flex flex-col justify-between gap-6 rounded-3xl border border-border bg-surface p-6 shadow-sm"
		>
			<div class="space-y-2">
				<h3 class="flex items-center gap-2 font-bold text-foreground">
					<RefreshCwIcon size={18} class="text-frost-400" />
					Verify Integrity
				</h3>
				<p class="text-xs leading-relaxed text-muted">
					Scans local files against the manifest to ensure no data is corrupted or missing. This
					will pause any active downloads for this title.
				</p>
			</div>

			<button
				onclick={handleRecheck}
				disabled={isProcessing}
				class="btn-secondary flex w-full items-center justify-center gap-2 rounded-2xl py-3 text-sm font-bold text-foreground transition-all hover:border-frost-500/50 hover:bg-frost-500/5 disabled:opacity-50"
			>
				{#if isProcessing}
					<LoaderIcon size={16} class="animate-spin" />
				{:else}
					<RefreshCwIcon size={16} />
				{/if}
				Recheck Files
			</button>
		</div>

		<!-- Delete/Redownload Card -->
		<div
			class="flex flex-col justify-between gap-6 rounded-3xl border border-red-500/10 bg-surface p-6 shadow-sm"
		>
			<div class="space-y-2">
				<h3 class="flex items-center gap-2 font-bold text-red-400">
					<Trash2Icon size={18} />
					Danger Zone
				</h3>
				<p class="text-xs leading-relaxed text-muted">
					Completely removes all local data for this game and re-queues it for a fresh download. Use
					this if a recheck fails to fix issues.
				</p>
			</div>

			<button
				onclick={() => (showDeleteConfirm = true)}
				disabled={isProcessing}
				class="flex w-full items-center justify-center gap-2 rounded-full border border-red-500/20 bg-red-500/10 py-3 text-sm font-bold text-red-400 transition-all hover:bg-red-500 hover:text-white disabled:opacity-50"
			>
				<Trash2Icon size={16} />
				Delete and Redownload
			</button>
		</div>
	</main>
</div>

<!-- DELETE CONFIRMATION MODAL -->
{#if showDeleteConfirm}
	<div
		class="fixed inset-0 z-[200] flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<div
			class="absolute inset-0 bg-background/80 backdrop-blur-md"
			onclick={() => !isProcessing && (showDeleteConfirm = false)}
		></div>

		<div
			class="relative w-full max-w-sm space-y-6 rounded-[2rem] border border-red-500/20 bg-surface p-8 shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<div class="flex flex-col items-center gap-4 text-center">
				<div class="rounded-full border border-red-500/20 bg-red-500/10 p-4 text-red-400">
					<AlertTriangleIcon size={40} strokeWidth={1.5} />
				</div>
				<div class="space-y-2">
					<h2 class="text-xl font-bold text-foreground">Are you sure?</h2>
					<p class="text-xs leading-relaxed text-muted">
						This will permanently delete all local files for <span class="font-bold text-foreground"
							>{game?.Meta?.Name}</span
						>. A new download will start immediately after.
					</p>
				</div>
			</div>

			<div class="flex flex-col gap-2">
				<button
					onclick={handleDeleteRedownload}
					disabled={isProcessing}
					class="flex w-full items-center justify-center gap-2 rounded-full bg-red-500 py-3 text-sm font-bold text-white transition-all hover:bg-red-400"
				>
					{#if isProcessing}
						<LoaderIcon size={16} class="animate-spin" />
					{:else}
						Purge and Redownload
					{/if}
				</button>
				<button
					onclick={() => (showDeleteConfirm = false)}
					disabled={isProcessing}
					class="btn-secondary w-full rounded-2xl py-3 text-sm font-bold text-muted transition-all hover:text-foreground"
				>
					Cancel
				</button>
			</div>
		</div>
	</div>
{/if}
