// Apply the saved theme before the app paints. Dark remains the default.
try {
	const saved = localStorage.getItem('glacier-theme');
	document.documentElement.dataset.theme = saved === 'light' ? 'light' : 'dark';
} catch {
	document.documentElement.dataset.theme = 'dark';
}
