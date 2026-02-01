import { Glacier } from '$lib/api/api';


export function checkAuth() {
	let loading = $state(false);
	let error = $state('');
	let value = $state<boolean>(false);

	const runner = async () => {
		loading = true;
		error = '';

		try {
			const sd = await fetch(`${Glacier.base}/ping`);
			value = sd.ok;
		} catch (e: any) {
			value = false;
			console.error(e.message || 'An error occurred');
			error = e.message || 'An error occurred';
		} finally {
			loading = false;

			console.log(`Auth status ${value}`);
		}
	};

	return {
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get value() {
			return value;
		},
		runner
	};
}
