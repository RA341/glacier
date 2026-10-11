<script lang="ts">
	import {
		AlertTriangleIcon,
		CornerLeftUpIcon,
		FileIcon,
		FolderIcon,
		HardDriveIcon,
		LoaderIcon,
		RefreshCwIcon,
		SearchXIcon,
		Trash2Icon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { type Game, LibraryService } from '#lib/gen/library/v1/library_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { formatBytes } from '#lib/api/byte-math.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	let { game }: { game: Game } = $props();

	const libSrv = glacierCli(LibraryService);

	let base = $state('');
	let downloaded = $state(false);
	let pendingDeletePath = $state<string | null>(null);

	let fileRpc = createRPCRunner(() =>
		libSrv.listFiles({
			GameId: game.ID,
			BasePath: base,
			Downloaded: downloaded
		})
	);

	$effect(() => {
		fileRpc.runner();
	});

	function navigateTo(dirName: string) {
		base = base === '' ? dirName : `${base}/${dirName}`;
	}

	function navigateUp() {
		const parts = base.split('/');
		parts.pop();
		base = parts.join('/');
	}

	const breadcrumbs = $derived(base === '' ? [] : base.split('/'));
	const sm = getSnackbarCtx();

	async function confirmDelete() {
		if (!pendingDeletePath) return;

		const { err } = await callRPC(() =>
			libSrv.deleteFile({
				Path: `${base}/${pendingDeletePath}`,
				Downloaded: downloaded,
				GameId: game.ID
			})
		);

		if (err) {
			sm.push(`Error deleting file: ${err}`, 'error');
		} else {
			sm.push('File deleted successfully', 'success');
		}

		pendingDeletePath = null;
		await fileRpc.runner();
	}

	async function markAsExe(RelPath: string) {
		const { err } = await callRPC(() =>
			libSrv.edit({
				game: {
					ID: game.ID,
					file: {
						Exe: RelPath
					}
				}
			})
		);

		if (err) {
			sm.push(`Error marking exe: ${err}`, 'error');
		}

		await fileRpc.runner();
	}
</script>

<div class="relative flex h-full flex-col overflow-hidden">
	<!-- Compact Toolbar -->
	<div class="flex items-center justify-between border-b border-border bg-panel/50 p-2">
		<div class="flex gap-1 rounded-lg border border-border bg-surface p-0.5">
			<button
				onclick={() => (downloaded = false)}
				class="rounded-md px-3 py-1 text-[10px] font-bold transition-all {!downloaded
					? 'bg-panel text-frost-400 shadow-sm'
					: 'text-muted hover:text-foreground'}"
			>
				Complete
			</button>
			<button
				onclick={() => (downloaded = true)}
				class="rounded-md px-3 py-1 text-[10px] font-bold transition-all {downloaded
					? 'bg-panel text-frost-400 shadow-sm'
					: 'text-muted hover:text-foreground'}"
			>
				Incomplete
			</button>
		</div>

		<button
			onclick={() => fileRpc.runner()}
			class="p-1.5 text-muted transition-colors hover:text-frost-400"
		>
			<RefreshCwIcon size={14} class={fileRpc.loading ? 'animate-spin' : ''} />
		</button>
	</div>

	<!-- Mini Breadcrumbs -->
	<div
		class="no-scrollbar flex items-center gap-1 overflow-x-auto border-b border-border bg-panel/20 px-3 py-1.5 text-[10px] whitespace-nowrap"
	>
		<button onclick={() => (base = '')} class="text-muted hover:text-frost-400">
			<HardDriveIcon size={12} />
		</button>
		{#if base !== ''}<span class="text-muted/30">/</span>{/if}
		{#each breadcrumbs as part, i}
			<button
				onclick={() => (base = breadcrumbs.slice(0, i + 1).join('/'))}
				class="font-bold text-muted hover:text-foreground">{part}</button
			>
			{#if i < breadcrumbs.length - 1}<span class="text-muted/30">/</span>{/if}
		{/each}
	</div>

	<!-- Scrollable File Area -->
	<div class="custom-scrollbar flex-1 overflow-y-auto bg-black/5">
		{#if fileRpc.loading && !fileRpc.value}
			<div class="flex h-full items-center justify-center gap-2 text-muted">
				<LoaderIcon class="animate-spin text-frost-500" size={14} />
				<span class="text-[10px] font-bold tracking-widest uppercase">Loading...</span>
			</div>
		{:else if fileRpc.error}
			<div class="p-4 text-center text-[10px] font-bold text-red-400 uppercase">
				{fileRpc.error}
			</div>
		{:else}
			<div class="flex flex-col">
				{#if base !== ''}
					<button
						onclick={navigateUp}
						class="flex items-center gap-2 border-b border-border/30 px-4 py-2 text-[11px] font-bold text-frost-400 hover:bg-frost-500/5"
					>
						<CornerLeftUpIcon size={14} />
						..
					</button>
				{/if}

				{#if fileRpc.value?.files && fileRpc.value.files.length > 0}
					{#each fileRpc.value.files as file}
						{@const isSelected = game.file?.Exe === file.RelPath}
						<div
							class="group flex items-center justify-between border-b border-border/20 px-4 py-2 transition-colors last:border-0
                            {isSelected
								? 'border-l-2 border-l-frost-400 bg-frost-500/10 hover:bg-frost-500/15'
								: 'hover:bg-panel/50'}
                        "
						>
							<button
								onclick={() => (file.IsDir ? navigateTo(file.RelPath.split('/').pop()!) : null)}
								class="flex min-w-0 flex-1 items-center gap-2 text-left"
							>
								<span
									class="shrink-0 {isSelected
										? 'text-frost-400'
										: 'text-muted group-hover:text-frost-400'}"
								>
									{#if file.IsDir}<FolderIcon size={14} />{:else}<FileIcon size={14} />{/if}
								</span>
								<span
									class="truncate text-[11px] {file.IsDir
										? 'font-bold text-foreground'
										: isSelected
											? 'text-frost-300'
											: 'text-muted'}"
								>
									{file.RelPath.split('/').pop()}
								</span>
								<span
									class="truncate text-[11px] {file.IsDir
										? 'font-bold text-foreground'
										: isSelected
											? 'text-frost-300'
											: 'text-muted'}"
								>
									{isSelected ? 'EXE' : ''}
								</span>
							</button>

							{#if !isSelected}
								<button
									onclick={() => markAsExe(file.RelPath)}
									class="rounded-lg p-1.5 text-muted opacity-0 transition-all group-hover:opacity-100 hover:text-frost-500"
									title="Mark Exe"
								>
									Mark as Exe
								</button>
							{/if}

							<span class="mr-2 shrink-0 font-mono text-[9px] text-muted/40"
								>{formatBytes(file.Size)}</span
							>

							<!-- Trash Icon Button -->
							<button
								onclick={() => (pendingDeletePath = file.RelPath)}
								class="rounded-lg p-1.5 text-muted opacity-0 transition-all group-hover:opacity-100 hover:bg-red-400/10 hover:text-red-400"
								title="Delete File"
							>
								<Trash2Icon size={14} />
							</button>
						</div>
					{/each}
				{:else}
					<div class="flex flex-col items-center justify-center gap-1 py-8 text-muted/20">
						<SearchXIcon size={24} />
						<span class="text-[10px] font-bold uppercase">Empty</span>
					</div>
				{/if}
			</div>
		{/if}
	</div>

	<!-- DELETE CONFIRMATION DIALOG -->
	{#if pendingDeletePath}
		<div
			class="absolute inset-0 z-50 flex items-center justify-center p-4"
			transition:fade={{ duration: 100 }}
		>
			<!-- Local Backdrop -->
			<!-- svelte-ignore a11y_click_events_have_key_events -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div
				class="absolute inset-0 bg-background/60 backdrop-blur-sm"
				onclick={() => (pendingDeletePath = null)}
			></div>

			<div
				class="relative w-full max-w-xs space-y-4 rounded-2xl border border-red-500/20 bg-surface p-5 shadow-2xl"
				transition:fly={{ y: 10, duration: 200 }}
			>
				<div class="flex items-center gap-3 text-red-400">
					<div class="rounded-lg bg-red-500/10 p-2">
						<AlertTriangleIcon size={18} />
					</div>
					<span class="text-sm font-bold">Confirm Deletion</span>
				</div>

				<p class="text-xs leading-relaxed text-muted">
					Are you sure you want to delete <span class="font-mono text-foreground"
						>{pendingDeletePath}</span
					>? This action cannot be undone.
				</p>

				<div class="flex justify-end gap-2">
					<button
						onclick={() => (pendingDeletePath = null)}
						class="rounded-lg border border-border bg-panel px-3 py-1.5 text-[11px] font-bold text-muted transition-colors hover:text-foreground"
					>
						Cancel
					</button>
					<button
						onclick={confirmDelete}
						class="rounded-lg bg-red-500 px-4 py-1.5 text-[11px] font-bold text-white shadow-lg shadow-red-500/20 transition-colors hover:bg-red-400"
					>
						Delete
					</button>
				</div>
			</div>
		</div>
	{/if}
</div>

<style>
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}

	.custom-scrollbar::-webkit-scrollbar {
		width: 4px;
	}

	.custom-scrollbar::-webkit-scrollbar-track {
		background: transparent;
	}

	.custom-scrollbar::-webkit-scrollbar-thumb {
		background: rgba(255, 255, 255, 0.05);
		border-radius: 10px;
	}
</style>
