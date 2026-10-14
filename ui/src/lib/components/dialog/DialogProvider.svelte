<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { dialogIcons, dialogManager, type DialogType, setDialogCtx } from './dialog.svelte';

	let { children } = $props();

	setDialogCtx();

	const themes: Record<DialogType, string> = {
		success: 'text-green-400 border-green-500/20 bg-green-500/5',
		info: 'text-frost-400 border-frost-500/20 bg-frost-500/5',
		warn: 'text-amber-400 border-amber-500/20 bg-amber-500/5',
		error: 'text-red-400 border-red-500/20 bg-red-500/5'
	};
</script>

{@render children()}

{#if dialogManager.active}
	<div
		class="fixed inset-0 z-300 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<!-- Backdrop -->
		<div
			class="absolute inset-0 bg-background/80 backdrop-blur-md"
			onclick={() => dialogManager.active?.resolve(false)}
		></div>

		<!-- Dialog Box -->
		<div
			class="relative w-full max-w-sm overflow-hidden rounded-[2rem] border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<div class="space-y-6 p-8">
				<!-- Icon & Header -->
				<div class="flex flex-col items-center gap-4 text-center">
					<div class="rounded-full border p-4 {themes[dialogManager.active.type]}">
						{#if dialogManager?.active?.type}
							{@const Icon = dialogIcons[dialogManager.active.type]}
							<Icon size={32} strokeWidth={1.5} />
						{/if}
					</div>

					<div class="space-y-2">
						<h2 class="text-xl font-black tracking-tight text-foreground uppercase">
							{dialogManager.active.title}
						</h2>
						<p class="text-sm leading-relaxed whitespace-pre-wrap text-muted">
							{dialogManager.active.body}
						</p>
					</div>
				</div>

				<!-- Actions -->
				<div class="flex gap-3 pt-2">
					{#if dialogManager.active.variant === 'confirm'}
						<button
							onclick={() => dialogManager.active?.resolve(false)}
							class="btn-secondary flex-1 rounded-2xl py-3 text-sm font-bold text-muted transition-all hover:text-foreground"
						>
							Cancel
						</button>
					{/if}

					<button
						onclick={() => dialogManager.active?.resolve(true)}
						class="btn-frost flex-2 rounded-full py-3 text-sm font-black tracking-widest text-foreground uppercase shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-95"
					>
						Confirm
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
