<script lang="ts">
    import {frostCli, glacierCli, isFrost} from "$lib/api/api";
    import {ConfigService} from "$lib/gen/config/v1/config_pb";
    import type {Client} from "@connectrpc/connect";
    import {createRPCRunner} from "$lib/api/rpc.svelte";
    import {
        FolderIcon,
        ChevronRightIcon,
        ArrowLeftIcon,
        XIcon,
        CheckIcon,
        SearchXIcon,
        LoaderIcon,
        FolderOpenIcon
    } from "@lucide/svelte";
    import {fade, fly} from "svelte/transition";

    interface FolderPickerProps {
        value?: string;
        placeholder?: string;
        label?: string;
        disabled?: boolean;
    }

    let {
        value = $bindable(""),
        placeholder = "Select a folder...",
        label = "",
        disabled = false
    }: FolderPickerProps = $props();

    const configSrv = isFrost ? frostCli(ConfigService) : glacierCli(ConfigService);

    let isOpen = $state(false);
    let base = $state("");
    const fileListRpc = createRPCRunner(() => configSrv.listFiles({base: base}));

    $effect(() => {
        if (isOpen) {
            fileListRpc.runner();
        }
    });

    function openPicker() {
        if (disabled) return;
        base = value || "";
        isOpen = true;
    }

    const breadcrumbs = $derived(
        base === "/" || base === ""
            ? []
            : base.split('/').filter(Boolean)
    );

    function navigateTo(absolutePath: string) {
        base = absolutePath;
    }

    function navigateUp() {
        if (base === "/" || base === "") return;
        const parts = base.split('/').filter(Boolean);
        parts.pop();
        base = "/" + parts.join("/");
    }

    function jumpToBreadcrumb(index: number) {
        const parts = base.split('/').filter(Boolean);
        base = "/" + parts.slice(0, index + 1).join("/");
    }

    function selectFolder() {
        value = base || "/";
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

<div class="flex flex-col gap-1.5 w-full">
    {#if label}
        <label class="text-[10px] font-bold text-muted uppercase tracking-widest ml-1">
            {label}
        </label>
    {/if}

    <div class="flex items-center group relative">
        <input
                type="text"
                bind:value={value}
                {placeholder}
                class="flex-1 min-w-0 bg-panel border border-border text-foreground text-sm rounded-l-xl px-4 py-2.5 outline-none focus:border-frost-500 transition-all"
        />

        <button
                {disabled}
                type="button"
                onclick={openPicker}
                class="bg-panel border border-l-0 border-border px-4 py-2.5 rounded-r-xl hover:bg-surface text-muted hover:text-frost-400 transition-all flex items-center justify-center active:scale-95 disabled:opacity-50"
                title="Browse folder"
        >
            <FolderOpenIcon size={18}/>
        </button>
    </div>
</div>

<!-- FOLDER PICKER DIALOG -->
{#if isOpen}
    <div class="fixed inset-0 z-200 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={() => isOpen = false}></div>

        <div class="relative w-full max-w-lg bg-surface border border-border rounded-4xl shadow-2xl flex flex-col overflow-hidden h-125"
             transition:fly={{ y: 20, duration: 300 }}
        >
            <div bind:this={scrollContainer}
                 class="flex items-center gap-1 p-2 bg-background border border-border rounded-xl overflow-x-auto no-scrollbar whitespace-nowrap text-xs scroll-smooth"
            >
                <button onclick={() => base = "/"}
                        class="px-2 py-1 shrink-0 hover:bg-panel rounded text-muted hover:text-frost-400 font-medium transition-colors"
                >
                    /
                </button>

                {#each breadcrumbs as part, i}
                    <span class="text-muted/30 shrink-0">/</span>
                    <button onclick={() => jumpToBreadcrumb(i)}
                            class="px-2 py-1 shrink-0 hover:bg-panel rounded text-muted hover:text-foreground font-bold transition-colors"
                    >
                        {part}
                    </button>
                {/each}
            </div>

            <!-- Folder List Area -->
            <div class="flex-1 overflow-y-auto p-2 bg-panel/10">
                {#if fileListRpc.loading}
                    <div class="flex flex-col items-center justify-center h-full gap-3 text-muted">
                        <LoaderIcon class="animate-spin text-frost-500" size={24}/>
                        <span class="text-[10px] font-bold uppercase tracking-widest opacity-50">Scanning Storage...</span>
                    </div>
                {:else if fileListRpc.error}
                    <div class="flex flex-col items-center justify-center h-full p-8 text-center text-red-400 gap-2">
                        <SearchXIcon size={32}/>
                        <p class="text-xs font-bold uppercase truncate w-full">{fileListRpc.error}</p>
                    </div>
                {:else}
                    <div class="flex flex-col gap-1">
                        {#if base !== "/" && base !== ""}
                            <button
                                    onclick={navigateUp}
                                    class="flex items-center gap-3 p-3 rounded-xl hover:bg-frost-500/5 text-frost-400 transition-colors text-left font-bold text-sm"
                            >
                                <ArrowLeftIcon size={18}/>
                                <span>.. (Parent Directory)</span>
                            </button>
                        {/if}

                        {#if fileListRpc.value?.files && fileListRpc.value.files.length > 0}
                            {#each fileListRpc.value.files as folder}
                                <button
                                        onclick={() => navigateTo(folder)}
                                        class="flex items-center justify-between p-3 rounded-xl hover:bg-panel transition-all group text-left"
                                >
                                    <div class="flex items-center gap-3 min-w-0">
                                        <FolderIcon size={18} class="text-muted group-hover:text-frost-400 shrink-0"
                                                    fill="currentColor" fill-opacity="0.1"/>
                                        <span class="text-sm font-medium truncate text-foreground/80 group-hover:text-foreground">
                                            {folder.split('/').pop() || folder}
                                        </span>
                                    </div>
                                    <ChevronRightIcon size={16} class="text-muted/20 group-hover:text-frost-400"/>
                                </button>
                            {/each}
                        {:else}
                            <div class="flex flex-col items-center justify-center py-20 text-muted/20 gap-2">
                                <FolderIcon size={48} strokeWidth={1}/>
                                <span class="text-[10px] font-bold uppercase tracking-widest">No accessible subfolders</span>
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>

            <!-- Footer -->
            <div class="p-4 border-t border-border bg-panel/30 flex justify-between items-center gap-4">
                <div class="flex-1 min-w-0">
                    <p class="text-[9px] font-bold text-muted uppercase tracking-widest mb-1 ml-1 opacity-50">Target
                        Directory</p>
                    <p class="text-xs font-mono text-frost-400 truncate bg-background p-2 rounded-lg border border-border">
                        {base || '/'}
                    </p>
                </div>
                <button
                        onclick={selectFolder}
                        class="px-6 py-3 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 transition-all shadow-lg shadow-frost-500/20 flex items-center gap-2 shrink-0 active:scale-95"
                >
                    <CheckIcon size={18}/>
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
