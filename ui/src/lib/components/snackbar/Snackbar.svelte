<script lang="ts">
    import {fly} from 'svelte/transition';
    import {XIcon} from '@lucide/svelte';
    import {setSnackbarCtx, type SnackbarType, toastIcons} from './snackbar-provider.svelte';

    let {children} = $props();

    const snackbarManager = setSnackbarCtx();

    const themes: Record<SnackbarType, string> = {
        success: 'border-green-500/20 bg-green-500/10 text-green-400',
        info: 'border-frost-500/20 bg-frost-500/10 text-frost-400',
        warn: 'border-amber-500/20 bg-amber-500/10 text-amber-400',
        error: 'border-red-500/20 bg-red-500/10 text-red-400'
    };
</script>

{@render children()}

<!-- SNACKBAR CONTAINER -->
<div class="fixed bottom-8 left-1/2 -translate-x-1/2 z-200 flex flex-col items-center gap-3 w-full max-w-md px-6 pointer-events-none">
    {#each snackbarManager.toasts as toast (toast.id)}
        {@const Icon = toastIcons[toast.type]}
        <div
                transition:fly={{ y: 20, duration: 400 }}
                class="pointer-events-auto flex items-center gap-3 px-4 py-3 rounded-2xl border backdrop-blur-md shadow-2xl w-full {themes[toast.type]}"
        >
            <!-- Dynamic Icon -->
            <Icon size={18} class="shrink-0"/>

            <p class="text-sm font-bold flex-1 leading-tight">
                {toast.message}
            </p>

            <!-- Close Button -->
            <button
                    onclick={() => snackbarManager.remove(toast.id)}
                    class="p-1 hover:bg-black/10 rounded-lg transition-colors text-current opacity-60 hover:opacity-100"
            >
                <XIcon size={16}/>
            </button>
        </div>
    {/each}
</div>