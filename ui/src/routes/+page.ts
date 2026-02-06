import { redirect } from '@sveltejs/kit';
import { Glacier } from '$lib/api/api';

export async function load({ url }) {
	const isAuthPage = url.pathname.startsWith('/auth');
	const isRoot = url.pathname === '/';

	if (isAuthPage) return;

	try {
		const res = await fetch(`${Glacier.base}/ping`, {
			credentials: 'include'
		});

		if (!res.ok) {
			throw redirect(307, '/auth');
		}

		if (isRoot) {
			throw redirect(307, '/library/');
		}
	} catch (e: any) {
		if (e.status && e.status >= 300 && e.status < 400) throw e;

		console.log('error fetching ping', e);
		throw redirect(307, '/auth');
	}
}
