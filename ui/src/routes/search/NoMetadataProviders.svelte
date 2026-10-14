<script lang="ts">
	import { DatabaseIcon, SettingsIcon } from '@lucide/svelte';
</script>

<div
	class="flex min-h-80 flex-col items-center justify-center gap-4 rounded-3xl border-2 border-dashed border-border p-6 text-center"
>
	<DatabaseIcon size={48} strokeWidth={1.5} class="text-frost-400" />
	<div class="space-y-2">
		<h2 class="text-xl font-bold">No metadata providers set up</h2>
		<p class="text-sm text-muted">Set up a metadata provider to search for games.</p>
	</div>
	<a
		href="/settings/glacier/metadata/"
		class="inline-flex items-center gap-2 rounded-xl bg-frost-500 px-5 py-3 text-sm font-bold text-background transition-colors hover:bg-frost-400"
	>
		<SettingsIcon size={18} />
		Set up metadata providers
	</a>
</div>
