import { DownloadIcon, LibraryIcon, SettingsIcon, SearchIcon, UserIcon } from '@lucide/svelte';

export interface NavigationGuards {
	isOmni: boolean;
	isFrost: boolean;
}

const hasConfig = ({ isOmni, isFrost }: NavigationGuards) => isOmni || isFrost;

// Guards control navigation visibility; authorization remains the server's responsibility.
export const routes = {
	library: { label: 'Library', href: '/library', icon: LibraryIcon },
	search: { label: 'Search', href: '/search', icon: SearchIcon },
	downloads: { label: 'Downloads', href: '/downloads', icon: DownloadIcon },
	config: { label: 'Config', href: '/settings', icon: SettingsIcon, guard: hasConfig },
	profile: { label: 'Profile', href: '/profile', icon: UserIcon }
} as const;
