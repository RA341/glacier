<script lang="ts">
    export let label = ""; // "Default" or "ENV"
    export let value = ""; // The actual value to copy

    let copied = false;

    async function handleCopy() {
        try {
            await navigator.clipboard.writeText(value);
            copied = true;
            // Reset the icon after 2 seconds
            setTimeout(() => (copied = false), 2000);
        } catch (err) {
            console.error("Failed to copy!", err);
        }
    }

    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter' || event.key === ' ') {
            event.preventDefault();
            handleCopy();
        }
    }
</script>

<div
        class="flex items-center gap-2 group min-h-5 cursor-pointer select-none"
        role="button"
        tabindex="0"
        on:click={handleCopy}
        on:keydown={handleKeyDown}
>
    <!-- The Code Block -->
    <code class="text-[12px] font-mono bg-zinc-100 dark:bg-zinc-800 px-1.5 py-0.5 rounded border border-zinc-200 dark:border-zinc-700 text-zinc-700 dark:text-zinc-300">
        {#if label}<span class="opacity-50 font-sans mr-1">{label}:</span>{/if}{value}
    </code>

    <!-- The Copy Button (appears on hover) -->
    <button
            type="button"
            on:click={handleCopy}
            class="opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity p-1 rounded hover:bg-zinc-100 dark:hover:bg-zinc-800"
            title="Copy to clipboard"
    >
        {#if copied}
            <!-- Checkmark Icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-green-500" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
        {:else}
            <!-- Copy Icon -->
            <svg xmlns="http://www.w3.org/2000/svg" class="w-3 h-3 text-zinc-400" viewBox="0 0 24 24" fill="none"
                 stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
        {/if}
    </button>
</div>