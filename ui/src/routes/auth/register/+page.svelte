<script lang="ts">
	import { fade, fly } from 'svelte/transition';
	import { KeyIcon, LoaderIcon, TriangleAlert, UserIcon } from '@lucide/svelte';
	import { glacierPubCli } from '#lib/api/api.js';
	import { AuthService } from '#lib/gen/auth/v1/auth_pb.js';
	import { createRPCRunner } from '#lib/api/rpc.svelte.js';
	import { goto } from '$app/navigation';

	let username = $state('');
	let password = $state('');
	let passwordVerify = $state('');

	const authSrv = glacierPubCli(AuthService);
	const regRpc = createRPCRunner(() =>
		authSrv.register({
			username: username,
			password: password,
			passwordVerify: passwordVerify
		})
	);

	async function handleRegister(e: Event) {
		e.preventDefault();
		await regRpc.runner();
		if (!regRpc.error) {
			await goto('/auth/login');
		}
	}
</script>

<div class="space-y-6" in:fade={{ duration: 300 }}>
	<div class="space-y-1">
		<h2 class="text-xl font-bold text-foreground">Create account</h2>
		<!--        <p class="text-sm text-muted">Join this server.</p>-->
	</div>

	<form onsubmit={handleRegister} class="space-y-4">
		<div class="space-y-2">
			<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">Username</label
			>
			<div class="relative">
				<UserIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
				<input
					type="text"
					bind:value={username}
					placeholder="John117"
					required
					class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
				/>
			</div>
		</div>

		<div class="space-y-2">
			<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase">Password</label
			>
			<div class="relative">
				<KeyIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
				<input
					type="password"
					bind:value={password}
					placeholder="••••••••"
					required
					class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
				/>
			</div>
		</div>
		<div class="space-y-2">
			<label class="ml-1 text-[10px] font-bold tracking-widest text-muted uppercase"
				>Verify Password</label
			>
			<div class="relative">
				<KeyIcon size={18} class="absolute top-1/2 left-4 -translate-y-1/2 text-muted/50" />
				<input
					type="password"
					bind:value={passwordVerify}
					placeholder="••••••••"
					required
					class="w-full rounded-[0.75rem] border border-border bg-panel py-3.5 pr-4 pl-12 text-sm transition-all outline-none focus:border-frost-500"
				/>
			</div>
		</div>

		<button
			type="submit"
			disabled={regRpc.loading}
			class="btn-frost flex w-full items-center justify-center gap-2 rounded-full py-4 font-bold text-foreground shadow-lg shadow-frost-500/20 transition-all hover:bg-frost-500 hover:text-background active:scale-[0.98]"
		>
			{#if regRpc.loading}
				<LoaderIcon size={20} class="animate-spin" />
			{:else}
				Create Account
			{/if}
		</button>
	</form>
</div>

{#if regRpc.error}
	<div
		class="fixed inset-0 z-110 flex items-center justify-center p-6"
		transition:fade={{ duration: 100 }}
	>
		<!-- Darker backdrop for the error alert -->
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div class="absolute inset-0 bg-black/40 backdrop-blur-sm" onclick={regRpc.clear}></div>

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
				{regRpc.error}
			</p>

			<div class="flex justify-end pt-2">
				<button
					onclick={regRpc.clear}
					class="btn-secondary rounded-2xl px-5 py-2 text-sm font-bold transition-all hover:text-foreground"
				>
					Dismiss
				</button>
			</div>
		</div>
	</div>
{/if}
