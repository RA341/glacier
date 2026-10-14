<script lang="ts">
	import { fly } from 'svelte/transition';
	import { XIcon, CopyIcon, CheckIcon } from '@lucide/svelte'; // Import Copy icons
	import { setSnackbarCtx, type SnackbarType, toastIcons } from './snackbar-provider.svelte';

	let { children } = $props();
	const snackbarManager = setSnackbarCtx();

	let copiedId = $state<string | null>(null);

	async function copyToClipboard(id: string, text: string) {
		await navigator.clipboard.writeText(text);
		copiedId = id;
		setTimeout(() => {
			if (copiedId === id) copiedId = null;
		}, 2000);
	}

	const themes: Record<SnackbarType, string> = {
		success: 'border-green-500/20 bg-green-500/10 text-green-400',
		info: 'border-frost-500/20 bg-frost-500/10 text-frost-400',
		warn: 'border-amber-500/20 bg-amber-500/10 text-amber-400',
		error: 'border-red-500/20 bg-red-500/10 text-red-400'
	};
</script>

{@render children()}

<div
	class="pointer-events-none fixed bottom-8 left-1/2 z-[200] flex w-full max-w-md -translate-x-1/2 flex-col items-center gap-3 px-6"
>
	{#each snackbarManager.toasts as toast (toast.id)}
		{@const Icon = toastIcons[toast.type]}
		<div
			transition:fly={{ y: 20, duration: 400 }}
			onmouseenter={() => snackbarManager.pause(toast.id)}
			onmouseleave={() => snackbarManager.resume(toast.id)}
			role="alert"
			class="pointer-events-auto flex w-full items-start gap-3 rounded-2xl border px-4 py-3 shadow-2xl backdrop-blur-md {themes[
				toast.type
			]}"
		>
			<Icon size={18} class="mt-0.5 shrink-0" />

			<!-- line-clamp-3 handles the 3-line limit with ellipsis -->
			<p class="line-clamp-3 flex-1 text-sm leading-tight font-medium">
				{toast.message}
			</p>

			<div class="flex shrink-0 items-center gap-1">
				<!-- Copy Button -->
				<button
					onclick={() => copyToClipboard(toast.id, toast.message)}
					class="rounded-full p-1.5 text-current opacity-60 transition-colors hover:bg-black/10 hover:opacity-100"
					title="Copy message"
				>
					{#if copiedId === toast.id}
						<CheckIcon size={16} class="text-green-500" />
					{:else}
						<CopyIcon size={16} />
					{/if}
				</button>

				<!-- Close Button -->
				<button
					onclick={() => snackbarManager.remove(toast.id)}
					class="rounded-full p-1.5 text-current opacity-60 transition-colors hover:bg-black/10 hover:opacity-100"
				>
					<XIcon size={16} />
				</button>
			</div>
		</div>
	{/each}
</div>
