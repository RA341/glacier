<script lang="ts">
	import { fade } from 'svelte/transition';
	import FileManager from './FileManager.svelte';

	let { game = $bindable() } = $props();
</script>

<div class="h-150 overflow-hidden rounded-3xl border border-border bg-surface shadow-sm">
	<FileManager {game} />
</div>
