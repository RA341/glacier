<script lang="ts">
    import {
        GlobeIcon,
        ZapIcon,
        ChevronRightIcon,
        LoaderIcon,
        ShieldCheckIcon,
        CheckCircle2Icon
    } from '@lucide/svelte';
    import {fade, fly} from 'svelte/transition';

    let url = $state("");
    let isTesting = $state(false);
    let isProceeding = $state(false);
    let testSuccess = $state<boolean | null>(null);

    async function handleTest() {
        if (!url || isTesting) return;
        isTesting = true;
        testSuccess = null;



        testSuccess = true; // Placeholder result
        isTesting = false;
    }

    async function handleNext() {
        if (!url || isProceeding) return;
        isProceeding = true;

        // API PLACEHOLDER: implement save and navigate logic here
        await new Promise(r => setTimeout(r, 1000));

        isProceeding = false;
    }
</script>

<div
        class="w-full max-w-lg z-10 space-y-8"
        transition:fly={{ y: 20, duration: 600 }}
>

    <div class="space-y-4">
        <div class="flex justify-between items-end px-1">
            <label for="url" class="text-[10px] font-black text-white uppercase tracking-[0.2em]">Base Glacier
                URL</label>
            {#if testSuccess}
                        <span in:fade class="text-[10px] font-bold text-green-400 uppercase flex items-center gap-1">
                            <CheckCircle2Icon size={12}/> Connection Verified
                        </span>
            {/if}
        </div>

        <div class="relative group">
            <div class="absolute left-5 top-1/2 -translate-y-1/2 text-muted group-focus-within:text-frost-400 transition-colors">
                <GlobeIcon size={20}/>
            </div>
            <input
                    type="url"
                    id="url"
                    bind:value={url}
                    placeholder="https://glacier.your-domain.com"
                    class="w-full bg-panel border border-border rounded-2xl py-4 pl-14 pr-6 outline-none focus:border-frost-500 transition-all text-sm font-mono tracking-tight shadow-inner"
            />
        </div>
        <p class="text-[13px] text-white px-2 leading-relaxed italic">
            Enter the full address of your Glacier instance, including protocol (http/https).
        </p>
    </div>

    <!-- Action Buttons -->
    <div class="grid grid-cols-2 gap-4">
        <button
                onclick={handleTest}
                disabled={!url || isTesting || isProceeding}
                class="py-4 bg-panel border border-border text-foreground font-black uppercase tracking-widest text-[10px] rounded-2xl hover:bg-surface hover:border-frost-500/50 transition-all flex items-center justify-center gap-2 disabled:opacity-30 active:scale-95"
        >
            {#if isTesting}
                <LoaderIcon size={16} class="animate-spin text-frost-500"/>
                Validating...
            {:else}
                <ZapIcon size={16} class="text-frost-400"/>
                Test Link
            {/if}
        </button>

        <button
                onclick={handleNext}
                disabled={!url || isTesting || isProceeding}
                class="py-4 bg-frost-500 text-background font-black uppercase tracking-widest text-[10px] rounded-2xl hover:bg-frost-400 transition-all flex items-center justify-center gap-2 shadow-lg shadow-frost-500/20 disabled:opacity-30 active:scale-95"
        >
            {#if isProceeding}
                <LoaderIcon size={16} class="animate-spin"/>
                Connecting...
            {:else}
                Continue
                <ChevronRightIcon size={16}/>
            {/if}
        </button>
    </div>
</div>

<style>
    :global(body) {
        background-color: #0a0a0c;
    }

    input::placeholder {
        opacity: 0.3;
    }
</style>
