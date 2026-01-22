import { type Client, ConnectError, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import { createRPCRunner } from '$lib/api/svelte-api.svelte';
import type { DescService } from '@bufbuild/protobuf';
import { browser } from '$app/environment';

const mode = import.meta.env.MODE;

const devServer = 'http://localhost:6699';
const frostServer = 'http://localhost:9966';

const isFrostDev = mode === 'frostdev';
export const isFrost = mode === 'frost' || isFrostDev;
console.log(`Is frost ${isFrost}`);

const getBase = () => {
	if (mode === 'development') {
		return devServer;
	}

	if (isFrostDev) {
		return frostServer;
	}

	return '/';
};

export const HOST = getBase();
console.log("HOST", HOST);
////////////////////////////////////////////////////////////////////////////////////////////////////
// glacier client init
const GLACIER_API = "/api/server"
export const GLACIER_API_BASE = HOST === '/' ? GLACIER_API : `${HOST}${GLACIER_API}`;
console.log(`API url: ${GLACIER_API_BASE} `);

export const glacierTransport = createConnectTransport({
	baseUrl: GLACIER_API_BASE,
	useBinaryFormat: true,
	interceptors: []
});

export function glacierCli<T extends DescService>(service: T): Client<T> {
	return createClient(service, glacierTransport);
}

////////////////////////////////////////////////////////////////////////////////////////////////////
// frost client init

const FROST_API = `/api/frost`;
export const FROST_API_BASE = HOST === '/' ? FROST_API : `${HOST}${FROST_API}`;
console.log(`API url: ${FROST_API_BASE} `);

export const frostTransport = createConnectTransport({
	baseUrl: FROST_API_BASE,
	useBinaryFormat: true,
	interceptors: []
});

export function frostCli<T extends DescService>(service: T): Client<T> {
	return createClient(service, frostTransport);
}

////////////////////////////////////////////////////////////////////////////////////////////////////

export async function callRPC<T>(exec: () => Promise<T>): Promise<{ val: T | null; err: string }> {
	try {
		const val = await exec();
		return { val, err: '' };
	} catch (error: unknown) {
		if (error instanceof ConnectError) {
			console.error(`Error: ${error.message}`);
			// todo maybe ?????
			// if (error.code == Code.Unauthenticated) {
			//     nav("/")
			//
			return { val: null, err: `${error.rawMessage}` };
		}

		return { val: null, err: `Unknown error while calling api: ${(error as Error).toString()}` };
	}
}

export async function pingWithAuth() {
	try {
		console.log('Checking authentication status with server...');
		const response = await fetch('/ping', {
			redirect: 'follow'
		});

		if (response.status == 302) {
			const location = await response.text();
			console.log(`oidc is enabled redirecting to oidc auth: ${location}`);
			if (browser) {
				window.location.assign(location);
			}

			return false;
		}

		console.log(`Server response isOK: ${response.ok}`);
		return response.ok;
	} catch (error) {
		console.error('Authentication check failed:', error);
		return false;
	}
}

export function formatDate(timestamp: bigint | number | string) {
	const numericTimestamp =
		typeof timestamp === 'bigint'
			? // convert to ms from seconds
				Number(timestamp) * 1000
			: timestamp;
	return new Date(numericTimestamp).toLocaleDateString('en-US', {
		year: 'numeric',
		month: 'short',
		day: 'numeric',
		hour: '2-digit',
		minute: '2-digit'
	});
}
