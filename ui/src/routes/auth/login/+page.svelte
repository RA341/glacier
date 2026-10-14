<script lang="ts">
	import { TriangleAlert, ArrowRightIcon, KeyIcon, LoaderIcon, User } from '@lucide/svelte';
	import { fade, fly } from 'svelte/transition';
	import { glacierPubCli } from '#lib/api/api.js';
	import { AuthService } from '#lib/gen/auth/v1/auth_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { goto } from '$app/navigation';

	let username = $state('');
	let password = $state('');

	const authSrv = glacierPubCli(AuthService);
	const loginRpc = createRPCRunner(() =>
		authSrv.login({
			username: username,
			password: password
		})
	);

	async function handleLogin(e: Event) {
		e.preventDefault();
		await loginRpc.runner();

		if (!loginRpc.error) {
			await goto('/library', { replaceState: true });
		}
	}
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
	<div class="space-y-1">
		<h2 class="text-xl font-bold text-foreground">Welcome back</h2>
		<p class="text-sm text-muted">Please enter your credentials to continue.</p>
	</div>

	<form onsubmit={handleLogin} class="space-y-4">
		<div class="space-y-2">
			<label for="email" class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">
				Username
			</label>
			<div class="relative">
				<User size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
				<input
					type="text"
					id="username"
					bind:value={username}
					placeholder="Username"
					required
					class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<div class="flex items-center justify-between px-1">
				<label for="password" class="text-[10px] font-bold tracking-widest text-muted uppercase"
					>Password</label
				>
				<a
					href="/auth/forgot"
					class="text-[10px] font-bold tracking-widest text-frost-500 uppercase hover:text-frost-400"
					>Forgot?</a
				>
			</div>
			<div class="relative">
				<KeyIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
				<input
					type="password"
					id="password"
					bind:value={password}
					placeholder="••••••••"
					required
					class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
				/>
			</div>
		</div>

		<button
			type="submit"
			disabled={loginRpc.loading}
			class="btn-frost flex w-full items-center justify-center gap-2 rounded-full py-4 font-bold text-foreground transition-all hover:bg-frost-500 hover:text-background active:scale-[0.98] disabled:opacity-50"
		>
			{#if loginRpc.loading}
				<LoaderIcon size={20} class="animate-spin" />
			{:else}
				Sign In
				<ArrowRightIcon size={18} />
			{/if}
		</button>
	</form>

	<button
		type="submit"
		onclick={() => (window.location.href = '/api/server/public/auth/oidc')}
		class="btn-frost flex w-full items-center justify-center gap-2 rounded-2xl py-4 font-bold text-foreground transition-all hover:bg-frost-500 hover:text-background active:scale-[0.98] disabled:opacity-50"
	>
		{#if loginRpc.loading}
			<LoaderIcon size={20} class="animate-spin" />
		{:else}
			Use OIDC
			<ArrowRightIcon size={18} />
		{/if}
	</button>
</div>

{#if loginRpc.error}
	<div
		class="fixed inset-0 z-110 flex items-center justify-center p-6"
		transition:fade={{ duration: 100 }}
	>
		<!-- Darker backdrop for the error alert -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="absolute inset-0 bg-black/40 backdrop-blur-sm" onclick={loginRpc.clear}></div>

		<div
			class="relative flex w-full max-w-sm flex-col gap-4 rounded-2xl border border-red-500/30 bg-surface p-6 shadow-2xl"
			transition:fly={{ y: 10, duration: 200 }}
		>
			<div class="flex items-center gap-3 text-red-400">
				<div class="rounded-lg bg-red-500/10 p-2">
					<TriangleAlert size={20} />
				</div>
				<h3 class="font-bold">Action Failed</h3>
			</div>

			<p class="text-sm leading-relaxed text-muted">
				{loginRpc.error}
			</p>

			<div class="flex justify-end pt-2">
				<button
					onclick={loginRpc.clear}
					class="btn-secondary rounded-2xl px-5 py-2 text-sm font-bold transition-all hover:text-foreground"
				>
					Dismiss
				</button>
			</div>
		</div>
	</div>
{/if}
