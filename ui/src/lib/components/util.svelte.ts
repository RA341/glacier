export function extractValues(schema: any) {
	const result: any = {};
	for (const key in schema) {
		const item = schema[key];
		if (item.IsStruct && item.Nested) {
			result[key] = extractValues(item.Nested);
		} else {
			result[key] = item.Value;
		}
	}
	return result;
}
