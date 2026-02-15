<script lang="ts">
    import {
        ExternalLinkIcon,
        Video,
        ImageIcon,
        LinkIcon,
        LoaderIcon,
        PlusIcon,
        Trash2Icon,
        UploadIcon,
        XIcon,
        EyeIcon,
        ServerIcon,
    } from "@lucide/svelte";
    import {fade, fly} from 'svelte/transition';
    import type {Asset} from "$lib/gen/search/v1/search_pb";
    import {getAssetPath, uploadAsset} from "$lib/api/assets";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {Glacier} from "$lib/api/api";
    import {onMount} from "svelte";
    import {splitCamelCase, trimPrefix} from "$lib/api/strings";
    import AssetImg from "$lib/components/assets/AssetImg.svelte";
    import AssetVideo from "$lib/components/assets/AssetVideo.svelte";

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
    let uploadType = $state("poster");

    let currentAsset = $derived(assets?.[selectedIndex] || null);

    const base = `${Glacier.base}/assets`

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
            sm.push("Asset processed successfully", 'success');
            isUploadOpen = false;
            await refresh();
        }

        isUploading = false;
    }

    let assetTypes = $state<string[]>([])

    async function getAssetTypes(): Promise<void> {
        const final = `${base}/types`
        try {
            const resp = await fetch(final)
            if (resp.ok) {
                assetTypes = JSON.parse(await resp.text()) as string[];
            } else {
                sm.push(`Failed to load asset types: ${resp.statusText}`, 'error');
            }
        } catch (e) {
            sm.push(`Failed to load asset types: ${e}`, 'error');
        }
    }

    async function deleteAssetApi(assetId: bigint): Promise<void> {
        const final = `${base}/delete/${assetId}`
        try {
            const resp = await fetch(final)
            if (resp.ok) {
                sm.push("Asset deleted successfully", 'success');
            } else {
                sm.push(`Failed to delete asset: ${resp.statusText}`, 'error');
            }
        } catch (e) {
            sm.push(`Could not delete asset: ${e}`, 'error');
        }
    }

    onMount(async () => {
        await getAssetTypes()
        uploadType = assetTypes.at(1) ?? ""
    })

    async function deleteAsset(assetId: bigint, index: number) {
        if (confirm("Remove this asset from metadata?")) {
            await deleteAssetApi(assetId)
            await refresh()
        }
    }

    // Reset tab when switching assets
    $effect(() => {
        selectedIndex;
        activeTab = 'preview';
    });
</script>

