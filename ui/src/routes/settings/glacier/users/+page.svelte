<script lang="ts">
	import {
		CircleAlert,
		EllipsisVertical,
		LoaderIcon,
		PlusIcon,
		RefreshCcw,
		SearchIcon,
		ShieldCheckIcon,
		Trash2Icon,
		UserCogIcon,
		UserIcon,
		UserMinusIcon
	} from '@lucide/svelte';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { type User, UserService } from '#lib/gen/user/v1/user_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import UserDialog from './UserDialog.svelte';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	const userSrv = glacierCli(UserService);
	let query = $state('');

	let listRpc = createRPCRunner(() => userSrv.list({ query: query }));

	function refreshList() {
		listRpc.runner();
	}

	onMount(() => {
		refreshList();
	});

	function getRoleTheme(role: string = '') {
		const r = role.toLowerCase();
		if (r === 'admin') return 'text-frost-400 bg-frost-400/10 border-frost-400/20';
		return 'text-muted bg-panel border-border';
	}

	let isOpen = $state(false);

	let user = $state<User | undefined>();

	function editUser(input: User) {
		isOpen = true;
		user = input;
	}

	let sm = getSnackbarCtx();

	async function deleteUser(input: User) {
		const { err } = await callRPC(() => userSrv.delete({ id: input.id }));
		sm.push(`error deleting user: ${err}`, 'error');
		refreshList();
	}

	function onClose() {
		isOpen = false;
		refreshList();
	}
</script>

<div class="max-w-1l mx-auto space-y-6">
	<!-- HEADER & SEARCH -->
	<header class="flex flex-col justify-between gap-4 px-2 sm:flex-row sm:items-center">
		<div class="flex items-center gap-3">
			<button
				onclick={() => {
					isOpen = true;
					user = undefined;
				}}
				class="flex items-center gap-2 rounded-xl bg-frost-500 px-6 py-2.5 text-sm font-bold text-background shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-400 active:scale-95"
			>
				<PlusIcon size={18} />
				Add
			</button>

			<div class="group relative">
				<SearchIcon
					size={16}
					class="absolute top-1/2 left-4 -translate-y-1/2 text-muted transition-colors group-focus-within:text-frost-400"
				/>
				<input
					type="text"
					bind:value={query}
					onkeydown={(e) => e.key === 'Enter' && refreshList()}
					placeholder="Filter by username..."
					class="w-full rounded-xl border border-border bg-panel py-2 pr-4 pl-11 text-sm transition-all outline-none focus:border-frost-500 sm:w-64"
				/>
			</div>
			<button
				onclick={refreshList}
				class="rounded-xl border border-border bg-panel p-2 text-muted transition-all hover:text-frost-400"
			>
				<RefreshCcw size={18} class={listRpc.loading ? 'animate-spin' : ''} />
			</button>
		</div>
	</header>

	<!-- LIST AREA -->
	<main class="min-h-100">
		{#if listRpc.loading && !listRpc.value}
			<div class="flex h-64 flex-col items-center justify-center gap-4 text-muted">
				<LoaderIcon class="animate-spin text-frost-500" size={32} />
				<p class="animate-pulse text-sm font-medium">Fetching users...</p>
			</div>
		{:else if listRpc.error}
			<div
				class="flex h-64 flex-col items-center justify-center gap-3 rounded-3xl border border-red-500/10 bg-red-500/5 text-red-400"
			>
				<CircleAlert size={32} />
				<p class="text-sm font-medium">{listRpc.error}</p>
			</div>
		{:else if !listRpc.value?.users || listRpc.value.users.length === 0}
			<div
				class="flex h-64 flex-col items-center justify-center rounded-3xl border-2 border-dashed border-border text-muted/30"
			>
				<UserMinusIcon size={48} strokeWidth={1} class="mb-2" />
				<p class="text-sm font-medium">No users found</p>
			</div>
		{:else}
			<!-- THE LIST STACK -->
			<div class="flex flex-col gap-2">
				{#each listRpc.value.users as user (user.Username)}
					<div
						role="button"
						tabindex="0"
						onclick={() => editUser(user)}
						onkeydown={(e) => {
							if (e.key === 'Enter' || e.key === ' ') {
								e.preventDefault();
								editUser(user);
							}
						}}
						class="group flex cursor-pointer items-center justify-between rounded-xl border border-border bg-panel/20 p-3 pl-4 transition-all hover:border-frost-500/30 hover:bg-panel/40"
						in:fade={{ duration: 200 }}
					>
						<div class="flex min-w-0 items-center gap-4">
							<!-- Role-based Icon -->
							<div
								class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-border bg-surface text-muted transition-colors group-hover:text-frost-400"
							>
								{#if user.Role?.toLowerCase() === 'admin'}
									<ShieldCheckIcon size={20} />
								{:else}
									<UserIcon size={20} />
								{/if}
							</div>

							<div class="flex min-w-0 flex-col">
								<span class="truncate text-sm font-bold text-foreground">{user.Username}</span>
							</div>
						</div>

						<div class="flex items-center gap-6">
							<!-- Role Badge -->
							<span
								class="hidden rounded border px-2 py-0.5 text-[10px] font-bold tracking-widest uppercase sm:inline-block {getRoleTheme(
									user.Role
								)}"
							>
								{user.Role || 'User'}
							</span>

							<!-- Actions -->
							<div
								class="flex translate-x-1 items-center gap-1 pr-1 opacity-0 transition-all group-hover:translate-x-0 group-hover:opacity-100"
							>
								<button
									onclick={(e) => {
										e.preventDefault();
										e.stopPropagation();
										editUser(user);
									}}
									class="rounded-lg p-2 text-muted transition-all hover:bg-frost-500/10 hover:text-frost-400"
									title="Edit Permissions"
								>
									<UserCogIcon size={18} />
								</button>
								<button
									onclick={(e) => {
										e.preventDefault();
										e.stopPropagation();
										deleteUser(user);
									}}
									class="rounded-lg p-2 text-muted transition-all hover:bg-red-500/10 hover:text-red-400"
									title="Remove User"
								>
									<Trash2Icon size={18} />
								</button>
							</div>

							<!-- Small indicator for mobile or non-hover -->
							<div class="text-muted/20 group-hover:hidden">
								<EllipsisVertical size={18} />
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</main>
</div>

<UserDialog open={isOpen} editUser={user} {onClose} />

<style>
	/* Prevent text selection on names */
	span {
		user-select: none;
	}
</style>
