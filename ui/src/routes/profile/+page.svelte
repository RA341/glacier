<script lang="ts">
	import {
		PenLineIcon,
		FingerprintPattern,
		KeyIcon,
		LoaderIcon,
		LockIcon,
		SaveIcon,
		ShieldCheckIcon,
		UserIcon
	} from '@lucide/svelte';
	import { callRPC, glacierCli } from '#lib/api/api.js';
	import { UserService } from '#lib/gen/user/v1/user_pb.js';
	import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';
	import { getUserCtx } from '#lib/components/user/provider.svelte.js';

	const userSrv = glacierCli(UserService);

	let newUsername = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');

	let isUpdatingUsername = $state(false);
	let isUpdatingPassword = $state(false);

	const user = getUserCtx();

	const sm = getSnackbarCtx();

	async function upUser(username?: string, pass?: string) {
		if (!user.user?.user?.id) {
			sm.push('Cannot update, empty user info', 'warn');
			return;
		}

		const { err } = await callRPC(() =>
			userSrv.edit({
				user: {
					id: user.user!.user!.id,
					Role: user.user!.user!.Role,
					Username: username ?? '',
					Password: pass ?? ''
				}
			})
		);

		if (err) {
			sm.push(`error occurred while updating: ${err}`, 'error');
			return false;
		}
		return true;
	}

	async function handleUpdateUsername() {
		if (!newUsername || newUsername === user.user?.user?.Username) return;

		isUpdatingUsername = true;
		const ok = await upUser(newUsername);
		isUpdatingUsername = false;
		if (!ok) {
			return;
		}
		sm.push('Username updated successfully', 'success');

		await user.refresh();
	}

	async function handleUpdatePassword() {
		if (!newPassword || newPassword !== confirmPassword) {
			sm.push('Passwords do not match', 'warn');
			return;
		}

		isUpdatingPassword = true;
		const ok = await upUser('', newPassword);
		isUpdatingPassword = false;
		if (!ok) {
			return;
		}

		newPassword = '';
		confirmPassword = '';
		sm.push('Password updated successfully', 'success');

		await user.refresh();
	}
</script>

<div class="mx-auto space-y-8 p-10">
	<!-- Header -->
	<header class="flex items-center gap-4 px-2">
		<div class="rounded-2xl border border-border bg-panel p-3 text-frost-400">
			<FingerprintPattern size={28} />
		</div>
		<div>
			<h1 class="headline-xl text-foreground">My Profile</h1>
			<p class="text-sm text-muted">Manage your account</p>
		</div>
	</header>

	{#if user.user && !user.user}
		<div class="flex h-64 flex-col items-center justify-center gap-4 text-muted">
			<LoaderIcon class="animate-spin text-frost-500" size={32} />
			<p class="animate-pulse text-sm">Retrieving account details...</p>
		</div>
	{:else if user.user}
		<div class="grid grid-cols-1 gap-8 md:grid-cols-3">
			<!-- LEFT: Profile Preview -->
			<div class="md:col-span-1">
				<div
					class="sticky top-8 flex flex-col items-center gap-4 rounded-4xl border border-border bg-surface p-8 text-center shadow-xl"
				>
					<div
						class="flex h-24 w-24 items-center justify-center rounded-3xl border-4 border-surface bg-panel text-frost-400 shadow-2xl"
					>
						{#if user.user.user?.Role?.toLowerCase() === 'admin'}
							<ShieldCheckIcon size={48} strokeWidth={1.5} />
						{:else}
							<UserIcon size={48} strokeWidth={1.5} />
						{/if}
					</div>

					<div>
						<h2 class="max-w-full truncate px-2 text-2xl font-black text-foreground">
							{user.user.user?.Username}
						</h2>
						<span
							class="mt-1 inline-block rounded-full border border-frost-500/20 bg-frost-500/10 px-3 py-1 text-[10px] font-bold tracking-widest text-frost-400 uppercase"
						>
							{user.user.user?.Role || 'Unknown'}
						</span>
					</div>
				</div>
			</div>

			<!-- RIGHT: Settings Forms -->
			<div class="space-y-6 md:col-span-2">
				<!-- Section 1: General Settings (Username) -->
				<div class="space-y-6 rounded-4xl border border-border bg-surface p-8 shadow-xl">
					<div class="flex items-center gap-3">
						<div class="rounded-xl bg-panel p-2 text-muted">
							<PenLineIcon size={20} />
						</div>
						<h3 class="text-lg font-bold">General Settings</h3>
					</div>

					<div class="space-y-4">
						<div class="space-y-2">
							<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
								>Display Username</label
							>
							<div class="relative">
								<UserIcon
									size={18}
									class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/40"
								/>
								<input
									type="text"
									bind:value={newUsername}
									class="w-full rounded-[0.75rem] border border-border bg-panel py-4 pr-4 pl-12 text-sm font-bold transition-all outline-none focus:border-frost-500"
								/>
							</div>
						</div>

						<div class="flex justify-end">
							<button
								onclick={handleUpdateUsername}
								disabled={isUpdatingUsername || newUsername === user.user?.user?.Username}
								class="btn-secondary flex items-center gap-2 rounded-2xl px-8 py-3 text-xs font-bold text-foreground transition-all hover:border-frost-500 disabled:opacity-30"
							>
								{#if isUpdatingUsername}
									<LoaderIcon size={14} class="animate-spin" />
								{:else}
									<SaveIcon size={14} />
								{/if}
								Save Username
							</button>
						</div>
					</div>
				</div>

				<!-- Section 2: Security Settings (Password) -->
				<div class="space-y-6 rounded-4xl border border-border bg-surface p-8 shadow-xl">
					<div class="flex items-center gap-3">
						<div class="rounded-xl bg-panel p-2 text-muted">
							<LockIcon size={20} />
						</div>
						<h3 class="text-lg font-bold">Security</h3>
					</div>

					<div class="space-y-5">
						<div class="space-y-2">
							<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
								>New Password</label
							>
							<div class="relative">
								<KeyIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/40" />
								<input
									type="password"
									bind:value={newPassword}
									placeholder="••••••••"
									class="w-full rounded-[0.75rem] border border-border bg-panel py-4 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
								/>
							</div>
						</div>

						<div class="space-y-2">
							<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
								>Confirm Password</label
							>
							<div class="relative">
								<LockIcon
									size={18}
									class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/40"
								/>
								<input
									type="password"
									bind:value={confirmPassword}
									placeholder="••••••••"
									class="w-full rounded-[0.75rem] border border-border bg-panel py-4 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
								/>
							</div>
						</div>

						<button
							onclick={handleUpdatePassword}
							disabled={isUpdatingPassword || !newPassword}
							class="btn-frost flex w-full items-center justify-center gap-2 rounded-full py-4 font-bold text-foreground transition-all hover:bg-frost-500 hover:text-background disabled:opacity-30"
						>
							{#if isUpdatingPassword}
								<LoaderIcon size={20} class="animate-spin" />
							{:else}
								<SaveIcon size={18} />
								Update Password
							{/if}
						</button>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
