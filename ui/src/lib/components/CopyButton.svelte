<script lang="ts">
	export let label = ''; // "Default" or "ENV"
	export let value = ''; // The actual value to copy

	let copied = false;

	async function handleCopy() {
		try {
			await navigator.clipboard.writeText(value);
			copied = true;
			// Reset the icon after 2 seconds
			setTimeout(() => (copied = false), 2000);
		} catch (err) {
			console.error('Failed to copy!', err);
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
	class="group flex min-h-5 cursor-pointer items-center gap-2 select-none"
	role="button"
	tabindex="0"
	on:click={handleCopy}
	on:keydown={handleKeyDown}
>
	<!-- The Code Block -->
	<code
		class="rounded-full border border-border bg-panel px-2 py-1 font-hud text-xs text-muted hover:border-frost-400/40 hover:text-frost-300"
	>
		{#if label}<span class="mr-1 font-sans opacity-50">{label}:</span>{/if}{value}
	</code>

	<!-- The Copy Button (appears on hover) -->
	<button
		type="button"
		on:click={handleCopy}
		class="rounded p-1 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-zinc-100 focus:opacity-100 dark:hover:bg-zinc-800"
		title="Copy to clipboard"
	>
		{#if copied}
			<!-- Checkmark Icon -->
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-3 w-3 text-green-500"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="3"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<polyline points="20 6 9 17 4 12"></polyline>
			</svg>
		{:else}
			<!-- Copy Icon -->
			<svg
				xmlns="http://www.w3.org/2000/svg"
				class="h-3 w-3 text-zinc-400"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
				<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
			</svg>
		{/if}
	</button>
</div>
