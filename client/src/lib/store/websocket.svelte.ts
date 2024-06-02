import { browser, dev } from '$app/environment';
import { getContext, setContext } from 'svelte';

type wsState = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';

type msgType<T> = {
	type: string;
	payload?: T;
};

class WebSocketStore {
	state = $state<wsState>('connecting');
	websocket = $state<WebSocket>();
	isAuthenticated = $state(false);
	constructor() {
		if (!browser) return;
		const url = dev ? `ws://${import.meta.env.VITE_API_URL}/ws` : `wss://${import.meta.env.VITE_API_URL}/ws`;
		this.websocket = new WebSocket(url);

		this.websocket.addEventListener('open', () => {
			this.state = 'connected';
		});
	}

	send(msg: msgType<unknown>) {
		if (!this.websocket) return;
		this.websocket.send(JSON.stringify(msg));
	}

	close() {
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
}

export function setWebsocketStore() {
	const store = new WebSocketStore();
	setContext('websocket', store);
	return store;
}

export function getWebsocketStore() {
	return getContext<WebSocketStore>('websocket');
}
