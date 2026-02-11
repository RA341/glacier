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
