import { Glacier } from '$lib/api/api';
import * as sea from 'node:sea';

interface UploadAssetParams {
	id: bigint;
	gameId: bigint;
	type: string;
	file: File;
}

// todo maybe use use frost base on frost clietns
export function getAssetPath({
	assetType,
	assetPath,
	gameId
}: {
	gameId?: bigint;
	assetPath?: string;
	assetType?: string;
}): string {
	const assetBase = `${Glacier.base}/assets/${gameId}`;

	if (assetType) {
		return `${assetBase}/t/${assetType}`;
	}

	return `${assetBase}/${assetPath}`;
}

export async function uploadAsset({
	file,
	type,
	id,
	gameId
}: UploadAssetParams): Promise<string | null> {
	try {
		const formData = new FormData();

		formData.append('file', file);
		formData.append('type', type);
		formData.append('id', id.toString());
		formData.append('gameId', gameId.toString());

		const response = await fetch(`${Glacier.base}/assets/upload`, {
			method: 'POST',
			body: formData
		});

		if (!response.ok) {
			const errorText = await response.text();
			return errorText || `Upload failed with status ${response.status}`;
		}

		return null;
	} catch (err) {
		if (err instanceof Error) {
			return err.message;
		}
		return 'An unknown error occurred during upload';
	}
}
