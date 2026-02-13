<script lang="ts">
    import {callRPC, frostCli} from "$lib/api/api";
    import {FrostLibraryService} from "$lib/gen/frost_library/v1/frost_library_pb";
    import {getSnackbarCtx} from "$lib/components/snackbar/snackbar-provider.svelte";
    import {CheckIcon, GaugeIcon} from '@lucide/svelte';

    const frostLib = frostCli(FrostLibraryService);
    const sm = getSnackbarCtx();

    let presets = [1, 5, 10, 100, 500];
    let customValue = $state<number | null>(null);
    let activeLimit = $state<number | null>(null);
    let isUpdating = $state(false);

    async function setLimit(limit: number) {
        if (isUpdating) return;
        isUpdating = true;

        const {err} = await callRPC(() => frostLib.throttleSpeed({
            limit: limit * 1024 * 1024
        }));

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

<div class="bg-panel border border-border rounded-2xl p-4 flex flex-col sm:flex-row items-center gap-6 shadow-sm">
    <!-- Icon & Label -->
    <div class="flex items-center gap-3 shrink-0">
        <div class="p-2 bg-frost-500/10 rounded-lg text-frost-400">
            <GaugeIcon size={18}/>
        </div>
        <div class="flex flex-col">
            <span class="text-[10px] font-bold text-muted uppercase tracking-widest">Throttle</span>
            <span class="text-xs font-bold text-foreground">Speed Limit</span>
        </div>
    </div>

    <!-- Presets Section -->
    <div class="flex gap-1 bg-surface p-1 rounded-xl border border-border">
        {#each presets as p}
            <button
                    onclick={() => setLimit(p)}
                    disabled={isUpdating}
                    class="px-4 py-1.5 rounded-lg text-[11px] font-bold transition-all
                {activeLimit === p
                    ? 'bg-frost-500 text-background shadow-lg shadow-frost-500/20'
                    : 'text-muted hover:text-foreground hover:bg-panel'}"
            >
                {p}M
            </button>
        {/each}

        <button
                onclick={() => setLimit(0)}
                disabled={isUpdating}
                class="px-4 py-1.5 rounded-lg text-[11px] font-bold transition-all
            {activeLimit === 0
                ? 'bg-frost-500 text-background shadow-lg shadow-frost-500/20'
                : 'text-muted hover:text-amber-400 hover:bg-panel'}"
                title="Remove Limit"
        >
            Unlimited
        </button>
    </div>

    <!-- Divider (Hidden on small mobile) -->
    <div class="hidden sm:block w-px h-8 bg-border"></div>

    <!-- Custom Input -->
    <div class="relative flex-1 min-w-35 group">
        <input
                type="number"
                bind:value={customValue}
                onkeydown={handleCustomSubmit}
                placeholder="Custom (MB/s)"
                class="w-full bg-surface border border-border rounded-xl py-2 pl-4 pr-10 text-xs font-bold outline-none focus:border-frost-500 transition-all placeholder:text-muted/30"
        />
        <button
                onclick={() => customValue && setLimit(customValue)}
                disabled={!customValue || isUpdating}
                class="absolute right-1.5 top-1.5 bottom-1.5 px-2 bg-panel border border-border rounded-lg text-muted hover:text-frost-400 transition-all disabled:opacity-0"
        >
            <CheckIcon size={14}/>
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

    input[type=number] {
        -moz-appearance: textfield;
    }
</style>
