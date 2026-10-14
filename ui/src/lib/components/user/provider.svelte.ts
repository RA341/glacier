import { getContext, setContext } from 'svelte';
import { goto } from '$app/navigation';
import { page } from '$app/state';
import { glacierCli } from '#lib/api/api.js';
import { UserService } from '#lib/gen/user/v1/user_pb.js';
import { createRPCRunner } from '#lib/api/rpc.svelte.js';
import { AuthService } from '#lib/gen/auth/v1/auth_pb.js';
import { getSnackbarCtx } from '#lib/components/snackbar/snackbar-provider.svelte.js';

class UserInfoManager {
	private llSrv = glacierCli(UserService);
	private authSrv = glacierCli(AuthService);
	private rpc = createRPCRunner(() => this.llSrv.self({}));

	user = $derived(this.rpc.value);
	isLoading = $derived(this.rpc.loading);
	error = $derived(this.rpc.error);
	role = $derived(this.user?.user?.Role);
	isOmni = $derived(this.role === 'Omnissiah');

	async check() {
		const isAuthPage = page.url.pathname.startsWith('/auth');
		const isRoot = page.url.pathname === '/';

		await this.rpc.runner();

		// If fetch failed or user is null
		if (!this.user) {
			if (!isAuthPage) {
				await goto('/auth/login', { replace: true });
			}
		} else {
			// If logged in and on /auth or / redirect to library
			if (isAuthPage || isRoot) {
				await goto('/library', { replace: true });
			}
		}
	}

	async refresh() {
		await this.check();
	}

	async logout() {
		await this.authSrv.logout({});
		this.rpc.clear();
		await goto('/auth/login');
	}
}

const CTX_KEY = Symbol('user');

export function setUserCtx() {
	const manager = new UserInfoManager();
	return setContext(CTX_KEY, manager);
}

export function getUserCtx() {
	const ctx = getContext<UserInfoManager>(CTX_KEY);
	if (!ctx) throw new Error('getUserCtx must be used within a provider');
	return ctx;
}
