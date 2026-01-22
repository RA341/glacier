import { redirect } from '@sveltejs/kit';
import { isFrost } from '$lib/api/api';

export function load() {
	throw redirect(307, isFrost ? '/downloads/frost' : '/downloads/glacier');
}
