<script lang="ts">
	import {
		MessageSquareIcon,
		ArrowLeftIcon,
		ShieldAlertIcon,
		ChevronLeftIcon
	} from '@lucide/svelte';
	import { fade } from 'svelte/transition';
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
	<!-- Header Section -->
	<div class="space-y-2">
		<h2 class="text-xl leading-tight font-bold text-foreground">Password Recovery</h2>
		<!--        <p class="text-sm text-muted">Self-service recovery is disabled for security.</p>-->
	</div>

	<!-- Instructions Card -->
	<div
		class="flex flex-col items-center gap-4 rounded-2xl border border-border bg-panel/40 p-6 text-center"
	>
		<div class="rounded-full border border-amber-500/20 bg-amber-500/10 p-4 text-amber-500">
			<MessageSquareIcon size={32} strokeWidth={1.5} />
		</div>

		<div class="space-y-2">
			<p class="text-sm font-bold text-foreground">Administrator Action Required</p>
			<p class="px-2 text-xs leading-relaxed text-muted">
				Please contact your administrator to request a credential reset.
			</p>
		</div>
	</div>

	<!-- Navigation -->
	<div class="border-t border-border pt-4">
		<a
			href="/auth/login"
			class="group flex w-full items-center justify-center gap-2 rounded-2xl border border-border bg-panel py-4 font-bold text-foreground transition-all hover:border-frost-500/50 hover:bg-surface active:scale-[0.98]"
		>
			<ChevronLeftIcon size={18} class="transition-transform group-hover:-translate-x-1" />
			Back to Login
		</a>
	</div>
</div>
