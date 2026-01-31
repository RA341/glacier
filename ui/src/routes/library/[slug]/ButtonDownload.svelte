<script lang="ts">
    import { DownloadIcon, LoaderIcon, AlertTriangleIcon, XIcon } from '@lucide/svelte';
    import { callRPC, frostCli } from "$lib/api/api";
    import { FrostLibraryService } from "$lib/gen/frost_library/v1/frost_library_pb";
    import { fade, fly } from 'svelte/transition';

    let { gameId }: { gameId: bigint } = $props();

    const llService = frostCli(FrostLibraryService);

    let isSending = $state(false);
    let errorMessage = $state<string | null>(null);

    async function download() {
        if (isSending) return;

        isSending = true;
        errorMessage = null;

        const { err } = await callRPC(() => llService.download({
            gameId: BigInt(gameId!),
            downloadFolder: "./downloads"
        }));

        if (err) {
            errorMessage = String(err);
        } else {
            // Optional: Show success state or close page
        }

        isSending = false;
    }
</script>

<button
        onclick={download}
        disabled={isSending}
        class="px-8 py-2 bg-frost-500 text-background rounded-xl text-sm font-bold hover:bg-frost-400 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2 shadow-lg shadow-frost-500/20 active:scale-95"
>
    {#if isSending}
        <LoaderIcon size={16} class="animate-spin" />
        Processing...
    {:else}
        <DownloadIcon size={16}/>
        Download
    {/if}
</button>

{#if errorMessage}
    <div class="fixed inset-0 z-100 flex items-center justify-center p-6" transition:fade={{ duration: 150 }}>
        <!-- Backdrop -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={() => errorMessage = null}></div>

        <!-- Dialog Box -->
        <div
                class="relative w-full max-w-sm bg-surface border border-red-500/30 rounded-2xl shadow-2xl p-6 flex flex-col gap-4"
                transition:fly={{ y: 20, duration: 300 }}
        >
            <div class="flex items-center justify-between">
                <div class="flex items-center gap-3 text-red-400">
                    <div class="p-2 bg-red-500/10 rounded-lg">
                        <AlertTriangleIcon size={20} />
                    </div>
                    <h3 class="font-bold">Download Failed</h3>
                </div>
                <button onclick={() => errorMessage = null} class="text-muted hover:text-foreground">
                    <XIcon size={20} />
                </button>
            </div>

            <p class="text-sm text-muted leading-relaxed">
                {errorMessage}
            </p>

            <div class="flex justify-end mt-2">
                <button
                        onclick={() => errorMessage = null}
                        class="px-6 py-2 bg-panel border border-border rounded-xl text-sm font-bold hover:text-foreground transition-all active:scale-95"
                >
                    Dismiss
                </button>
            </div>
        </div>
    </div>
{/if}
