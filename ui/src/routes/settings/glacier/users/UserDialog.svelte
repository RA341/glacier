<script lang="ts">
	import {
		XIcon,
		UserIcon,
		LockIcon,
		ShieldIcon,
		SaveIcon,
		UserPlusIcon,
		UserCogIcon,
		ChevronDownIcon
	} from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { type User, UserService } from '#lib/gen/user/v1/user_pb.js';
	import { glacierCli } from '#lib/api/api.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { callRPC } from '#lib/api/api.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

	interface Props {
		open: boolean;
		editUser?: User;
		onClose?: () => void;
	}

	let { onClose, editUser, open = $bindable() }: Props = $props();

	let username = $state('');
	let password = $state('');
	let role = $state('User');

	$effect(() => {
		if (open) {
			rolesRpc.runner();
			username = '';
			password = '';
			role = 'User';
		}

		if (editUser) {
			username = editUser.Username;
			password = '';
			role = editUser.Role || 'User';
		}
	});

	function handleClose() {
		open = false;
		if (onClose) onClose();
	}

	const userSrv = glacierCli(UserService);

	let rolesRpc = createRPCRunner(() => userSrv.listRoles({}));
	const snackbarManager = getSnackbarCtx();

	async function handleSave() {
		const user = {
			id: BigInt(editUser?.id ?? 0),
			Username: username,
			Password: password,
			Role: role
		};

		if (editUser) {
			const { err } = await callRPC(() => userSrv.edit({ user: user }));
			if (err) {
				snackbarManager.push(`Error creating ${err}`, 'error');
			}
		} else {
			const { err } = await callRPC(() => userSrv.new({ user: user }));
			if (err) {
				snackbarManager.push(`Error editing ${err}`, 'error');
			}
		}

		handleClose();
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-100 flex items-center justify-center p-6"
		transition:fade={{ duration: 150 }}
	>
		<!-- Backdrop -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="absolute inset-0 bg-background/80 backdrop-blur-md" onclick={handleClose}></div>

		<!-- Dialog Content -->
		<div
			class="relative flex w-full max-w-md flex-col overflow-hidden rounded-4xl border border-border bg-surface shadow-2xl"
			transition:fly={{ y: 20, duration: 300 }}
		>
			<!-- Header -->
			<div class="flex items-center justify-between border-b border-border bg-panel/30 p-6">
				<div class="flex items-center gap-3">
					<div class="rounded-xl bg-frost-500/10 p-2.5 text-frost-400">
						{#if editUser}
							<UserCogIcon size={20} />
						{:else}
							<UserPlusIcon size={20} />
						{/if}
					</div>
					<h2 class="text-xl font-bold text-foreground">
						{editUser ? 'Edit User' : 'Create User'}
					</h2>
				</div>
				<button
					onclick={handleClose}
					class="p-1 text-muted transition-colors hover:text-foreground"
				>
					<XIcon size={20} />
				</button>
			</div>

			<!-- Form Body -->
			<div class="space-y-6 p-8">
				<!-- Username -->
				<div class="space-y-2">
					<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
						>Username</label
					>
					<div class="relative">
						<UserIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
						<input
							type="text"
							bind:value={username}
							placeholder="johndoe"
							class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
						/>
					</div>
				</div>

				<!-- Password -->
				<div class="space-y-2">
					<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">
						{'Password'}
					</label>
					<div class="relative">
						<LockIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
						<input
							type="password"
							bind:value={password}
							placeholder="••••••••"
							class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
						/>
					</div>
					{#if editUser}
						<p class="px-1 text-[10px] text-muted/50 italic">
							Leave blank to keep existing password
						</p>
					{/if}
				</div>

				<!-- Role -->
				<div class="space-y-2">
					<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
						>System Role</label
					>
					<div class="relative">
						<ShieldIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
						<select
							bind:value={role}
							class="w-full cursor-pointer appearance-none rounded-[0.75rem] border border-border bg-panel py-3.5 pr-10 pl-12 text-sm font-bold transition-all outline-none focus:border-frost-500"
						>
							{#each rolesRpc.value?.roles as r}
								<option value={r.Name}>{r.Name}</option>
							{/each}
						</select>
						<div
							class="pointer-events-none absolute top-1/2 right-4 -translate-y-1/2 text-muted/50"
						>
							<ChevronDownIcon size={16} />
						</div>
					</div>
				</div>
			</div>

			<!-- Footer -->
			<div class="flex justify-end gap-3 border-t border-border bg-panel/30 p-6">
				<button
					onclick={handleClose}
					class="rounded-full border border-border px-6 py-2.5 text-sm font-bold text-muted transition-all hover:bg-panel hover:text-foreground"
				>
					Cancel
				</button>
				<button
					onclick={handleSave}
					class="btn-frost flex items-center gap-2 rounded-full px-8 py-2.5 text-sm font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-95"
				>
					<SaveIcon size={16} />
					{editUser ? 'Save Changes' : 'Create User'}
				</button>
			</div>
		</div>
	</div>
{/if}

<style>
	select {
		background-image: none;
	}
</style>
