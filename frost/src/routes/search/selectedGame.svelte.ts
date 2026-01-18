import type { GameMetadata } from '$lib/gen/search/v1/search_pb';

export const transferStore = $state({
	data: null as GameMetadata | null
});