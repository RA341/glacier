export function trimPrefix(toTrim: string, trim: string): string {
	if (!toTrim || !trim) {
		return toTrim;
	}
	const index = toTrim.indexOf(trim);
	if (index !== 0) {
		return toTrim;
	}
	return toTrim.substring(trim.length);
}

export function splitCamelCase(str: string): string {
	return str.replace(/([A-Z])/g, (match, letter, offset, original) =>
		offset === 0 || /[A-Z]/.test(original[offset - 1]) ? letter : ' ' + letter
	);
}
