/**
 * Connects to a WebSocket and handles updates for a process status.
 * @param url The websocket URL (e.g., "ws://localhost:8080/process/running?exe=app.exe")
 * @param onMessage Callback receiving the parsed JSON data
 * @param onError Callback receiving the error message string
 * @param onDone
 */
export function handleProcessStatus({
	url,
	onConnect,
	onDone,
	onMessage,
	onError
}: {
	url: string;
	onMessage?: (data: object) => void;
	onError?: (err: string) => void;
	onConnect?: () => void;
	onDone?: () => void;
}): () => void {
	const socket = new WebSocket(url);

	socket.onopen = (ev) => {
		onConnect?.();
	};

	socket.onmessage = (event: MessageEvent) => {
		try {
			console.log(event.data);

			// Parse the JSON string from the Go server into an object
			const data = JSON.parse(event.data);
			onMessage?.(data);
		} catch (e) {
			onError?.('Failed to parse message: ' + e);
		}
	};

	socket.onerror = (event: Event) => {
		onError?.('WebSocket error occurred');
		console.error('WS Error:', event);
	};

	socket.onclose = (event: CloseEvent) => {
		if (!event.wasClean) {
			onError?.(`Connection lost: ${event.reason || 'Unknown reason'}`);
		}

		onDone?.();
	};

	// Return a cleanup function to close the connection
	return () => {
		if (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING) {
			socket.close();
		}
	};
}
