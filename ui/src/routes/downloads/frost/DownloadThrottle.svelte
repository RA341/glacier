<script lang="ts">
	import { callRPC, frostCli } from '#lib/api/api.js';
	import { FrostLibraryService } from '#lib/gen/frost_library/v1/frost_library_pb.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { CheckIcon, GaugeIcon } from '@lucide/svelte';

	const frostLib = frostCli(FrostLibraryService);
	const sm = getSnackbarCtx();

	let presets = [1, 5, 10, 100, 500];
	let customValue = $state<number | null>(null);
	let activeLimit = $state<number | null>(null);
	let isUpdating = $state(false);

	async function setLimit(limit: number) {
		if (isUpdating) return;
		isUpdating = true;

		const { err } = await callRPC(() =>
			frostLib.throttleSpeed({
				limit: limit * 1024 * 1024
			})
		);

		if (err) {
			sm.push(`Could not update limit: ${err}`, 'error');
		} else {
			activeLimit = limit;
			sm.push(`Speed limit set to ${limit} MB/s`, 'success');
		}
		isUpdating = false;
	}

	function handleCustomSubmit(e: KeyboardEvent) {
		if (e.key === 'Enter' && customValue) {
			setLimit(customValue);
		}
	}
</script>

<div
	class="flex flex-col items-center gap-6 rounded-2xl border border-border bg-panel p-4 shadow-sm sm:flex-row"
>
	<!-- Icon & Label -->
	<div class="flex shrink-0 items-center gap-3">
		<div class="rounded-lg bg-frost-500/10 p-2 text-frost-400">
			<GaugeIcon size={18} />
		</div>
		<div class="flex flex-col">
			<span class="text-[10px] font-bold tracking-widest text-muted uppercase">Throttle</span>
			<span class="text-xs font-bold text-foreground">Speed Limit</span>
		</div>
	</div>

	<!-- Presets Section -->
	<div class="flex gap-1 rounded-xl border border-border bg-surface p-1">
		{#each presets as p}
			<button
				onclick={() => setLimit(p)}
				disabled={isUpdating}
				class="rounded-full px-4 py-1.5 text-[11px] font-bold transition-all
                {activeLimit === p
					? 'btn-frost text-foreground shadow-lg shadow-frost-500/20'
					: 'text-muted hover:bg-panel hover:text-foreground'}"
			>
				{p}M
			</button>
		{/each}

		<button
			onclick={() => setLimit(0)}
			disabled={isUpdating}
			class="rounded-full px-4 py-1.5 text-[11px] font-bold transition-all
            {activeLimit === 0
				? 'btn-frost text-foreground shadow-lg shadow-frost-500/20'
				: 'text-muted hover:bg-panel hover:text-amber-400'}"
			title="Remove Limit"
		>
			Unlimited
		</button>
	</div>

	<!-- Divider (Hidden on small mobile) -->
	<div class="hidden h-8 w-px bg-border sm:block"></div>

	<!-- Custom Input -->
	<div class="group relative min-w-35 flex-1">
		<input
			type="number"
			bind:value={customValue}
			onkeydown={handleCustomSubmit}
			placeholder="Custom (MB/s)"
			class="w-full rounded-xl border border-border bg-surface py-2 pr-10 pl-4 text-xs font-bold transition-all outline-none placeholder:text-muted/30 focus:border-frost-500"
		/>
		<button
			onclick={() => customValue && setLimit(customValue)}
			disabled={!customValue || isUpdating}
			class="btn-secondary absolute top-1.5 right-1.5 bottom-1.5 rounded-2xl px-2 text-muted transition-all hover:text-frost-400 disabled:opacity-0"
		>
			<CheckIcon size={14} />
		</button>
	</div>
</div>

<style>
	/* Remove arrows from number input */
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}

	input[type='number'] {
		-moz-appearance: textfield;
	}
</style>