<div class="bg-surface border border-border rounded-3xl overflow-hidden shadow-sm flex h-137.5">
    <!-- SIDEBAR (Left) -->
    <div class="w-64 border-r border-border flex flex-col bg-panel/20">
        <div class="p-4 border-b border-border bg-surface">
            <h2 class="text-[10px] font-black text-frost-400 uppercase tracking-[0.2em] flex items-center gap-2">
                <ImageIcon size={12}/>
                Media Assets
            </h2>
        </div>

        <div class="p-2 border-t border-border bg-surface">
            <button
                    onclick={openAdd}
                    class="w-full py-2 bg-panel border border-border rounded-xl text-[10px] font-bold uppercase tracking-widest text-muted hover:text-frost-400 hover:border-frost-500/50 transition-all flex items-center justify-center gap-2"
            >
                <PlusIcon size={14}/>
                Add Asset
            </button>
        </div>

        <div class="flex-1 overflow-y-auto p-2 space-y-1 custom-scrollbar">
            {#each assets as asset, i}
                {@const assetTy = trimPrefix(asset.Type, "Asset")}

                <button
                        onclick={() => selectedIndex = i}
                        class="w-full flex items-center gap-3 p-3 rounded-xl transition-all text-left group
                        {selectedIndex === i ? 'bg-frost-500/10 border-frost-500/20 text-frost-400' : 'hover:bg-panel text-muted'}"
                >
                    <div class="shrink-0">
                        {#if asset.Type.toLowerCase().includes('video') || asset.Type.toLowerCase().includes('trailer')}
                            <Video size={16}/>
                        {:else}
                            <ImageIcon size={16}/>
                        {/if}
                    </div>
                    <div class="flex-1 min-w-0">
                        <p class="text-xs font-bold truncate capitalize">{asset.LocalPath}</p>
                        <span class="px-2 py-0.5 rounded bg-panel border border-border text-[9px] font-black uppercase text-muted">
                            {assetTy}
                        </span>
                    </div>
                </button>
            {/each}
        </div>
    </div>

    <!-- PREVIEW & EDIT AREA (Right) -->
    <div class="flex-1 bg-background/50 flex flex-col min-w-0">
        {#if currentAsset}
            <!-- Toolbar -->
            <div class="p-3 border-b border-border flex justify-between items-center gap-4 bg-surface flex-wrap">
                <!-- Asset Type Picker -->
                <div class="flex items-center gap-2 flex-1 min-w-0">
                    <span class="text-[9px] font-bold text-muted uppercase tracking-widest shrink-0">Type</span>
                    <select
                            bind:value={currentAsset.Type}
                            class="bg-panel border border-border rounded-xl py-1.5 px-3 outline-none focus:border-frost-500 text-xs font-bold appearance-none truncate max-w-48"
                    >
                        {#each assetTypes as assetT}
                            <option value={assetT}>{splitCamelCase(trimPrefix(assetT, "Asset"))}</option>
                        {/each}
                    </select>
                </div>

                <!-- Action Buttons -->
                <div class="flex gap-2 shrink-0">
                    <button
                            onclick={() => openEdit(currentAsset)}
                            class="p-2 bg-panel border border-border rounded-lg text-muted hover:text-frost-400 transition-all"
                            title="Upload Replacement"
                    >
                        <UploadIcon size={16}/>
                    </button>
                    <button
                            onclick={() => deleteAsset(currentAsset.ID, selectedIndex)}
                            class="p-2 bg-panel border border-border rounded-lg text-muted hover:text-red-400 transition-all"
                            title="Delete Asset"
                    >
                        <Trash2Icon size={16}/>
                    </button>
                </div>
            </div>

            <!-- Tabs -->
            <div class="flex border-b border-border bg-surface">
                <button
                        onclick={() => activeTab = 'preview'}
                        class="flex items-center gap-2 px-5 py-2.5 text-[10px] font-black uppercase tracking-widest border-b-2 transition-all
                        {activeTab === 'preview' ? 'border-frost-500 text-frost-400' : 'border-transparent text-muted hover:text-foreground'}"
                >
                    <EyeIcon size={12}/>
                    Preview
                </button>
                <button
                        onclick={() => activeTab = 'remote'}
                        class="flex items-center gap-2 px-5 py-2.5 text-[10px] font-black uppercase tracking-widest border-b-2 transition-all
                        {activeTab === 'remote' ? 'border-frost-500 text-frost-400' : 'border-transparent text-muted hover:text-foreground'}"
                >
                    <LinkIcon size={12}/>
                    Remote URL
                </button>
            </div>

            {#if activeTab === 'preview'}
                <!-- Visual Preview -->
                <div class="flex-1 p-6 flex flex-col items-center justify-center overflow-hidden bg-[radial-gradient(circle_at_center,var(--tw-gradient-stops))] from-panel/20 via-transparent to-transparent gap-4 min-h-0">
                    {#if currentAsset.LocalPath || currentAsset.RemoteUrl}
                        {@const assetT = currentAsset.Type.toLowerCase()}
                        {@const isVideo = assetT.includes('video') || assetT.includes('trailer')}
                        {@const localUrl = getAssetPath({gameId: gameId, assetPath: currentAsset.LocalPath})}

                        <!-- Media Preview -->
                        <div class="relative group w-full flex-1 min-h-0 shadow-2xl rounded-2xl overflow-hidden border border-border bg-panel">
                            {#if isVideo}
                                <AssetVideo
                                        src={localUrl}
                                        autoplay={true}
                                        muted={true}
                                        class="w-full h-full"
                                >
                                    {#snippet loadingSlot()}
                                        <div class="flex h-full w-full items-center justify-center bg-black">
                                            <LoaderIcon class="animate-spin text-white/50"/>
                                        </div>
                                    {/snippet}

                                    {#snippet errorSlot()}
                                        <div class="flex flex-col h-full w-full items-center justify-center bg-surface text-muted-foreground">
                                            <PlusIcon size={48} class="rotate-45 opacity-50"/>
                                            <p class="text-xs font-bold uppercase tracking-widest">Video Unavailable</p>
                                        </div>
                                    {/snippet}
                                </AssetVideo>
                            {:else}
                                <AssetImg src={localUrl} class="w-full h-full">
                                    {#snippet loadingSlot()}
                                        <div class="flex h-full w-full items-center justify-center bg-surface animate-pulse">
                                            <LoaderIcon class="animate-spin text-muted"/>
                                        </div>
                                    {/snippet}

                                    {#snippet errorSlot()}
                                        <div class="flex h-full w-full items-center justify-center bg-surface text-muted/20">
                                            <ImageIcon size={48}/>
                                        </div>
                                    {/snippet}
                                </AssetImg>
                            {/if}

                            <div class="absolute inset-0 bg-black/60 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-4">
                                <a href={localUrl} target="_blank"
                                   class="p-3 bg-panel border border-border text-foreground rounded-full hover:text-frost-400 transition-colors">
                                    <ExternalLinkIcon size={20}/>
                                </a>
                            </div>
                        </div>

                        <!-- Local Path (below preview) -->
                        <div class="w-full space-y-1">
                            <p class="text-[9px] font-bold text-muted uppercase tracking-widest ml-1 opacity-50 flex items-center gap-1">
                                <ServerIcon size={10}/>
                                Local Server Path
                            </p>
                            <p class="text-xs font-mono text-muted/80 truncate px-3 py-2 bg-background/50 rounded-xl border border-border/50 italic">
                                {currentAsset.LocalPath || 'No local file uploaded'}
                            </p>
                        </div>
                    {:else}
                        <div class="text-center space-y-2 opacity-10">
                            <ImageIcon size={64} strokeWidth={1}/>
                            <p class="text-xs font-black uppercase tracking-widest">No Content Found</p>
                        </div>
                    {/if}
                </div>

            {:else if activeTab === 'remote'}
                <!-- Remote URL Tab -->
                <div class="flex-1 p-6 flex flex-col justify-start gap-4">
                    <div class="space-y-2">
                        <p class="text-[9px] font-bold text-muted uppercase tracking-widest ml-1 opacity-50">Remote
                            Source URL</p>

                        {#if !currentAsset.RemoteUrl}
                            <div class="flex items-center gap-3 px-4 py-3 bg-panel/50 border border-border/50 rounded-xl text-muted/40 italic">
                                <LinkIcon size={14}/>
                                <span class="text-xs">No remote URL found</span>
                            </div>
                        {/if}

                        <div class="relative group/input">
                            <div class="absolute left-3 top-1/2 -translate-y-1/2 text-muted/30 group-focus-within/input:text-frost-500 transition-colors">
                                <LinkIcon size={14}/>
                            </div>
                            <input
                                    type="text"
                                    bind:value={currentAsset.RemoteUrl}
                                    placeholder="https://..."
                                    class="w-full bg-background border border-border rounded-xl py-2.5 pl-10 pr-4 text-xs font-mono outline-none focus:border-frost-500 transition-all shadow-inner"
                            />
                        </div>
                        <p class="text-[8px] text-muted/40 ml-1">
                            Changes to URL will be staged until you Save the Game configuration.
                        </p>
                    </div>
                </div>
            {/if}

        {:else}
            <div class="flex-1 flex flex-col items-center justify-center text-muted/20 gap-4">
                <ImageIcon size={80} strokeWidth={1}/>
                <p class="text-sm font-black uppercase tracking-[0.3em]">Select an Asset</p>
            </div>
        {/if}
    </div>
</div>

<!-- UPLOAD DIALOG -->
{#if isUploadOpen}
    <div class="fixed inset-0 z-100 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md"
             onclick={() => !isUploading && (isUploadOpen = false)}></div>

        <div
                class="relative w-full max-w-md bg-surface border border-border rounded-[2rem] shadow-2xl overflow-hidden"
                transition:fly={{ y: 20, duration: 300 }}
        >
            <div class="p-6 border-b border-border bg-panel/30 flex justify-between items-center">
                <h3 class="font-bold text-foreground flex items-center gap-2">
                    <UploadIcon size={18} class="text-frost-400"/>
                    {editingAssetId ? 'Replace Asset Content' : 'Upload New Asset'}
                </h3>
                <button onclick={() => isUploadOpen = false} class="text-muted hover:text-foreground">
                    <XIcon size={20}/>
                </button>
            </div>

            <div class="p-8 space-y-6">
                <div class="space-y-2">
                    <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">Asset Type</label>
                    <select bind:value={uploadType}
                            class="w-full bg-panel border border-border rounded-xl py-3 px-4 outline-none focus:border-frost-500 text-sm font-bold appearance-none"
                    >
                        {#each assetTypes as assetT}
                            <option value={assetT}>{splitCamelCase(trimPrefix(assetT, "Asset"))}</option>
                        {/each}
                    </select>
                </div>

                <div class="space-y-2">
                    <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">File Source</label>
                    <input
                            type="file"
                            onchange={(e) => uploadFile = e.currentTarget.files?.[0] || null}
                            class="w-full text-xs text-muted file:mr-4 file:py-2 file:px-4 file:rounded-xl file:border-0 file:text-[10px] file:font-black file:uppercase file:bg-frost-500 file:text-background hover:file:bg-frost-400 file:cursor-pointer"
                    />
                </div>

                {#if editingAssetId}
                    <div class="p-3 bg-amber-500/5 border border-amber-500/20 rounded-xl flex gap-3 items-center">
                        <LoaderIcon size={14} class="text-amber-500"/>
                        <p class="text-[10px] text-amber-200/60 leading-tight">
                            Editing Asset ID: {editingAssetId}. This will replace the existing file on the server.
                        </p>
                    </div>
                {/if}
            </div>

            <div class="p-6 border-t border-border bg-panel/30 flex gap-3">
                <button
                        onclick={() => isUploadOpen = false}
                        class="flex-1 py-3 bg-panel border border-border rounded-xl text-xs font-bold text-muted"
                >
                    Cancel
                </button>
                <button
                        onclick={handleUpload}
                        disabled={!uploadFile || isUploading}
                        class="flex-2 py-3 bg-frost-500 text-background rounded-xl text-xs font-black uppercase tracking-widest flex items-center justify-center gap-2 disabled:opacity-50"
                >
                    {#if isUploading}
                        <LoaderIcon size={16} class="animate-spin"/>
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
