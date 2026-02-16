<script lang="ts">
    import {setContext} from 'svelte';
    import {fade, fly} from 'svelte/transition';
    import {dialogManager, dialogIcons, type DialogType, setDialogCtx} from './dialog.svelte';

    let {children} = $props();

    setDialogCtx()

    const themes: Record<DialogType, string> = {
        success: 'text-green-400 border-green-500/20 bg-green-500/5',
        info: 'text-frost-400 border-frost-500/20 bg-frost-500/5',
        warn: 'text-amber-400 border-amber-500/20 bg-amber-500/5',
        error: 'text-red-400 border-red-500/20 bg-red-500/5'
    };
</script>

{@render children()}

{#if dialogManager.active}
    <div class="fixed inset-0 z-300 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md"
             onclick={() => dialogManager.active?.resolve(false)}></div>

        <!-- Dialog Box -->
        <div
                class="relative w-full max-w-sm bg-surface border border-border rounded-[2rem] shadow-2xl overflow-hidden"
                transition:fly={{ y: 20, duration: 300 }}
        >
            <div class="p-8 space-y-6">
                <!-- Icon & Header -->
                <div class="flex flex-col items-center text-center gap-4">

                    <div class="p-4 rounded-full border {themes[dialogManager.active.type]}">
                        {#if dialogManager?.active?.type}
                            {@const Icon = dialogIcons[dialogManager.active.type]}
                            <Icon size={32} strokeWidth={1.5}/>
                        {/if}
                    </div>

                    <div class="space-y-2">
                        <h2 class="text-xl font-black uppercase tracking-tight text-foreground">
                            {dialogManager.active.title}
                        </h2>
                        <p class="text-sm text-muted leading-relaxed">
                            {dialogManager.active.body}
                        </p>
                    </div>
                </div>

                <!-- Actions -->
                <div class="flex gap-3 pt-2">
                    {#if dialogManager.active.variant === 'confirm'}
                        <button
                                onclick={() => dialogManager.active?.resolve(false)}
                                class="flex-1 py-3 bg-panel border border-border rounded-xl text-sm font-bold text-muted hover:text-foreground transition-all"
                        >
                            Cancel
                        </button>
                    {/if}

                    <button
                            onclick={() => dialogManager.active?.resolve(true)}
                            class="flex-2 py-3 bg-frost-500 text-background rounded-xl text-sm font-black uppercase tracking-widest hover:bg-frost-400 active:scale-95 transition-all shadow-lg shadow-frost-500/20"
                    >
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}