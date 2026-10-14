export type Theme = 'dark' | 'light';
export const theme = $state<{ current: Theme }>({ current: 'dark' });
const storageKey = 'glacier-theme';

function applyTheme(value: Theme) {
	theme.current = value;
	document.documentElement.dataset.theme = value;
}

export function initializeTheme() {
	applyTheme(document.documentElement.dataset.theme === 'light' ? 'light' : 'dark');
	const onStorage = (event: StorageEvent) => {
		if (event.key === storageKey || event.key === null) {
			applyTheme(event.newValue === 'light' ? 'light' : 'dark');
		}
	};
	window.addEventListener('storage', onStorage);
	return () => window.removeEventListener('storage', onStorage);
}

export function toggleTheme() {
	applyTheme(theme.current === 'dark' ? 'light' : 'dark');
	try {
		localStorage.setItem(storageKey, theme.current);
	} catch {
		// Theme switching still works when browser storage is unavailable.
	}
}
