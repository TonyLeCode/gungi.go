import { browser, dev } from '$app/environment';
import { getContext, setContext } from 'svelte';

type wsState = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';

type msgType<T> = {
	type: string;
	payload?: T;
};

//TODO gracefully try to reconnect
const maxAttempts = 4;

class WebSocketStore {
	state = $state<wsState>('connecting');
	websocket = $state<WebSocket>();
	isAuthenticated = $state(false);
	// reconnectionAttempts = $state(0);

	constructor() {
		if (!browser) return;
		this.connect(0);
	}

	connect(attempts: number) {
		const url = dev ? `ws://${import.meta.env.VITE_API_URL}/ws` : `wss://${import.meta.env.VITE_API_URL}/ws`;
		const newWebSocket = new WebSocket(url);

		newWebSocket.addEventListener('open', () => {
			this.websocket = newWebSocket;
			this.state = 'connected';
		});

		// Do not yet need error event listener

		// We don't add notifications here because the store should be separate from the notifications store
		newWebSocket.addEventListener('close', () => {
			if (attempts == 0) {
				console.log('Attempting reconnection...');
			}
			this.websocket?.close();
			this.isAuthenticated = false;
			this.attemptReconnection(attempts);
		});

		newWebSocket.addEventListener('message', (event) => {
			const data = JSON.parse(event.data);
			if (data?.type == 'auth') {
				if (data?.payload == 'success') {
					this.isAuthenticated = true;
				}
			}
		});
	}

	send(msg: msgType<unknown>) {
		if (!this.websocket) return;
		this.websocket.send(JSON.stringify(msg));
	}

	close() {
		this.isAuthenticated = false;
		if (!this.websocket) return;
		this.websocket.close();
	}

	addMsgListener(fn: (event?: MessageEvent) => void) {
		if (this.websocket === undefined) return;
		this.websocket.addEventListener('message', fn);
		return () => {
			if (this.websocket === undefined) return;
			this.websocket.removeEventListener('message', fn);
		};
	}

	removeMsgListener(fn: (event?: MessageEvent) => void) {
		if (this.websocket === undefined) return;
		this.websocket.removeEventListener('message', fn);
	}

	authenticate(token: string) {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		this.send(msg);
	}

	attemptReconnection(attempts: number) {
		if (this.state !== 'reconnecting') this.state = 'reconnecting';
		if (attempts < maxAttempts) {
			const timeout = 1200 * 2 ** attempts + Math.random() * 1000; // exponential backoff with jitter
			console.log('Reconnecting in ' + timeout + 'ms');
			setTimeout(() => {
				this.connect(attempts + 1);
			}, timeout);
		} else {
			console.error('Maximum reconnection attempts reached');
			this.state = 'closed';
			this.websocket = undefined;
		}
	}
}

export function setWebsocketStore() {
	const store = new WebSocketStore();
	setContext('websocket', store);
	return store;
}

export function getWebsocketStore() {
	return getContext<WebSocketStore>('websocket');
}
