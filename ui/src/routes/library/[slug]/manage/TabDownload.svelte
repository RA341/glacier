<script lang="ts">
	import { ShieldCheckIcon } from '@lucide/svelte';
	import { type Game, LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	let { game = $bindable() }: { game: Game; refresh: () => Promise<void> } = $props();

	const lib = glacierCli(LibraryService);
	const sm = getSnackbarCtx();

	async function handleRedownload() {
		if (!game?.ID) {
			return;
		}

		const { err } = await callRPC(() => lib.redownload({ gameId: game.ID }));
		if (err) {
			sm.push(`Could not redownload ${err}`, 'error');
		}
		// await refresh()
	}
</script>

<!-- Download Section -->
<section class="space-y-4">
	<h2 class="px-2 text-sm font-bold tracking-widest text-muted uppercase">Download</h2>
	<div
		class="grid grid-cols-1 gap-6 rounded-3xl border border-border bg-surface p-6 lg:grid-cols-12"
	>
		<!-- State Info -->
		<div class="grid grid-cols-2 gap-x-4 gap-y-6 border-r border-border pr-6 lg:col-span-5">
			<div>
				<button onclick={handleRedownload} class="p-1 text-frost-400 hover:text-foreground">
					Retry
				</button>
			</div>
			<div>
				<p class="text-[9px] font-bold text-muted uppercase">Client Type</p>
				<p class="text-sm font-bold">{game?.DownloadState?.Client}</p>
			</div>
			<div>
				<p class="text-[9px] font-bold text-muted uppercase">Download ID</p>
				<p class="truncate font-mono text-sm opacity-60">{game?.DownloadState?.DownloadId}</p>
			</div>
			<div>
				<p class="text-[9px] font-bold text-muted uppercase">State</p>
				<p class="text-sm font-bold text-frost-400">{game?.DownloadState?.State}</p>
			</div>
			<div>
				<p class="text-[9px] font-bold text-muted uppercase">Progress</p>
				<p class="text-sm font-bold">{game?.DownloadState?.Progress}</p>
			</div>

			<div>
				<p class="text-[9px] font-bold text-muted uppercase">Download Url</p>
				<p class="overflow-hidden text-sm font-bold text-frost-400">
					{game?.DownloadState?.DownloadUrl}
				</p>
			</div>
			<div>
				<p class="text-[9px] font-bold text-muted uppercase">Path</p>
				<p class="text-sm font-bold wrap-break-word">{game?.DownloadState?.DownloadPath}</p>
			</div>
		</div>

		<div class="col-span-2 pt-4">
			<p class="mb-2 flex items-center gap-1 text-[9px] font-bold text-muted uppercase">
				<ShieldCheckIcon size={12} />
				Virus Scan
			</p>
			<p class="text-xs font-medium text-green-400">TODO</p>
		</div>
	</div>
</section>
