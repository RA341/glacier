<script lang="ts">
	import { frostCli, glacierCli, isFrost } from '#lib/api/api.js';
	import { ConfigService } from '#lib/gen/config/v1/config_pb.js';
	import type { Client } from '@connectrpc/connect';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import {
		FolderIcon,
		ChevronRightIcon,
		ArrowLeftIcon,
		XIcon,
		CheckIcon,
		SearchXIcon,
		LoaderIcon,
		FolderOpenIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';

	interface FolderPickerProps {
		value?: string;
		placeholder?: string;
		label?: string;
		disabled?: boolean;
	}

	let {
		value = $bindable(''),
		placeholder = 'Select a folder...',
		label = '',
		disabled = false
	}: FolderPickerProps = $props();

	const configSrv = isFrost ? frostCli(ConfigService) : glacierCli(ConfigService);

	let isOpen = $state(false);
	let base = $state('');
	const fileListRpc = createRPCRunner(() => configSrv.listFiles({ base: base }));

	$effect(() => {
		if (isOpen) {
			fileListRpc.runner();
		}
	});

	function openPicker() {
		if (disabled) return;
		base = value || '';
		isOpen = true;
	}

	const breadcrumbs = $derived(base === '/' || base === '' ? [] : base.split('/').filter(Boolean));

	function navigateTo(absolutePath: string) {
		base = absolutePath;
	}

	function navigateUp() {
		if (base === '/' || base === '') return;
		const parts = base.split('/').filter(Boolean);
		parts.pop();
		base = '/' + parts.join('/');
	}

	function jumpToBreadcrumb(index: number) {
		const parts = base.split('/').filter(Boolean);
		base = '/' + parts.slice(0, index + 1).join('/');
	}

	function selectFolder() {
		value = base || '/';
		isOpen = false;
	}

	let scrollContainer = $state<HTMLDivElement>();
	$effect(() => {
		if (breadcrumbs && scrollContainer) {
			setTimeout(() => {
				scrollContainer!.scrollTo({
					left: scrollContainer!.scrollWidth,
					behavior: 'smooth'
				});
			}, 0);
		}
	});
</script>

<div class="flex w-full flex-col gap-1.5">
	{#if label}
		<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">
			{label}
		</label>
	{/if}

	<div class="group relative flex items-center">
		<input
			type="text"
			bind:value
			{placeholder}
			class="min-w-0 flex-1 rounded-l-xl border border-border bg-panel px-4 py-2.5 text-sm text-foreground transition-all outline-none focus:border-frost-500"
		/>

		<button
			{disabled}
			type="button"
			onclick={openPicker}
			class="flex items-center justify-center rounded-r-xl border border-l-0 border-border bg-panel px-4 py-2.5 text-muted transition-all hover:bg-surface hover:text-frost-400 active:scale-95 disabled:opacity-50"
			title="Browse folder"
		>
			<FolderOpenIcon size={18} />
		</button>
	</div>
</div>

<!-- FOLDER PICKER DIALOG -->
{#if isOpen}
	<div
		class="fixed inset-0 z-200 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<div
			class="absolute inset-0 bg-background/80 backdrop-blur-md"
			onclick={() => (isOpen = false)}
		></div>

		<div
			class="relative flex h-125 w-full max-w-lg flex-col overflow-hidden rounded-4xl border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<div
				bind:this={scrollContainer}
				class="no-scrollbar flex items-center gap-1 overflow-x-auto scroll-smooth rounded-xl border border-border bg-background p-2 text-xs whitespace-nowrap"
			>
				<button
					onclick={() => (base = '/')}
					class="shrink-0 rounded px-2 py-1 font-medium text-muted transition-colors hover:bg-panel hover:text-frost-400"
				>
					/
				</button>

				{#each breadcrumbs as part, i}
					<span class="shrink-0 text-muted/30">/</span>
					<button
						onclick={() => jumpToBreadcrumb(i)}
						class="shrink-0 rounded px-2 py-1 font-bold text-muted transition-colors hover:bg-panel hover:text-foreground"
					>
						{part}
					</button>
				{/each}
			</div>

			<!-- Folder List Area -->
			<div class="flex-1 overflow-y-auto bg-panel/10 p-2">
				{#if fileListRpc.loading}
					<div class="flex h-full flex-col items-center justify-center gap-3 text-muted">
						<LoaderIcon class="animate-spin text-frost-500" size={24} />
						<span class="text-[10px] font-bold tracking-widest uppercase opacity-50"
							>Scanning Storage...</span
						>
					</div>
				{:else if fileListRpc.error}
					<div
						class="flex h-full flex-col items-center justify-center gap-2 p-8 text-center text-red-400"
					>
						<SearchXIcon size={32} />
						<p class="w-full truncate text-xs font-bold uppercase">{fileListRpc.error}</p>
					</div>
				{:else}
					<div class="flex flex-col gap-1">
						{#if base !== '/' && base !== ''}
							<button
								onclick={navigateUp}
								class="flex items-center gap-3 rounded-full p-3 text-left text-sm font-bold text-frost-400 transition-colors hover:bg-frost-500/5"
							>
								<ArrowLeftIcon size={18} />
								<span>.. (Parent Directory)</span>
							</button>
						{/if}

						{#if fileListRpc.value?.files && fileListRpc.value.files.length > 0}
							{#each fileListRpc.value.files as folder}
								<button
									onclick={() => navigateTo(folder)}
									class="group flex items-center justify-between rounded-full p-3 text-left transition-all hover:bg-panel"
								>
									<div class="flex min-w-0 items-center gap-3">
										<FolderIcon
											size={18}
											class="shrink-0 text-muted group-hover:text-frost-400"
											fill="currentColor"
											fill-opacity="0.1"
										/>
										<span
											class="truncate text-sm font-medium text-foreground/80 group-hover:text-foreground"
										>
											{folder.split('/').pop() || folder}
										</span>
									</div>
									<ChevronRightIcon size={16} class="text-muted/20 group-hover:text-frost-400" />
								</button>
							{/each}
						{:else}
							<div class="flex flex-col items-center justify-center gap-2 py-20 text-muted/20">
								<FolderIcon size={48} strokeWidth={1} />
								<span class="text-[10px] font-bold tracking-widest uppercase"
									>No accessible subfolders</span
								>
							</div>
						{/if}
					</div>
				{/if}
			</div>

			<!-- Footer -->
			<div class="flex items-center justify-between gap-4 border-t border-border bg-panel/30 p-4">
				<div class="min-w-0 flex-1">
					<p class="mb-1 ml-1 text-[9px] font-bold tracking-widest text-muted uppercase opacity-50">
						Target Directory
					</p>
					<p
						class="truncate rounded-lg border border-border bg-background p-2 font-mono text-xs text-frost-400"
					>
						{base || '/'}
					</p>
				</div>
				<button
					onclick={selectFolder}
					class="btn-frost flex shrink-0 items-center gap-2 rounded-full px-6 py-3 text-sm font-bold text-foreground transition-all hover:bg-frost-500 hover:text-background active:scale-95"
				>
					<CheckIcon size={18} />
					Select
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	.no-scrollbar::-webkit-scrollbar {
		display: none;
	}

	.no-scrollbar {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
</style>
