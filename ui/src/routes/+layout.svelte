<script lang="ts">
	import './layout.css';
	import { page } from '$app/state';
	import { appName } from '#lib/api/api.js';
	import Snackbar from '#lib/components/snackbar/Snackbar.svelte';
	import DesktopSidebar from '#lib/navigation/DesktopSidebar.svelte';
	import MobileTopBar from '#lib/navigation/MobileTopBar.svelte';
	import MobileBottomNav from '#lib/navigation/MobileBottomNav.svelte';
	import User from '#lib/components/user/User.svelte';
	import DialogProvider from '#lib/components/dialog/DialogProvider.svelte';
	import { derived } from 'svelte/store';

	let { children } = $props();

	let isNotAuth = !page.url.pathname.startsWith('/auth')
</script>

<svelte:head>
	<title>{appName}</title>
	<link rel="icon" href={'/favicon.svg'} />
</svelte:head>

<Snackbar>
	<DialogProvider>
		<User>
			<div class="flex h-dvh min-h-0 w-full flex-col md:h-full md:flex-row">
				{#if isNotAuth}
					<DesktopSidebar />
					<MobileTopBar />
				{/if}
				<main class="relative min-h-0 min-w-0 flex-1 overflow-y-auto">
					{@render children()}
				</main>
				{#if isNotAuth}
					<MobileBottomNav />
				{/if}
			</div>
		</User>
	</DialogProvider>
</Snackbar>

<style>
	::-webkit-scrollbar {
		width: 6px;
	}

	::-webkit-scrollbar-track {
		background: transparent;
	}

	::-webkit-scrollbar-thumb {
		background: var(--border);
		border-radius: 10px;
	}

	::-webkit-scrollbar-thumb:hover {
		background: #444;
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
