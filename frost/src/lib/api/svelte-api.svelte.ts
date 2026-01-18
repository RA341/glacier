
export function createRPCRunner<T>(exec: () => Promise<T>) {
	let loading = $state(false);
	let error = $state('');
	let value = $state<T | null>(null);

	const runner = async () => {
		loading = true;
		error = '';

		try {
			value = await exec();
		} catch (e: any) {
			value = null;
			console.error(e.message || 'An error occurred');
			error = e.message || 'An error occurred';
		} finally {
			loading = false;
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
