import type { GameMetadata } from '#lib/gen/search/v1/search_pb.js';

export const transferStore = $state({
	data: null as GameMetadata | null
});
