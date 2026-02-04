import { type Client, ConnectError, createClient } from '@connectrpc/connect';
import { createConnectTransport } from '@connectrpc/connect-web';
import type { DescService } from '@bufbuild/protobuf';
import { browser } from '$app/environment';

const SERVICES = {
	GLACIER: '/api/server/protected',
	GLACIERPUB: '/api/server/public',
	FROST: '/api/frost'
} as const;

const SERVERS = {
	development: 'http://localhost:6699',
	frostdev: 'http://localhost:9966'
} as const;

const getBaseHost = () => {
	if (mode === 'development') return SERVERS.development;
	if (mode === 'frostdev') return SERVERS.frostdev;
	if (browser) {
		return window.location.origin;
	}
	return "/"
};

const buildBaseUrl = (servicePrefix: string) => {
	const host = getBaseHost();
	// Ensure no double slashes and no trailing slash
	return `${host.replace(/\/$/, '')}${servicePrefix}`;
};

function createServiceContext(prefix: string) {
	const baseUrl = buildBaseUrl(prefix);
	const transport = createConnectTransport({
		baseUrl: baseUrl,
		useBinaryFormat: true,
		interceptors: []
	});

	return {
		transport,
		base: baseUrl,
		get: <T extends DescService>(service: T): Client<T> => createClient(service, transport)
	};
}

const mode = import.meta.env.MODE;

const isFrostDev = mode === 'frostdev';
export const isFrost = mode === 'frost' || isFrostDev;
console.log(`Is frost ${isFrost}`);

export const appName = isFrost ? 'Frost' : 'Glacier';

export const Glacier = createServiceContext(SERVICES.GLACIER);
export const glacierCli = Glacier.get;

export const GlacierPub = createServiceContext(SERVICES.GLACIERPUB);
export const glacierPubCli = GlacierPub.get;

export const Frost = createServiceContext(SERVICES.FROST);
export const frostCli = Frost.get;

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
