<script lang="ts">
    import {
        ChevronDownIcon,
        FolderIcon,
        CheckCircle2Icon,
        ClockIcon
    } from "@lucide/svelte";
    import { slide } from 'svelte/transition';
    import type { FolderProgress } from "$lib/gen/frost_library/v1/frost_library_pb";

    let { detail, ID }: { ID: string; detail: FolderProgress } = $props();

    let isExpanded = $state(false);

    const totalComplete = $derived(Number(detail.Complete));
    const totalLeft = $derived(Number(detail.Left));
    const totalSize = $derived(totalComplete + totalLeft);
    const overallProgress = $derived(totalSize > 0 ? (totalComplete / totalSize) * 100 : 0);

    // Human readable byte converter
    function formatBytes(bytes: number, decimals = 2) {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const dm = decimals < 0 ? 0 : decimals;
        const sizes = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
    }

    const getFileProgress = (complete: bigint | number, left: bigint | number) => {
        const c = Number(complete);
        const l = Number(left);
        return (c + l) > 0 ? (c / (c + l)) * 100 : 0;
    };
</script>

<div class="bg-panel/30 border border-border rounded-2xl overflow-hidden transition-all duration-300 {isExpanded ? 'border-frost-500/30 ring-1 ring-frost-500/10' : ''}">

    <!-- ACCORDION HEADER (Summary) -->
    <button
            onclick={() => isExpanded = !isExpanded}
            class="w-full flex items-center gap-4 p-4 hover:bg-panel/50 transition-colors text-left"
    >
        <div class="p-2.5 bg-surface border border-border rounded-xl text-muted group-hover:text-frost-400">
            <FolderIcon size={20} class={overallProgress === 100 ? 'text-green-500' : 'text-frost-400'}/>
        </div>

        <div class="flex-1 min-w-0 space-y-1.5">
            <div class="flex justify-between items-end">
                <span class="font-bold text-sm text-foreground truncate">{ID}</span>
                <span class="text-[10px] font-mono font-bold text-muted">
                    {overallProgress.toFixed(1)}%
                </span>
            </div>

            <!-- Overall Progress Bar -->
            <div class="h-1.5 w-full bg-panel rounded-full overflow-hidden">
                <div
                        class="h-full bg-frost-500 transition-all duration-500"
                        style="width: {overallProgress}%"
                ></div>
            </div>
        </div>

        <div class="flex items-center gap-3 ml-2">
            <div class="hidden sm:flex flex-col items-end text-[10px] font-bold uppercase tracking-widest text-muted/50 whitespace-nowrap">
                <span>{formatBytes(totalComplete)} Done</span>
                <span>{formatBytes(totalLeft)} Left</span>
            </div>
            <ChevronDownIcon
                    size={18}
                    class="text-muted transition-transform duration-300 {isExpanded ? 'rotate-180 text-frost-400' : ''}"
            />
        </div>
    </button>

    <!-- ACCORDION CONTENT (File List) -->
    {#if isExpanded}
        <div transition:slide={{ duration: 300 }} class="border-t border-border bg-black/10">
            <div class="p-2 space-y-1">
                {#each detail.files as file}
                    {@const fileComplete = Number(file.Complete)}
                    {@const fileTotal = Number(file.Complete) + Number(file.Left)}
                    {@const filePct = getFileProgress(file.Complete, file.Left)}
                    {@const isDone = Number(file.Left) === 0}

                    <div class="flex items-center gap-3 p-3 rounded-xl hover:bg-panel/40 transition-colors group">
                        <div class="text-muted/40 group-hover:text-frost-400/50 transition-colors">
                            {#if isDone}
                                <CheckCircle2Icon size={14} class="text-green-500/50"/>
                            {:else}
                                <ClockIcon size={14}/>
                            {/if}
                        </div>

                        <div class="flex-1 min-w-0">
                            <div class="flex justify-between items-center mb-1">
                                <p class="text-xs font-medium text-muted group-hover:text-foreground transition-colors truncate pr-4">
                                    {file.Name}
                                </p>
                                <span class="text-[9px] font-mono text-muted/40 whitespace-nowrap">
                                    {formatBytes(fileComplete)} / {formatBytes(fileTotal)}
                                </span>
                            </div>

                            <!-- Individual File Progress -->
                            <div class="h-1 w-full bg-panel/50 rounded-full overflow-hidden">
                                <div
                                        class="h-full {isDone ? 'bg-green-500/40' : 'bg-frost-500/40'} transition-all duration-700"
                                        style="width: {filePct}%"
                                ></div>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}
</div>

<style>
    button {
        user-select: none;
    }
</style>
