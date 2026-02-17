export const formatReleaseDate = (dateStr: string) => {
	if (!dateStr) return { formatted: 'Unknown', relative: '' };

	const date = new Date(dateStr);

	// Formatted: "23rd March 2023"
	const formatted = date.toLocaleDateString('en-GB', {
		day: 'numeric',
		month: 'long',
		year: 'numeric'
	});

	// Relative time: "2 years ago"
	const now = new Date();
	// https://github.com/microsoft/TypeScript/issues/58400
	// @ts-ignore
	const diffMs = now - date;
	const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

	let relative = '';
	if (diffDays < 1) relative = 'Today';
	else if (diffDays < 7) relative = `${diffDays} day${diffDays > 1 ? 's' : ''} ago`;
	else if (diffDays < 30)
		relative = `${Math.floor(diffDays / 7)} week${Math.floor(diffDays / 7) > 1 ? 's' : ''} ago`;
	else if (diffDays < 365)
		relative = `${Math.floor(diffDays / 30)} month${Math.floor(diffDays / 30) > 1 ? 's' : ''} ago`;
	else
		relative = `${Math.floor(diffDays / 365)} year${Math.floor(diffDays / 365) > 1 ? 's' : ''} ago`;

	return { formatted, relative };
};
