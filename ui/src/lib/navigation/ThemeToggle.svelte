<script lang="ts">
	import { MoonIcon, SunIcon } from '@lucide/svelte';
	import { theme, toggleTheme } from './theme.svelte.js';

	const label = $derived(`Switch to ${theme.current === 'dark' ? 'light' : 'dark'} theme`);
</script>

<button
	type="button"
	onclick={toggleTheme}
	title={label}
	aria-label={label}
	class="flex h-11 w-11 items-center justify-center rounded-full text-muted transition-colors hover:bg-panel hover:text-frost-400"
>
	{#if theme.current === 'dark'}
		<SunIcon size={22} strokeWidth={2} />
	{:else}
		<MoonIcon size={22} strokeWidth={2} />
	{/if}
</button>
