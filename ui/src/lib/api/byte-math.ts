export function calculatePercent(complete: any, left: any): number {
	const c = Number(complete || 0);
	const l = Number(left || 0);
	const total = c + l;
	if (total === 0) return 0;
	return Math.min(100, Math.max(0, (c / total) * 100));
}

export function formatBytes(bytes: number | bigint, decimals = 2) {
	const b = BigInt(bytes);
	if (b === 0n) return '0 B';

	const k = 1024;
	const dm = decimals < 0 ? 0 : decimals;
	const sizes = ['B', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];

	const i = Math.floor(Math.log(Number(b)) / Math.log(k));
	const index = Math.min(i, sizes.length - 1);
	const val = Number(b) / Math.pow(k, index);

	return `${val.toFixed(dm)} ${sizes[index]}`;
}
