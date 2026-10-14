import { page } from '$app/state';
import { isFrost } from '#lib/api/api.js';
import { routes, type NavigationGuards } from './routes.js';

interface GuardedRoute {
	href: string;
	guard?: (guards: NavigationGuards) => boolean;
}

// Create per component/provider so user state is never shared across requests.
export function createNavigationStore(user: { readonly isOmni: boolean }) {
	const guards = $derived({ isOmni: user.isOmni, isFrost });
	const visible = (route: GuardedRoute) => !route.guard || route.guard(guards);
	const main = [routes.library, routes.search, routes.downloads];
	const footer = [routes.config, routes.profile];

	return {
		routes,
		get main() {
			return main;
		},
		get footer() {
			return footer.filter(visible);
		},
		isActive(href: string) {
			const path = page.url.pathname;
			return path === href || path.startsWith(`${href}/`);
		}
	};
}
