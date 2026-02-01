import { redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';
import { page } from '$app/state';
import { Glacier } from '$lib/api/api';

export async function load() {
	// Skip check for the login page itself
	if (browser && page.url.pathname.startsWith('/auth')) return;

	try {
		const res = await fetch(`${Glacier.base}/ping`);
		if (!res.ok) {
			throw redirect(307, '/auth');
		}
	} catch (e) {
		throw redirect(307, '/auth');
	}

	throw redirect(307, '/library/');
}
