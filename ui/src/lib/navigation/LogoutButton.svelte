<script lang="ts">
	import { goto } from '$app/navigation';
	import { LogOutIcon } from '@lucide/svelte';
	import { glacierPubCli } from '#lib/api/api.js';
	import { AuthService } from '#lib/gen/auth/v1/auth_pb.js';
	import { getDialogCtx } from '#lib/components/dialog/dialog.svelte.js';
	import { resolve } from '$app/paths';

	const auth = glacierPubCli(AuthService);
	const dm = getDialogCtx();

	async function logout() {
		const ok = await dm.confirm('Logout', 'Did you grab everything? Phone. Keys. Wallet.', 'info');
		if (!ok) return;
		await auth.logout({});
		await goto(resolve('/auth/login'), { replace: true });
	}
</script>

<button
	type="button"
	onclick={logout}
	title="Logout"
	aria-label="Logout"
	class="flex h-11 w-11 items-center justify-center rounded-xl text-muted transition-colors hover:bg-panel hover:text-foreground"
>
	<LogOutIcon size={22} strokeWidth={2} />
</button>
