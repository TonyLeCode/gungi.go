import { writable, readable } from 'svelte/store';
import { AddNotification, type notificationType } from './notification';
import { nanoid } from 'nanoid';
import { browser } from '$app/environment';

type wsConnStateType = 'connecting' | 'connected' | 'reconnecting' | 'closed' | 'error';
// type wsRouteType = 'overview' | 'game' | 'roomList' | null;
interface msgType {
	type: string;
	payload: unknown;
}

export const ws = writable<WebSocket>();
export const wsConnState = writable<wsConnStateType>('closed');
// export const wsRoute = writable<wsRouteType>();

// export function changeRoute(route: wsRouteType) {
// 	//TODO send route change message to server
// 	wsRoute.set(route);
// }

export function websocketConnect(url: string, token: string) {
	const newWS = new WebSocket(url);

	wsConnState.set('connecting');

	newWS.addEventListener('open', () => {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		// wsConnState.set('connected');
		newWS.send(JSON.stringify(msg));
	});

	newWS.addEventListener('message', (event) => {
		try {
			const data = JSON.parse(event.data);
			if (data.type == 'auth') {
				data.payload == '1' ? wsConnState.set('connected') : newWS.close();
			}
		} catch (err) {
			console.error('Error: ', err);
		}
	});

	newWS.addEventListener('error', (event) => {
		console.error('Error: ', event);
		wsConnState.set('error');
	});

	newWS.addEventListener('close', (event) => {
		console.log(event);
		wsConnState.set('closed');
		AddNotification({
			id: nanoid(),
			title: 'Disconnected',
			type: 'default',
			msg: 'Please refresh or come back later',
		} as notificationType);
	});

	ws.set(newWS);
}

function createWsStore() {
	if (!browser) {
		return;
	}
	const newSocket = new WebSocket(`ws://${import.meta.env.VITE_API_URL}/ws`);
	const { subscribe } = readable<wsConnStateType>('closed', (set) => {
		newSocket.addEventListener('open', () => {
			set('connecting');
		});
		newSocket.addEventListener('error', () => {
			set('error');
		});
		newSocket.addEventListener('close', () => {
			set('closed');
			AddNotification({
				id: nanoid(),
				title: 'Disconnected',
				type: 'default',
				msg: 'Please refresh or come back later',
			} as notificationType);
		});
		newSocket.addEventListener('message', (event) => {
			try {
				const data = JSON.parse(event.data);
				if (data.type == 'auth') {
					data.payload == '1' ? set('connected') : newSocket.close();
				}
			} catch (err) {
				console.error('Error: ', err);
			}
		});
	});

	function send(msg: msgType) {
		newSocket.send(JSON.stringify(msg));
	}

	function authenticate(token: string) {
		const msg = {
			type: 'auth',
			payload: `Bearer ${token}`,
		};
		send(msg);
	}

	function close() {
		newSocket.close();
	}

	function addMsgListener(fn: (event?: MessageEvent) => void) {
		newSocket.addEventListener('message', fn);
	}
	function removeMsgListener(fn: (event?: MessageEvent) => void) {
		newSocket.removeEventListener('message', fn);
	}

	console.log('create ws store');
	return {
		subscribe,
		send,
		authenticate,
		close,
		addMsgListener,
		removeMsgListener,
	};
}

export const df = createWsStore();
