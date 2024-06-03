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

		this.websocket.addEventListener('error', () => {
			this.state = 'error';
		});

		this.websocket.addEventListener('close', () => {
			this.state = 'closed';
			// const topNotificationStore = getTopNotificationStore();
			// topNotificationStore.SetNotification('You are disconnected! Please refresh or try again later.');
			console.error("disconnected from websocket");
		});

		this.websocket.addEventListener('message', (event) => {
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

	authenticate(token: string) {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		this.send(msg);
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
