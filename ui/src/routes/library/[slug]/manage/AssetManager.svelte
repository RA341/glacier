<script lang="ts">
	import {
		Save,
		ExternalLinkIcon,
		EyeIcon,
		ImageIcon,
		LinkIcon,
		LoaderIcon,
		PlusIcon,
		ServerIcon,
		Trash2Icon,
		UploadIcon,
		Video,
		XIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { getAssetPath, uploadAsset } from '#lib/api/assets.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { callRPC, Glacier, glacierCli } from '#lib/api/api.js';
	import { onMount } from 'svelte';
	import { splitCamelCase, trimPrefix } from '#lib/api/strings.js';
	import AssetImg from '#lib/components/assets/AssetImg.svelte';
	import AssetVideo from '#lib/components/assets/AssetVideo.svelte';
	import { type Asset, AssetService } from '#lib/gen/assets/v1/assets_pb.js';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';

	let {
		assets = $bindable(),
		gameId,
		refresh
	}: {
		gameId?: bigint;
		assets: Asset[];
		refresh: () => Promise<void>;
	} = $props();

	const sm = getSnackbarCtx();

	let selectedIndex = $state(0);
	let isUploadOpen = $state(false);
	let isUploading = $state(false);
	let editingAssetId = $state<bigint | null>(null);
	let activeTab = $state<'preview' | 'remote'>('preview');

	let uploadFile = $state<File | null>(null);
	let uploadType = $state('poster');

	let currentAsset = $derived(assets?.[selectedIndex] ?? null);

	const base = `${Glacier.base}/assets`;

	function openAdd() {
		editingAssetId = null;
		uploadFile = null;
		isUploadOpen = true;
	}

	function openEdit(asset: Asset) {
		editingAssetId = asset.ID;
		uploadType = asset.Type;
		uploadFile = null;
		isUploadOpen = true;
	}

	async function handleUpload() {
		if (!uploadFile || !gameId) return;

		isUploading = true;
		const err = await uploadAsset({
			file: uploadFile,
			type: uploadType,
			id: editingAssetId ?? BigInt(0),
			gameId: gameId
		});

		if (err) {
			sm.push(`Upload failed: ${err}`, 'error');
		} else {
			sm.push('Asset processed successfully', 'success');
			isUploadOpen = false;
			await refresh();
		}

		isUploading = false;
	}

	let assetTypes = $state<string[]>([]);

	async function getAssetTypes(): Promise<void> {
		const final = `${base}/types`;
		try {
			const resp = await fetch(final);
			if (resp.ok) {
				assetTypes = JSON.parse(await resp.text()) as string[];
			} else {
				sm.push(`Failed to load asset types: ${resp.statusText}`, 'error');
			}
		} catch (e) {
			sm.push(`Failed to load asset types: ${e}`, 'error');
		}
	}

	onMount(async () => {
		await getAssetTypes();
		uploadType = assetTypes.at(1) ?? '';
	});

	const dm = getDialogCtx();

	async function deleteAsset() {
		if (!currentAsset) {
			return;
		}

		const ok = await dm.confirm('Are you sure you want to delete this asset ?', 'warning');
		if (!ok) {
			return;
		}

		const { err } = await callRPC(() => assetSrv.delete({ ID: currentAsset.ID }));
		if (err) {
			sm.push(`Could not update asset: ${err}`, 'error');
		}

		await refresh();
	}

	const assetSrv = glacierCli(AssetService);

	async function editAsset() {
		if (!currentAsset) {
			return;
		}

		const { err } = await callRPC(() => assetSrv.edit({ asset: currentAsset }));
		if (err) {
			sm.push(`Could not update asset: ${err}`, 'error');
		}
	}

	// Reset tab when switching assets
	$effect(() => {
		selectedIndex;
		activeTab = 'preview';
	});
</script>

<div class="flex h-137.5 overflow-hidden rounded-3xl border border-border bg-surface shadow-sm">
	<!-- SIDEBAR (Left) -->
	<div class="flex w-64 flex-col border-r border-border bg-panel/20">
		<div class="border-b border-border bg-surface p-4">
			<h2
				class="flex items-center gap-2 text-[10px] font-black tracking-[0.2em] text-frost-400 uppercase"
			>
				<ImageIcon size={12} />
				Media Assets
			</h2>
		</div>

		<div class="border-t border-border bg-surface p-2">
			<button
				onclick={openAdd}
				class="flex w-full items-center justify-center gap-2 rounded-xl border border-border bg-panel py-2 text-[10px] font-bold tracking-widest text-muted uppercase transition-all hover:border-frost-500/50 hover:text-frost-400"
			>
				<PlusIcon size={14} />
				Add Asset
			</button>
		</div>

		<div class="custom-scrollbar flex-1 space-y-1 overflow-y-auto p-2">
			{#each assets as asset, i}
				{@const assetTy = trimPrefix(asset.Type, 'Asset')}

				<button
					onclick={() => (selectedIndex = i)}
					class="group flex w-full items-center gap-3 rounded-xl p-3 text-left transition-all
                        {selectedIndex === i
						? 'border-frost-500/20 bg-frost-500/10 text-frost-400'
						: 'text-muted hover:bg-panel'}"
				>
					<div class="shrink-0">
						{#if asset.Type.toLowerCase().includes('video') || asset.Type.toLowerCase().includes('trailer')}
							<Video size={16} />
						{:else}
							<ImageIcon size={16} />
						{/if}
					</div>
					<div class="min-w-0 flex-1">
						<p class="truncate text-xs font-bold capitalize">{asset.LocalPath}</p>
						<span
							class="rounded border border-border bg-panel px-2 py-0.5 text-[9px] font-black text-muted uppercase"
						>
							{assetTy}
						</span>
					</div>
				</button>
			{/each}
		</div>
	</div>

	<!-- PREVIEW & EDIT AREA (Right) -->
	<div class="flex min-w-0 flex-1 flex-col bg-background/50">
		{#if currentAsset}
			<div class="flex flex-nowrap items-center gap-4 border-b border-border bg-surface p-3">
				<div class="flex min-w-0 flex-1 items-center gap-2">
					<button
						onclick={editAsset}
						class="rounded-lg border border-border bg-panel p-2 text-muted transition-all hover:text-red-400"
						title="Delete Asset"
					>
						<Save size={16} />
					</button>
					<button
						onclick={() => openEdit(currentAsset)}
						class="rounded-lg border border-border bg-panel p-2 text-muted transition-all hover:text-frost-400"
						title="Upload Replacement"
					>
						<UploadIcon size={16} />
					</button>
					<button
						onclick={() => deleteAsset(currentAsset.ID, selectedIndex)}
						class="rounded-lg border border-border bg-panel p-2 text-muted transition-all hover:text-red-400"
						title="Delete Asset"
					>
						<Trash2Icon size={16} />
					</button>

					<span class="ml-5 shrink-0 text-[9px] font-bold tracking-widest text-muted uppercase"
						>Type</span
					>
					<select
						bind:value={currentAsset.Type}
						class="max-w-48 min-w-0 flex-1 appearance-none truncate rounded-xl border border-border bg-panel px-3 py-1.5 text-xs font-bold outline-none focus:border-frost-500"
					>
						{#each assetTypes as assetT}
							<option value={assetT}>{splitCamelCase(trimPrefix(assetT, 'Asset'))}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Tabs -->
			<div class="flex border-b border-border bg-surface">
				<button
					onclick={() => (activeTab = 'preview')}
					class="flex items-center gap-2 border-b-2 px-5 py-2.5 text-[10px] font-black tracking-widest uppercase transition-all
                        {activeTab === 'preview'
						? 'border-frost-500 text-frost-400'
						: 'border-transparent text-muted hover:text-foreground'}"
				>
					<EyeIcon size={12} />
					Preview
				</button>
				<button
					onclick={() => (activeTab = 'remote')}
					class="flex items-center gap-2 border-b-2 px-5 py-2.5 text-[10px] font-black tracking-widest uppercase transition-all
                        {activeTab === 'remote'
						? 'border-frost-500 text-frost-400'
						: 'border-transparent text-muted hover:text-foreground'}"
				>
					<LinkIcon size={12} />
					Remote URL
				</button>
			</div>

			{#if activeTab === 'preview'}
				<!-- Visual Preview -->
				<div
					class="flex min-h-0 flex-1 flex-col items-center justify-center gap-4 overflow-hidden bg-[radial-gradient(circle_at_center,var(--tw-gradient-stops))] from-panel/20 via-transparent to-transparent p-6"
				>
					{#if currentAsset.LocalPath || currentAsset.RemoteUrl}
						{@const assetT = currentAsset.Type.toLowerCase()}
						{@const isVideo = assetT.includes('video') || assetT.includes('trailer')}
						{@const localUrl = getAssetPath({ gameId: gameId, assetPath: currentAsset.LocalPath })}

						<!-- Media Preview -->
						<div
							class="group relative min-h-0 w-full flex-1 overflow-hidden rounded-2xl border border-border bg-panel shadow-2xl"
						>
							{#if isVideo}
								<AssetVideo src={localUrl} autoplay={true} muted={true} class="h-full w-full">
									{#snippet loadingSlot()}
										<div class="flex h-full w-full items-center justify-center bg-black">
											<LoaderIcon class="animate-spin text-white/50" />
										</div>
									{/snippet}

									{#snippet errorSlot()}
										<div
											class="text-muted-foreground flex h-full w-full flex-col items-center justify-center bg-surface"
										>
											<PlusIcon size={48} class="rotate-45 opacity-50" />
											<p class="text-xs font-bold tracking-widest uppercase">Video Unavailable</p>
										</div>
									{/snippet}
								</AssetVideo>
							{:else}
								<AssetImg src={localUrl} class="h-full w-full">
									{#snippet loadingSlot()}
										<div
											class="flex h-full w-full animate-pulse items-center justify-center bg-surface"
										>
											<LoaderIcon class="animate-spin text-muted" />
										</div>
									{/snippet}

									{#snippet errorSlot()}
										<div
											class="flex h-full w-full items-center justify-center bg-surface text-muted/20"
										>
											<ImageIcon size={48} />
										</div>
									{/snippet}
								</AssetImg>
							{/if}
						</div>

						<!-- Local Path (below preview) -->
						<div class="w-full space-y-1">
							<!-- Label -->
							<p
								class="ml-1 flex items-center gap-1 text-[9px] font-bold tracking-widest text-muted uppercase opacity-50"
							>
								<ServerIcon size={10} />
								Cached Path
							</p>

							<!-- Inline pill -->
							<div
								class="flex items-center gap-2 rounded-xl border border-border/50 bg-background/50 px-3 py-2"
							>
								<ServerIcon size={12} class="shrink-0 text-muted/60" />
								<p class="min-w-0 flex-1 truncate font-mono text-xs text-muted/80 italic">
									{currentAsset.LocalPath || 'No local file uploaded'}
								</p>
								<a
									href={localUrl}
									target="_blank"
									class="shrink-0 text-muted transition-colors hover:text-frost-400"
									title="Open in browser"
								>
									<ExternalLinkIcon size={12} />
								</a>
							</div>
						</div>
					{:else}
						<div class="space-y-2 text-center opacity-10">
							<ImageIcon size={64} strokeWidth={1} />
							<p class="text-xs font-black tracking-widest uppercase">No Content Found</p>
						</div>
					{/if}
				</div>
			{:else if activeTab === 'remote'}
				<!-- Remote URL Tab -->
				<div class="flex flex-1 flex-col justify-start gap-4 p-6">
					<div class="space-y-2">
						<p class="ml-1 text-[9px] font-bold tracking-widest text-muted uppercase opacity-50">
							Remote Source URL
						</p>

						{#if !currentAsset.RemoteUrl}
							<div
								class="flex items-center gap-3 rounded-xl border border-border/50 bg-panel/50 px-4 py-3 text-muted/40 italic"
							>
								<LinkIcon size={14} />
								<span class="text-xs">No remote URL found</span>
							</div>
						{/if}

						<div class="group/input relative">
							<div
								class="absolute top-1/2 left-3 -translate-y-1/2 text-muted/30 transition-colors group-focus-within/input:text-frost-500"
							>
								<LinkIcon size={14} />
							</div>
							<input
								type="text"
								bind:value={currentAsset.RemoteUrl}
								placeholder="https://..."
								class="w-full rounded-xl border border-border bg-background py-2.5 pr-4 pl-10 font-mono text-xs shadow-inner transition-all outline-none focus:border-frost-500"
							/>
						</div>
						<p class="ml-1 text-[8px] text-muted/40">
							Changes to URL will be staged until you Save the Game configuration.
						</p>
					</div>
				</div>
			{/if}
		{:else}
			<div class="flex flex-1 flex-col items-center justify-center gap-4 text-muted/20">
				<ImageIcon size={80} strokeWidth={1} />
				<p class="text-sm font-black tracking-[0.3em] uppercase">Select an Asset</p>
			</div>
		{/if}
	</div>
</div>

<!-- UPLOAD DIALOG -->
{#if isUploadOpen}
	<div
		class="fixed inset-0 z-100 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<div
			class="absolute inset-0 bg-background/80 backdrop-blur-md"
			onclick={() => !isUploading && (isUploadOpen = false)}
		></div>

		<div
			class="relative w-full max-w-md overflow-hidden rounded-[2rem] border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<div class="flex items-center justify-between border-b border-border bg-panel/30 p-6">
				<h3 class="flex items-center gap-2 font-bold text-foreground">
					<UploadIcon size={18} class="text-frost-400" />
					{editingAssetId ? 'Replace Asset Content' : 'Upload New Asset'}
				</h3>
				<button onclick={() => (isUploadOpen = false)} class="text-muted hover:text-foreground">
					<XIcon size={20} />
				</button>
			</div>

			<div class="space-y-6 p-8">
				<div class="space-y-2">
					<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
						>Asset Type</label
					>
					<select
						bind:value={uploadType}
						class="w-full appearance-none rounded-xl border border-border bg-panel px-4 py-3 text-sm font-bold outline-none focus:border-frost-500"
					>
						{#each assetTypes as assetT}
							<option value={assetT}>{splitCamelCase(trimPrefix(assetT, 'Asset'))}</option>
						{/each}
					</select>
				</div>

				<div class="space-y-2">
					<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
						>File Source</label
					>
					<input
						type="file"
						onchange={(e) => (uploadFile = e.currentTarget.files?.[0] || null)}
						class="w-full text-xs text-muted file:mr-4 file:cursor-pointer file:rounded-xl file:border-0 file:bg-frost-500 file:px-4 file:py-2 file:text-[10px] file:font-black file:text-background file:uppercase hover:file:bg-frost-400"
					/>
				</div>

				{#if editingAssetId}
					<div
						class="flex items-center gap-3 rounded-xl border border-amber-500/20 bg-amber-500/5 p-3"
					>
						<LoaderIcon size={14} class="text-amber-500" />
						<p class="text-[10px] leading-tight text-amber-200/60">
							Editing Asset ID: {editingAssetId}. This will replace the existing file on the server.
						</p>
					</div>
				{/if}
			</div>

			<div class="flex gap-3 border-t border-border bg-panel/30 p-6">
				<button
					onclick={() => (isUploadOpen = false)}
					class="flex-1 rounded-xl border border-border bg-panel py-3 text-xs font-bold text-muted"
				>
					Cancel
				</button>
				<button
					onclick={handleUpload}
					disabled={!uploadFile || isUploading}
					class="flex flex-2 items-center justify-center gap-2 rounded-xl bg-frost-500 py-3 text-xs font-black tracking-widest text-background uppercase disabled:opacity-50"
				>
					{#if isUploading}
						<LoaderIcon size={16} class="animate-spin" />
						Uploading...
					{:else}
						Start Upload
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
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
